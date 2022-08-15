// Code generated by mockery v2.14.0. DO NOT EDIT.

package extensions

import (
	mock "github.com/stretchr/testify/mock"
	jnigi "tekao.net/jnigi"
)

// MockIHandler is an autogenerated mock type for the IHandler type
type MockIHandler struct {
	mock.Mock
}

// ToGoRepresentation provides a mock function with given fields: env, object, dest
func (_m *MockIHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest interface{}) error {
	ret := _m.Called(env, object, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*jnigi.Env, *jnigi.ObjectRef, interface{}) error); ok {
		r0 = rf(env, object, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ToJniRepresentation provides a mock function with given fields: env, value
func (_m *MockIHandler) ToJniRepresentation(env *jnigi.Env, value interface{}) (*jnigi.ObjectRef, error) {
	ret := _m.Called(env, value)

	var r0 *jnigi.ObjectRef
	if rf, ok := ret.Get(0).(func(*jnigi.Env, interface{}) *jnigi.ObjectRef); ok {
		r0 = rf(env, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jnigi.ObjectRef)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*jnigi.Env, interface{}) error); ok {
		r1 = rf(env, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockIHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockIHandler creates a new instance of MockIHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockIHandler(t mockConstructorTestingTNewMockIHandler) *MockIHandler {
	mock := &MockIHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}