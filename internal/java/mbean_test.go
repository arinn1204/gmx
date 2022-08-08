package java

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var java *Java

func TestMain(m *testing.M) {

	java, _ = CreateJvm()

	if os.Getenv("TEST_ENV") == "IT" {
		m.Run()
	}

	java.ShutdownJvm()
}

func TestCanInitializeConnectionToRemoteJVM(t *testing.T) {
	mbean := &MBean{
		Java: java,
	}
	err := mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	defer mbean.Close()
	assert.Nil(t, err)
}

func TestCanInitializeTheJVMMultipleTimes(t *testing.T) {
	mbean := &MBean{
		Java: java,
	}
	mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	mbean.Close()

	mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	mbean.Close()
}

func TestCanCallIntoJmxAndGetResult(t *testing.T) {
	mbean := &MBean{
		Java: java,
	}
	mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	defer mbean.Close()

	type testData struct {
		value         string
		className     string
		operationName string
	}

	type testDataContainer struct {
		initialData *testData
		readData    *testData
		testName    string
		expectedVal any
	}

	container := []testDataContainer{
		{
			initialData: &testData{value: "fan369", className: "java.lang.String", operationName: "putString"},
			readData:    &testData{value: "messi", operationName: "getString"},
			testName:    "StringTesting",
			expectedVal: "fan369",
		},
	}

	for _, data := range container {
		t.Run(data.testName, func(t *testing.T) {
			initialData := data.initialData
			insertData(initialData.value, initialData.className, initialData.operationName, t, mbean)
			result := readData(data.readData.value, data.readData.operationName, t, mbean)
			assert.Equal(t, "fan369", result)
		})
	}
}

func readData(value string, operationName string, t *testing.T, mbean *MBean) any {

	operation := MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: operationName,
		Args: []MBeanOperationArgs{
			{
				Value: value,
				Type:  "java.lang.String",
			},
		},
	}

	result, err := mbean.Execute(operation)
	assert.Nil(t, err)

	return result
}

func insertData(value string, className string, operationName string, t *testing.T, mbean *MBean) {
	operation := MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: operationName,
		Args: []MBeanOperationArgs{
			{
				Value: "messi",
				Type:  "java.lang.String",
			},
			{
				Value: value,
				Type:  className,
			},
		},
	}

	_, err := mbean.Execute(operation)
	assert.Nil(t, err)
}

func TestOnConnectionErrors(t *testing.T) {
	mbean := &MBean{
		Java: java,
	}
	err := mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9999/jmxrmi")

	expected := "failed to create a JMX connection Factory::java.io.IOException: Failed to retrieve RMIServer stub: javax.naming.ServiceUnavailableException [Root exception is java.rmi.ConnectException: Connection refused to host: 127.0.0.1; nested exception is: \n\tjava.net.ConnectException: Connection refused]"
	assert.Equal(t, expected, err.Error())
}
