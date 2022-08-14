package class

import (
	"errors"
	"strconv"

	"tekao.net/jnigi"
)

const (
	FLOAT = "java/lang/Float"
)

type floatHandler struct{}

func (handler *floatHandler) toJniRepresentation(env *jnigi.Env, value float32) (*jnigi.ObjectRef, error) {
	stringifiedValue := strconv.FormatFloat(float64(value), 'f', -1, 32)
	strHandler := &stringHandler{}

	strRef, err := strHandler.toJniRepresentation(env, stringifiedValue)

	if err != nil {
		return nil, err
	}

	defer env.DeleteLocalRef(strRef)

	floatRef := jnigi.NewObjectRef(FLOAT)

	if err = env.CallStaticMethod(FLOAT, "valueOf", floatRef, strRef); err != nil {
		return nil, errors.New("failed to convert to a float::" + err.Error())
	}

	return floatRef, nil
}

func (handler *floatHandler) toGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *float32) error {
	if err := object.CallMethod(env, "floatValue", dest); err != nil {
		return errors.New("failed to create a float::" + err.Error())
	}

	return nil
}
