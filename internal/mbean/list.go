package mbean

import (
	"tekao.net/jnigi"
)

func (mbean *Client) createGoArrayFromList(param *jnigi.ObjectRef, env *jnigi.Env, dest *[]any) error {
	iterator, err := getIterator(env, param)

	if err != nil {
		return err
	}

	for iterator.hasNext(env) {
		value, err := iterator.getNext(env)
		if err != nil {
			return err
		}
		if err = mbean.fromJava(value, env, dest); err != nil {
			return err
		}
	}

	return nil
}
