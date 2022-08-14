// Code generated by mockery v2.14.0. DO NOT EDIT.

package extensions

import (
	mock "github.com/stretchr/testify/mock"
	jnigi "tekao.net/jnigi"
)

// MockIClassHandler is an autogenerated mock type for the IClassHandler type
type MockIClassHandler[T interface{}] struct {
	mock.Mock
}

// toGoRepresentation provides a mock function with given fields: env, object, dest
func (_m *MockIClassHandler[T]) toGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest *T) error {
	ret := _m.Called(env, object, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*jnigi.Env, *jnigi.ObjectRef, *T) error); ok {
		r0 = rf(env, object, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// toJniRepresentation provides a mock function with given fields: env, value
func (_m *MockIClassHandler[T]) toJniRepresentation(env *jnigi.Env, value T) (*jnigi.ObjectRef, error) {
	ret := _m.Called(env, value)

	var r0 *jnigi.ObjectRef
	if rf, ok := ret.Get(0).(func(*jnigi.Env, T) *jnigi.ObjectRef); ok {
		r0 = rf(env, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jnigi.ObjectRef)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*jnigi.Env, T) error); ok {
		r1 = rf(env, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockIClassHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockIClassHandler creates a new instance of MockIClassHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockIClassHandler[T interface{}](t mockConstructorTestingTNewMockIClassHandler) *MockIClassHandler[T] {
	mock := &MockIClassHandler[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
