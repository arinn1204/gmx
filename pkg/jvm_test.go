package gmx

import (
	"errors"
	"sync"
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

	client := &client{}

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

	client := &client{}

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

	client := &client{}

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

	client := &client{
		mbeans: sync.Map{},
	}

	id, err := client.Connect("localhost", 9001)

	assert.Nil(t, err)
	mockJava.AssertCalled(t, "CreateMBeanConnection", "service:jmx:rmi:///jndi/rmi://localhost:9001/jmxrmi")

	exec, ok := client.mbeans.Load(*id)

	assert.True(t, ok)
	assert.Equal(t, exec, mockExecutor)
}

func TestConnect_ConnectFails(t *testing.T) {
	mockJava := &jvm.MockIJava{}
	mockJava.On("CreateMBeanConnection", "service:jmx:rmi:///jndi/rmi://localhost:9001/jmxrmi").Return(nil, errors.New("something went wrong"))

	java = mockJava

	client := &client{
		mbeans: sync.Map{},
	}

	id, err := client.Connect("localhost", 9001)

	assert.Nil(t, id)

	assert.Equal(t, errors.New("failed to create a connection::something went wrong"), err)
}

func TestClose_WithConnections(t *testing.T) {
	mockJVM := &jvm.MockIJava{}

	mockJVM.On("ShutdownJvm").Return(nil)

	java = mockJVM

	client := &client{
		mbeans:              sync.Map{},
		numberOfConnections: 1,
	}

	id := uuid.New()
	mockBean := mbean.MockBeanExecutor{}

	mockBean.On("Close").Once()

	client.mbeans.Store(id, &mockBean)

	client.Close()

	exec, ok := client.mbeans.Load(id)
	assert.False(t, ok)
	assert.Nil(t, exec)

	mockBean.AssertCalled(t, "Close")
	mockJVM.AssertCalled(t, "ShutdownJvm")
}
