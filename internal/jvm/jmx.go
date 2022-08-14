package jvm

import (
	"errors"
	"fmt"

	"github.com/arinn1204/gmx/internal/mbean"
	"github.com/arinn1204/gmx/pkg/extensions"

	"tekao.net/jnigi"
)

// CreateMBeanConnection is the factory method responsible for creating MBean connections based on the provided URI
// They can be created multiple per thread, or in parallel threads. They are bound to OS threads
// They should be used with caution for this reason
func (java *Java) CreateMBeanConnection(uri string) (mbean.BeanExecutor, error) {

	env := java.Attach()

	jmxConnector, err := buildJMXConnector(env, uri)

	if err != nil {
		if jmxConnector != nil {
			jmxConnector.CallMethod(env, "close", nil)
		}
		return nil, err
	}

	mbean := &mbean.Client{
		JmxConnection: jmxConnector,
		Env:           env,
		ClassHandlers: make(map[string]extensions.IHandler),
	}

	return mbean, err
}

func buildJMXConnector(env *jnigi.Env, jndiURI string) (*jnigi.ObjectRef, error) {
	stringRef, err := env.NewObject("java/lang/String", []byte(jndiURI))

	if err != nil {
		return nil, fmt.Errorf("failed to create a string from %s::%s", jndiURI, err.Error())
	}

	jmxURL, err := env.NewObject("javax/management/remote/JMXServiceURL", stringRef)
	if err != nil {
		return nil, errors.New("failed to create JMXServiceURL::" + err.Error())
	}

	if err != nil {
		return nil, errors.New("failed to create a blank map::" + err.Error())
	}

	jmxConnector := jnigi.NewObjectRef("javax/management/remote/JMXConnector")

	connectorFactory := "javax/management/remote/JMXConnectorFactory"
	if err = env.CallStaticMethod(connectorFactory, "connect", jmxConnector, jmxURL); err != nil {
		return nil, errors.New("failed to create a JMX connection Factory::" + err.Error())
	}

	return jmxConnector, nil
}
