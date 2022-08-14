package handlers

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

const (
	JNI_STRING       = "java/lang/String"
	STRING_CLASSPATH = "java.lang.String"
)

type StringHandler struct{}

func (handler *StringHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	stringRef, err := env.NewObject(JNI_STRING, []byte(value.(string)))
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", value, err.Error())
	}

	return stringRef, nil
}

// ToGoRepresentiation is the processing to go from a string java class to a string go class
// Dest must be a string pointer to receive the newly created string
func (handler *StringHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	var strBytes []byte
	if err := object.CallMethod(env, "getBytes", &strBytes); err != nil {
		return errors.New("failed to create a string::" + err.Error())
	}

	(*dest.(*string)) = string(strBytes)

	return nil
}
