package gmx

import (
	"github.com/arinn1204/gmx/internal/mbean"

	"github.com/google/uuid"
)

// Client is the main mbean client.
// This is responsible for creating the JVM, creating individual MBean Clients, and cleaning it all up
// The client is also responsible for orchestrating the JMX operations
type Client struct {
	mbeans map[uuid.UUID]mbean.BeanExecutor // The map of underlying clients. The map is identified as id -> client
}

// MBeanOperator is an interface that describes the functions needed to fully operate against MBeans over JMXRMI
type MBeanOperator interface {
	// This will initialize the JVM if needed (only once) and an MBean connection
	Initialize() (*Client, error)
	// This will close out the JVM and free up any clients that are remaining
	Close()
	// This will initialize a new MBean connection
	Connect(hostname string, port int) (*uuid.UUID, error)
	// This will execute the given operation against every MBean that has already been created.
	// It will return a mapping of results and errors based on the ID of the bean
	ExecuteAgainstAll(domain string, name string, operation string, args ...MBeanArgs) (map[uuid.UUID]any, map[uuid.UUID]error)
	// This will execute the given operation against the spefied bean
	ExecuteAgainstID(id uuid.UUID, domain string, name string, operation string, args ...MBeanArgs) (any, error)
}

// MBeanArgs is the container for the operation arguments.
// If you have an MBean defined as the following
//
//	getValue(String name)
//
// then you will want to structure your args like
//
//	  MBeanArgs{
//	    Value: "theNameOfTheStringYouAreFetching",
//		JavaType: "java.lang.String"
//	  }
type MBeanArgs struct {
	Value    any
	JavaType string
}

// ExecuteAgainstAll will execute a single command against every mbean that is currently registered.
// This will return a mapping of all results and errors, based on the UUID that the connection has been assigned.
func (client *Client) ExecuteAgainstAll(domain string, name string, operation string, args ...MBeanArgs) (map[uuid.UUID]any, map[uuid.UUID]error) {
	result := make(map[uuid.UUID]any)
	receivedErrors := make(map[uuid.UUID]error)

	for id := range client.mbeans {
		res, err := client.ExecuteAgainstID(id, domain, name, operation, args...)
		result[id] = res
		receivedErrors[id] = err
	}

	return result, receivedErrors
}

// ExecuteAgainstID is a method that will take a given operation and MBean ID and make the JMX request.
// It will return whatever is returned downstream, errors and all
func (client *Client) ExecuteAgainstID(id uuid.UUID, domain string, name string, operation string, args ...MBeanArgs) (any, error) {
	bean := client.mbeans[id]

	operationArgs := make([]mbean.OperationArgs, 0)

	for _, arg := range args {
		operationArgs = append(
			operationArgs,
			mbean.OperationArgs{
				Value: arg.Value,
				Type:  arg.JavaType,
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