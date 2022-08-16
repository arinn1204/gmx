package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

// These are the constants for the Set classpath and JNI representation
const (
	SetJniRepresentation = "java/util/Set"
	SetClassPath         = "java.util.Set"
)

// SetHandler is the type that will be able to convert lists to and from go arrays
type SetHandler struct {
	ClassHandlers *map[string]extensions.IHandler
}

// ToJniRepresentation is the implementation that will convert from a go type
// to a JNI representation of that type
func (handler *SetHandler) ToJniRepresentation(env *jnigi.Env, elementType string, value any) (*jnigi.ObjectRef, error) {
	byteStr := []byte(value.(string))

	switch elementType {
	case StringClasspath:
		arr := make([]string, 0)
		json.Unmarshal(byteStr, &arr)
		if set, err := createJavaSet(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, set, arr, handler.ClassHandlers)
		}
	case DoubleClasspath:
		arr := make([]float64, 0)
		json.Unmarshal(byteStr, &arr)
		if set, err := createJavaSet(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, set, arr, handler.ClassHandlers)
		}
	case FloatClasspath:
		arr := make([]float32, 0)
		json.Unmarshal(byteStr, &arr)
		if set, err := createJavaSet(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, set, arr, handler.ClassHandlers)
		}
	case IntClasspath:
		arr := make([]int, 0)
		json.Unmarshal(byteStr, &arr)
		if set, err := createJavaSet(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, set, arr, handler.ClassHandlers)
		}
	case LongClasspath:
		arr := make([]int64, 0)
		json.Unmarshal(byteStr, &arr)
		if set, err := createJavaSet(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, set, arr, handler.ClassHandlers)
		}
	case BoolClasspath:
		arr := make([]bool, 0)
		json.Unmarshal(byteStr, &arr)
		if set, err := createJavaSet(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, set, arr, handler.ClassHandlers)
		}
	}

	return nil, fmt.Errorf("there are no registered handlers for %s<%s>", ListClassPath, elementType)
}

// ToGoRepresentation will convert from a JNI type to a go type
func (handler *SetHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	iterator, err := getIterator(env, object, handler.ClassHandlers)

	if err != nil {
		return err
	}

	defer env.DeleteLocalRef(iterator.iterable)
	for iterator.hasNext(env) {
		value, err := iterator.getNext(env)
		if err != nil {
			return err
		}
		defer env.DeleteLocalRef(value)
		val, err := iterator.fromJava(value, env)
		if err != nil {
			return err
		} else {
			*dest.(*[]any) = append(*dest.(*[]any), val)
		}
	}

	return nil
}

func createJavaSet[T any](env *jnigi.Env, arr []T, handlers *map[string]extensions.IHandler) (*iterableRef[T], error) {
	size := len(arr)
	arrayList, err := env.NewObject("java/util/HashSet", size)
	if err != nil {
		return nil, errors.New("failed to create an arraylist::" + err.Error())
	}

	return &iterableRef[T]{iterable: arrayList, classHandlers: handlers}, nil
}
