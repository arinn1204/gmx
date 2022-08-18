package gmx

import (
	"fmt"
	"testing"

	"github.com/arinn1204/gmx/internal/jvm"
	"github.com/arinn1204/gmx/internal/mbean"
	"tekao.net/jnigi"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExecuteAgainstID(t *testing.T) {
	id := uuid.New()

	mbeans := make(map[uuid.UUID]mbean.BeanExecutor)
	executor := &mbean.MockBeanExecutor{}

	mockJava := jvm.MockIJava{}

	env := &jnigi.Env{}

	mockJava.On("Attach").Return(env)
	mockJava.On("Detach").Return(nil)

	java = &mockJava

	for i := range []int{0, 1, 2} {

		operationArgs := make([]mbean.OperationArgs, 0)
		args := make([]MBeanArgs, 0)

		for j := 0; j < i; j++ {
			operationArgs = append(operationArgs, mbean.OperationArgs{
				Value:    fmt.Sprintf("value%d", j),
				JavaType: "java.lang.String",
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
		executor.On("WithEnvironment", env).Return(executor)

		mbeans[id] = executor

		client := client{
			mbeans: mbeans,
		}

		operator := client.GetOperator()

		val, err := operator.ExecuteAgainstID(id, "org.example", "game", "getValue", args...)

		assert.Nil(t, err)
		assert.Equal(t, "hello", val)

		executor.AssertCalled(t, "Execute", expected)

	}
}

func TestExecuteAgainstAll(t *testing.T) {
	mockJava := jvm.MockIJava{}

	env := &jnigi.Env{}

	mockJava.On("Attach").Return(env)
	mockJava.On("Detach").Return(nil)

	java = &mockJava
	expectedOperation := make(map[uuid.UUID]mbean.Operation)
	operation := mbean.Operation{
		Domain:    "com.google",
		Name:      "spyware",
		Operation: "getLocation",
		Args:      make([]mbean.OperationArgs, 0),
	}

	locationID := uuid.New()
	gameID := uuid.New()

	expectedOperation[locationID] = operation

	gameExecutor := mbean.MockBeanExecutor{}
	gameExecutor.On("WithEnvironment", env).Return(&gameExecutor)
	gameExecutor.On("Execute", operation).Return("NV", nil)

	locationExecutor := mbean.MockBeanExecutor{}
	locationExecutor.On("WithEnvironment", env).Return(&locationExecutor)
	locationExecutor.On("Execute", operation).Return("CA", nil)

	mbeans := make(map[uuid.UUID]mbean.BeanExecutor)
	mbeans[locationID] = &locationExecutor
	mbeans[gameID] = &gameExecutor

	client := client{
		mbeans: mbeans,
	}
	operator := client.GetOperator()

	res, err := operator.ExecuteAgainstAll("com.google", "spyware", "getLocation")

	assert.Nil(t, err[gameID])
	assert.Nil(t, err[locationID])

	assert.Equal(t, "NV", res[gameID])
	assert.Equal(t, "CA", res[locationID])

}
