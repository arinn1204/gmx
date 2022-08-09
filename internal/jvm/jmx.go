package jvm

import (
	"errors"
	"fmt"
	"gmx/internal/mbean"
	"runtime"

	"tekao.net/jnigi"
)

// CreateMBeanConnection is the factory method responsible for creating MBean connections based on the provided URI
// They can be created multiple per thread, or in parallel threads. They are bound to OS threads
// They should be used with caution for this reason
func CreateMBeanConnection(java *Java, uri string) (*mbean.Client, error) {

	runtime.LockOSThread()
	env := java.jvm.AttachCurrentThread()
	configureEnvironment(env)

	jmxConnector, err := buildJMXConnector(env, uri)

	if err != nil {
		if jmxConnector != nil {
			jmxConnector.CallMethod(env, "close", nil)
		}
		return nil, err
	}

	mBeanServerConnector := jnigi.NewObjectRef("javax/management/MBeanServerConnection")
	if err = jmxConnector.CallMethod(env, "getMBeanServerConnection", mBeanServerConnector); err != nil {
		return nil, errors.New("failed to create the mbean server connection::" + err.Error())
	}

	mbean := &mbean.Client{
		JmxConnection: jmxConnector,
		Env:           env,
	}

	return mbean, err
}

func buildJMXConnector(env *jnigi.Env, jndiUri string) (*jnigi.ObjectRef, error) {
	stringRef, err := env.NewObject("java/lang/String", []byte(jndiUri))

	if err != nil {
		return nil, fmt.Errorf("failed to create a string from %s::%s", jndiUri, err.Error())
	}

	jmxUrl, err := env.NewObject("javax/management/remote/JMXServiceURL", stringRef)
	if err != nil {
		return nil, errors.New("failed to create JMXServiceURL::" + err.Error())
	}

	if err != nil {
		return nil, errors.New("failed to create a blank map::" + err.Error())
	}

	jmxConnector := jnigi.NewObjectRef("javax/management/remote/JMXConnector")

	connectorFactory := "javax/management/remote/JMXConnectorFactory"
	if err = env.CallStaticMethod(connectorFactory, "connect", jmxConnector, jmxUrl); err != nil {
		return nil, errors.New("failed to create a JMX connection Factory::" + err.Error())
	}

	return jmxConnector, nil
}
