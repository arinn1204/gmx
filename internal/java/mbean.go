package java

import (
	"errors"
	"fmt"
	"gmx/internal/jvm"
	"strings"

	"tekao.net/jnigi"
)

type MBean struct {
	Java             *jvm.Java
	serverConnection *jnigi.ObjectRef
	jmxConnection    *jnigi.ObjectRef
}

type MBeanOperation struct {
	Domain    string
	Name      string
	Operation string
	Args      []MBeanOperationArgs
}

type MBeanOperationArgs struct {
	Value any
	Type  string
}

type BeanExecutor interface {
	Execute(operation MBeanOperation) (any, error)
	InitializeMBeanConnection(uri string) error
	Close()
}

func (mbean *MBean) InitializeMBeanConnection(uri string) error {

	jmxConnector, err := buildJMXConnector(mbean.Java, uri)

	if err != nil {
		if jmxConnector != nil {
			closeReferences(mbean.Java.Env, jmxConnector)
		}
		return err
	}

	mBeanServerConnector := jnigi.NewObjectRef("javax/management/MBeanServerConnection")
	if err = jmxConnector.CallMethod(mbean.Java.Env, "getMBeanServerConnection", mBeanServerConnector); err != nil {
		return errors.New("failed to create the mbean server connection::" + err.Error())
	}

	mbean.serverConnection = mBeanServerConnector
	mbean.jmxConnection = jmxConnector

	return err
}

func (mbean *MBean) Execute(operation MBeanOperation) (any, error) {

	returnString := jnigi.NewObjectRef(jvm.OBJECT)
	if err := invoke(operation, mbean, returnString); err != nil {
		return "", err
	}

	return toGoString(mbean, returnString, jvm.STRING)
}

func (mbean *MBean) Close() {
	if mbean.Java == nil {
		return
	}
	closeReferences(mbean.Java.Env, mbean.jmxConnection)
}

func invoke(operation MBeanOperation, mbean *MBean, outParam *jnigi.ObjectRef) error {
	mbeanName := fmt.Sprintf("%s:name=%s", operation.Domain, operation.Name)
	objectParam, err := mbean.Java.CreateString(mbeanName)

	defer deleteReference(mbean, objectParam)
	if err != nil {
		return err
	}

	objectName, err := mbean.Java.Env.NewObject("javax/management/ObjectName", objectParam)
	defer deleteReference(mbean, objectName)
	if err != nil {
		return errors.New("failed to create ObjectName::" + err.Error())
	}

	names, err := getNameArray(mbean.Java, operation.Args)
	defer deleteReference(mbean, names)
	if err != nil {
		return err
	}

	types, err := getTypeArray(mbean.Java, operation.Args)
	defer deleteReference(mbean, types)
	if err != nil {
		return err
	}

	operationRef, err := mbean.Java.CreateString(operation.Operation)
	defer deleteReference(mbean, operationRef)
	if err != nil {
		return err
	}

	if err = mbean.serverConnection.CallMethod(mbean.Java.Env, "invoke", outParam, objectName, operationRef, names, types); err != nil {
		return errors.New("failed to call invoke::" + err.Error())
	}

	return nil
}

func getNameArray(java *jvm.Java, args []MBeanOperationArgs) (*jnigi.ObjectRef, error) {
	values := make([]any, 0)
	classes := make([]string, 0)

	for _, arg := range args {
		values = append(values, arg.Value)
		classes = append(classes, arg.Type)
	}

	return getArray(java, values, classes, jvm.OBJECT)
}

func getTypeArray(java *jvm.Java, args []MBeanOperationArgs) (*jnigi.ObjectRef, error) {
	types := make([]any, 0)
	classes := make([]string, 0)

	for _, arg := range args {
		types = append(types, arg.Type)
		classes = append(classes, jvm.STRING)
	}

	return getArray(java, types, classes, jvm.STRING)
}

func getArray(java *jvm.Java, values []any, classes []string, className string) (*jnigi.ObjectRef, error) {

	types := make([]*jnigi.ObjectRef, 0)
	for i, value := range values {
		var err error
		var obj *jnigi.ObjectRef

		jniClassPath := strings.Replace(classes[i], ".", "/", -1)

		if jniClassPath == jvm.STRING {
			obj, err = java.CreateString(value.(string))
		} else {
			obj, err = java.CreateJavaNative(value, jniClassPath)
		}

		if err != nil {
			return nil, err
		}

		types = append(types, obj)
	}

	return java.Env.ToObjectArray(types, className), nil
}
