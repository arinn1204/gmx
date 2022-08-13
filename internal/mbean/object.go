package mbean

import (
	"encoding/json"
	"fmt"
	"strconv"

	"tekao.net/jnigi"
)

func createObjectReference(env *jnigi.Env, value string, classPath string) (*jnigi.ObjectRef, error) {
	if classPath == STRING {
		return createString(env, value)
	} else if classPath == FLOAT {
		return createFloat(env, value)
	} else if classPath == DOUBLE {
		return createDouble(env, value)
	} else if classPath == LIST {
		return createList(value, env)
	} else if classPath == MAP {
		return createMap(value, env)
	} else {
		return createJavaNative(env, value, classPath)
	}
}

func createList(value string, env *jnigi.Env) (*jnigi.ObjectRef, error) {
	dest := make(map[any]any)
	if err := json.Unmarshal([]byte(value), &dest); err != nil {
		return nil, fmt.Errorf("failed to convert %s to a map::%s", value, err)
	}

	return createJavaList(env, value)
}

func createMap(value string, env *jnigi.Env) (*jnigi.ObjectRef, error) {
	dest := make(map[any]any)
	if err := json.Unmarshal([]byte(value), &dest); err != nil {
		return nil, fmt.Errorf("failed to convert %s to a map::%s", value, err)
	}

	return createJavaMap(env, dest)
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
