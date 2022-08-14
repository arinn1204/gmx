package extensions

import "tekao.net/jnigi"

type IClassHandler[T any] interface {
	toJniRepresentation(env *jnigi.Env, value T) (*jnigi.ObjectRef, error)
	toGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *T) error
}

type IInterfaceHandler[T comparable] interface {
	toJniRepresentation(env *jnigi.Env, value T) (*jnigi.Env, error)
	toGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *T) error
}
