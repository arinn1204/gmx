package gmx

import (
	"log"
	"sync"

	"github.com/arinn1204/gmx/internal/mbean"
	"github.com/arinn1204/gmx/pkg/extensions"

	"github.com/google/uuid"
)

// Client is the main mbean client.
// This is responsible for creating the JVM, creating individual MBean Clients, and cleaning it all up
// The client is also responsible for orchestrating the JMX operations
type client struct {
	maxNumberOfGoRoutines uint                             // The maximum number of goroutines to be used when doing parallel operations
	mbeans                map[uuid.UUID]mbean.BeanExecutor // The map of underlying clients. The map is identified as id -> client
	classHandlers         map[string]extensions.IHandler   // The map of type handlers to be used
	interfaceHandlers     map[string]extensions.InterfaceHandler
}

type attributeManager struct {
	maxNumberOfGoRoutines uint
	mbeans                *map[uuid.UUID]mbean.BeanExecutor
	classHandlers         *map[string]extensions.IHandler // The map of type handlers to be used
	interfaceHandlers     *map[string]extensions.InterfaceHandler
}

type operator struct {
	maxNumberOfGoRoutines uint
	mbeans                *map[uuid.UUID]mbean.BeanExecutor
	classHandlers         *map[string]extensions.IHandler // The map of type handlers to be used
	interfaceHandlers     *map[string]extensions.InterfaceHandler
}

// MBeanClient is an interface that describes the functions needed to fully operate against MBeans over JMXRMI
type MBeanClient interface {
	// Initialize will initialize the client:
	// This starts the JVM and registers all basic class and interface handlers
	// Basic handlers are: Integer, Double, String, Float, Boolean, Long, List, Set, Map<String, Object>
	Initialize() error

	// Close will close all connections that have been created and then shut down the JVM
	Close()

	// RegisterClassHandler will register a new class handler
	// This is needed if executing operations or retrieving/updating attributes
	// that require more complex objects. The default handlers only include primitives and strings
	RegisterClassHandler(typeName string, handler extensions.IHandler)

	// RegisterInterfaceHandler will register a new interface handler
	// In the event no class handler is found for a given type, we will then scan the available
	// interfaces implemented by the type. This will then check all the included handlers to determine
	// how to convert between go and jni
	RegisterInterfaceHandler(typeName string, handler extensions.InterfaceHandler)

	// Connect will create a new mbean connection defined by the hostname and port
	// The reference to this connection is stored for the life of the operator
	Connect(hostname string, port int) (*uuid.UUID, error)

	// This will return the type that is responsible for executing operations
	GetOperator() MBeanOperator

	// This will return the attribute manager
	GetAttributeManager() MBeanAttributeManager
}

// MBeanAttributeManager is a type that will be responsible for managing attributes of a given mbean
type MBeanAttributeManager interface {

	// Get will fetch an attribute by a given name for the given bean across all connections
	// The args are required in order to be able to deserialize lists from attributes
	//
	// JavaType and JavaContainerType are only required for Lists/Sets/Maps and any other generic collections
	// Value is not used here
	Get(domain string, beanName string, attributeName string, args ...MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error)

	// Put will change the given attribute across all connections
	// The args are required in order to be able to serialize data being sent to the attribute
	Put(domain string, beanName string, attributeName string, args ...MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error)

	// GetById will execute Get against the one given ID.
	// The args are required in order to be able to deserialize lists from attributes
	//
	// JavaType and JavaContainerType are only required for Lists/Sets/Maps and any other generic collections
	// Value is not used here
	GetById(id uuid.UUID, domain string, beanName string, attributeName string, args ...MBeanArgs) (string, error)

	// PutById will execute Put against the one given ID.
	// The args are required in order to be able to serialize data being sent to the attribute
	PutById(id uuid.UUID, domain string, beanName string, attributeName string, args ...MBeanArgs) (string, error)
}

// MBeanOperator is a type that is responsible for executing operations against a defined mbean
type MBeanOperator interface {

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

// CreateClient is a method that will create an unbound MBeanClient
// This means that it will consume as many native threads as there are connected mbeans
func CreateClient() MBeanClient {
	client := &client{}
	if err := client.Initialize(); err != nil {
		log.Fatal(err)
	}

	return client
}

// CreateClient is a method that will create a bound MBeanClient
// This will only use as many native threads as provided by limit
func CreateClientWithLimitation(limit uint) MBeanClient {
	client := &client{
		maxNumberOfGoRoutines: limit,
	}

	if err := client.Initialize(); err != nil {
		log.Fatal(err)
	}

	return client
}

func internalExecuteAgainstAll(mbeans *map[uuid.UUID]mbean.BeanExecutor, maxNumberOfGoRoutines uint, execution func(uuid.UUID) (string, error)) (map[uuid.UUID]string, map[uuid.UUID]error) {
	results := make(chan batchExecutionResult, len(*mbeans))
	wg := &sync.WaitGroup{}

	maxLimitOfConcurrency := len(*mbeans)

	if maxNumberOfGoRoutines > 0 {
		maxLimitOfConcurrency = int(maxNumberOfGoRoutines)
	}

	guard := make(chan struct{}, maxLimitOfConcurrency)

	for id := range *mbeans {

		// this will block if there are no available threads
		guard <- struct{}{}
		wg.Add(1)
		go func(id uuid.UUID) {
			defer wg.Done()
			res, err := execution(id)
			result := batchExecutionResult{
				id:     id,
				result: res,
				err:    err,
			}

			results <- result
			<-guard // this will release this threads guard
		}(id)

	}

	wg.Wait()
	result := make(map[uuid.UUID]string)
	receivedErrors := make(map[uuid.UUID]error)

	for i := 0; i < len(*mbeans); i++ {
		res := <-results

		result[res.id] = res.result
		receivedErrors[res.id] = res.err
	}

	return result, receivedErrors
}

// GetOperator is a function that returns an mbean operator
// This will be used to execute operations against any given mbean that is
// registered with the client
func (client *client) GetOperator() MBeanOperator {
	return &operator{
		maxNumberOfGoRoutines: client.maxNumberOfGoRoutines,
		mbeans:                &client.mbeans,
		classHandlers:         &client.classHandlers,
		interfaceHandlers:     &client.interfaceHandlers,
	}
}

// GetAttributeManager is a function that returns an mbean attribute manager
// This will be used to read and update attributes for any/all mbeans that are associated
// with the client that creates the manager
func (client *client) GetAttributeManager() MBeanAttributeManager {
	return &attributeManager{
		maxNumberOfGoRoutines: client.maxNumberOfGoRoutines,
		mbeans:                &client.mbeans,
		classHandlers:         &client.classHandlers,
		interfaceHandlers:     &client.interfaceHandlers,
	}
}
