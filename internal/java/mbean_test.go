package java

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var java *Java

func TestMain(m *testing.M) {

	java, _ = CreateJvm()

	m.Run()

	java.ShutdownJvm()
}

func TestCanInitializeConnectionToRemoteJVM(t *testing.T) {
	if os.Getenv("TEST_ENV") != "IT" {
		return
	}

	mbean := &MBean{
		Java: java,
	}
	err := mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	defer mbean.Close()
	assert.Nil(t, err)
}

func TestCanInitializeTheJVMMultipleTimes(t *testing.T) {
	if os.Getenv("TEST_ENV") != "IT" {
		return
	}
	mbean := &MBean{
		Java: java,
	}
	mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	mbean.Close()

	mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	mbean.Close()
}

func TestCanCallIntoJmxAndGetResult(t *testing.T) {
	if os.Getenv("TEST_ENV") != "IT" {
		return
	}

	mbean := &MBean{
		Java: java,
	}
	mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	defer mbean.Close()

	operation := MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: "putString",
		Args: []MBeanOperationArgs{
			{
				Value: "messi",
				Type:  "java.lang.String",
			},
			{
				Value: "fan369",
				Type:  "java.lang.String",
			},
		},
	}

	_, err := mbean.Execute(operation)

	assert.Nil(t, err)

	operation = MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: "getString",
		Args: []MBeanOperationArgs{
			{
				Value: "messi",
				Type:  "java.lang.String",
			},
		},
	}

	result, err := mbean.Execute(operation)
	assert.Nil(t, err)
	assert.Equal(t, "fan369", result)
}

func TestOnConnectionErrors(t *testing.T) {
	if os.Getenv("TEST_ENV") != "IT" {
		return
	}

	mbean := &MBean{
		Java: java,
	}
	err := mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9999/jmxrmi")

	expected := "failed to create a JMX connection Factory::java.io.IOException: Failed to retrieve RMIServer stub: javax.naming.ServiceUnavailableException [Root exception is java.rmi.ConnectException: Connection refused to host: 127.0.0.1; nested exception is: \n\tjava.net.ConnectException: Connection refused]"
	assert.Equal(t, expected, err.Error())
}
