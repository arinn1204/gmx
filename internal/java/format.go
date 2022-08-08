package java

import (
	"errors"
	"strings"

	"tekao.net/jnigi"
)

func toGoString(mbean *MBean, param *jnigi.ObjectRef, outputType string) (string, error) {
	if param.IsNil() {
		return "NIL", nil
	}

	var bytes []byte

	clazz, err := getClass(param, mbean)

	if err != nil {
		return "", err
	}

	if strings.EqualFold(clazz, "String") {
		if err := fromJavaString(param, mbean, &bytes); err != nil {
			return "", err
		}
	} else {
		return "", nil
	}

	return string(bytes), nil
}

func fromJavaString(param *jnigi.ObjectRef, mbean *MBean, dest *[]byte) error {
	if err := param.CallMethod(mbean.Java.env, "getBytes", dest); err != nil {
		return errors.New("failed to convert response to a byte array::" + err.Error())
	}

	return nil
}

func getClass(param *jnigi.ObjectRef, mbean *MBean) (string, error) {

	cls := jnigi.NewObjectRef("java/lang/Class")
	name := jnigi.NewObjectRef(STRING)

	if err := param.CallMethod(mbean.Java.env, "getClass", cls); err != nil {
		return "", errors.New("failed to call getClass::" + err.Error())
	}

	if err := cls.CallMethod(mbean.Java.env, "getSimpleName", name); err != nil {
		return "", errors.New("failed to get class name::" + err.Error())
	}

	var bytes []byte
	if err := name.CallMethod(mbean.Java.env, "getBytes", &bytes); err != nil {
		return "", errors.New("failed to get byte representation::" + err.Error())
	}

	return string(bytes), nil
}
