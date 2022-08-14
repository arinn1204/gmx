package handlers

import (
	"errors"
	"strconv"

	"tekao.net/jnigi"
)

const (
	JNI_FLOAT = "java/lang/Float"
)

type FloatHandler struct{}

func (handler *FloatHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	stringifiedValue := strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
	strHandler := &StringHandler{}

	strRef, err := strHandler.ToJniRepresentation(env, stringifiedValue)

	if err != nil {
		return nil, err
	}

	defer env.DeleteLocalRef(strRef)

	floatRef := jnigi.NewObjectRef(JNI_FLOAT)

	if err = env.CallStaticMethod(JNI_FLOAT, "valueOf", floatRef, strRef); err != nil {
		return nil, errors.New("failed to convert to a float::" + err.Error())
	}

	return floatRef, nil
}

func (handler *FloatHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *any) error {
	val := float32(0)
	if err := object.CallMethod(env, "floatValue", &val); err != nil {
		return errors.New("failed to create a float::" + err.Error())
	}

	*dest = val

	return nil
}
