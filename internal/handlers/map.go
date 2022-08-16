package handlers

import (
	"errors"

	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

const (
	MapClassPath = "java.util.Map"
)

// MapHandler is the type that will be able to convert maps to and from go arrays
type MapHandler struct {
	ClassHandlers *map[string]extensions.IHandler
}

// ToJniRepresentation is the ability to translate from a go map to a java map
func (mh *MapHandler) ToJniRepresentation(env *jnigi.Env, elementType string, value any) (*jnigi.ObjectRef, error) {
	return nil, errors.New("translating from go map to JNI is not supported at this time")
}

// ToGoRepresentation is the handling method to translate from a java map type to a go map
func (mh *MapHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	entrySet := jnigi.NewObjectRef(SetJniRepresentation)
	if err := object.CallMethod(env, "entrySet", entrySet); err != nil {
		return err
	}

	defer env.DeleteLocalRef(entrySet)

	iterator, err := getIterator(env, entrySet, mh.ClassHandlers)

	if err != nil {
		return err
	}

	defer env.DeleteLocalRef(iterator.iterable)

	returnedMap := make(map[string]any)

	for iterator.hasNext(env) {
		entry, err := iterator.getNext(env)
		if err != nil {
			return err
		}
		key, value, err := getKeyAndValue(env, iterator, entry)

		if err != nil {
			return err
		}

		returnedMap[key] = value
	}

	*dest.(*map[string]any) = returnedMap

	return nil
}

func getKeyAndValue(env *jnigi.Env, iterator *iterableRef[any], entry *jnigi.ObjectRef) (string, any, error) {

	keyRef, valueRef, err := getKeyAndValueReferences(env, entry)

	if err != nil {
		return "", nil, err
	}

	defer env.DeleteLocalRef(keyRef)
	defer env.DeleteLocalRef(valueRef)

	key, err := iterator.fromJava(keyRef, env)
	if err != nil {
		return "", nil, err
	}

	value, err := iterator.fromJava(valueRef, env)

	if err != nil {
		return "", nil, err
	}

	stringKey, err := anyToString(key)

	if err != nil {
		return "", nil, err
	}

	return stringKey, value, nil
}

func getKeyAndValueReferences(env *jnigi.Env, entry *jnigi.ObjectRef) (*jnigi.ObjectRef, *jnigi.ObjectRef, error) {

	defer env.DeleteLocalRef(entry)

	keyRef := jnigi.NewObjectRef(ObjectJniRepresentation)
	valueRef := jnigi.NewObjectRef(ObjectJniRepresentation)

	if err := entry.CallMethod(env, "getKey", keyRef); err != nil {
		return nil, nil, err
	}

	if err := entry.CallMethod(env, "getValue", valueRef); err != nil {
		return nil, nil, err
	}

	return keyRef, valueRef, nil
}
