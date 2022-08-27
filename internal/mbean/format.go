package mbean

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

var stringHandler extensions.IHandler

func init() {
	stringHandler = &handlers.StringHandler{}
}

func (mbean *Client) toGoString(env *jnigi.Env, param *jnigi.ObjectRef) (string, error) {
	if param.IsNil() {
		return "", nil
	}

	clazz, err := getClassName(param, env)

	defer env.DeleteLocalRef(param)
	if err != nil {
		return "", err
	}

	if handler, exists := mbean.ClassHandlers.Load(clazz); exists {
		_ = handler
		return fromJava(clazz, env, param, handler.(extensions.IHandler))
	}

	val, err := handlers.CheckForKnownInterfaces(env, param, clazz, &mbean.InterfaceHandlers)

	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(val)

	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func fromJava(classPath string, env *jnigi.Env, parameter *jnigi.ObjectRef, handler extensions.IHandler) (string, error) {
	switch classPath {
	case handlers.StringClasspath:
		var str string

		if err := handler.ToGoRepresentation(env, parameter, &str); err != nil {
			return "", err
		}

		return str, nil
	case handlers.BoolClasspath:
		res := false

		if err := handler.ToGoRepresentation(env, parameter, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%t", res), nil
	case handlers.LongClasspath:
		res := int64(0)

		if err := handler.ToGoRepresentation(env, parameter, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", res), nil
	case handlers.IntClasspath:
		res := 0

		if err := handler.ToGoRepresentation(env, parameter, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", res), nil
	case handlers.FloatClasspath:
		res := float32(0)

		if err := handler.ToGoRepresentation(env, parameter, &res); err != nil {
			return "", err
		}

		return strconv.FormatFloat(float64(res), 'f', -1, 32), nil
	case handlers.DoubleClasspath:
		res := float64(0)

		if err := handler.ToGoRepresentation(env, parameter, &res); err != nil {
			return "", err
		}

		return strconv.FormatFloat(res, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("no handler exists for %s", classPath)
	}
}

func getClass(param *jnigi.ObjectRef, env *jnigi.Env) (*jnigi.ObjectRef, error) {
	cls := jnigi.NewObjectRef("java/lang/Class")

	if err := param.CallMethod(env, "getClass", cls); err != nil {
		return nil, errors.New("failed to call getClass::" + err.Error())
	}

	return cls, nil
}

func getClassName(param *jnigi.ObjectRef, env *jnigi.Env) (string, error) {
	name := jnigi.NewObjectRef(STRING)
	defer env.DeleteLocalRef(name)

	cls, err := getClass(param, env)
	defer env.DeleteLocalRef(cls)

	if err != nil {
		return "", err
	}

	if err := cls.CallMethod(env, "getName", name); err != nil {
		return "", errors.New("failed to get class name::" + err.Error())
	}

	var strName string
	if err = stringHandler.ToGoRepresentation(env, name, &strName); err != nil {
		return "", err
	}

	return strName, nil
}

func toTypeFromString(value string, className string) (any, error) {
	var val any
	var err error
	switch className {
	case handlers.BoolClasspath:
		val, err = strconv.ParseBool(value)
	case handlers.DoubleClasspath:
		val, err = strconv.ParseFloat(value, 64)
	case handlers.FloatClasspath:
		val, err = strconv.ParseFloat(value, 32)
		val = float32(val.(float64))
	case handlers.IntClasspath:
		val, err = strconv.ParseInt(value, 10, 32)
		val = int(val.(int64))
	case handlers.LongClasspath:
		val, err = strconv.ParseInt(value, 10, 64)
	case handlers.StringClasspath, handlers.StringJniRepresentation:
		val, err = value, nil
	}

	return val, err
}
