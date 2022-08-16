package handlers

import (
	"errors"
	"fmt"

	"github.com/arinn1204/gmx/pkg/extensions"
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

func (iterator iterableRef[T]) fromJava(param *jnigi.ObjectRef, env *jnigi.Env) (any, error) {
	cls, err := getClassName(param, env)
	if err != nil {
		return nil, err
	}

	handlers := *iterator.classHandlers

	// always go class and then interface
	if handler, exists := handlers[cls]; exists {
		return getFromClassHandler(cls, handler, env, param)
	}

	return CheckForKnownInterfaces(env, param, cls, iterator.interfaceHandlers)
}

func getFromClassHandler(cls string, handler extensions.IHandler, env *jnigi.Env, param *jnigi.ObjectRef) (any, error) {
	var err error
	switch cls {
	case StringClasspath:
		var valDest string
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return nil, err
		}

		return valDest, nil
	case IntClasspath:
		var valDest int
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return nil, err
		}

		return valDest, nil
	case LongClasspath:
		var valDest int64
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return nil, err
		}

		return valDest, nil
	case FloatClasspath:
		var valDest float32
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return nil, err
		}

		return valDest, nil
	case DoubleClasspath:
		var valDest float64
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return nil, err
		}

		return valDest, nil
	case BoolClasspath:
		var valDest bool
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return nil, err
		}

		return valDest, nil
	default:
		return nil, fmt.Errorf("no known formatter for %s", cls)
	}
}
