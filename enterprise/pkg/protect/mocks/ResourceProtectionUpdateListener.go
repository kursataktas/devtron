// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	protect "github.com/devtron-labs/devtron/enterprise/pkg/protect"
	mock "github.com/stretchr/testify/mock"
)

// ResourceProtectionUpdateListener is an autogenerated mock type for the ResourceProtectionUpdateListener type
type ResourceProtectionUpdateListener struct {
	mock.Mock
}

// OnStateChange provides a mock function with given fields: appId, envId, state, userId
func (_m *ResourceProtectionUpdateListener) OnStateChange(appId int, envId int, state protect.ProtectionState, userId int32) {
	_m.Called(appId, envId, state, userId)
}

type mockConstructorTestingTNewResourceProtectionUpdateListener interface {
	mock.TestingT
	Cleanup(func())
}

// NewResourceProtectionUpdateListener creates a new instance of ResourceProtectionUpdateListener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewResourceProtectionUpdateListener(t mockConstructorTestingTNewResourceProtectionUpdateListener) *ResourceProtectionUpdateListener {
	mock := &ResourceProtectionUpdateListener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
