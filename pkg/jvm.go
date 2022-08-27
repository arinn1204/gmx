package gmx

import (
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

// RegisterClassHandler is the method to use when wanting to register additional handlers
// By default this client will handle everything in internal/handlers
func (client *client) RegisterClassHandler(typeName string, handler extensions.IHandler) {
	client.classHandlers[typeName] = handler

	client.mbeans.Range(func(key, value any) bool {
		value.(mbean.BeanExecutor).RegisterClassHandler(typeName, handler)
		return true
	})
}

// RegisterInterfaceHandler is the method to use when wanting to register
// additional handlers that will blanket apply to a Java interface
func (client *client) RegisterInterfaceHandler(typeName string, handler extensions.InterfaceHandler) {
	client.interfaceHandlers[typeName] = handler

	client.mbeans.Range(func(key, value any) bool {
		value.(mbean.BeanExecutor).RegisterInterfaceHandler(typeName, handler)
		return true
	})
}

func (client *client) registerNewBean(id uuid.UUID, bean mbean.BeanExecutor) {
	for typeName, handler := range client.classHandlers {
		bean.RegisterClassHandler(typeName, handler)
	}

	for typeName, handler := range client.interfaceHandlers {
		bean.RegisterInterfaceHandler(typeName, handler)
	}

	client.inc()
	client.mbeans.Store(id, bean)
}

// Initialize is the initial method to create a GMX client.
// This will initialize the JVM if necessary as well as setting up the object
func (client *client) Initialize() error {
	startJvm()

	client.mbeans = sync.Map{}
	client.classHandlers = make(map[string]extensions.IHandler)
	client.interfaceHandlers = make(map[string]extensions.InterfaceHandler)

	client.RegisterClassHandler(handlers.BoolClasspath, &handlers.BoolHandler{})
	client.RegisterClassHandler(handlers.DoubleClasspath, &handlers.DoubleHandler{})
	client.RegisterClassHandler(handlers.FloatClasspath, &handlers.FloatHandler{})
	client.RegisterClassHandler(handlers.IntClasspath, &handlers.IntHandler{})
	client.RegisterClassHandler(handlers.LongClasspath, &handlers.LongHandler{})
	client.RegisterClassHandler(handlers.StringClasspath, &handlers.StringHandler{})

	client.RegisterInterfaceHandler(handlers.ListClassPath, &handlers.ListHandler{ClassHandlers: &client.classHandlers, InterfaceHandlers: &client.interfaceHandlers})
	client.RegisterInterfaceHandler(handlers.SetClassPath, &handlers.SetHandler{ClassHandlers: &client.classHandlers, InterfaceHandlers: &client.interfaceHandlers})
	client.RegisterInterfaceHandler(handlers.MapClassPath, &handlers.MapHandler{ClassHandlers: &client.classHandlers, InterfaceHandlers: &client.interfaceHandlers})

	return nil
}

func (client *client) inc() {
	lock.Lock()
	client.numberOfConnections++
	lock.Unlock()
}

func (client *client) dec() {
	lock.Lock()
	client.numberOfConnections--
	lock.Unlock()
}

// Connect is the initializing method for the MBean itself. It will
// connect to the remote server and assign the given connection a UUID.
// The GMX client will store references to MBean clients, the UUID's will be
// helpful if wanting to be able to tell which MBeans go to which location
func (client *client) Connect(hostname string, port int) (*uuid.UUID, error) {
	jmxURI := fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%d/jmxrmi", hostname, port)

	bean := &mbean.Client{
		JmxURI:            jmxURI,
		ClassHandlers:     make(map[string]extensions.IHandler),
		InterfaceHandlers: make(map[string]extensions.InterfaceHandler),
	}

	id := uuid.New()
	client.registerNewBean(id, bean)

	return &id, nil
}

// Close is a method that will close the connection. It will free up any resources
// that the GMX client is still holding on as well as shutting down the JVM
func (client *client) Close() {

	client.mbeans.Range(func(key, value any) bool {
		client.mbeans.Delete(key)
		client.dec()
		return true

	})

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
