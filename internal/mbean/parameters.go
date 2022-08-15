package mbean

import (
	"errors"
	"fmt"

	"github.com/arinn1204/gmx/internal/handlers"
	"tekao.net/jnigi"
)

func deleteLocalArray(env *jnigi.Env, arr []*jnigi.ObjectRef) {
	if arr == nil {
		return
	}

	for _, ob := range arr {
		env.DeleteLocalRef(ob)
	}
}

func getOperationParameterTypes(env *jnigi.Env,
	objectName *jnigi.ObjectRef,
	serverConnection *jnigi.ObjectRef,
	operationName string) (*jnigi.ObjectRef, []string, error) {

	operations, err := getBeanOperations(env, objectName, serverConnection)
	if err != nil {
		return nil, nil, err
	}

	defer deleteLocalArray(env, operations)

	operatonRef, err := getOperation(env, operations, operationName)

	if err != nil {
		return nil, nil, err
	}

	defer env.DeleteLocalRef(operatonRef)

	types, err := getSignature(env, operatonRef)

	if err != nil {
		return nil, nil, err
	}

	typeReferences := make([]*jnigi.ObjectRef, 0)

	for _, paramType := range types {
		ref, err := stringHandler.ToJniRepresentation(env, paramType)

		if err != nil {
			defer deleteLocalArray(env, typeReferences)
			return nil, nil, err
		}

		typeReferences = append(typeReferences, ref)
	}

	return env.ToObjectArray(typeReferences, handlers.StringJniRepresentation), types, nil
}

func getSignature(env *jnigi.Env, operationInfo *jnigi.ObjectRef) ([]string, error) {
	parameterInfo := jnigi.NewObjectArrayRef("javax/management/MBeanParameterInfo")

	if err := operationInfo.CallMethod(env, "getSignature", parameterInfo); err != nil {
		return nil, err
	}

	defer env.DeleteLocalRef(parameterInfo)

	parameters := env.FromObjectArray(parameterInfo)

	parameterTypes := make([]string, 0)

	for _, parameter := range parameters {
		stringRef := jnigi.NewObjectRef(handlers.StringJniRepresentation)

		if err := parameter.CallMethod(env, "getType", stringRef); err != nil {
			return nil, err
		}

		defer env.DeleteLocalRef(stringRef)

		var parameterType string
		if err := stringHandler.ToGoRepresentation(env, stringRef, &parameterType); err != nil {
			return nil, errors.New("failed to convert to a go string::" + err.Error())
		}

		parameterTypes = append(parameterTypes, parameterType)
	}

	return parameterTypes, nil
}

func getOperation(env *jnigi.Env, operations []*jnigi.ObjectRef, operationName string) (*jnigi.ObjectRef, error) {
	var mbeanOperationInfo *jnigi.ObjectRef

	for _, operation := range operations {
		nameReference := jnigi.NewObjectRef(handlers.StringJniRepresentation)

		if err := operation.CallMethod(env, "getName", nameReference); err != nil {
			continue
		}

		defer env.DeleteLocalRef(nameReference)

		var name string
		if err := stringHandler.ToGoRepresentation(env, nameReference, &name); err != nil {
			continue
		}

		if name == operationName {
			mbeanOperationInfo = operation
			break
		} else {
			env.DeleteLocalRef(operation)
		}
	}

	var err error

	if mbeanOperationInfo == nil {
		err = fmt.Errorf("failed to find an operation matching name %s", operationName)
	}

	return mbeanOperationInfo, err
}

func getBeanOperations(env *jnigi.Env, objectName *jnigi.ObjectRef, serverConnection *jnigi.ObjectRef) ([]*jnigi.ObjectRef, error) {
	beanInfo := jnigi.NewObjectRef("javax/management/MBeanInfo")

	if err := serverConnection.CallMethod(env, "getMBeanInfo", beanInfo, objectName); err != nil {
		return nil, errors.New("failed to get the mbean info::" + err.Error())
	}

	defer env.DeleteLocalRef(beanInfo)

	operationInfoArray := jnigi.NewObjectArrayRef("javax/management/MBeanOperationInfo")

	if err := beanInfo.CallMethod(env, "getOperations", operationInfoArray); err != nil {
		return nil, errors.New("failed to get operations for bean::" + err.Error())
	}

	return env.FromObjectArray(operationInfoArray), nil
}
