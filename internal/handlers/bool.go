package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	// BoolJniRepresentation is the representation of the boolean type in the JNI
	BoolJniRepresentation = "java/lang/Boolean"

	// BoolClasspath is the fully qualified boolean name
	BoolClasspath = "java.lang.Boolean"
)

// BoolHandler is the implementation of the IHandler that handles boolean conversions
type BoolHandler struct{}

// ToJniRepresentation is the implementation that will convert from a go type
// to a JNI representation of that type
func (handler *BoolHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(BoolJniRepresentation, value.(bool))

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

// ToGoRepresentation will convert from a JNI type to a go type
func (handler *BoolHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	if err := object.CallMethod(env, "booleanValue", dest); err != nil {
		return errors.New("failed to create a bool::" + err.Error())
	}
	return nil
}
