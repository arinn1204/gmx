package class

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	STRING = "java/lang/String"
)

type StringHandler struct{}

func (handler *StringHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	stringRef, err := env.NewObject(STRING, []byte(value.(string)))
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", value, err.Error())
	}

	return stringRef, nil
}

func (handler *StringHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *any) error {
	strBytes := make([]byte, 0)
	if err := object.CallMethod(env, "getBytes", strBytes); err != nil {
		return errors.New("failed to create a string::" + err.Error())
	}

	*dest = string(strBytes)

	return nil
}
