package extensions

import "tekao.net/jnigi"

// IHandler is the extension interface that can be used to provide custom handling of objects
// this will be registered with the BeanExporter when sending and receiving information from java.
// It is intended to drive the conversion logic between the JNI and go
type IHandler interface {
	ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error)
	ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error
}

// InterfaceHandler is the extension interface that will be used to translate to and from JNI representations
// of an object that is defined by an interface. An example would be a list, a set, or a map
type InterfaceHandler interface {
	ToJniRepresentation(env *jnigi.Env, elementType string, value any) (*jnigi.ObjectRef, error)
	ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error
}
