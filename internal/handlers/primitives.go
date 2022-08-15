package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

func getClass(param *jnigi.ObjectRef, env *jnigi.Env) (*jnigi.ObjectRef, error) {
	cls := jnigi.NewObjectRef("java/lang/Class")

	if err := param.CallMethod(env, "getClass", cls); err != nil {
		return nil, errors.New("failed to call getClass::" + err.Error())
	}

	return cls, nil
}

func getClassName(param *jnigi.ObjectRef, env *jnigi.Env) (string, error) {
	name := jnigi.NewObjectRef(StringJniRepresentation)
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
	if err = strHandler.ToGoRepresentation(env, name, &strName); err != nil {
		return "", err
	}

	return strName, nil
}

func (iterator iterableRef[T]) fromJava(param *jnigi.ObjectRef, env *jnigi.Env, dest *[]any) error {
	cls, err := getClassName(param, env)
	if err != nil {
		return err
	}

	handler := iterator.classHandlers[cls]

	switch cls {
	case StringClasspath:
		var valDest string
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, string(valDest))
	case IntClasspath:
		var valDest int
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case LongClasspath:
		var valDest int64
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case FloatClasspath:
		var valDest float32
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case DoubleClasspath:
		var valDest float64
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case BoolClasspath:
		var valDest bool
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	default:
		return fmt.Errorf("no known formatter for %s", cls)
	}

	return nil
}
