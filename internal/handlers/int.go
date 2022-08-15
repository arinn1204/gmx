package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

// These are the constants for the Integer classpath and JNI representation
const (
	IntegerJniRepresentation = "java/lang/Integer"
	IntClasspath             = "java.lang.Integer"
)

// IntHandler is the type that can convert to and from java.lang.Integer
type IntHandler struct{}

// ToJniRepresentation is the implementation that will convert from a go type
// to a JNI representation of that type
func (handler *IntHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	intref, err := env.NewObject(IntegerJniRepresentation, value.(int))

	if err != nil {
		return nil, fmt.Errorf("failed to create integer from %d::%s", value, err)
	}

	return intref, err
}

// ToGoRepresentation will convert from a JNI type to a go type
func (handler *IntHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	if err := object.CallMethod(env, "intValue", dest); err != nil {
		return errors.New("failed to create a integer::" + err.Error())
	}

	return nil
}
