package gmx

import (
	"fmt"
	"testing"

	"github.com/arinn1204/gmx/internal/mbean"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExecuteAgainstID(t *testing.T) {
	id := uuid.New()

	mbeans := make(map[uuid.UUID]mbean.BeanExecutor)
	executor := &mbean.MockBeanExecutor{}

	for i := range []int{0, 1, 2} {

		operationArgs := make([]mbean.OperationArgs, 0)
		args := make([]MBeanArgs, 0)

		for j := 0; j < i; j++ {
			operationArgs = append(operationArgs, mbean.OperationArgs{
				Value: fmt.Sprintf("value%d", j),
				Type:  "java.lang.String",
			})

			args = append(args, MBeanArgs{
				Value:    fmt.Sprintf("value%d", j),
				JavaType: "java.lang.String",
			})
		}

		expected := mbean.Operation{
			Domain:    "org.example",
			Name:      "game",
			Operation: "getValue",
			Args:      operationArgs,
		}

		executor.On("Execute", expected).Return("hello", nil)

		mbeans[id] = executor

		client := Client{
			mbeans: mbeans,
		}

		val, err := client.ExecuteAgainstID(id, "org.example", "game", "getValue", args...)

		assert.Nil(t, err)
		assert.Equal(t, "hello", val)

		executor.AssertCalled(t, "Execute", expected)

	}
}

func TestExecuteAgainstAll(t *testing.T) {
	expectedOperation := make(map[uuid.UUID]mbean.Operation)
	operation := mbean.Operation{
		Domain:    "com.google",
		Name:      "spyware",
		Operation: "getLocation",
		Args:      make([]mbean.OperationArgs, 0),
	}

	locationId := uuid.New()
	gameId := uuid.New()

	expectedOperation[locationId] = operation

	gameExecutor := mbean.MockBeanExecutor{}
	gameExecutor.On("Execute", operation).Return("NV", nil)

	locationExecutor := mbean.MockBeanExecutor{}
	locationExecutor.On("Execute", operation).Return("CA", nil)

	mbeans := make(map[uuid.UUID]mbean.BeanExecutor)
	mbeans[locationId] = &locationExecutor
	mbeans[gameId] = &gameExecutor

	client := Client{
		mbeans: mbeans,
	}

	res, err := client.ExecuteAgainstAll("com.google", "spyware", "getLocation")

	assert.Nil(t, err[gameId])
	assert.Nil(t, err[locationId])

	assert.Equal(t, "NV", res[gameId])
	assert.Equal(t, "CA", res[locationId])

}
