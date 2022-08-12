package gmx

import (
	"testing"

	"github.com/arinn1204/gmx/internal/jvm"
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
