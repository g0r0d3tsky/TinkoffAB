// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	domain "homework/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// DeviceUseCase is an autogenerated mock type for the DeviceUseCase type
type DeviceUseCase struct {
	mock.Mock
}

// CreateDevice provides a mock function with given fields: d
func (_m *DeviceUseCase) CreateDevice(d domain.Device) error {
	ret := _m.Called(d)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.Device) error); ok {
		r0 = rf(d)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDevice provides a mock function with given fields: _a0
func (_m *DeviceUseCase) DeleteDevice(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDevice provides a mock function with given fields: _a0
func (_m *DeviceUseCase) GetDevice(_a0 string) (domain.Device, error) {
	ret := _m.Called(_a0)

	var r0 domain.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.Device, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) domain.Device); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(domain.Device)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateDevice provides a mock function with given fields: _a0
func (_m *DeviceUseCase) UpdateDevice(_a0 domain.Device) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.Device) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDeviceUseCase creates a new instance of DeviceUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeviceUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeviceUseCase {
	mock := &DeviceUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
