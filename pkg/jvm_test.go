package gmx

import (
	"errors"
	"testing"

	"github.com/arinn1204/gmx/internal/jvm"
	"github.com/arinn1204/gmx/internal/mbean"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInitialize_CalledFirstTime(t *testing.T) {
	mockJava := &jvm.MockIJava{}

	mockJava.On("IsStarted").Return(false)
	mockJava.On("CreateJVM").Return(mockJava, nil)

	java = mockJava

	client := &Client{}

	client.Initialize()

	mockJava.AssertNumberOfCalls(t, "IsStarted", 2)
	mockJava.AssertNumberOfCalls(t, "CreateJVM", 1)
}

func TestInitialize_SimulatedRaceCondition(t *testing.T) {
	mockJava := &jvm.MockIJava{}

	mockJava.On("CreateJVM").Return(mockJava, nil)

	mockJava.On("IsStarted").Return(false).Once()
	mockJava.On("IsStarted").Return(true).Once()

	java = mockJava

	client := &Client{}

	client.Initialize()

	mockJava.AssertNumberOfCalls(t, "IsStarted", 2)
	mockJava.AssertNotCalled(t, "CreateJVM")
}

func TestInitialize_OnlyEverCalledOnce(t *testing.T) {
	mockJava := &jvm.MockIJava{}

	mockJava.On("CreateJVM").Return(mockJava, nil)

	mockJava.On("IsStarted").Return(false).Once()
	mockJava.On("IsStarted").Return(false).Once()
	mockJava.On("IsStarted").Return(true).Once()
	java = mockJava

	client := &Client{}

	//can't really do true parallel testing with a mock jvm
	client.Initialize()
	client.Initialize()

	mockJava.AssertNumberOfCalls(t, "IsStarted", 3)
	mockJava.AssertNumberOfCalls(t, "CreateJVM", 1)
}

func TestConnect_HappyPath(t *testing.T) {
	mockJava := &jvm.MockIJava{}
	mockExecutor := &mbean.MockBeanExecutor{}
	mockJava.On("CreateMBeanConnection", "service:jmx:rmi:///jndi/rmi://localhost:9001/jmxrmi").Return(mockExecutor, nil)

	java = mockJava

	client := &Client{
		mbeans: make(map[uuid.UUID]mbean.BeanExecutor),
	}

	id, err := client.Connect("localhost", 9001)

	assert.Nil(t, err)
	mockJava.AssertCalled(t, "CreateMBeanConnection", "service:jmx:rmi:///jndi/rmi://localhost:9001/jmxrmi")

	assert.Equal(t, client.mbeans[*id], mockExecutor)
}

func TestConnect_ConnectFails(t *testing.T) {
	mockJava := &jvm.MockIJava{}
	mockJava.On("CreateMBeanConnection", "service:jmx:rmi:///jndi/rmi://localhost:9001/jmxrmi").Return(nil, errors.New("something went wrong"))

	java = mockJava

	client := &Client{
		mbeans: make(map[uuid.UUID]mbean.BeanExecutor),
	}

	id, err := client.Connect("localhost", 9001)

	assert.Nil(t, id)

	assert.Equal(t, errors.New("failed to create a connection::something went wrong"), err)
}
