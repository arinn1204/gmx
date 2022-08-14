package extensions

import "tekao.net/jnigi"

type IClassHandler interface {
	ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error)
	ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *any) error
}
