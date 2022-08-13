// Code generated by mockery v2.14.0. DO NOT EDIT.

package mbean

import (
	mock "github.com/stretchr/testify/mock"
	jnigi "tekao.net/jnigi"
)

// MockBeanExecutor is an autogenerated mock type for the BeanExecutor type
type MockBeanExecutor struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockBeanExecutor) Close() {
	_m.Called()
}

// Execute provides a mock function with given fields: operation
func (_m *MockBeanExecutor) Execute(operation Operation) (string, error) {
	ret := _m.Called(operation)

	var r0 string
	if rf, ok := ret.Get(0).(func(Operation) string); ok {
		r0 = rf(operation)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(Operation) error); ok {
		r1 = rf(operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEnv provides a mock function with given fields:
func (_m *MockBeanExecutor) GetEnv() *jnigi.Env {
	ret := _m.Called()

	var r0 *jnigi.Env
	if rf, ok := ret.Get(0).(func() *jnigi.Env); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jnigi.Env)
		}
	}

	return r0
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
