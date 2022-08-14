package mbean

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"tekao.net/jnigi"
)

func toGoString(env *jnigi.Env, param *jnigi.ObjectRef) (string, error) {
	if param.IsNil() {
		return "", nil
	}

	var bytes []byte

	clazz, err := getClassName(param, env)

	defer env.DeleteLocalRef(param)
	if err != nil {
		return "", err
	}

	if strings.EqualFold(clazz, "String") {
		if err := fromJavaString(param, env, &bytes); err != nil {
			return "", err
		}
		return string(bytes), nil
	} else if strings.EqualFold(clazz, "Long") {
		res := int64(0)

		if err := fromJavaLong(param, env, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", res), nil
	} else if strings.EqualFold(clazz, "Integer") {
		res := 0

		if err := fromJavaInteger(param, env, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", res), nil
	} else if strings.EqualFold(clazz, "Double") {
		res := float64(0)

		if err := fromJavaDouble(param, env, &res); err != nil {
			return "", err
		}

		return strconv.FormatFloat(res, 'f', -1, 64), nil
	} else if strings.EqualFold(clazz, "Float") {
		res := float32(0)

		if err := fromJavaFloat(param, env, &res); err != nil {
			return "", err
		}

		return strconv.FormatFloat(float64(res), 'f', -1, 32), nil
	} else if strings.EqualFold(clazz, "Boolean") {
		res := false

		if err := fromJavaBoolean(param, env, &res); err != nil {
			return "", err
		}

		return fmt.Sprintf("%t", res), nil
	} else {
		return checkForKnownInterfaces(env, param, clazz)
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

	if err := cls.CallMethod(env, "getSimpleName", name); err != nil {
		return "", errors.New("failed to get class name::" + err.Error())
	}

	var bytes []byte
	if err := name.CallMethod(env, "getBytes", &bytes); err != nil {
		return "", errors.New("failed to get byte representation::" + err.Error())
	}

	return string(bytes), nil
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
