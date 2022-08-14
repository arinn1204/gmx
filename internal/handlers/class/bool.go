package class

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	BOOLEAN = "java/lang/Boolean"
)

type boolHandler struct{}

func (handler *boolHandler) toJniRepresentation(env *jnigi.Env, value bool) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(BOOLEAN, value)

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

func (handler *boolHandler) toGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *bool) error {
	if err := object.CallMethod(env, "boolValue", dest); err != nil {
		return errors.New("failed to create a bool::" + err.Error())
	}

	return nil
}
