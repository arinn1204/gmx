package handlers

import (
	"errors"
	"strconv"

	"tekao.net/jnigi"
)

const (
	DoubleJniRepresentation = "java/lang/Double"
	DoubleClasspath         = "java.lang.Double"
)

type DoubleHandler struct{}

func (handler *DoubleHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	stringifiedValue := strconv.FormatFloat(value.(float64), 'f', -1, 64)
	strHandler := &StringHandler{}

	strRef, err := strHandler.ToJniRepresentation(env, stringifiedValue)

	if err != nil {
		return nil, err
	}

	defer env.DeleteLocalRef(strRef)

	floatRef := jnigi.NewObjectRef(DoubleJniRepresentation)

	if err = env.CallStaticMethod(DoubleJniRepresentation, "valueOf", floatRef, strRef); err != nil {
		return nil, errors.New("failed to convert to a double::" + err.Error())
	}

	return floatRef, nil
}

func (handler *DoubleHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	if err := object.CallMethod(env, "doubleValue", dest); err != nil {
		return errors.New("failed to create a float::" + err.Error())
	}

	return nil
}
