// Code generated by mockery v2.14.0. DO NOT EDIT.

package extensions

import (
	mock "github.com/stretchr/testify/mock"
	jnigi "tekao.net/jnigi"
)

// MockInterfaceHandler is an autogenerated mock type for the InterfaceHandler type
type MockInterfaceHandler struct {
	mock.Mock
}

// ToGoRepresentation provides a mock function with given fields: env, object, dest
func (_m *MockInterfaceHandler) ToGoRepresentation(env *jnigi.Env, object *jnigi.ObjectRef, dest interface{}) error {
	ret := _m.Called(env, object, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*jnigi.Env, *jnigi.ObjectRef, interface{}) error); ok {
		r0 = rf(env, object, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ToJniRepresentation provides a mock function with given fields: env, elementType, value
func (_m *MockInterfaceHandler) ToJniRepresentation(env *jnigi.Env, elementType string, value interface{}) (*jnigi.ObjectRef, error) {
	ret := _m.Called(env, elementType, value)

	var r0 *jnigi.ObjectRef
	if rf, ok := ret.Get(0).(func(*jnigi.Env, string, interface{}) *jnigi.ObjectRef); ok {
		r0 = rf(env, elementType, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jnigi.ObjectRef)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*jnigi.Env, string, interface{}) error); ok {
		r1 = rf(env, elementType, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockInterfaceHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockInterfaceHandler creates a new instance of MockInterfaceHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockInterfaceHandler(t mockConstructorTestingTNewMockInterfaceHandler) *MockInterfaceHandler {
	mock := &MockInterfaceHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
