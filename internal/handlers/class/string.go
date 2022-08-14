package class

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	STRING = "java/lang/String"
)

type stringHandler struct{}

func (handler *stringHandler) toJniRepresentation(env *jnigi.Env, value string) (*jnigi.ObjectRef, error) {
	stringRef, err := env.NewObject(STRING, []byte(value))
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", value, err.Error())
	}

	return stringRef, nil
}

func (handler *stringHandler) toGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *string) error {
	strBytes := make([]byte, 0)
	if err := object.CallMethod(env, "getBytes", strBytes); err != nil {
		return errors.New("failed to create a string::" + err.Error())
	}

	*dest = string(strBytes)

	return nil
}
