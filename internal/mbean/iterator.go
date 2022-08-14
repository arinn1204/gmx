package mbean

import (
	"errors"

	"tekao.net/jnigi"
)

type iterableRef[T any] struct {
	iterable *jnigi.ObjectRef
}

func getIterator(env *jnigi.Env, param *jnigi.ObjectRef) (*iterableRef[any], error) {
	iterator := jnigi.NewObjectRef("java/util/Iterator")

	if err := param.CallMethod(env, "iterator", iterator); err != nil {
		return nil, errors.New("failed to construct iterator::" + err.Error())
	}

	return &iterableRef[any]{iterable: iterator}, nil
}

func (iterable *iterableRef[T]) getNext(env *jnigi.Env) (*jnigi.ObjectRef, error) {
	next := jnigi.NewObjectRef(OBJECT)

	if err := iterable.iterable.CallMethod(env, "next", next); err != nil {
		return nil, errors.New("failed to construct iterator::" + err.Error())
	}

	return next, nil
}

func (iterable *iterableRef[T]) hasNext(env *jnigi.Env) bool {
	var next bool
	if err := iterable.iterable.CallMethod(env, "hasNext", &next); err != nil {
		return false
	}

	return next
}
