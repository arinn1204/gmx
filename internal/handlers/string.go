package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

// These are the constants for the String classpath and JNI representation
const (
	StringJniRepresentation = "java/lang/String"
	StringClasspath         = "java.lang.String"
)

// StringHandler is the type that can convert to and from java.lang.String
type StringHandler struct{}

// ToJniRepresentation is the implementation that will convert from a go type
// to a JNI representation of that type
func (handler *StringHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	stringRef, err := env.NewObject(StringJniRepresentation, []byte(value.(string)))
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", value, err.Error())
	}

	return stringRef, nil
}

// ToGoRepresentation is the processing to go from a string java class to a string go class
// Dest must be a string pointer to receive the newly created string
func (handler *StringHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	var strBytes []byte
	if err := object.CallMethod(env, "getBytes", &strBytes); err != nil {
		return errors.New("failed to create a string::" + err.Error())
	}

	(*dest.(*string)) = string(strBytes)

	return nil
}
