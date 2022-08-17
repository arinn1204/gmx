// Code generated by mockery v2.14.0. DO NOT EDIT.

package gmx

import (
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// MockMBeanOperator is an autogenerated mock type for the MBeanOperator type
type MockMBeanOperator struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockMBeanOperator) Close() {
	_m.Called()
}

// Connect provides a mock function with given fields: hostname, port
func (_m *MockMBeanOperator) Connect(hostname string, port int) (*uuid.UUID, error) {
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

// ExecuteAgainstAll provides a mock function with given fields: domain, name, operation, args
func (_m *MockMBeanOperator) ExecuteAgainstAll(domain string, name string, operation string, args ...MBeanArgs) (map[uuid.UUID]string, map[uuid.UUID]error) {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, domain, name, operation)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 map[uuid.UUID]string
	if rf, ok := ret.Get(0).(func(string, string, string, ...MBeanArgs) map[uuid.UUID]string); ok {
		r0 = rf(domain, name, operation, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uuid.UUID]string)
		}
	}

	var r1 map[uuid.UUID]error
	if rf, ok := ret.Get(1).(func(string, string, string, ...MBeanArgs) map[uuid.UUID]error); ok {
		r1 = rf(domain, name, operation, args...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[uuid.UUID]error)
		}
	}

	return r0, r1
}

// ExecuteAgainstID provides a mock function with given fields: id, domain, name, operation, args
func (_m *MockMBeanOperator) ExecuteAgainstID(id uuid.UUID, domain string, name string, operation string, args ...MBeanArgs) (string, error) {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, id, domain, name, operation)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	if rf, ok := ret.Get(0).(func(uuid.UUID, string, string, string, ...MBeanArgs) string); ok {
		r0 = rf(id, domain, name, operation, args...)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string, string, string, ...MBeanArgs) error); ok {
		r1 = rf(id, domain, name, operation, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: domain, beanName, attributeName
func (_m *MockMBeanOperator) Get(domain string, beanName string, attributeName string) (string, error) {
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

// Initialize provides a mock function with given fields:
func (_m *MockMBeanOperator) Initialize() (*Client, error) {
	ret := _m.Called()

	var r0 *Client
	if rf, ok := ret.Get(0).(func() *Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Client)
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

// Put provides a mock function with given fields: domain, beanName, attributeName, value
func (_m *MockMBeanOperator) Put(domain string, beanName string, attributeName string, value interface{}) (string, error) {
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

type mockConstructorTestingTNewMockMBeanOperator interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockMBeanOperator creates a new instance of MockMBeanOperator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockMBeanOperator(t mockConstructorTestingTNewMockMBeanOperator) *MockMBeanOperator {
	mock := &MockMBeanOperator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
