package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	IntegerJniRepresentation = "java/lang/Integer"
	IntClasspath             = "java.lang.Integer"
)

type IntHandler struct{}

func (handler *IntHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(IntegerJniRepresentation, value.(int))

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

func (handler *IntHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	if err := object.CallMethod(env, "intValue", dest); err != nil {
		return errors.New("failed to create a integer::" + err.Error())
	}

	return nil
}
