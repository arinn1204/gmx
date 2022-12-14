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
	client.classHandlers.Store(typeName, handler)

	client.mbeans.Range(func(key, value any) bool {
		value.(mbean.BeanExecutor).RegisterClassHandler(typeName, handler)
		return true
	})
}

// RegisterInterfaceHandler is the method to use when wanting to register
// additional handlers that will blanket apply to a Java interface
func (client *client) RegisterInterfaceHandler(typeName string, handler extensions.InterfaceHandler) {
	client.interfaceHandlers.Store(typeName, handler)

	client.mbeans.Range(func(key, value any) bool {
		value.(mbean.BeanExecutor).RegisterInterfaceHandler(typeName, handler)
		return true
	})
}

func (client *client) registerNewBean(id uuid.UUID, bean mbean.BeanExecutor) {
	client.classHandlers.Range(func(key, value any) bool {
		typeName := key.(string)
		handler := value.(extensions.IHandler)
		bean.RegisterClassHandler(typeName, handler)
		return true
	})

	client.interfaceHandlers.Range(func(key, value any) bool {
		typeName := key.(string)
		handler := value.(extensions.InterfaceHandler)
		bean.RegisterInterfaceHandler(typeName, handler)
		return true
	})

	client.inc()
	client.mbeans.Store(id, bean)
}

// Initialize is the initial method to create a GMX client.
// This will initialize the JVM if necessary as well as setting up the object
func (client *client) Initialize() error {
	startJvm()

	client.mbeans = sync.Map{}
	client.classHandlers = sync.Map{}
	client.interfaceHandlers = sync.Map{}

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

// RegisterBean is the initializing method for the MBean itself. It will store the
// address provided and register the bean with the client to be executed later.
func (client *client) RegisterConnection(hostname string, port int) (*uuid.UUID, error) {
	jmxURI := fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%d/jmxrmi", hostname, port)

	bean := &mbean.Client{
		JmxURI:            jmxURI,
		ClassHandlers:     sync.Map{},
		InterfaceHandlers: sync.Map{},
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
