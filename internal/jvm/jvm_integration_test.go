package jvm

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
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

func TestCanConnectToMultipleMBeansSynchronously(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration tests when running short mode")
	}
	lockCurrentThread(java)
	defer unlockCurrentThread(java)

	mbean1 := &mbean.Client{
		JmxURI:            "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi",
		ClassHandlers:     sync.Map{},
		InterfaceHandlers: sync.Map{},
		Env:               java.Env,
	}

	registerHandlers(mbean1)
	mbean2 := &mbean.Client{
		JmxURI:            "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi",
		ClassHandlers:     sync.Map{},
		InterfaceHandlers: sync.Map{},
		Env:               java.Env,
	}

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

		mbean := &mbean.Client{
			JmxURI:            "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi",
			ClassHandlers:     sync.Map{},
			InterfaceHandlers: sync.Map{},
			Env:               java.Env,
		}

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

		mbean := &mbean.Client{
			JmxURI:            "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi",
			ClassHandlers:     sync.Map{},
			InterfaceHandlers: sync.Map{},
			Env:               java.Env,
		}

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

func toString(value any, t *testing.T) string {

	switch value := value.(type) {
	case int32:
		return fmt.Sprintf("%d", value)
	case int64:
		return fmt.Sprintf("%d", value)
	case bool:
		return fmt.Sprintf("%t", value)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(float64(value), 'f', -1, 64)
	case string:
		return value
	case []string, []float64, []float32, []bool, []int64, []int32:
		b, err := json.Marshal(value)
		assert.Nil(t, err)
		return string(b)
	default:
		name := reflect.TypeOf(value).Name()
		assert.Fail(t, fmt.Sprintf("unkown type %s", name))
		return ""
	}
}

func TestCanCallIntoJmxAndGetResultWithMapsThatHaveInterfaceValues(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration tests when running short mode")
	}

	collections := []string{"List", "Set"}

	valueMapping := map[string]any{
		"List": []int32{int32(rand.Int31()), int32(rand.Int31())},
		"Set":  []string{"hello", "world", "!!"},
	}

	typeMapping := map[string]string{
		"List": "Integer",
		"Set":  "String",
	}

	for _, collectionType := range collections {
		innerType := typeMapping[collectionType]
		t.Run(fmt.Sprintf("TestJmxAndGetResultsFor_AdvancedMap<String, %s>", innerType), func(t *testing.T) {

			lockCurrentThread(java)
			defer unlockCurrentThread(java)
			values := valueMapping[collectionType]
			str := toString(values, t)

			data := testData{value: str, className: fmt.Sprintf("java.lang.%s", innerType), containerName: fmt.Sprintf("java.util.%s", collectionType), operationName: fmt.Sprintf("put%s", collectionType)}

			mbean := &mbean.Client{
				JmxURI:            "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi",
				ClassHandlers:     sync.Map{},
				InterfaceHandlers: sync.Map{},
				Env:               java.Env,
			}

			registerHandlers(mbean)

			insertData(java.Env, data, t, mbean)

			data = testData{value: collectionType, operationName: "getMap"}

			stringData := readData(java.Env, data, t, mbean)

			if innerType == "Integer" {
				dict := make(map[string][]int32)
				err := json.Unmarshal([]byte(stringData), &dict)
				assert.Nil(t, err)

				arrayEqual(t, values.([]int32), dict["messi"])

			} else if innerType == "String" {
				dict := make(map[string][]string)
				err := json.Unmarshal([]byte(stringData), &dict)
				assert.Nil(t, err)

				assert.Equal(t, len(values.([]string)), len(dict["messi"]))
			}

		})
	}
}

func arrayEqual(t *testing.T, left any, right any) {
	assert.Equal(t, reflect.TypeOf(left), reflect.TypeOf(right))

	switch left := left.(type) {
	case []int32:
		assert.Equal(t, len(left), len(right.([]int32)))
		foundIdentical := 0
		for _, l := range left {
			for _, r := range right.([]int32) {
				if l == r {
					foundIdentical++
				}
			}
		}

		assert.Equal(t, len(left), foundIdentical)
	case []string:
		assert.Equal(t, len(left), len(right.([]string)))
		foundIdentical := 0
		for _, l := range left {
			for _, r := range right.([]string) {
				if l == r {
					foundIdentical++
				}
			}
		}

		assert.Equal(t, len(left), foundIdentical)
	}

}

func TestCanCallIntoJmxAndGetResultWithBasicMaps(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration tests when running short mode")
	}
	floatValues := rand.Float32()
	doubleValues := rand.Float64()
	intValues := int32(rand.Int31())
	longValues := int64(rand.Int63())
	boolValues := true
	stringValues := "hello"

	valueMapping := map[string]any{
		"Integer": intValues,
		"Float":   floatValues,
		"Double":  doubleValues,
		"Long":    longValues,
		"Boolean": boolValues,
		"String":  stringValues,
	}
	primitiveTypes := []string{"Integer", "Long", "Float", "Double", "Boolean", "String"}

	for _, valueType := range primitiveTypes {
		t.Run(fmt.Sprintf("TestJmxAndGetResultsFor_Map<String,%s>", valueType), func(t *testing.T) {
			lockCurrentThread(java)
			defer unlockCurrentThread(java)

			values := valueMapping[valueType]

			str := toString(values, t)

			className := fmt.Sprintf("java.lang.%s", valueType)

			data := testData{value: str, className: className, operationName: fmt.Sprintf("put%s", valueType)}

			mbean := &mbean.Client{
				JmxURI:            "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi",
				ClassHandlers:     sync.Map{},
				InterfaceHandlers: sync.Map{},
				Env:               java.Env,
			}

			registerHandlers(mbean)
			var err error

			insertData(java.Env, data, t, mbean)

			data = testData{value: valueType, operationName: "getMap"}

			stringData := readData(java.Env, data, t, mbean)

			dest := make(map[string]any)
			switch className {
			case handlers.FloatClasspath:
				var typedDest map[string]float32
				err = json.Unmarshal([]byte(stringData), &typedDest)
				for k, v := range typedDest {
					dest[k] = v
				}
			case handlers.LongClasspath:
				var typedDest map[string]int64
				err = json.Unmarshal([]byte(stringData), &typedDest)
				for k, v := range typedDest {
					dest[k] = v
				}
			case handlers.IntClasspath:
				var typedDest map[string]int32
				err = json.Unmarshal([]byte(stringData), &typedDest)
				for k, v := range typedDest {
					dest[k] = v
				}
			case handlers.StringClasspath:
				var typedDest map[string]string
				err = json.Unmarshal([]byte(stringData), &typedDest)
				for k, v := range typedDest {
					dest[k] = v
				}
			case handlers.BoolClasspath:
				var typedDest map[string]bool
				err = json.Unmarshal([]byte(stringData), &typedDest)
				for k, v := range typedDest {
					dest[k] = v
				}
			case handlers.DoubleClasspath:
				var typedDest map[string]float64
				err = json.Unmarshal([]byte(stringData), &typedDest)
				for k, v := range typedDest {
					dest[k] = v
				}
			}

			assert.Nil(t, err)

			expected := make(map[string]any)
			expected["messi"] = values

			if valueType == "Long" {
				expected["hello"] = int64(3141592)
			}

			assert.Equal(t, expected, dest)
		})

	}

}

func TestCanCallIntoJmxAndGetResultWithCollections(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration tests when running short mode")
	}
	floatValues := []any{rand.Float32(), rand.Float32(), rand.Float32()}
	doubleValues := []any{rand.Float64(), rand.Float64(), rand.Float64()}
	intValues := []any{int32(rand.Int31()), int32(rand.Int31()), int32(rand.Int31())}
	longValues := []any{int64(rand.Int63()), int64(rand.Int63()), int64(rand.Int63())}
	boolValues := []any{true, false}
	stringValues := []any{"hello", "world", "whatsgoinonyo"}

	valueMapping := map[string][]any{
		"Integer": intValues,
		"Float":   floatValues,
		"Double":  doubleValues,
		"Long":    longValues,
		"Boolean": boolValues,
		"String":  stringValues,
	}
	collectionTypes := []string{"List", "Set"}
	primitiveTypes := []string{"Integer", "Long", "Float", "Double", "Boolean", "String"}

	// {
	// 	initialData: &testData{value: "[1, 2, 3]", className: "java.lang.Integer", containerName: "java.util.List", operationName: "putList"},
	// 	readData:    &testData{value: "messi", operationName: "getList"},
	// 	testName:    "IntListTesting",
	// 	expectedVal: "[1,2,3]",
	// },

	for _, collection := range collectionTypes {
		for _, primitiveType := range primitiveTypes {
			t.Run(fmt.Sprintf("TestJmxAndGetResultsFor_%s<%s>", collection, primitiveType), func(t *testing.T) {
				lockCurrentThread(java)
				defer unlockCurrentThread(java)

				values := valueMapping[primitiveType]

				strBytes, err := json.Marshal(values)
				assert.Nil(t, err)

				className := fmt.Sprintf("java.lang.%s", primitiveType)
				containerName := fmt.Sprintf("java.util.%s", collection)

				data := testData{value: string(strBytes), className: className, containerName: containerName, operationName: fmt.Sprintf("put%s", collection)}

				mbean := &mbean.Client{
					JmxURI:            "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi",
					ClassHandlers:     sync.Map{},
					InterfaceHandlers: sync.Map{},
					Env:               java.Env,
				}

				registerHandlers(mbean)

				insertData(java.Env, data, t, mbean)

				data = testData{value: "messi", operationName: fmt.Sprintf("get%s", collection)}

				stringData := readData(java.Env, data, t, mbean)

				dest := make([]any, 0)
				switch className {
				case handlers.FloatClasspath:
					var typedDest []float32
					err = json.Unmarshal([]byte(stringData), &typedDest)
					for _, f := range typedDest {
						dest = append(dest, f)
					}
				case handlers.LongClasspath:
					var typedDest []int64
					err = json.Unmarshal([]byte(stringData), &typedDest)
					for _, f := range typedDest {
						dest = append(dest, f)
					}
				case handlers.IntClasspath:
					var typedDest []int32
					err = json.Unmarshal([]byte(stringData), &typedDest)
					for _, f := range typedDest {
						dest = append(dest, f)
					}
				case handlers.StringClasspath:
					var typedDest []string
					err = json.Unmarshal([]byte(stringData), &typedDest)
					for _, f := range typedDest {
						dest = append(dest, f)
					}
				case handlers.BoolClasspath:
					var typedDest []bool
					err = json.Unmarshal([]byte(stringData), &typedDest)
					for _, f := range typedDest {
						dest = append(dest, f)
					}
				case handlers.DoubleClasspath:
					var typedDest []float64
					err = json.Unmarshal([]byte(stringData), &typedDest)
					for _, f := range typedDest {
						dest = append(dest, f)
					}
				}

				assert.Nil(t, err)

				assert.Equal(t, len(values), len(dest))

				containsCounter := 0

				for _, value := range values {
					for _, item := range dest {
						if item == value {
							containsCounter++
						}
					}
				}

				if len(values) != containsCounter {
					assert.Fail(t, fmt.Sprintf("expected '%s' to be equal to '%s'", string(strBytes), stringData))
				}
			})
		}
	}

}

func TestCanCallIntoJmxAndGetResultWithPrimitiveTypes(t *testing.T) {
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
	}

	lockCurrentThread(java)
	defer unlockCurrentThread(java)

	for _, data := range container {
		t.Run(data.testName, func(t *testing.T) {
			lockCurrentThread(java)
			defer unlockCurrentThread(java)

			bean := &mbean.Client{
				JmxURI:            "service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi",
				ClassHandlers:     sync.Map{},
				InterfaceHandlers: sync.Map{},
				Env:               java.Env,
			}

			registerHandlers(bean)

			insertData(java.Env, *data.initialData, t, bean)
			result := readData(java.Env, *data.readData, t, bean)
			assert.Equal(t, data.expectedVal, result)
		})
	}
}

func readData(env *jnigi.Env, data testData, t *testing.T, bean mbean.BeanExecutor) string {

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
	bean.RegisterInterfaceHandler(handlers.MapClassPath, &handlers.MapHandler{ClassHandlers: &client.ClassHandlers, InterfaceHandlers: &client.InterfaceHandlers})
	bean.RegisterInterfaceHandler(handlers.ListClassPath, &handlers.ListHandler{ClassHandlers: &client.ClassHandlers, InterfaceHandlers: &client.InterfaceHandlers})
	bean.RegisterInterfaceHandler(handlers.SetClassPath, &handlers.SetHandler{ClassHandlers: &client.ClassHandlers, InterfaceHandlers: &client.InterfaceHandlers})
}
