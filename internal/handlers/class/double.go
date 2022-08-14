package class

import (
	"errors"
	"strconv"

	"tekao.net/jnigi"
)

const (
	DOUBLE = "java/lang/Double"
)

type doubleHandler struct{}

func (handler *doubleHandler) toJniRepresentation(env *jnigi.Env, value float64) (*jnigi.ObjectRef, error) {
	stringifiedValue := strconv.FormatFloat(value, 'f', -1, 32)
	strHandler := &stringHandler{}

	strRef, err := strHandler.toJniRepresentation(env, stringifiedValue)

	if err != nil {
		return nil, err
	}

	defer env.DeleteLocalRef(strRef)

	floatRef := jnigi.NewObjectRef(DOUBLE)

	if err = env.CallStaticMethod(DOUBLE, "valueOf", floatRef, strRef); err != nil {
		return nil, errors.New("failed to convert to a double::" + err.Error())
	}

	return floatRef, nil
}

func (handler *doubleHandler) toGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *float64) error {
	if err := object.CallMethod(env, "doubleValue", dest); err != nil {
		return errors.New("failed to create a float::" + err.Error())
	}

	return nil
}
