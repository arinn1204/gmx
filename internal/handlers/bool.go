package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	JNI_BOOLEAN    = "java/lang/Boolean"
	BOOL_CLASSPATH = "java.lang.Boolean"
)

type BoolHandler struct{}

func (handler *BoolHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(JNI_BOOLEAN, value.(bool))

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

func (handler *BoolHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *any) error {
	val := false
	if err := object.CallMethod(env, "boolValue", &val); err != nil {
		return errors.New("failed to create a bool::" + err.Error())
	}

	*dest = val

	return nil
}
