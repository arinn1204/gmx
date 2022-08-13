package mbean

import (
	"strconv"

	"tekao.net/jnigi"
)

func fromGoAny(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {

	switch value.(type) {
	case float32, float64:
		return createFloat(env, strconv.FormatFloat(value.(float64), 'f', -1, 64))
	}

	return nil, nil
}
func getIterator(env *jnigi.Env, obj *jnigi.ObjectRef) (*jnigi.ObjectRef, error) {
	return nil, nil
}

func getNext(env *jnigi.Env, obj *jnigi.ObjectRef) (*jnigi.ObjectRef, error) {
	return nil, nil
}

func hasNext(env *jnigi.Env, obj *jnigi.ObjectRef) bool {
	return false
}
