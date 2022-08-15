package handlers

import (
	"errors"

	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

const (
	MapClassPath = "java.util.Map"
)

// MapHandler is the type that will be able to convert maps to and from go arrays
type MapHandler struct {
	ClassHandlers map[string]extensions.IHandler
}

// ToJniRepresentation is the ability to translate from a go map to a java map
func (mh *MapHandler) ToJniRepresentation(env *jnigi.Env, elementType string, value any) (*jnigi.ObjectRef, error) {
	return nil, errors.New("translating from go map to JNI is not supported at this time")
}

// ToGoRepresentation is the handling method to translate from a java map type to a go map
func (mh *MapHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	return nil
}
