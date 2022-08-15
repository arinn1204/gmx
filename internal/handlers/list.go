package handlers

import "tekao.net/jnigi"

// These are the constants for the List classpath and JNI representation
const (
	ListJniRepresentation = "java/util/List"
	ListClassPath         = "java.util.List"
)

// ListHandler is the type that will be able to convert lists to and from go arrays
type ListHandler struct{}

// ToJniRepresentation is the implementation that will convert from a go type
// to a JNI representation of that type
func (handler *ListHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	return nil, nil
}

// ToGoRepresentation will convert from a JNI type to a go type
func (handler *ListHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	return nil
}
