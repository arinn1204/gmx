package java

import "tekao.net/jnigi"

func deleteReference(mbean *MBean, param *jnigi.ObjectRef) {
	if param != nil {
		mbean.Java.env.DeleteLocalRef(param)
	}
}
