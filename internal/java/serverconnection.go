package java

import (
	"errors"
	"fmt"

	"tekao.net/jnigi"
)

func buildJMXConnector(java *Java, jndiUri string) (*jnigi.ObjectRef, error) {
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

	return jmxConnector, nil
}
