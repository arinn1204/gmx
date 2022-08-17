// Code generated by mockery v2.14.0. DO NOT EDIT.

package gmx

import (
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// MockAttributeHandler is an autogenerated mock type for the AttributeHandler type
type MockAttributeHandler struct {
	mock.Mock
}

// Get provides a mock function with given fields: domain, beanName, attributeName
func (_m *MockAttributeHandler) Get(domain string, beanName string, attributeName string) (string, error) {
	ret := _m.Called(domain, beanName, attributeName)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string) string); ok {
		r0 = rf(domain, beanName, attributeName)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(domain, beanName, attributeName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id, domain, beanName, attributeName
func (_m *MockAttributeHandler) GetById(id uuid.UUID, domain string, beanName string, attributeName string) (string, error) {
	ret := _m.Called(id, domain, beanName, attributeName)

	var r0 string
	if rf, ok := ret.Get(0).(func(uuid.UUID, string, string, string) string); ok {
		r0 = rf(id, domain, beanName, attributeName)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string, string, string) error); ok {
		r1 = rf(id, domain, beanName, attributeName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Put provides a mock function with given fields: domain, beanName, attributeName, value
func (_m *MockAttributeHandler) Put(domain string, beanName string, attributeName string, value interface{}) (string, error) {
	ret := _m.Called(domain, beanName, attributeName, value)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, interface{}) string); ok {
		r0 = rf(domain, beanName, attributeName, value)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, interface{}) error); ok {
		r1 = rf(domain, beanName, attributeName, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutById provides a mock function with given fields: id, domain, beanName, attributeName, value
func (_m *MockAttributeHandler) PutById(id uuid.UUID, domain string, beanName string, attributeName string, value interface{}) (string, error) {
	ret := _m.Called(id, domain, beanName, attributeName, value)

	var r0 string
	if rf, ok := ret.Get(0).(func(uuid.UUID, string, string, string, interface{}) string); ok {
		r0 = rf(id, domain, beanName, attributeName, value)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string, string, string, interface{}) error); ok {
		r1 = rf(id, domain, beanName, attributeName, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockAttributeHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAttributeHandler creates a new instance of MockAttributeHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAttributeHandler(t mockConstructorTestingTNewMockAttributeHandler) *MockAttributeHandler {
	mock := &MockAttributeHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
