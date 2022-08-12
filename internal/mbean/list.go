package mbean

import "tekao.net/jnigi"

func createJavaList(env *jnigi.Env, list []any) (*jnigi.ObjectRef, error) {
	return nil, nil
}

func createGoArray(list *jnigi.ObjectRef, env *jnigi.Env, dest *[]any) error {
	return nil
}
