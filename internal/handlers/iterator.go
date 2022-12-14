package handlers

import (
	"errors"
	"sync"

	"tekao.net/jnigi"
)

const (
	// ObjectJniRepresentation is the JNI representation of an object
	ObjectJniRepresentation = "java/lang/Object"

	// IteratorJniRepresentation is the JNI representation of an iterator
	IteratorJniRepresentation = "java/util/Iterator"
)

type iterableRef[T any] struct {
	iterable          *jnigi.ObjectRef
	classHandlers     *sync.Map
	interfaceHandlers *sync.Map
}

func getIterator(env *jnigi.Env, param *jnigi.ObjectRef, classHandlers *sync.Map, interfaceHandlers *sync.Map) (*iterableRef[any], error) {
	iterator := jnigi.NewObjectRef(IteratorJniRepresentation)

	if err := param.CallMethod(env, "iterator", iterator); err != nil {
		return nil, errors.New("failed to construct iterator::" + err.Error())
	}

	return &iterableRef[any]{
		iterable:          iterator,
		classHandlers:     classHandlers,
		interfaceHandlers: interfaceHandlers,
	}, nil
}

func (iterable *iterableRef[T]) getNext(env *jnigi.Env) (*jnigi.ObjectRef, error) {
	next := jnigi.NewObjectRef(ObjectJniRepresentation)

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
