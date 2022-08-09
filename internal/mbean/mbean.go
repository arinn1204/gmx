package mbean

import (
	"errors"
	"fmt"
	"gmx/internal/jniwrapper"
	"strings"

	"tekao.net/jnigi"
)

type MBean struct {
	JmxConnection *jnigi.ObjectRef
	Env           *jnigi.Env
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
	Execute(env *jnigi.Env, operation MBeanOperation) (any, error)
}

func (mbean *MBean) Execute(env *jnigi.Env, operation MBeanOperation) (any, error) {

	returnString := jnigi.NewObjectRef(jniwrapper.OBJECT)
	if err := invoke(env, operation, mbean, returnString); err != nil {
		return "", err
	}

	return toGoString(env, returnString, jniwrapper.STRING)
}

func invoke(env *jnigi.Env, operation MBeanOperation, mbean *MBean, outParam *jnigi.ObjectRef) error {
	mbeanName := fmt.Sprintf("%s:name=%s", operation.Domain, operation.Name)
	objectParam, err := jniwrapper.CreateString(env, mbeanName)

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

	operationRef, err := jniwrapper.CreateString(env, operation.Operation)
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

func getNameArray(env *jnigi.Env, args []MBeanOperationArgs) (*jnigi.ObjectRef, error) {
	values := make([]any, 0)
	classes := make([]string, 0)

	for _, arg := range args {
		values = append(values, arg.Value)
		classes = append(classes, arg.Type)
	}

	return getArray(env, values, classes, jniwrapper.OBJECT)
}

func getTypeArray(env *jnigi.Env, args []MBeanOperationArgs) (*jnigi.ObjectRef, error) {
	types := make([]any, 0)
	classes := make([]string, 0)

	for _, arg := range args {
		types = append(types, arg.Type)
		classes = append(classes, jniwrapper.STRING)
	}

	return getArray(env, types, classes, jniwrapper.STRING)
}

func getArray(env *jnigi.Env, values []any, classes []string, className string) (*jnigi.ObjectRef, error) {

	types := make([]*jnigi.ObjectRef, 0)
	for i, value := range values {
		var err error
		var obj *jnigi.ObjectRef

		jniClassPath := strings.Replace(classes[i], ".", "/", -1)

		if jniClassPath == jniwrapper.STRING {
			obj, err = jniwrapper.CreateString(env, value.(string))
		} else {
			obj, err = jniwrapper.CreateJavaNative(env, value, jniClassPath)
		}

		if err != nil {
			return nil, err
		}

		types = append(types, obj)
	}

	return env.ToObjectArray(types, className), nil
}
