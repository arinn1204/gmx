package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	LongJniRepresentation = "java/lang/Long"
	LongClasspath         = "java.lang.Long"
)

type LongHandler struct{}

func (handler *LongHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(LongJniRepresentation, value.(int64))

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

func (handler *LongHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	if err := object.CallMethod(env, "longValue", dest); err != nil {
		return errors.New("failed to create a long::" + err.Error())
	}

	return nil
}
