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

	if strings.EqualFold(clazz, handlers.STRING_CLASSPATH) {
		var str string

		if err = handler.ToGoRepresentation(env, param, &str); err != nil {
			return "", err
		}

		return str, nil
	} else if strings.EqualFold(clazz, handlers.LONG_CLASSPATH) {
		res := int64(0)

		if err = handler.ToGoRepresentation(env, param, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", res), nil
	} else if strings.EqualFold(clazz, handlers.INT_CLASSPATH) {
		res := 0

		if err = handler.ToGoRepresentation(env, param, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", res), nil
	} else if strings.EqualFold(clazz, handlers.DOUBLE_CLASSPATH) {
		res := float64(0)

		if err = handler.ToGoRepresentation(env, param, &res); err != nil {
			return "", err
		}

		return strconv.FormatFloat(res, 'f', -1, 64), nil
	} else if strings.EqualFold(clazz, handlers.FLOAT_CLASSPATH) {
		res := float32(0)

		if err = handler.ToGoRepresentation(env, param, &res); err != nil {
			return "", err
		}

		return strconv.FormatFloat(float64(res), 'f', -1, 32), nil
	} else if strings.EqualFold(clazz, handlers.BOOL_CLASSPATH) {
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

func createFloat(env *jnigi.Env, str string) (*jnigi.ObjectRef, error) {
	return createFloatingPointValue(env, str, FLOAT)
}

func createDouble(env *jnigi.Env, str string) (*jnigi.ObjectRef, error) {
	return createFloatingPointValue(env, str, DOUBLE)
}

func createFloatingPointValue(env *jnigi.Env, str string, class string) (*jnigi.ObjectRef, error) {
	stringifiedFloat, err := createString(env, str)

	if err != nil {
		return nil, err
	}

	floatRef := jnigi.NewObjectRef(class)
	if err = env.CallStaticMethod(class, "valueOf", floatRef, stringifiedFloat); err != nil {
		return nil, fmt.Errorf("failed to create a %s from stringref::%s", class, err)
	}

	return floatRef, nil
}

func toTypeFromString(value string, className string) (any, error) {
	var val any
	var err error
	switch className {
	case handlers.BOOL_CLASSPATH:
		val, err = strconv.ParseBool(value)
	case handlers.DOUBLE_CLASSPATH:
		val, err = strconv.ParseFloat(value, 64)
	case handlers.FLOAT_CLASSPATH:
		val, err = strconv.ParseFloat(value, 32)
		val = float32(val.(float64))
	case handlers.INT_CLASSPATH:
		val, err = strconv.ParseInt(value, 10, 32)
		val = int(val.(int64))
	case handlers.LONG_CLASSPATH:
		val, err = strconv.ParseInt(value, 10, 64)
	case handlers.STRING_CLASSPATH, handlers.JNI_STRING:
		val, err = value, nil
	}

	return val, err
}
