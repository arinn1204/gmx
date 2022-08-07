package java

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

func buildMBeanServerConnection(java *Java, jndiUri string) (*jnigi.ObjectRef, error) {
	stringRef, err := java.env.NewObject("java/lang/String", []byte(jndiUri))

	if err != nil {
		return nil, fmt.Errorf("failed to create a string from %s::%s", jndiUri, err.Error())
	}

	jmxUrl, err := java.env.NewObject("javax/management/remote/JMXServiceURL", stringRef)
	if err != nil {
		return nil, errors.New("failed to create JMXServiceURL::" + err.Error())
	}

	if err != nil {
		return nil, errors.New("failed to create a blank map::" + err.Error())
	}

	jmxConnector := jnigi.NewObjectRef("javax/management/remote/JMXConnector")

	err = java.env.CallStaticMethod(
		"javax/management/remote/JMXConnectorFactory",
		"connect",
		jmxConnector,
		jmxUrl)

	if err != nil {
		return nil, errors.New("failed to create a JMX connection Factory::" + err.Error())
	}

	mBeanServerConnector := jnigi.NewObjectRef("javax/management/MBeanServerConnection")
	err = jmxConnector.CallMethod(
		java.env,
		"getMBeanServerConnection",
		mBeanServerConnector)

	if err != nil {
		return nil, errors.New("failed to create the mbean server connection::" + err.Error())
	}

	return mBeanServerConnector, nil
}
