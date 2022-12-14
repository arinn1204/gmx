// Code generated by mockery v2.14.0. DO NOT EDIT.

package gmx

import (
	extensions "github.com/arinn1204/gmx/pkg/extensions"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockMBeanClient is an autogenerated mock type for the MBeanClient type
type MockMBeanClient struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockMBeanClient) Close() {
	_m.Called()
}

// GetAttributeManager provides a mock function with given fields:
func (_m *MockMBeanClient) GetAttributeManager() MBeanAttributeManager {
	ret := _m.Called()

	var r0 MBeanAttributeManager
	if rf, ok := ret.Get(0).(func() MBeanAttributeManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(MBeanAttributeManager)
		}
	}

	return r0
}

// GetOperator provides a mock function with given fields:
func (_m *MockMBeanClient) GetOperator() MBeanOperator {
	ret := _m.Called()

	var r0 MBeanOperator
	if rf, ok := ret.Get(0).(func() MBeanOperator); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(MBeanOperator)
		}
	}

	return r0
}

// Initialize provides a mock function with given fields:
func (_m *MockMBeanClient) Initialize() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterClassHandler provides a mock function with given fields: typeName, handler
func (_m *MockMBeanClient) RegisterClassHandler(typeName string, handler extensions.IHandler) {
	_m.Called(typeName, handler)
}

// RegisterConnection provides a mock function with given fields: hostname, port
func (_m *MockMBeanClient) RegisterConnection(hostname string, port int) (*uuid.UUID, error) {
	ret := _m.Called(hostname, port)

	var r0 *uuid.UUID
	if rf, ok := ret.Get(0).(func(string, int) *uuid.UUID); ok {
		r0 = rf(hostname, port)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(hostname, port)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterInterfaceHandler provides a mock function with given fields: typeName, handler
func (_m *MockMBeanClient) RegisterInterfaceHandler(typeName string, handler extensions.InterfaceHandler) {
	_m.Called(typeName, handler)
}

type mockConstructorTestingTNewMockMBeanClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockMBeanClient creates a new instance of MockMBeanClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockMBeanClient(t mockConstructorTestingTNewMockMBeanClient) *MockMBeanClient {
	mock := &MockMBeanClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
