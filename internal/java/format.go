package java

import (
	"errors"

	"tekao.net/jnigi"
)

func (mbean *MBean) toGoString(param *jnigi.ObjectRef, outputType string) (string, error) {
	if param.IsNil() {
		return "NIL", nil
	}

	var bytes []byte

	if err := param.CallMethod(mbean.Java.env, "getBytes", &bytes); err != nil {
		return "", errors.New("failed to convert response to a byte array::" + err.Error())
	}

	return string(bytes), nil
}
