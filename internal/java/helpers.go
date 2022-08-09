package java

import "tekao.net/jnigi"

func deleteReference(mbean *MBean, param *jnigi.ObjectRef) {
	if param != nil {
		mbean.Java.Env.DeleteLocalRef(param)
	}
}

func closeReferences(env *jnigi.Env, reference *jnigi.ObjectRef) {
	reference.CallMethod(env, "close", nil)
}
