package mbean

import (
	"errors"

	"tekao.net/jnigi"
)

func getIterator(env *jnigi.Env, obj *jnigi.ObjectRef) (*jnigi.ObjectRef, error) {
	iterator := jnigi.NewObjectRef("java/util/Iterator")

	if err := obj.CallMethod(env, "iterator", iterator); err != nil {
		return nil, errors.New("failed to construct iterator::" + err.Error())
	}

	return iterator, nil
}

func getNext(env *jnigi.Env, obj *jnigi.ObjectRef) (*jnigi.ObjectRef, error) {
	next := jnigi.NewObjectRef(OBJECT)

	if err := obj.CallMethod(env, "next", next); err != nil {
		return nil, errors.New("failed to construct iterator::" + err.Error())
	}

	return next, nil
}

func hasNext(env *jnigi.Env, obj *jnigi.ObjectRef) bool {
	var next bool
	if err := obj.CallMethod(env, "hasNext", &next); err != nil {
		return false
	}

	return next
}
