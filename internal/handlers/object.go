package handlers

import (
	"errors"
	"reflect"

	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

func populateGenericContainer[T any](env *jnigi.Env, collection *iterableRef[T], arr []T, handlers map[string]extensions.IHandler) (*jnigi.ObjectRef, error) {
	for _, item := range arr {
		if err := collection.add(env, item); err != nil {
			return nil, err
		}
	}

	return collection.iterable, nil
}

func getClassPathFromType(env *jnigi.Env, value any) (string, error) {
	switch value.(type) {
	case int:
		return IntClasspath, nil
	case int64:
		return LongClasspath, nil
	case float32:
		return FloatClasspath, nil
	case float64:
		return DoubleClasspath, nil
	case bool:
		return BoolClasspath, nil
	case string:
		return StringClasspath, nil
	}

	return "", errors.New("no defined translater for value " + reflect.TypeOf(value).Name())
}
