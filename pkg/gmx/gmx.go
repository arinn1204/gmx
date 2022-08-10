package gmx

import (
	"gmx/internal/mbean"
)

// Client is the main mbean client.
// This is responsible for creating the JVM, creating individual MBean Clients, and cleaning it all up
// The client is also responsible for orchestrating out the JMX operations
type Client struct {
	Hostname string          // The hostname/ip address of the JMXRMI server
	Port     int             // The port that RMI is configured to listen on
	client   []*mbean.Client // The array of underying clients, each new client that gets created will store its reference here
}

// MBeanOperator is an interface that describes the functions needed to fully operate against MBeans over JMXRMI
type MBeanOperator interface {
	Initialize() (*Client, error)      // This will initialize the JVM if needed (only once) and an MBean connection
	Close()                            // This will close out the JVM and free up any clients that are remaining
	Connect(hostname string, port int) // This will initialize a new MBean connection
}
