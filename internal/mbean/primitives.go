package mbean

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

func fromJava(param *jnigi.ObjectRef, env *jnigi.Env, dest *[]any) error {
	cls, err := getClassName(param, env)
	if err != nil {
		return err
	}

	switch cls {
	case "String":
		var valDest []byte
		if err = fromJavaString(param, env, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, string(valDest))
	case "Integer":
		var valDest int
		if err = fromJavaInteger(param, env, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case "Long":
		var valDest int64
		if err = fromJavaLong(param, env, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case "Float":
		var valDest float32
		if err = fromJavaFloat(param, env, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case "Double":
		var valDest float64
		if err = fromJavaDouble(param, env, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case "Boolean":
		var valDest bool
		if err = fromJavaBoolean(param, env, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	default:
		return fmt.Errorf("no known formatter for %s", cls)
	}

	return nil
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
