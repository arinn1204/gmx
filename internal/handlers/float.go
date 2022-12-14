package handlers

import (
	"errors"
	"strconv"

	"github.com/arinn1204/gmx/pkg/extensions"
	"tekao.net/jnigi"
)

// These are the constants for the Float classpath and JNI representation
const (
	FloatJniRepresentation = "java/lang/Float"
	FloatClasspath         = "java.lang.Float"
)

// FloatHandler is the type that can convert to and from java.lang.Float
type FloatHandler struct{}

var strHandler extensions.IHandler

func init() {
	strHandler = &StringHandler{}
}

// ToJniRepresentation is the implementation that will convert from a go type
// to a JNI representation of that type
func (handler *FloatHandler) ToJniRepresentation(env *jnigi.Env, value any) (*jnigi.ObjectRef, error) {
	stringifiedValue := strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)

	strRef, err := strHandler.ToJniRepresentation(env, stringifiedValue)

	if err != nil {
		return nil, err
	}

	defer env.DeleteLocalRef(strRef)

	floatRef := jnigi.NewObjectRef(FloatJniRepresentation)

	if err = env.CallStaticMethod(FloatJniRepresentation, "valueOf", floatRef, strRef); err != nil {
		return nil, errors.New("failed to convert to a float::" + err.Error())
	}

	return floatRef, nil
}

// ToGoRepresentation will convert from a JNI type to a go type
func (handler *FloatHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest any) error {
	if err := object.CallMethod(env, "floatValue", dest); err != nil {
		return errors.New("failed to create a float::" + err.Error())
	}
	return nil
}
