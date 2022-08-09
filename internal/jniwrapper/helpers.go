package jniwrapper

import (
	"fmt"

	"tekao.net/jnigi"
)

// the commonly used types
const (
	STRING  = "java/lang/String"
	OBJECT  = "java/lang/Object"
	LONG    = "java/lang/Long"
	INTEGER = "java/lang/Integer"
	BOOLEAN = "java/lang/Boolean"
	FLOAT   = "java/lang/Float"
	DOUBLE  = "java/lang/Double"
)

func CreateString(env *jnigi.Env, str string) (*jnigi.ObjectRef, error) {
	stringRef, err := env.NewObject(STRING, []byte(str))
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", str, err.Error())
	}

	return stringRef, nil
}

func DeleteLocalReference(env *jnigi.Env, param *jnigi.ObjectRef) {
	if param != nil {
		env.DeleteLocalRef(param)
	}
}

func CreateJavaNative(env *jnigi.Env, obj any, typeName string) (*jnigi.ObjectRef, error) {
	ref, err := env.NewObject(typeName, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into an object::%s", obj, err.Error())
	}

	return ref, nil
}
