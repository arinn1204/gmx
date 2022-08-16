package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

// These are the constants for the List classpath and JNI representation
const (
	ListJniRepresentation = "java/util/List"
	ListClassPath         = "java.util.List"
)

// ListHandler is the type that will be able to convert lists to and from go arrays
type ListHandler struct {
	ClassHandlers *map[string]extensions.IHandler
}

// ToJniRepresentation is the implementation that will convert from a go type
// to a JNI representation of that type
func (handler *ListHandler) ToJniRepresentation(env *jnigi.Env, elementType string, value any) (*jnigi.ObjectRef, error) {
	byteStr := []byte(value.(string))

	switch elementType {
	case StringClasspath:
		arr := make([]string, 0)
		json.Unmarshal(byteStr, &arr)
		if lst, err := createJavaList(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, lst, arr, handler.ClassHandlers)
		}
	case DoubleClasspath:
		arr := make([]float64, 0)
		json.Unmarshal(byteStr, &arr)
		if lst, err := createJavaList(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, lst, arr, handler.ClassHandlers)
		}
	case FloatClasspath:
		arr := make([]float32, 0)
		json.Unmarshal(byteStr, &arr)
		if lst, err := createJavaList(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, lst, arr, handler.ClassHandlers)
		}
	case IntClasspath:
		arr := make([]int, 0)
		json.Unmarshal(byteStr, &arr)
		if lst, err := createJavaList(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, lst, arr, handler.ClassHandlers)
		}
	case LongClasspath:
		arr := make([]int64, 0)
		json.Unmarshal(byteStr, &arr)
		if lst, err := createJavaList(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, lst, arr, handler.ClassHandlers)
		}
	case BoolClasspath:
		arr := make([]bool, 0)
		json.Unmarshal(byteStr, &arr)
		if lst, err := createJavaList(env, arr, handler.ClassHandlers); err == nil {
			return populateGenericContainer(env, lst, arr, handler.ClassHandlers)
		}
	}

	return nil, fmt.Errorf("there are no registered handlers for %s<%s>", ListClassPath, elementType)
}

// ToGoRepresentation will convert from a JNI type to a go type
func (handler *ListHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
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

func createJavaList[T any](env *jnigi.Env, arr []T, handlers *map[string]extensions.IHandler) (*iterableRef[T], error) {
	size := len(arr)
	arrayList, err := env.NewObject("java/util/ArrayList", size)
	if err != nil {
		return nil, errors.New("failed to create an arraylist::" + err.Error())
	}

	return &iterableRef[T]{iterable: arrayList, classHandlers: handlers}, nil
}

func (iterable *iterableRef[T]) add(env *jnigi.Env, item T) error {
	res := false
	classPath, err := getClassPathFromType(env, item)

	if err != nil {
		return err
	}

	if handler, exists := (*iterable.classHandlers)[classPath]; exists {
		param, err := handler.ToJniRepresentation(env, item)

		if err != nil {
			return err
		}

		env.PrecalculateSignature("(Ljava/lang/Object;)Z") //since we don't have type params for the list
		return iterable.iterable.CallMethod(env, "add", &res, param)
	}

	return fmt.Errorf("there is no class handler defined for %s", reflect.TypeOf(item).Name())
}
