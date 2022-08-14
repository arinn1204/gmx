package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	JNI_LONG       = "java/lang/Long"
	LONG_CLASSPATH = "java.lang.Long"
)

type LongHandler struct{}

func (handler *LongHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(JNI_LONG, value.(int64))

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

func (handler *LongHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *any) error {
	val := int64(0)

	if err := object.CallMethod(env, "longValue", &val); err != nil {
		return errors.New("failed to create a long::" + err.Error())
	}

	*dest = val

	return nil
}
