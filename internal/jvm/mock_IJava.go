// Code generated by mockery v2.14.0. DO NOT EDIT.

package jvm

import (
	mock "github.com/stretchr/testify/mock"
	jnigi "tekao.net/jnigi"
)

// MockIJava is an autogenerated mock type for the IJava type
type MockIJava struct {
	mock.Mock
}

// Attach provides a mock function with given fields:
func (_m *MockIJava) Attach() *jnigi.Env {
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

// CreateJVM provides a mock function with given fields:
func (_m *MockIJava) CreateJVM() (IJava, error) {
	ret := _m.Called()

	var r0 IJava
	if rf, ok := ret.Get(0).(func() IJava); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(IJava)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Detach provides a mock function with given fields:
func (_m *MockIJava) Detach() {
	_m.Called()
}

// IsStarted provides a mock function with given fields:
func (_m *MockIJava) IsStarted() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ShutdownJvm provides a mock function with given fields:
func (_m *MockIJava) ShutdownJvm() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockIJava interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockIJava creates a new instance of MockIJava. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockIJava(t mockConstructorTestingTNewMockIJava) *MockIJava {
	mock := &MockIJava{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
