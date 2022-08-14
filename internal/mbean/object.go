package mbean

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"tekao.net/jnigi"
)

func createObjectReferenceFromValue(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	switch valueType := value.(type) {
	case int:
		return createObjectReference(env, fmt.Sprintf("%d", value), INTEGER)
	case int64:
		return createObjectReference(env, fmt.Sprintf("%d", value), LONG)
	case float32:
		return createObjectReference(env, fmt.Sprintf("%f", value), FLOAT)
	case float64:
		return createObjectReference(env, fmt.Sprintf("%f", value), DOUBLE)
	case bool:
		return createObjectReference(env, fmt.Sprintf("%t", value), BOOLEAN)
	case string:
		return createObjectReference(env, valueType, STRING)
	}

	return nil, errors.New("no defined translater for value " + reflect.TypeOf(value).Name())
}

func createObjectReference(env *jnigi.Env, value string, classPath string) (*jnigi.ObjectRef, error) {
	if classPath == STRING {
		return createString(env, value)
	} else if classPath == FLOAT {
		return createFloat(env, value)
	} else if classPath == DOUBLE {
		return createDouble(env, value)
	} else {
		return createJavaNative(env, value, classPath)
	}
}

func createContainerReference(env *jnigi.Env, value string, elementTypePath string, containerTypePath string) (*jnigi.ObjectRef, error) {
	byteStr := []byte(value)

	switch elementTypePath {
	case STRING:
		arr := make([]string, 0)
		json.Unmarshal(byteStr, &arr)
		return createGenericContainerReference(env, containerTypePath, arr)
	case DOUBLE:
		arr := make([]float64, 0)
		json.Unmarshal(byteStr, &arr)
		return createGenericContainerReference(env, containerTypePath, arr)
	case FLOAT:
		arr := make([]float32, 0)
		json.Unmarshal(byteStr, &arr)
		return createGenericContainerReference(env, containerTypePath, arr)
	case INTEGER:
		arr := make([]int, 0)
		json.Unmarshal(byteStr, &arr)
		return createGenericContainerReference(env, containerTypePath, arr)
	case LONG:
		arr := make([]int64, 0)
		json.Unmarshal(byteStr, &arr)
		return createGenericContainerReference(env, containerTypePath, arr)
	case BOOLEAN:
		arr := make([]bool, 0)
		json.Unmarshal(byteStr, &arr)
		return createGenericContainerReference(env, containerTypePath, arr)
	}

	return nil, fmt.Errorf("there are no registered handlers for %s<%s>", containerTypePath, elementTypePath)
}

func createGenericContainerReference[T any](env *jnigi.Env, containerType string, arr []T) (*jnigi.ObjectRef, error) {
	list, err := createJavaList[T](env, len(arr))

	if err != nil {
		return nil, err
	}

	for _, item := range arr {
		if err = list.add(env, item); err != nil {
			return nil, err
		}
	}

	return list.toObjectReference(), nil
}

// CreateJavaNative is a helper used to turn a primitive go type
// (int, int64, float32/64, bool) into the corresponding java types
func createJavaNative(env *jnigi.Env, input string, typeName string) (*jnigi.ObjectRef, error) {
	var obj any
	var err error

	switch typeName {
	case STRING:
		obj = input
	case INTEGER:
		obj, err = strconv.ParseInt(input, 10, 4*8)
		obj = int(obj.(int64))
	case LONG:
		obj, err = strconv.ParseInt(input, 10, 8*8)
	case BOOLEAN:
		obj, err = strconv.ParseBool(input)
	}

	if err != nil {
		return nil, err
	}

	ref, err := env.NewObject(typeName, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", obj, err.Error())
	}

	return ref, nil
}
