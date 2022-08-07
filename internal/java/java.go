package java

import (
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

func BuildMbeanConnection(uri string) (*MBean, error) {
	java, err := CreateJvm()
	if err != nil {
		return nil, err
	}

	connection, err := buildMBeanServerConnection(java, uri)

	if err != nil {
		return nil, err
	}

	bean := &MBean{
		ServerConnection: connection,
		Java:             java,
	}

	return bean, nil
}

func (mbean *MBean) Close() {
	mbean.Java.ShutdownJvm()
}
