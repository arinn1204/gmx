package mbean

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"tekao.net/jnigi"
)

func toGoString(env *jnigi.Env, param *jnigi.ObjectRef, outputType string) (string, error) {
	if param.IsNil() {
		return "", nil
	}

	var bytes []byte

	clazz, err := getClass(param, env)

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
	} else if strings.EqualFold(clazz, "List") {
		res := make([]any, 0)

		if err := createGoArrayFromList(param, env, &res); err != nil {
			return "", err
		}

		return "", nil
	} else if strings.EqualFold(clazz, "Map") {
		res := make(map[any]any)

		if err := createGoMap(param, env, &res); err != nil {
			return "", err
		}

		return "", nil
	} else {
		return "", fmt.Errorf("type of %s does not have a defined handler", clazz)
	}
}

func fromJavaString(param *jnigi.ObjectRef, env *jnigi.Env, dest *[]byte) error {
	if err := param.CallMethod(env, "getBytes", dest); err != nil {
		return errors.New("failed to convert response to a byte array::" + err.Error())
	}

	return nil
}

func fromJavaLong(param *jnigi.ObjectRef, env *jnigi.Env, dest *int64) error {
	if err := param.CallMethod(env, "longValue", dest); err != nil {
		return errors.New("failed to create a long::" + err.Error())
	}

	return nil
}

func fromJavaDouble(param *jnigi.ObjectRef, env *jnigi.Env, dest *float64) error {
	if err := param.CallMethod(env, "doubleValue", dest); err != nil {
		return errors.New("failed to create a long::" + err.Error())
	}

	return nil
}

func fromJavaFloat(param *jnigi.ObjectRef, env *jnigi.Env, dest *float32) error {
	if err := param.CallMethod(env, "floatValue", dest); err != nil {
		return errors.New("failed to create a long::" + err.Error())
	}

	return nil
}

func fromJavaBoolean(param *jnigi.ObjectRef, env *jnigi.Env, dest *bool) error {
	if err := param.CallMethod(env, "booleanValue", dest); err != nil {
		return errors.New("failed to create a long::" + err.Error())
	}

	return nil
}

func fromJavaInteger(param *jnigi.ObjectRef, env *jnigi.Env, dest *int) error {
	if err := param.CallMethod(env, "intValue", dest); err != nil {
		return errors.New("failed to create a integer::" + err.Error())
	}

	return nil
}

func getClass(param *jnigi.ObjectRef, env *jnigi.Env) (string, error) {

	cls := jnigi.NewObjectRef("java/lang/Class")
	name := jnigi.NewObjectRef(STRING)

	defer env.DeleteLocalRef(name)
	defer env.DeleteLocalRef(cls)

	if err := param.CallMethod(env, "getClass", cls); err != nil {
		return "", errors.New("failed to call getClass::" + err.Error())
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
	return createFloatingPoiintValue(env, str, FLOAT)
}

func createDouble(env *jnigi.Env, str string) (*jnigi.ObjectRef, error) {
	return createFloatingPoiintValue(env, str, DOUBLE)
}

func createFloatingPoiintValue(env *jnigi.Env, str string, class string) (*jnigi.ObjectRef, error) {
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
