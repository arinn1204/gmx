package class

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	INTEGER = "java/lang/Integer"
)

type intHandler struct{}

func (handler *intHandler) toJniRepresentation(env *jnigi.Env, value int) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(INTEGER, value)

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

func (handler *intHandler) toGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *int) error {
	if err := object.CallMethod(env, "intValue", dest); err != nil {
		return errors.New("failed to create a integer::" + err.Error())
	}

	return nil
}
