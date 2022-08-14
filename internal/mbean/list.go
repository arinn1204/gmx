package mbean

import (
	"errors"

	"tekao.net/jnigi"
)

type listRef[T any] struct {
	list *jnigi.ObjectRef
}

func createJavaList[T any](env *jnigi.Env, size int) (*listRef[T], error) {
	arrayList, err := env.NewObject("java/util/ArrayList", size)
	if err != nil {
		return nil, errors.New("failed to create an arraylist::" + err.Error())
	}

	return &listRef[T]{list: arrayList}, nil
}

func (list *listRef[T]) add(env *jnigi.Env, item T) error {
	res := false
	param, err := createObjectReferenceFromValue(env, item)

	if err != nil {
		return err
	}

	env.PrecalculateSignature("(Ljava/lang/Object;)Z") //since we don't have type params for the list
	return list.list.CallMethod(env, "add", &res, param)
}

func (list *listRef[T]) toObjectReference() *jnigi.ObjectRef {
	return list.list
}

func createGoArrayFromList(param *jnigi.ObjectRef, env *jnigi.Env, dest *[]any) error {
	iterator, err := getIterator(env, param)

	if err != nil {
		return err
	}

	for hasNext(env, iterator) {
		value, err := getNext(env, iterator)
		if err != nil {
			return err
		}
		if err = fromJava(value, env, dest); err != nil {
			return err
		}
	}

	return nil
}
