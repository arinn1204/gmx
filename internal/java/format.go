package java

import (
	"errors"
	"fmt"
	"strings"

	"tekao.net/jnigi"
)

func toGoString(mbean *MBean, param *jnigi.ObjectRef, outputType string) (any, error) {
	if param.IsNil() {
		return "NIL", nil
	}

	var bytes []byte

	clazz, err := getClass(param, mbean)

	defer deleteReference(mbean, param)
	if err != nil {
		return "", err
	}

	var result any

	if strings.EqualFold(clazz, "String") {
		if err := fromJavaString(param, mbean, &bytes); err != nil {
			return "", err
		}
		result = string(bytes)
	} else if strings.EqualFold(clazz, "Long") {
		res := int64(0)

		if err := fromJavaLong(param, mbean, &res); err != nil {
			return "", err
		}

		result = res
	} else if strings.EqualFold(clazz, "Integer") {
		res := 0

		if err := fromJavaInteger(param, mbean, &res); err != nil {
			return "", err
		}

		result = res
	} else {
		return "", fmt.Errorf("type of %s does not have a defined handler", clazz)
	}

	return result, nil
}

func fromJavaString(param *jnigi.ObjectRef, mbean *MBean, dest *[]byte) error {
	if err := param.CallMethod(mbean.Java.env, "getBytes", dest); err != nil {
		return errors.New("failed to convert response to a byte array::" + err.Error())
	}

	return nil
}

func fromJavaLong(param *jnigi.ObjectRef, mbean *MBean, dest *int64) error {
	if err := param.CallMethod(mbean.Java.env, "longValue", dest); err != nil {
		return errors.New("failed to create a long::" + err.Error())
	}

	return nil
}

func fromJavaInteger(param *jnigi.ObjectRef, mbean *MBean, dest *int) error {
	if err := param.CallMethod(mbean.Java.env, "intValue", dest); err != nil {
		return errors.New("failed to create a integer::" + err.Error())
	}

	return nil
}

func getClass(param *jnigi.ObjectRef, mbean *MBean) (string, error) {

	cls := jnigi.NewObjectRef("java/lang/Class")
	name := jnigi.NewObjectRef(STRING)

	defer deleteReference(mbean, name)
	defer deleteReference(mbean, cls)

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
