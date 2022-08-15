package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	BoolJniRepresentation = "java/lang/Boolean"
	BoolClasspath         = "java.lang.Boolean"
)

type BoolHandler struct{}

func (handler *BoolHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(BoolJniRepresentation, value.(bool))

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

func (handler *BoolHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	if err := object.CallMethod(env, "booleanValue", dest); err != nil {
		return errors.New("failed to create a bool::" + err.Error())
	}
	return nil
}
