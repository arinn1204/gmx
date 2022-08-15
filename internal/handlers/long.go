package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

// These are the constants for the Long classpath and JNI representation
const (
	LongJniRepresentation = "java/lang/Long"
	LongClasspath         = "java.lang.Long"
)

// LongHandler is the type that can convert to and from java.lang.Long
type LongHandler struct{}

// ToJniRepresentation is the implementation that will convert from a go type
// to a JNI representation of that type
func (handler *LongHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(LongJniRepresentation, value.(int64))

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

// ToGoRepresentation will convert from a JNI type to a go type
func (handler *LongHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	if err := object.CallMethod(env, "longValue", dest); err != nil {
		return errors.New("failed to create a long::" + err.Error())
	}

	return nil
}
