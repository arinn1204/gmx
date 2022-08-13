package mbean

import (
	"errors"
	"fmt"
	"reflect"

	"tekao.net/jnigi"
)

func createJavaList(env *jnigi.Env, list any) (*jnigi.ObjectRef, error) {
	switch list := list.(type) {
	case []int:
		return genericCreateJavaList(env, list)
	}

	return nil, fmt.Errorf("no existing handler for %s", reflect.TypeOf(list))
}

func genericCreateJavaList[T any](env *jnigi.Env, list []T) (*jnigi.ObjectRef, error) {
	arrayList, err := env.NewObject("java/util/ArrayList", len(list))
	if err != nil {
		return nil, errors.New("failed to create an arraylist::" + err.Error())
	}

	res := false
	for _, val := range list {
		env.PrecalculateSignature("(Ljava/lang/Object;)Z") //since we don't have type params for the list
		err = arrayList.CallMethod(env, "add", &res, val)

		if !res || err != nil {
			return nil, fmt.Errorf("failed to add %s to the list::%s", val, err.Error())
		}
	}

	return arrayList, nil
}

func createGoArray(list *jnigi.ObjectRef, env *jnigi.Env, dest *[]any) error {
	return nil
}
