package mbean

import (
	"encoding/json"
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

func getInterfaces(env *jnigi.Env, param *jnigi.ObjectRef) ([]*jnigi.ObjectRef, error) {
	cls, err := getClass(param, env)

	if err != nil {
		return nil, err
	}

	defer env.DeleteLocalRef(cls)

	interfaceRef := jnigi.NewObjectArrayRef("java/lang/Class")

	defer env.DeleteLocalRef(interfaceRef)

	if err = cls.CallMethod(env, "getInterfaces", interfaceRef); err != nil {
		return nil, fmt.Errorf("failed to get interfaces::%s", err)
	}

	return env.FromObjectArray(interfaceRef), nil
}

func (mbean *Client) checkForKnownInterfaces(env *jnigi.Env, param *jnigi.ObjectRef, clazz string) (string, error) {
	interfaces, err := getInterfaces(env, param)

	if err != nil {
		return "", fmt.Errorf("%s::%s", clazz, err)
	}

	for _, iface := range interfaces {
		name := jnigi.NewObjectRef(STRING)
		defer env.DeleteLocalRef(name)

		if err = iface.CallMethod(env, "getName", name); err != nil {
			return "", fmt.Errorf("failed to get name of interface::%s", err)
		}

		var dest string

		if err = stringHandler.ToGoRepresentation(env, name, &dest); err != nil {
			return "", err
		}

		if handler, exists := mbean.InterfaceHandlers[dest]; exists {
			dest := make([]any, 0)

			if err := handler.ToGoRepresentation(env, param, &dest); err != nil {
				return "", err
			}
			arr, err := json.Marshal(dest)

			if err != nil {
				return "", errors.New("failed to turn list into json array::" + err.Error())
			}
			return string(arr), nil
		}
	}

	return "", fmt.Errorf("type of %s does not have a defined handler", clazz)
}
