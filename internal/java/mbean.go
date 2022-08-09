package java

import (
	"errors"
	"fmt"
	"strings"

	"tekao.net/jnigi"
)

// the commonly used types
const (
	STRING  = "java/lang/String"
	OBJECT  = "java/lang/Object"
	LONG    = "java/lang/Long"
	INTEGER = "java/lang/Integer"
	BOOLEAN = "java/lang/Boolean"
)

type MBean struct {
	Java             *Java
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

func (mbean *MBean) Execute(operation MBeanOperation) (any, error) {

	returnString := jnigi.NewObjectRef(OBJECT)
	if err := invoke(operation, mbean, returnString); err != nil {
		return "", err
	}

	return toGoString(mbean, returnString, STRING)
}

func (mbean *MBean) Close() {
	if mbean.Java == nil {
		return
	}
	closeReferences(mbean.Java.env, mbean.jmxConnection)
}

func invoke(operation MBeanOperation, mbean *MBean, outParam *jnigi.ObjectRef) error {
	mbeanName := fmt.Sprintf("%s:name=%s", operation.Domain, operation.Name)
	objectParam, err := mbean.Java.createString(mbeanName)

	defer mbean.Java.env.DeleteLocalRef(objectParam)

	if err != nil {
		return err
	}

	objectName, err := mbean.Java.env.NewObject("javax/management/ObjectName", objectParam)
	defer mbean.Java.env.DeleteLocalRef(objectName)
	if err != nil {
		return errors.New("failed to create ObjectName::" + err.Error())
	}

	names, err := getNameArray(mbean.Java, operation.Args)
	defer mbean.Java.env.DeleteLocalRef(names)
	if err != nil {
		return err
	}

	types, err := getTypeArray(mbean.Java, operation.Args)
	defer mbean.Java.env.DeleteLocalRef(types)
	if err != nil {
		return err
	}

	operationRef, err := mbean.Java.createString(operation.Operation)
	defer mbean.Java.env.DeleteLocalRef(operationRef)
	if err != nil {
		return err
	}

	if err = mbean.serverConnection.CallMethod(mbean.Java.env, "invoke", outParam, objectName, operationRef, names, types); err != nil {
		return errors.New("failed to call invoke::" + err.Error())
	}

	return nil
}

func getNameArray(java *Java, args []MBeanOperationArgs) (*jnigi.ObjectRef, error) {
	values := make([]any, 0)
	classes := make([]string, 0)

	for _, arg := range args {
		values = append(values, arg.Value)
		classes = append(classes, arg.Type)
	}

	return getArray(java, values, classes, OBJECT)
}

func getTypeArray(java *Java, args []MBeanOperationArgs) (*jnigi.ObjectRef, error) {
	types := make([]any, 0)
	classes := make([]string, 0)

	for _, arg := range args {
		types = append(types, arg.Type)
		classes = append(classes, STRING)
	}

	return getArray(java, types, classes, STRING)
}

func getArray(java *Java, values []any, classes []string, className string) (*jnigi.ObjectRef, error) {

	types := make([]*jnigi.ObjectRef, 0)
	for i, value := range values {
		var err error
		var obj *jnigi.ObjectRef

		jniClassPath := strings.Replace(classes[i], ".", "/", -1)

		if jniClassPath == STRING {
			obj, err = java.createString(value.(string))
		} else if jniClassPath == LONG {
			obj, err = java.createLong(value.(int64))
		}

		if err != nil {
			return nil, err
		}

		types = append(types, obj)
	}

	return java.env.ToObjectArray(types, className), nil
}
