package mbean

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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

	handler := mbean.ClassHandlers[clazz]

	if strings.EqualFold(clazz, handlers.StringClasspath) {
		var str string

		if err = handler.ToGoRepresentation(env, param, &str); err != nil {
			return "", err
		}

		return str, nil
	} else if strings.EqualFold(clazz, handlers.LongClasspath) {
		res := int64(0)

		if err = handler.ToGoRepresentation(env, param, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", res), nil
	} else if strings.EqualFold(clazz, handlers.IntClasspath) {
		res := 0

		if err = handler.ToGoRepresentation(env, param, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", res), nil
	} else if strings.EqualFold(clazz, handlers.DoubleClasspath) {
		res := float64(0)

		if err = handler.ToGoRepresentation(env, param, &res); err != nil {
			return "", err
		}

		return strconv.FormatFloat(res, 'f', -1, 64), nil
	} else if strings.EqualFold(clazz, handlers.FloatClasspath) {
		res := float32(0)

		if err = handler.ToGoRepresentation(env, param, &res); err != nil {
			return "", err
		}

		return strconv.FormatFloat(float64(res), 'f', -1, 32), nil
	} else if strings.EqualFold(clazz, handlers.BoolClasspath) {
		res := false

		if err = handler.ToGoRepresentation(env, param, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%t", res), nil
	} else {
		return mbean.checkForKnownInterfaces(env, param, clazz)
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

func createString(env *jnigi.Env, str string) (*jnigi.ObjectRef, error) {
	stringRef, err := env.NewObject(STRING, []byte(str))
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", str, err.Error())
	}

	return stringRef, nil
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
