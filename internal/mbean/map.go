package mbean

import (
	"encoding/json"
	"fmt"

	"tekao.net/jnigi"
)

func createMapFromJSON(value string, env *jnigi.Env) (*jnigi.ObjectRef, error) {
	dest := make(map[any]any)
	if err := json.Unmarshal([]byte(value), &dest); err != nil {
		return nil, fmt.Errorf("failed to convert %s to a map::%s", value, err)
	}

	return createJavaMap(env, dest)
}

func createJavaMap(env *jnigi.Env, dict map[any]any) (*jnigi.ObjectRef, error) {
	return nil, nil
}

func createGoMap(obj *jnigi.ObjectRef, env *jnigi.Env, dest *map[any]any) error {
	return nil
}
