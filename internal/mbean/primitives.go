package mbean

import (
	"fmt"

	"github.com/arinn1204/gmx/internal/handlers"
	"tekao.net/jnigi"
)

func (mbean *Client) fromJava(param *jnigi.ObjectRef, env *jnigi.Env, dest *[]any) error {
	cls, err := getClassName(param, env)
	if err != nil {
		return err
	}

	handler := mbean.ClassHandlers[cls]

	switch cls {
	case handlers.StringClasspath:
		var valDest string
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, string(valDest))
	case handlers.IntClasspath:
		var valDest int
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case handlers.LongClasspath:
		var valDest int64
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case handlers.FloatClasspath:
		var valDest float32
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case handlers.DoubleClasspath:
		var valDest float64
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	case handlers.BoolClasspath:
		var valDest bool
		if err = handler.ToGoRepresentation(env, param, &valDest); err != nil {
			return err
		}

		*dest = append(*dest, valDest)
	default:
		return fmt.Errorf("no known formatter for %s", cls)
	}

	return nil
}
