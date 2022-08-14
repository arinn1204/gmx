package class

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	LONG = "java/lang/Long"
)

type longHandler struct{}

func (handler *longHandler) toJniRepresentation(env *jnigi.Env, value int64) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(LONG, value)

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

func (handler *longHandler) toGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *int64) error {
	if err := object.CallMethod(env, "longValue", dest); err != nil {
		return errors.New("failed to create a long::" + err.Error())
	}

	return nil
}
