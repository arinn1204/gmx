package mbean

import (
	"errors"
	"fmt"
	"strings"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

// a collection of JNI representations of java primitive types
// these will be the boxed representations, not the true primitive
const (
	STRING  = "java/lang/String"
	OBJECT  = "java/lang/Object"
	LONG    = "java/lang/Long"
	INTEGER = "java/lang/Integer"
	BOOLEAN = "java/lang/Boolean"
	FLOAT   = "java/lang/Float"
	DOUBLE  = "java/lang/Double"
	LIST    = "java/util/List"
	MAP     = "java/util/Map"
)

// Client is the overarching type that will facilitate JMX connections
// JmxConnection is the living connection that was created when `CreateMBeanConnection` was called to create the Client
// Env is the environment that belongs to this bean, this will not always match the JVM env!
type Client struct {
	JmxConnection     *jnigi.ObjectRef
	Env               *jnigi.Env
	ClassHandlers     map[string]extensions.IHandler
	InterfaceHandlers map[string]extensions.IHandler
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
	RegisterInterfaceHandler(typeName string, handler extensions.IHandler) error
	Execute(operation Operation) (string, error)
	WithEnvironment(env *jnigi.Env) BeanExecutor
	GetEnv() *jnigi.Env
	Close()
}

// RegisterClassHandler will register the given class handlers
// For a class handler to be valid it must implement a form of IClassHandler
func (mbean *Client) RegisterClassHandler(typeName string, handler extensions.IHandler) error {
	mbean.ClassHandlers[typeName] = handler
	return nil
}

// RegisterInterfaceHandler will register the given class handlers
// For a class handler to be valid it must implement a form of IClassHandler
func (mbean *Client) RegisterInterfaceHandler(typeName string, handler extensions.IHandler) error {
	mbean.InterfaceHandlers[typeName] = handler
	return nil
}

// WithEnvironment allos the client to spin up a new client using a new environment
// This is handy when using the same JmxConnection in sub threads
func (mbean *Client) WithEnvironment(env *jnigi.Env) BeanExecutor {
	return &Client{
		JmxConnection: mbean.JmxConnection,
		Env:           env,
	}
}

// GetEnv will expose the underlying environment that the client is associated with
func (mbean *Client) GetEnv() *jnigi.Env {
	return mbean.Env
}

// Close is a method that will clean up all of the MBeans resources
// It will close the JMX method within the JVM as well as deleting the connection
// from the JNI resources
func (mbean *Client) Close() {
	defer mbean.Env.DeleteLocalRef(mbean.JmxConnection)
	mbean.JmxConnection.CallMethod(mbean.Env, "close", nil)
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
	mbeanName := fmt.Sprintf("%s:name=%s", operation.Domain, operation.Name)
	objectParam, err := createString(env, mbeanName)

	defer env.DeleteLocalRef(objectParam)
	if err != nil {
		return err
	}

	objectName, err := env.NewObject("javax/management/ObjectName", objectParam)
	defer env.DeleteLocalRef(objectName)
	if err != nil {
		return errors.New("failed to create ObjectName::" + err.Error())
	}

	names, err := mbean.getValueArray(env, operation.Args)
	if names != nil {
		defer env.DeleteLocalRef(names)
	}

	if err != nil {
		return err
	}

	types, err := mbean.getTypeArray(env, operation.Args)

	if types != nil {
		defer env.DeleteLocalRef(types)
	}

	if err != nil {
		return err
	}

	operationRef, err := createString(env, operation.Operation)
	defer env.DeleteLocalRef(operationRef)
	if err != nil {
		return err
	}

	mBeanServerConnector := jnigi.NewObjectRef("javax/management/MBeanServerConnection")
	defer env.DeleteLocalRef(mBeanServerConnector)
	if err = mbean.JmxConnection.CallMethod(env, "getMBeanServerConnection", mBeanServerConnector); err != nil {
		return errors.New("failed to create the mbean server connection::" + err.Error())
	}

	if err = mBeanServerConnector.CallMethod(env, "invoke", outParam, objectName, operationRef, names, types); err != nil {
		return errors.New("failed to call invoke::" + err.Error())
	}

	return nil
}

func (mbean *Client) getValueArray(env *jnigi.Env, args []OperationArgs) (*jnigi.ObjectRef, error) {
	values := make([]string, 0)
	classes := make([]string, 0)
	containerType := make([]string, 0)

	for _, arg := range args {
		values = append(values, arg.Value)
		classes = append(classes, arg.JavaType)
		containerType = append(containerType, arg.JavaContainerType)
	}

	return mbean.getArray(env, values, classes, containerType, OBJECT)
}

func (mbean *Client) getTypeArray(env *jnigi.Env, args []OperationArgs) (*jnigi.ObjectRef, error) {
	types := make([]string, 0)
	classes := make([]string, 0)
	containerType := make([]string, 0)

	for _, arg := range args {
		var paramType string
		if arg.JavaContainerType == "" {
			paramType = arg.JavaType
		} else {
			paramType = arg.JavaContainerType
		}

		types = append(types, paramType)
		classes = append(classes, handlers.StringClasspath)
		containerType = append(containerType, "") // types can't be arrays
	}

	return mbean.getArray(env, types, classes, containerType, STRING)
}

func (mbean *Client) getArray(env *jnigi.Env, values []string, classes []string, containerType []string, className string) (*jnigi.ObjectRef, error) {

	types := make([]*jnigi.ObjectRef, 0)
	for i, value := range values {

		var err error
		var obj *jnigi.ObjectRef

		jniClassPath := strings.Replace(classes[i], ".", "/", -1)
		if containerType[i] == "" {
			handler := mbean.ClassHandlers[classes[i]]
			var parsedVal any

			parsedVal, err = toTypeFromString(value, classes[i])

			if err != nil {
				return nil, err
			}

			obj, err = handler.ToJniRepresentation(env, parsedVal)
		} else {
			containerClassPath := strings.Replace(containerType[i], ".", "/", -1)
			obj, err = createContainerReference(env, value, jniClassPath, containerClassPath)
		}

		if err != nil {
			return nil, err
		}

		types = append(types, obj)
	}

	return env.ToObjectArray(types, className), nil
}
