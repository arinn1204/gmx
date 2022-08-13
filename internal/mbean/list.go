package mbean

import (
	"encoding/json"
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

func createListFromJSON(value string, env *jnigi.Env) (*jnigi.ObjectRef, error) {
	dest := make([]any, 0)
	if err := json.Unmarshal([]byte(value), &dest); err != nil {
		return nil, fmt.Errorf("failed to convert %s to a map::%s", value, err)
	}

	return createJavaList(env, dest)
}

func createJavaList(env *jnigi.Env, list []any) (*jnigi.ObjectRef, error) {
	arrayList, err := env.NewObject("java/util/ArrayList", len(list))
	if err != nil {
		return nil, errors.New("failed to create an arraylist::" + err.Error())
	}

	res := false
	for _, val := range list {
		ref, err := fromGoAny(env, val)

		if err != nil {
			return nil, err
		}

		env.PrecalculateSignature("(Ljava/lang/Object;)Z") //since we don't have type params for the list
		err = arrayList.CallMethod(env, "add", &res, ref)

		if !res || err != nil {
			return nil, fmt.Errorf("failed to add %s to the list::%s", val, err.Error())
		}
	}

	return arrayList, nil
}

func createGoArray(list *jnigi.ObjectRef, env *jnigi.Env, dest *[]any) error {
	return nil
}
