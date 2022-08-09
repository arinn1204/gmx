// Code generated by mockery v2.14.0. DO NOT EDIT.

package mbean

import mock "github.com/stretchr/testify/mock"

// MockBeanExecutor is an autogenerated mock type for the BeanExecutor type
type MockBeanExecutor struct {
	mock.Mock
}

// Execute provides a mock function with given fields: operation
func (_m *MockBeanExecutor) Execute(operation Operation) (interface{}, error) {
	ret := _m.Called(operation)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(Operation) interface{}); ok {
		r0 = rf(operation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(Operation) error); ok {
		r1 = rf(operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockBeanExecutor interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockBeanExecutor creates a new instance of MockBeanExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockBeanExecutor(t mockConstructorTestingTNewMockBeanExecutor) *MockBeanExecutor {
	mock := &MockBeanExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
