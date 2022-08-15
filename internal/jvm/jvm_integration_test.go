package jvm

import (
	"log"
	"runtime"
	"sync"
	"testing"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/arinn1204/gmx/internal/mbean"

	"github.com/stretchr/testify/assert"
	"tekao.net/jnigi"
)

var java *Java

func TestMain(m *testing.M) {
	java = &Java{}
	_, err := java.CreateJVM()

	if err != nil {
		log.Fatal(err)
	}

	m.Run()

	java.ShutdownJvm()
}

func TestCanInitializeConnectionToRemoteJVM(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration tests when running short mode")
	}

	lockCurrentThread(java)
	defer unlockCurrentThread(java)
	_, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:5001/jmxrmi")
	assert.Nil(t, err)
}

func TestOnConnectionErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration tests when running short mode")
	}
	lockCurrentThread(java)
	defer unlockCurrentThread(java)

	_, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9901/jmxrmi")

	expected := "failed to create a JMX connection Factory::java.io.IOException: Failed to retrieve RMIServer stub: javax.naming.ServiceUnavailableException [Root exception is java.rmi.ConnectException: Connection refused to host: 127.0.0.1; nested exception is: \n\tjava.net.ConnectException: Connection refused]"
	assert.Equal(t, expected, err.Error())
}

func TestCanConnectToMultipleMBeansSynchronously(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration tests when running short mode")
	}
	lockCurrentThread(java)
	defer unlockCurrentThread(java)

	var err error
	var mbean1 mbean.BeanExecutor
	var mbean2 mbean.BeanExecutor

	mbean1, err = java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	assert.Nil(t, err)
	registerHandlers(mbean1)

	mbean2, err = java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
	assert.Nil(t, err)
	registerHandlers(mbean2)

	testData := []testDataContainer{
		{
			initialData: &testData{value: "fan369", className: "java.lang.String", operationName: "putString"},
			readData:    &testData{value: "messi", operationName: "getString"},
			testName:    "StringTesting",
			expectedVal: "fan369",
		},
		{
			initialData: &testData{value: "2148493647", className: "java.lang.Long", operationName: "putLong"},
			readData:    &testData{value: "messi", operationName: "getLong"},
			testName:    "LongTesting",
			expectedVal: "2148493647",
		},
	}

	var res any

	insertData(mbean1.GetEnv(), *testData[0].initialData, t, mbean1)
	res = readData(mbean1.GetEnv(), *testData[0].readData, t, mbean1)

	assert.Equal(t, testData[0].expectedVal, res)

	insertData(mbean1.GetEnv(), *testData[1].initialData, t, mbean2)
	res = readData(mbean2.GetEnv(), *testData[1].readData, t, mbean2)

	assert.Equal(t, testData[1].expectedVal, res)
}

func TestCreatingFloats(t *testing.T) {
	lockCurrentThread(java)
	defer unlockCurrentThread(java)

	value := float32(3.1415)
	stringRef, err := java.Env.NewObject(mbean.STRING, []byte("3.1415"))

	assert.Nil(t, err)
	floatRef := jnigi.NewObjectRef(mbean.FLOAT)
	err = java.Env.CallStaticMethod(mbean.FLOAT, "valueOf", floatRef, stringRef)

	assert.Nil(t, err)

	res32 := float32(0)

	err = floatRef.CallMethod(java.Env, "floatValue", &res32)

	assert.Nil(t, err)
	assert.Equal(t, value, res32)
}

func TestCanConnectToMultipleMBeansAsynchronously(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration tests when running short mode")
	}
	wg := &sync.WaitGroup{}

	lockCurrentThread(java)
	defer unlockCurrentThread(java)

	wg.Add(1)
	go func(t *testing.T, wg *sync.WaitGroup) {
		defer wg.Done()
		lockCurrentThread(java)
		defer unlockCurrentThread(java)

		mbean, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
		assert.Nil(t, err)
		registerHandlers(mbean)

		testData := testDataContainer{
			initialData: &testData{value: "fan369", className: "java.lang.String", operationName: "putString"},
			readData:    &testData{value: "messi", operationName: "getString"},
			testName:    "StringTesting",
			expectedVal: "fan369",
		}

		insertData(mbean.GetEnv(), *testData.initialData, t, mbean)
		res := readData(mbean.GetEnv(), *testData.readData, t, mbean)

		assert.Equal(t, testData.expectedVal, res)
	}(t, wg)

	wg.Add(1)
	go func(t *testing.T, wg *sync.WaitGroup) {
		defer wg.Done()

		lockCurrentThread(java)
		defer unlockCurrentThread(java)
		mbean, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:5001/jmxrmi")
		assert.Nil(t, err)
		registerHandlers(mbean)

		testData := testDataContainer{
			initialData: &testData{value: "2148493647", className: "java.lang.Long", operationName: "putLong"},
			readData:    &testData{value: "messi", operationName: "getLong"},
			testName:    "LongTesting",
			expectedVal: "2148493647",
		}

		insertData(mbean.GetEnv(), *testData.initialData, t, mbean)
		res := readData(mbean.GetEnv(), *testData.readData, t, mbean)

		assert.Equal(t, testData.expectedVal, res)
	}(t, wg)

	wg.Wait()
}

type testData struct {
	value         string
	className     string
	containerName string
	operationName string
}

type testDataContainer struct {
	initialData *testData
	readData    *testData
	testName    string
	expectedVal any
}

func TestCanCallIntoJmxAndGetResult(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration tests when running short mode")
	}

	container := []testDataContainer{
		/*
			PRIMITIVE TESTING
		*/
		{
			initialData: &testData{value: "fan369", className: "java.lang.String", operationName: "putString"},
			readData:    &testData{value: "messi", operationName: "getString"},
			testName:    "StringTesting",
			expectedVal: "fan369",
		},
		{
			initialData: &testData{value: "2148493647", className: "java.lang.Long", operationName: "putLong"},
			readData:    &testData{value: "messi", operationName: "getLong"},
			testName:    "LongTesting",
			expectedVal: "2148493647",
		},
		{
			initialData: &testData{value: "214493647", className: "java.lang.Integer", operationName: "putInteger"},
			readData:    &testData{value: "messi", operationName: "getInteger"},
			testName:    "IntegerTesting",
			expectedVal: "214493647",
		},
		{
			initialData: &testData{value: "214493647.1431", className: "java.lang.Double", operationName: "putDouble"},
			readData:    &testData{value: "messi", operationName: "getDouble"},
			testName:    "DoubleTesting",
			expectedVal: "214493647.1431",
		},
		{
			initialData: &testData{value: "32.431", className: "java.lang.Float", operationName: "putFloat"},
			readData:    &testData{value: "messi", operationName: "getFloat"},
			testName:    "FloatTesting",
			expectedVal: "32.431",
		},
		{
			initialData: &testData{value: "true", className: "java.lang.Boolean", operationName: "putBoolean"},
			readData:    &testData{value: "messi", operationName: "getBoolean"},
			testName:    "BooleanTesting",
			expectedVal: "true",
		},
		/**
		LIST TESTING
		*/
		{
			initialData: &testData{value: "[1, 2, 3]", className: "java.lang.Integer", containerName: "java.util.List", operationName: "putList"},
			readData:    &testData{value: "messi", operationName: "getList"},
			testName:    "IntListTesting",
			expectedVal: "[1,2,3]",
		},
		{
			initialData: &testData{value: "[1, 2, 3]", className: "java.lang.Long", containerName: "java.util.List", operationName: "putList"},
			readData:    &testData{value: "messi", operationName: "getList"},
			testName:    "LongListTesting",
			expectedVal: "[1,2,3]",
		},
		{
			initialData: &testData{value: "[1.31,2.431,3.6543]", className: "java.lang.Float", containerName: "java.util.List", operationName: "putList"},
			readData:    &testData{value: "messi", operationName: "getList"},
			testName:    "FloatListTesting",
			expectedVal: "[1.31,2.431,3.6543]",
		},
		{
			initialData: &testData{value: "[1.31,2.431,3.6543]", className: "java.lang.Double", containerName: "java.util.List", operationName: "putList"},
			readData:    &testData{value: "messi", operationName: "getList"},
			testName:    "DoubleListTesting",
			expectedVal: "[1.31,2.431,3.6543]",
		},
		{
			initialData: &testData{value: "[true,false,true]", className: "java.lang.Boolean", containerName: "java.util.List", operationName: "putList"},
			readData:    &testData{value: "messi", operationName: "getList"},
			testName:    "BoolListTesting",
			expectedVal: "[true,false,true]",
		},
		/**
		SET TESTING
		*/
		// {
		// 	initialData: &testData{value: "[1, 2, 3]", className: "java.lang.Integer", containerName: "java.util.Set", operationName: "putSet"},
		// 	readData:    &testData{value: "messi", operationName: "getSet"},
		// 	testName:    "IntSetTesting",
		// 	expectedVal: "[1,2,3]",
		// },
		// {
		// 	initialData: &testData{value: "[1, 2, 3]", className: "java.lang.Long", containerName: "java.util.Set", operationName: "putSet"},
		// 	readData:    &testData{value: "messi", operationName: "getSet"},
		// 	testName:    "LongSetTesting",
		// 	expectedVal: "[1,2,3]",
		// },
		// {
		// 	initialData: &testData{value: "[1.31,2.431,3.6543]", className: "java.lang.Float", containerName: "java.util.Set", operationName: "putSet"},
		// 	readData:    &testData{value: "messi", operationName: "getSet"},
		// 	testName:    "FloatSetTesting",
		// 	expectedVal: "[1.31,2.431,3.6543]",
		// },
		// {
		// 	initialData: &testData{value: "[1.31,2.431,3.6543]", className: "java.lang.Double", containerName: "java.util.Set", operationName: "putSet"},
		// 	readData:    &testData{value: "messi", operationName: "getSet"},
		// 	testName:    "DoubleSetTesting",
		// 	expectedVal: "[1.31,2.431,3.6543]",
		// },
		// {
		// 	initialData: &testData{value: "[true,false,true]", className: "java.lang.Boolean", containerName: "java.util.Set", operationName: "putSet"},
		// 	readData:    &testData{value: "messi", operationName: "getSet"},
		// 	testName:    "BoolSetTesting",
		// 	expectedVal: "[true,false,true]",
		// },
	}

	lockCurrentThread(java)
	defer unlockCurrentThread(java)

	for _, data := range container {
		t.Run(data.testName, func(t *testing.T) {
			lockCurrentThread(java)
			defer unlockCurrentThread(java)

			mbean, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
			assert.Nil(t, err)
			registerHandlers(mbean)

			insertData(java.Env, *data.initialData, t, mbean)
			result := readData(java.Env, *data.readData, t, mbean)
			assert.Equal(t, data.expectedVal, result)
		})
	}
}

func readData(env *jnigi.Env, data testData, t *testing.T, bean mbean.BeanExecutor) any {

	operation := mbean.Operation{
		Domain:    "org.example",
		Name:      "game",
		Operation: data.operationName,
		Args: []mbean.OperationArgs{
			{
				Value:    data.value,
				JavaType: "java.lang.String",
			},
		},
	}

	result, err := bean.Execute(operation)
	assert.Nil(t, err)

	return result
}

func insertData(env *jnigi.Env, data testData, t *testing.T, bean mbean.BeanExecutor) {
	operation := mbean.Operation{
		Domain:    "org.example",
		Name:      "game",
		Operation: data.operationName,
		Args: []mbean.OperationArgs{
			{
				Value:    "messi",
				JavaType: "java.lang.String",
			},
			{
				Value:             data.value,
				JavaType:          data.className,
				JavaContainerType: data.containerName,
			},
		},
	}

	_, err := bean.Execute(operation)
	assert.Nil(t, err)
}

func lockCurrentThread(java *Java) {
	runtime.LockOSThread()
	env := java.jvm.AttachCurrentThread()
	configureEnvironment(env)
	java.Env = env
}

func unlockCurrentThread(java *Java) {
	java.jvm.DetachCurrentThread()
	runtime.UnlockOSThread()
}

func registerHandlers(bean mbean.BeanExecutor) {
	bean.RegisterClassHandler(handlers.BoolClasspath, &handlers.BoolHandler{})
	bean.RegisterClassHandler(handlers.DoubleClasspath, &handlers.DoubleHandler{})
	bean.RegisterClassHandler(handlers.FloatClasspath, &handlers.FloatHandler{})
	bean.RegisterClassHandler(handlers.IntClasspath, &handlers.IntHandler{})
	bean.RegisterClassHandler(handlers.LongClasspath, &handlers.LongHandler{})
	bean.RegisterClassHandler(handlers.StringClasspath, &handlers.StringHandler{})

	client := bean.(*mbean.Client)

	bean.RegisterInterfaceHandler(handlers.ListClassPath, &handlers.ListHandler{ClassHandlers: client.ClassHandlers})
}
