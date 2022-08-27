package mbean

import (
	"errors"
	"fmt"

	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

// a collection of JNI representations of java primitive types
// these will be the boxed representations, not the true primitive
const (
	STRING = "java/lang/String"
	OBJECT = "java/lang/Object"
	FLOAT  = "java/lang/Float"
)

// Client is the overarching type that will facilitate JMX connections
// JmxConnection is the living connection that was created when `CreateMBeanConnection` was called to create the Client
// Env is the environment that belongs to this bean, this will not always match the JVM env!
type Client struct {
	JmxURI            string
	Env               *jnigi.Env
	ClassHandlers     map[string]extensions.IHandler
	InterfaceHandlers map[string]extensions.InterfaceHandler
}

// Operation is the operation that is being performed
// Domain is the fully qualified name of the MBean `org.example`
// Name is the name of the mbean itself `game`
// Operation is the name of the operation that is attempted to be interacted with `getString`
// Args are the optional argument array that is for the operation
type Operation struct {
	Domain    string
	Name      string
	Operation string
	Args      []OperationArgs
}

// OperationArgs is the type that holds data about the arguments used for MBean operations
// Value is the value that is being entered in string form
// JavaType is the fully qualified java type `java.lang.String`
// JavaContainerType is the fully qualified type of the container that will be holding JavaType
type OperationArgs struct {
	Value             string
	JavaType          string
	JavaContainerType string
}

// BeanExecutor is the interface used around this package.
// This is how an execution is performed.
// This will always rely on the MBean Client's environment
// Close will handle any types of cleanup that is related directly to MBean operations
//
// for example: cleaning up the JMX connection and deleting the reference
type BeanExecutor interface {
	RegisterClassHandler(typeName string, handler extensions.IHandler) error
	RegisterInterfaceHandler(typeName string, handler extensions.InterfaceHandler) error
	Execute(operation Operation) (string, error)
	Get(domainName string, beanName string, attributeName string, args OperationArgs) (string, error)
	Put(domainName string, beanName string, attributeName string, args OperationArgs) (string, error)
	WithEnvironment(env *jnigi.Env) BeanExecutor
	GetEnv() *jnigi.Env
	OpenConnection(jndiURI string) (*jnigi.ObjectRef, error)
}

// RegisterClassHandler will register the given class handlers
// For a class handler to be valid it must implement a form of IClassHandler
func (mbean *Client) RegisterClassHandler(typeName string, handler extensions.IHandler) error {
	mbean.ClassHandlers[typeName] = handler
	return nil
}

// RegisterInterfaceHandler will register the given class handlers
// For a class handler to be valid it must implement a form of IClassHandler
func (mbean *Client) RegisterInterfaceHandler(typeName string, handler extensions.InterfaceHandler) error {
	mbean.InterfaceHandlers[typeName] = handler
	return nil
}

// WithEnvironment allos the client to spin up a new client using a new environment
// This is handy when using the same JmxConnection in sub threads
func (mbean *Client) WithEnvironment(env *jnigi.Env) BeanExecutor {
	return &Client{
		Env:               env,
		ClassHandlers:     mbean.ClassHandlers,
		InterfaceHandlers: mbean.InterfaceHandlers,
	}
}

// GetEnv will expose the underlying environment that the client is associated with
func (mbean *Client) GetEnv() *jnigi.Env {
	return mbean.Env
}

// Close is a method that will clean up all of the MBeans resources
// It will close the JMX method within the JVM as well as deleting the connection
// from the JNI resources
func Close(env *jnigi.Env, connection *jnigi.ObjectRef) {
	defer env.DeleteLocalRef(connection)
	connection.CallMethod(env, "close", nil)
}

// OpenConnection is a method that will establish a new connection against
// the given URI
func (mbean *Client) OpenConnection(jndiURI string) (*jnigi.ObjectRef, error) {
	stringRef, err := mbean.Env.NewObject("java/lang/String", []byte(jndiURI))

	if err != nil {
		return nil, fmt.Errorf("failed to create a string from %s::%s", jndiURI, err.Error())
	}

	jmxURL, err := mbean.Env.NewObject("javax/management/remote/JMXServiceURL", stringRef)
	if err != nil {
		return nil, errors.New("failed to create JMXServiceURL::" + err.Error())
	}

	if err != nil {
		return nil, errors.New("failed to create a blank map::" + err.Error())
	}

	jmxConnector := jnigi.NewObjectRef("javax/management/remote/JMXConnector")

	connectorFactory := "javax/management/remote/JMXConnectorFactory"
	if err = mbean.Env.CallStaticMethod(connectorFactory, "connect", jmxConnector, jmxURL); err != nil {
		return nil, errors.New("failed to create a JMX connection Factory::" + err.Error())
	}

	return jmxConnector, nil
}

// Execute is the orchestration for a JMX command execution.
func (mbean *Client) Execute(operation Operation) (string, error) {

	returnString := jnigi.NewObjectRef(OBJECT)
	if err := mbean.invoke(mbean.Env, operation, returnString); err != nil {
		return "", err
	}

	return mbean.toGoString(mbean.Env, returnString)
}

func (mbean *Client) invoke(env *jnigi.Env, operation Operation, outParam *jnigi.ObjectRef) error {
	objectName, err := getObjectName(env, operation)
	if err != nil {
		return errors.New("failed to create ObjectName::" + err.Error())
	}

	defer env.DeleteLocalRef(objectName)

	connection, err := mbean.OpenConnection(mbean.JmxURI)
	if err != nil {
		return err
	}

	defer Close(mbean.Env, connection)

	mBeanServerConnector, err := createMBeanServerConnection(env, connection)

	if err != nil {
		return errors.New("failed to create the mbean server connection::" + err.Error())
	}

	defer env.DeleteLocalRef(mBeanServerConnector)

	typeReferences, types, err := getOperationParameterTypes(env, objectName, mBeanServerConnector, operation.Operation)

	if err != nil {
		return err
	}

	defer env.DeleteLocalRef(typeReferences)
	names, err := mbean.getArray(env, operation.Args, types, OBJECT)
	if names != nil {
		defer env.DeleteLocalRef(names)
	}

	if err != nil {
		return err
	}

	operationRef, err := stringHandler.ToJniRepresentation(env, operation.Operation)
	defer env.DeleteLocalRef(operationRef)
	if err != nil {
		return err
	}

	if err = mBeanServerConnector.CallMethod(env, "invoke", outParam, objectName, operationRef, names, typeReferences); err != nil {
		return errors.New("failed to call invoke::" + err.Error())
	}

	return nil
}

func (mbean *Client) getArray(env *jnigi.Env, args []OperationArgs, methodTypes []string, className string) (*jnigi.ObjectRef, error) {

	types := make([]*jnigi.ObjectRef, 0)
	for i, arg := range args {

		obj, err := toJni(mbean, methodTypes[i], arg.JavaType, arg.JavaContainerType, arg.Value)

		if err != nil {
			return nil, err
		}

		types = append(types, obj)
	}

	return env.ToObjectArray(types, className), nil
}

func createMBeanServerConnection(env *jnigi.Env, connection *jnigi.ObjectRef) (*jnigi.ObjectRef, error) {
	mBeanServerConnector := jnigi.NewObjectRef("javax/management/MBeanServerConnection")
	if err := connection.CallMethod(env, "getMBeanServerConnection", mBeanServerConnector); err != nil {
		return nil, errors.New("failed to create the mbean server connection::" + err.Error())
	}

	return mBeanServerConnector, nil
}

func toJni(mbean *Client, methodType string, javaType string, containerType string, value string) (*jnigi.ObjectRef, error) {

	// the containers are always assumed to be interfaces
	if interfaceHandler, exists := mbean.InterfaceHandlers[methodType]; exists && containerType != "" {

		return interfaceHandler.ToJniRepresentation(mbean.Env, javaType, value)

	} else if handler, exists := mbean.ClassHandlers[methodType]; exists {
		var parsedVal any

		parsedVal, err := toTypeFromString(value, methodType)

		if err != nil {
			return nil, err
		}

		return handler.ToJniRepresentation(mbean.Env, parsedVal)
	}

	return nil, fmt.Errorf("no handlers exist for %s", javaType)
}
