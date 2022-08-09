package java

import (
	"errors"

	"tekao.net/jnigi"
)

func (mbean *MBean) InitializeMBeanConnection(uri string) error {

	jmxConnector, err := buildJMXConnector(mbean.Java, uri)

	if err != nil {
		if jmxConnector != nil {
			closeReferences(mbean.Java.Env, jmxConnector)
		}
		return err
	}

	mBeanServerConnector := jnigi.NewObjectRef("javax/management/MBeanServerConnection")
	err = jmxConnector.CallMethod(
		mbean.Java.Env,
		"getMBeanServerConnection",
		mBeanServerConnector)

	if err != nil {
		return errors.New("failed to create the mbean server connection::" + err.Error())
	}

	mbean.serverConnection = mBeanServerConnector
	mbean.jmxConnection = jmxConnector

	return err
}

func closeReferences(env *jnigi.Env, reference *jnigi.ObjectRef) {
	reference.CallMethod(env, "close", nil)
}
