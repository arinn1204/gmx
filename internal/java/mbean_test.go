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

func TestOnConnectionErrors(t *testing.T) {
	mbean := &MBean{
		Java: java,
	}
	err := mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9999/jmxrmi")

	expected := "failed to create a JMX connection Factory::java.io.IOException: Failed to retrieve RMIServer stub: javax.naming.ServiceUnavailableException [Root exception is java.rmi.ConnectException: Connection refused to host: 127.0.0.1; nested exception is: \n\tjava.net.ConnectException: Connection refused]"
	assert.Equal(t, expected, err.Error())
}

type testData struct {
	value         any
	className     string
	operationName string
}

type testDataContainer struct {
	initialData *testData
	readData    *testData
	testName    string
	expectedVal any
}

func TestCanCallIntoJmxAndGetResult(t *testing.T) {
	mbean := &MBean{
		Java: java,
	}
	mbean.InitializeMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	defer mbean.Close()

	container := []testDataContainer{
		{
			initialData: &testData{value: "fan369", className: "java.lang.String", operationName: "putString"},
			readData:    &testData{value: "messi", operationName: "getString"},
			testName:    "StringTesting",
			expectedVal: "fan369",
		},
		{
			initialData: &testData{value: int64(2148493647), className: "java.lang.Long", operationName: "putLong"},
			readData:    &testData{value: "messi", operationName: "getLong"},
			testName:    "LongTesting",
			expectedVal: int64(2148493647),
		},
		{
			initialData: &testData{value: 214493647, className: "java.lang.Integer", operationName: "putInteger"},
			readData:    &testData{value: "messi", operationName: "getInteger"},
			testName:    "IntegerTesting",
			expectedVal: 214493647,
		},
	}

	for _, data := range container {
		t.Run(data.testName, func(t *testing.T) {
			insertData(*data.initialData, t, mbean)
			result := readData(*data.readData, t, mbean)
			assert.Equal(t, data.expectedVal, result)
		})
	}
}

func readData(data testData, t *testing.T, mbean *MBean) any {

	operation := MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: data.operationName,
		Args: []MBeanOperationArgs{
			{
				Value: data.value,
				Type:  "java.lang.String",
			},
		},
	}

	result, err := mbean.Execute(operation)
	assert.Nil(t, err)

	return result
}

func insertData(data testData, t *testing.T, mbean *MBean) {
	operation := MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: data.operationName,
		Args: []MBeanOperationArgs{
			{
				Value: "messi",
				Type:  "java.lang.String",
			},
			{
				Value: data.value,
				Type:  data.className,
			},
		},
	}

	_, err := mbean.Execute(operation)
	assert.Nil(t, err)
}
