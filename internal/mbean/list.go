package mbean

import (
	"errors"

	"tekao.net/jnigi"
)

func createJavaList[T any](env *jnigi.Env, size int) (*iterableRef[T], error) {
	arrayList, err := env.NewObject("java/util/ArrayList", size)
	if err != nil {
		return nil, errors.New("failed to create an arraylist::" + err.Error())
	}

	return &iterableRef[T]{iterable: arrayList}, nil
}

func (iterable *iterableRef[T]) add(env *jnigi.Env, item T) error {
	res := false
	param, err := createObjectReferenceFromValue(env, item)

	if err != nil {
		return err
	}

	env.PrecalculateSignature("(Ljava/lang/Object;)Z") //since we don't have type params for the list
	return iterable.iterable.CallMethod(env, "add", &res, param)
}

func (iterable *iterableRef[T]) toObjectReference() *jnigi.ObjectRef {
	return iterable.iterable
}

func (mbean *Client) createGoArrayFromList(param *jnigi.ObjectRef, env *jnigi.Env, dest *[]any) error {
	iterator, err := getIterator(env, param)

	if err != nil {
		return err
	}

	for iterator.hasNext(env) {
		value, err := iterator.getNext(env)
		if err != nil {
			return err
		}
		if err = mbean.fromJava(value, env, dest); err != nil {
			return err
		}
	}

	return nil
}
