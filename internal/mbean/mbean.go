package mbean

import (
	"errors"
	"fmt"
	"strings"

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
)

// Client is the overarching type that will facilitate JMX connections
// JmxConnection is the living connection that was created when `CreateMBeanConnection` was called to create the Client
// Env is the environment that belongs to this bean, this will not always match the JVM env!
type Client struct {
	JmxConnection *jnigi.ObjectRef
	Env           *jnigi.Env
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
// Value is the value that is being entered
// Type is the fully qualified java type `java.lang.String`
type OperationArgs struct {
	Value any
	Type  string
}

// BeanExecutor is the interface used around this package.
// This is how an execution is performed.
// This will always rely on the MBean Client's environment
type BeanExecutor interface {
	Execute(operation Operation) (any, error)
}

// Execute is the orchestration for a JMX command execution.
func (mbean *Client) Execute(operation Operation) (any, error) {

	returnString := jnigi.NewObjectRef(OBJECT)
	if err := invoke(mbean.Env, operation, mbean, returnString); err != nil {
		return "", err
	}

	return toGoString(mbean.Env, returnString, STRING)
}

func invoke(env *jnigi.Env, operation Operation, mbean *Client, outParam *jnigi.ObjectRef) error {
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

	names, err := getNameArray(env, operation.Args)
	defer env.DeleteLocalRef(names)
	if err != nil {
		return err
	}

	types, err := getTypeArray(env, operation.Args)
	defer env.DeleteLocalRef(types)
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

func getNameArray(env *jnigi.Env, args []OperationArgs) (*jnigi.ObjectRef, error) {
	values := make([]any, 0)
	classes := make([]string, 0)

	for _, arg := range args {
		values = append(values, arg.Value)
		classes = append(classes, arg.Type)
	}

	return getArray(env, values, classes, OBJECT)
}

func getTypeArray(env *jnigi.Env, args []OperationArgs) (*jnigi.ObjectRef, error) {
	types := make([]any, 0)
	classes := make([]string, 0)

	for _, arg := range args {
		types = append(types, arg.Type)
		classes = append(classes, STRING)
	}

	return getArray(env, types, classes, STRING)
}

func getArray(env *jnigi.Env, values []any, classes []string, className string) (*jnigi.ObjectRef, error) {

	types := make([]*jnigi.ObjectRef, 0)
	for i, value := range values {
		var err error
		var obj *jnigi.ObjectRef

		jniClassPath := strings.Replace(classes[i], ".", "/", -1)

		if jniClassPath == STRING {
			obj, err = createString(env, value.(string))
		} else {
			obj, err = createJavaNative(env, value, jniClassPath)
		}

		if err != nil {
			return nil, err
		}

		types = append(types, obj)
	}

	return env.ToObjectArray(types, className), nil
}
