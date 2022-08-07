package java

import (
	"errors"
	"sync"

	"tekao.net/jnigi"
)

type Java struct {
	jvm     *jnigi.JVM
	env     *jnigi.Env
	lock    *sync.Mutex
	started bool
}

type IJava interface {
	CreateJvm() (*jnigi.Env, error)
	ShutdownJvm() error
	IsStarted() bool
}

func (mbean *MBean) InitializeMBeanConnection(uri string) error {
	java, err := CreateJvm()
	if err != nil {
		return err
	}

	jmxConnector, err := buildJMXConnector(java, uri)

	if err != nil {
		if jmxConnector != nil {
			closeReferences(mbean.java.env, jmxConnector)
		}
		return err
	}

	mBeanServerConnector := jnigi.NewObjectRef("javax/management/MBeanServerConnection")
	err = jmxConnector.CallMethod(
		java.env,
		"getMBeanServerConnection",
		mBeanServerConnector)

	if err != nil {
		return errors.New("failed to create the mbean server connection::" + err.Error())
	}

	mbean.java = java
	mbean.serverConnection = mBeanServerConnector
	mbean.jmxConnection = jmxConnector

	return err
}

func closeReferences(env *jnigi.Env, reference *jnigi.ObjectRef) {
	reference.CallMethod(env, "close", nil)
}
