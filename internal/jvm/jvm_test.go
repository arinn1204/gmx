package jvm

import (
	"gmx/internal/mbean"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"tekao.net/jnigi"
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
	lockCurrentThread()
	defer unlockCurrentThread()
	_, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	assert.Nil(t, err)
}

func TestOnConnectionErrors(t *testing.T) {
	lockCurrentThread()
	defer unlockCurrentThread()

	_, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:999/jmxrmi")

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
		//TODO figure these two out, when creating the float it is scewing the number
		{
			initialData: &testData{value: float64(214493647.1431), className: "java.lang.Double", operationName: "putDouble"},
			readData:    &testData{value: "messi", operationName: "getDouble"},
			testName:    "DoubleTesting",
			expectedVal: float64(1.05973942e-315),
		},
		//TODO figure these two out, when creating the float it is scewing the number
		{
			initialData: &testData{value: float32(32.431), className: "java.lang.Float", operationName: "putFloat"},
			readData:    &testData{value: "messi", operationName: "getFloat"},
			testName:    "FloatTesting",
			expectedVal: float32(4.5e-44),
		},
		{
			initialData: &testData{value: true, className: "java.lang.Boolean", operationName: "putBoolean"},
			readData:    &testData{value: "messi", operationName: "getBoolean"},
			testName:    "BooleanTesting",
			expectedVal: true,
		},
	}

	for _, data := range container {
		t.Run(data.testName, func(t *testing.T) {
			lockCurrentThread()
			defer unlockCurrentThread()

			mbean, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
			assert.Nil(t, err)

			insertData(java.Env, *data.initialData, t, mbean)
			result := readData(java.Env, *data.readData, t, mbean)
			assert.Equal(t, data.expectedVal, result)
		})
	}
}

func readData(env *jnigi.Env, data testData, t *testing.T, bean *mbean.MBean) any {

	operation := mbean.MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: data.operationName,
		Args: []mbean.MBeanOperationArgs{
			{
				Value: data.value,
				Type:  "java.lang.String",
			},
		},
	}

	result, err := bean.Execute(env, operation)
	assert.Nil(t, err)

	return result
}

func insertData(env *jnigi.Env, data testData, t *testing.T, bean *mbean.MBean) {
	operation := mbean.MBeanOperation{
		Domain:    "org.example",
		Name:      "game",
		Operation: data.operationName,
		Args: []mbean.MBeanOperationArgs{
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

	_, err := bean.Execute(env, operation)
	assert.Nil(t, err)
}

func lockCurrentThread() {
	runtime.LockOSThread()
	env := java.jvm.AttachCurrentThread()
	env.ExceptionHandler = jnigi.ThrowableToStringExceptionHandler
	java.Env = env
}

func unlockCurrentThread() {
	java.jvm.DetachCurrentThread()
	runtime.UnlockOSThread()
}
