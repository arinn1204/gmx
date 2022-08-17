package gmx

import (
	"sync"

	"github.com/arinn1204/gmx/internal/mbean"
	"github.com/arinn1204/gmx/pkg/extensions"

	"github.com/google/uuid"
)

// Client is the main mbean client.
// This is responsible for creating the JVM, creating individual MBean Clients, and cleaning it all up
// The client is also responsible for orchestrating the JMX operations
type Client struct {
	MaxNumberOfGoRoutines uint                             // The maximum number of goroutines to be used when doing parallel operations
	mbeans                map[uuid.UUID]mbean.BeanExecutor // The map of underlying clients. The map is identified as id -> client
	classHandlers         map[string]extensions.IHandler   // The map of type handlers to be used
	interfaceHandlers     map[string]extensions.InterfaceHandler
}

type Handler interface {
	RegisterClassHandler(typeName string, handler extensions.IHandler)
	RegisterInterfaceHandler(typeName string, handler extensions.InterfaceHandler)
}

type AttributeHandler interface {
	// Get will fetch an attribute by a given name for the given bean across all connections
	Get(domain string, beanName string, attributeName string) (string, error)

	// Put will change the given attribute across all connections
	Put(domain string, beanName string, attributeName string, value any) (string, error)

	// GetById will execute Get against the one given ID.
	GetById(id uuid.UUID, domain string, beanName string, attributeName string) (string, error)

	// PutById will execute Put against the one given ID.
	PutById(id uuid.UUID, domain string, beanName string, attributeName string, value any) (string, error)
}

// MBeanOperator is an interface that describes the functions needed to fully operate against MBeans over JMXRMI
type MBeanOperator interface {
	// Initialize will initialize the client:
	// This starts the JVM and registers all basic class and interface handlers
	// Basic handlers are: Integer, Double, String, Float, Boolean, Long, List, Set, Map<String, Object>
	Initialize() (*Client, error)

	// Close will close all connections that have been created and then shut down the JVM
	Close()

	// Connect will create a new mbean connection defined by the hostname and port
	// The reference to this connection is stored for the life of the operator
	Connect(hostname string, port int) (*uuid.UUID, error)

	// ExecuteAgainstID will execute an operation against the given id. This will only target the provided ID
	ExecuteAgainstID(id uuid.UUID, domain string, name string, operation string, args ...MBeanArgs) (string, error)

	// ExecuteAgainstAll will execute an operation against *all* connected beans.
	// These are ran in their own go routines. If there concerns/desired constraints please define MaxNumberOfGoRoutines
	ExecuteAgainstAll(domain string, name string, operation string, args ...MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error)
}

// MBeanArgs is the container for the operation arguments.
// If you have an MBean defined as the following
//
//	getValue(String name)
//
// Then you will want to structure your args like
//
//	  MBeanArgs{
//	    Value: "theNameOfTheStringYouAreFetching",
//		JavaType: "java.lang.String"
//	  }
//
// If the intent is to execute a command that accepts a list or a map, then the JavaContainerType will need to be defined
// The JavaType will be the inner type that is being sent, whereas the container type will be the container.
// For example:
//
//	  MBeanArgs{
//	    Value: "[\"foo\", \"bar\"]",
//		JavaType: "java.lang.String",
//		JavaContainerType: "java.util.List"
//	  }
type MBeanArgs struct {
	Value             string
	JavaType          string
	JavaContainerType string
}

type batchExecutionResult struct {
	id     uuid.UUID
	result string
	err    error
}

// ExecuteAgainstAll will execute a single command against every mbean that is currently registered.
// This will return a mapping of all results and errors, based on the UUID that the connection has been assigned.
//
// All executions will be run in separate go routines, so this needs to be planned for accordingly
func (client *Client) ExecuteAgainstAll(domain string, name string, operation string, args ...MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error) {
	results := make(chan batchExecutionResult, len(client.mbeans))
	wg := &sync.WaitGroup{}

	for id := range client.mbeans {

		wg.Add(1)
		go func(id uuid.UUID) {
			defer wg.Done()
			res, err := client.ExecuteAgainstID(id, domain, name, operation, args...)
			result := batchExecutionResult{
				id:     id,
				result: res,
				err:    err,
			}

			results <- result

		}(id)

	}

	wg.Wait()

	result := make(map[uuid.UUID]string)
	receivedErrors := make(map[uuid.UUID]error)

	for i := 0; i < len(client.mbeans); i++ {
		res := <-results

		result[res.id] = res.result
		receivedErrors[res.id] = res.err
	}

	return result, receivedErrors
}

// ExecuteAgainstID is a method that will take a given operation and MBean ID and make the JMX request.
// It will return whatever is returned downstream, errors and all
func (client *Client) ExecuteAgainstID(id uuid.UUID, domain string, name string, operation string, args ...MBeanArgs) (string, error) {
	env := java.Attach()
	defer java.Detach()

	bean := client.mbeans[id].WithEnvironment(env)

	operationArgs := make([]mbean.OperationArgs, 0)

	for _, arg := range args {
		operationArgs = append(
			operationArgs,
			mbean.OperationArgs{
				Value:             arg.Value,
				JavaType:          arg.JavaType,
				JavaContainerType: arg.JavaContainerType,
			},
		)
	}

	mbeanOp := mbean.Operation{
		Domain:    domain,
		Name:      name,
		Operation: operation,
		Args:      operationArgs,
	}

	return bean.Execute(mbeanOp)
}
