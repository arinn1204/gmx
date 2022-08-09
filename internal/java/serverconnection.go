package java

import (
	"errors"
	"fmt"
	"gmx/internal/jvm"

	"tekao.net/jnigi"
)

func buildJMXConnector(java *jvm.Java, jndiUri string) (*jnigi.ObjectRef, error) {
	stringRef, err := java.Env.NewObject("java/lang/String", []byte(jndiUri))

	if err != nil {
		return nil, fmt.Errorf("failed to create a string from %s::%s", jndiUri, err.Error())
	}

	jmxUrl, err := java.Env.NewObject("javax/management/remote/JMXServiceURL", stringRef)
	if err != nil {
		return nil, errors.New("failed to create JMXServiceURL::" + err.Error())
	}

	if err != nil {
		return nil, errors.New("failed to create a blank map::" + err.Error())
	}

	jmxConnector := jnigi.NewObjectRef("javax/management/remote/JMXConnector")

	connectorFactory := "javax/management/remote/JMXConnectorFactory"
	if err = java.Env.CallStaticMethod(connectorFactory, "connect", jmxConnector, jmxUrl); err != nil {
		return nil, errors.New("failed to create a JMX connection Factory::" + err.Error())
	}

	return jmxConnector, nil
}
