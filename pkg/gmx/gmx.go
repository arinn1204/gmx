package gmx

import (
	"gmx/internal/mbean"

	"github.com/google/uuid"
)

// Client is the main mbean client.
// This is responsible for creating the JVM, creating individual MBean Clients, and cleaning it all up
// The client is also responsible for orchestrating the JMX operations
type Client struct {
	mbeans map[uuid.UUID]*mbean.Client // The map of underlying clients. The map is identified as id -> client
}

// MBeanOperator is an interface that describes the functions needed to fully operate against MBeans over JMXRMI
type MBeanOperator interface {
	Initialize() (*Client, error)                         // This will initialize the JVM if needed (only once) and an MBean connection
	Close()                                               // This will close out the JVM and free up any clients that are remaining
	Connect(hostname string, port int) (uuid.UUID, error) // This will initialize a new MBean connection
}
