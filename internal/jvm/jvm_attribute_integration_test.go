package jvm

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/arinn1204/gmx/internal/mbean"
	"github.com/stretchr/testify/assert"
	"tekao.net/jnigi"
)

func TestCanReadAndSetPrimitiveAttributes(t *testing.T) {
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

	for _, primitive := range primitiveTypes {
		value := valueMapping[primitive]

		t.Run(fmt.Sprintf("TestCanReadAndSetPrimitiveAttributes_%s", primitive), func(t *testing.T) {
			lockCurrentThread(java)
			defer unlockCurrentThread(java)

			str := toString(value, t)

			className := fmt.Sprintf("java.lang.%s", primitive)

			data := testData{
				value:         str,
				className:     className,
				operationName: fmt.Sprintf("%sAttribute", primitive),
			}

			mbean, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
			assert.Nil(t, err)
			registerHandlers(mbean)
			defer mbean.Close()

			updateAttribute(java.Env, data, t, mbean)

			value, _ := readAttribute(java.Env, data.operationName, t, mbean)

			assert.Equal(t, str, value)
		})
	}
}

func TestCanReadAndWriteCollectionAttributes(t *testing.T) {
	collections := []string{"List", "Set"}
	floatValues := []float32{float32(rand.Float32()), float32(rand.Float32()), float32(rand.Float32())}
	stringValues := []string{"hello", "world", "whatsgoinonyo"}

	valueMapping := map[string]any{
		"List": stringValues,
		"Set":  floatValues,
	}

	typeMapping := map[string]any{
		"List": "String",
		"Set":  "Float",
	}

	for _, collection := range collections {
		value := valueMapping[collection]
		t.Run(fmt.Sprintf("TestCanReadAndWriteCollectionAttributes_%s", collection), func(t *testing.T) {
			lockCurrentThread(java)
			defer unlockCurrentThread(java)

			str := toString(value, t)

			javaType := fmt.Sprintf("java.lang.%s", typeMapping[collection])
			containerType := fmt.Sprintf("java.util.%s", collection)

			data := testData{
				value:         str,
				className:     javaType,
				containerName: containerType,
				operationName: fmt.Sprintf("%sAttribute", collection),
			}

			mbean, err := java.CreateMBeanConnection("service:jmx:rmi:///jndi/rmi://127.0.0.1:9001/jmxrmi")
			assert.Nil(t, err)
			registerHandlers(mbean)
			defer mbean.Close()

			updateAttribute(java.Env, data, t, mbean)

			strRes, _ := readAttribute(java.Env, data.operationName, t, mbean)

			if collection == "List" {
				var val []string
				json.Unmarshal([]byte(strRes), &val)
				arrayEqual(t, value, val)
			} else {
				var val []float32
				json.Unmarshal([]byte(strRes), &val)
				arrayEqual(t, value, val)
			}

		})
	}

}

func updateAttribute(env *jnigi.Env, data testData, t *testing.T, bean mbean.BeanExecutor) {
	_, err := bean.Put("org.example", "game", data.operationName, mbean.OperationArgs{
		Value:             data.value,
		JavaType:          data.className,
		JavaContainerType: data.containerName,
	})

	assert.Nil(t, err)
}

func readAttribute(env *jnigi.Env, attributeName string, t *testing.T, bean mbean.BeanExecutor) (string, error) {
	val, err := bean.Get("org.example", "game", attributeName, mbean.OperationArgs{})
	assert.Nil(t, err)

	return val, nil
}
