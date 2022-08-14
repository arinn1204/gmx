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
	mbeans   map[uuid.UUID]mbean.BeanExecutor // The map of underlying clients. The map is identified as id -> client
	handlers map[string]extensions.IHandler   // The map of type handlers to be used
}

// MBeanOperator is an interface that describes the functions needed to fully operate against MBeans over JMXRMI
type MBeanOperator interface {
	// This will initialize the JVM if needed (only once) and an MBean connection
	Initialize() (*Client, error)
	// This will close out the JVM and free up any clients that are remaining
	Close()
	RegisterHandler(typeName string, handler extensions.IHandler) // This will register additional handlers
	// This will initialize a new MBean connection
	Connect(hostname string, port int) (*uuid.UUID, error)
	// This will execute the given operation against every MBean that has already been created.
	// It will return a mapping of results and errors based on the ID of the bean
	ExecuteAgainstAll(domain string, name string, operation string, args ...MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error)
	// This will execute the given operation against the spefied bean
	ExecuteAgainstID(id uuid.UUID, domain string, name string, operation string, args ...MBeanArgs) (string, error)
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
