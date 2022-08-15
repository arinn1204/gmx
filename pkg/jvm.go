package gmx

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/arinn1204/gmx/internal/mbean"
	"github.com/arinn1204/gmx/pkg/extensions"

	"github.com/arinn1204/gmx/internal/jvm"

	"github.com/google/uuid"
)

var java jvm.IJava
var lock *sync.Mutex

func init() {
	lock = &sync.Mutex{}
	java = &jvm.Java{}
}

// RegisterHandler is the method to use when wanting to register additional handlers
// By default this client will handle everything in internal/handlers
func (client *Client) RegisterHandler(typeName string, handler extensions.IHandler) {
	client.handlers[typeName] = handler

	for _, bean := range client.mbeans {
		bean.RegisterClassHandler(typeName, handler)
	}
}

func (client *Client) registerNewBean(id uuid.UUID, bean mbean.BeanExecutor) {
	for typeName, handler := range client.handlers {
		bean.RegisterClassHandler(typeName, handler)
	}

	client.mbeans[id] = bean
}

// Initialize is the initial method to create a GMX client.
// This will initialize the JVM if necessary as well as setting up the object
func (client *Client) Initialize() error {
	startJvm()

	client.mbeans = make(map[uuid.UUID]mbean.BeanExecutor)
	client.handlers = make(map[string]extensions.IHandler)

	client.RegisterHandler(handlers.BoolClasspath, &handlers.BoolHandler{})
	client.RegisterHandler(handlers.DoubleClasspath, &handlers.DoubleHandler{})
	client.RegisterHandler(handlers.FloatClasspath, &handlers.FloatHandler{})
	client.RegisterHandler(handlers.IntClasspath, &handlers.IntHandler{})
	client.RegisterHandler(handlers.LongClasspath, &handlers.LongHandler{})
	client.RegisterHandler(handlers.StringClasspath, &handlers.StringHandler{})

	return nil
}

// Connect is the initializing method for the MBean itself. It will
// connect to the remote server and assign the given connection a UUID.
// The GMX client will store references to MBean clients, the UUID's will be
// helpful if wanting to be able to tell which MBeans go to which location
func (client *Client) Connect(hostname string, port int) (*uuid.UUID, error) {
	jmxURI := fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%d/jmxrmi", hostname, port)
	bean, err := java.CreateMBeanConnection(jmxURI)

	if err != nil {
		return nil, errors.New("failed to create a connection::" + err.Error())
	}

	id := uuid.New()

	client.registerNewBean(id, bean)

	return &id, nil
}

// Close is a method that will close the connection. It will free up any resources
// that the GMX client is still holding on as well as shutting down the JVM
func (client *Client) Close() {
	for uri := range client.mbeans {
		client.mbeans[uri].Close()
		client.mbeans[uri] = nil
	}

	java.ShutdownJvm()
}

func startJvm() {
	if java.IsStarted() {
		return
	}

	lock.Lock()
	defer lock.Unlock()
	if java.IsStarted() {
		return
	}

	var err error
	java, err = java.CreateJVM()

	if err != nil {
		log.Fatalf("Failed to create the JVM::" + err.Error())
	}
}
