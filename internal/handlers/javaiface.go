package handlers

import (
	"errors"
	"fmt"

	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

func getParent(param *jnigi.ObjectRef, env *jnigi.Env) (*jnigi.ObjectRef, error) {
	cls := jnigi.NewObjectRef("java/lang/Class")

	if err := param.CallMethod(env, "getSuperclass", cls); err != nil {
		return nil, errors.New("failed to call getClass::" + err.Error())
	}

	return cls, nil
}

func getInterfaces(env *jnigi.Env, cls *jnigi.ObjectRef) ([]*jnigi.ObjectRef, error) {
	var name string
	nameref := jnigi.NewObjectRef(StringJniRepresentation)

	cls.CallMethod(env, "getName", nameref)
	strHandler.ToGoRepresentation(env, nameref, &name)

	interfaceRef := jnigi.NewObjectArrayRef("java/lang/Class")

	defer env.DeleteLocalRef(interfaceRef)

	if err := cls.CallMethod(env, "getInterfaces", interfaceRef); err != nil {
		return nil, fmt.Errorf("failed to get interfaces::%s", err)
	}

	interfaces := env.FromObjectArray(interfaceRef)

	parent, err := getParent(cls, env)

	if err != nil || parent == nil || parent.IsNil() {
		return interfaces, err
	}

	defer env.DeleteLocalRef(parent)

	parentInterfaces, err := getInterfaces(env, parent)
	if err != nil {
		return nil, err
	}
	interfaces = append(interfaces, parentInterfaces...)

	return interfaces, nil
}

// CheckForKnownInterfaces is a function that will check if the object has an interface that can be
// processed by the given handlers
func CheckForKnownInterfaces(env *jnigi.Env, param *jnigi.ObjectRef, clazz string, interfaceHandlers *map[string]extensions.InterfaceHandler) (any, error) {
	cls, err := getClass(param, env)
	if err != nil {
		return nil, err
	}
	defer env.DeleteLocalRef(cls)
	interfaces, err := getInterfaces(env, cls)

	if err != nil {
		return "", fmt.Errorf("%s::%s", clazz, err)
	}

	for _, iface := range interfaces {
		name := jnigi.NewObjectRef(StringJniRepresentation)
		defer env.DeleteLocalRef(name)

		if err = iface.CallMethod(env, "getName", name); err != nil {
			return "", fmt.Errorf("failed to get name of interface::%s", err)
		}

		var cls string

		if err = strHandler.ToGoRepresentation(env, name, &cls); err != nil {
			return "", err
		}

		if handler, exists := (*interfaceHandlers)[cls]; exists {

			if cls == MapClassPath {
				dest := make(map[string]any)

				if err := handler.ToGoRepresentation(env, param, &dest); err != nil {
					return "", err
				}
				return dest, nil

			}
			dest := make([]any, 0)
			if err := handler.ToGoRepresentation(env, param, &dest); err != nil {
				return "", err
			}
			return dest, nil

		}
	}

	return "", fmt.Errorf("type of %s does not have a defined handler", clazz)
}
