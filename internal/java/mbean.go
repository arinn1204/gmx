package java

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

// the commonly used types
const (
	STRING = "java/lang/String"
	OBJECT = "java/lang/Object"
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
	Value string
	Type  string
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
	if err != nil {
		return err
	}

	objectName, err := mbean.Java.env.NewObject("javax/management/ObjectName", objectParam)

	if err != nil {
		return errors.New("failed to create ObjectName::" + err.Error())
	}

	names, err := getNameArray(mbean.Java, operation.Args)

	if err != nil {
		return err
	}

	types, err := getTypeArray(mbean.Java, operation.Args)

	if err != nil {
		return err
	}

	operationRef, err := mbean.Java.createString(operation.Operation)

	if err != nil {
		return err
	}

	if err = mbean.serverConnection.CallMethod(mbean.Java.env, "invoke", outParam, objectName, operationRef, names, types); err != nil {
		return errors.New("failed to call invoke::" + err.Error())
	}
	return nil
}

func getNameArray(java *Java, args []MBeanOperationArgs) (*jnigi.ObjectRef, error) {
	values := make([]string, 0)

	for _, arg := range args {
		values = append(values, arg.Value)
	}

	return getArray(java, values, OBJECT)
}

func getTypeArray(java *Java, args []MBeanOperationArgs) (*jnigi.ObjectRef, error) {
	types := make([]string, 0)

	for _, arg := range args {
		types = append(types, arg.Type)
	}

	return getArray(java, types, STRING)
}

func getArray(java *Java, values []string, className string) (*jnigi.ObjectRef, error) {

	types := make([]*jnigi.ObjectRef, 0)
	for _, value := range values {
		str, err := java.createString(value)
		if err != nil {
			return nil, err
		}

		types = append(types, str)
	}

	return java.env.ToObjectArray(types, className), nil
}
