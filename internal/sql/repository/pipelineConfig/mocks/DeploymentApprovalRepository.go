// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	mock "github.com/stretchr/testify/mock"
)

// DeploymentApprovalRepository is an autogenerated mock type for the DeploymentApprovalRepository type
type DeploymentApprovalRepository struct {
	mock.Mock
}

// ConsumeApprovalRequest provides a mock function with given fields: requestId
func (_m *DeploymentApprovalRepository) ConsumeApprovalRequest(requestId int) error {
	ret := _m.Called(requestId)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(requestId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchApprovalDataForArtifacts provides a mock function with given fields: artifactIds, pipelineId
func (_m *DeploymentApprovalRepository) FetchApprovalDataForArtifacts(artifactIds []int, pipelineId int) ([]*pipelineConfig.DeploymentApprovalRequest, error) {
	ret := _m.Called(artifactIds, pipelineId)

	var r0 []*pipelineConfig.DeploymentApprovalRequest
	var r1 error
	if rf, ok := ret.Get(0).(func([]int, int) ([]*pipelineConfig.DeploymentApprovalRequest, error)); ok {
		return rf(artifactIds, pipelineId)
	}
	if rf, ok := ret.Get(0).(func([]int, int) []*pipelineConfig.DeploymentApprovalRequest); ok {
		r0 = rf(artifactIds, pipelineId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.DeploymentApprovalRequest)
		}
	}

	if rf, ok := ret.Get(1).(func([]int, int) error); ok {
		r1 = rf(artifactIds, pipelineId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchApprovalDataForRequests provides a mock function with given fields: requestIds
func (_m *DeploymentApprovalRepository) FetchApprovalDataForRequests(requestIds []int) ([]*pipelineConfig.DeploymentApprovalUserData, error) {
	ret := _m.Called(requestIds)

	var r0 []*pipelineConfig.DeploymentApprovalUserData
	var r1 error
	if rf, ok := ret.Get(0).(func([]int) ([]*pipelineConfig.DeploymentApprovalUserData, error)); ok {
		return rf(requestIds)
	}
	if rf, ok := ret.Get(0).(func([]int) []*pipelineConfig.DeploymentApprovalUserData); ok {
		r0 = rf(requestIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.DeploymentApprovalUserData)
		}
	}

	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(requestIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchApprovedDataByApprovalId provides a mock function with given fields: approvalRequestId
func (_m *DeploymentApprovalRepository) FetchApprovedDataByApprovalId(approvalRequestId int) ([]*pipelineConfig.DeploymentApprovalUserData, error) {
	ret := _m.Called(approvalRequestId)

	var r0 []*pipelineConfig.DeploymentApprovalUserData
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]*pipelineConfig.DeploymentApprovalUserData, error)); ok {
		return rf(approvalRequestId)
	}
	if rf, ok := ret.Get(0).(func(int) []*pipelineConfig.DeploymentApprovalUserData); ok {
		r0 = rf(approvalRequestId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.DeploymentApprovalUserData)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(approvalRequestId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchById provides a mock function with given fields: requestId
func (_m *DeploymentApprovalRepository) FetchById(requestId int) (*pipelineConfig.DeploymentApprovalRequest, error) {
	ret := _m.Called(requestId)

	var r0 *pipelineConfig.DeploymentApprovalRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*pipelineConfig.DeploymentApprovalRequest, error)); ok {
		return rf(requestId)
	}
	if rf, ok := ret.Get(0).(func(int) *pipelineConfig.DeploymentApprovalRequest); ok {
		r0 = rf(requestId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pipelineConfig.DeploymentApprovalRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(requestId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchWithPipelineAndArtifactDetails provides a mock function with given fields: requestId
func (_m *DeploymentApprovalRepository) FetchWithPipelineAndArtifactDetails(requestId int) (*pipelineConfig.DeploymentApprovalRequest, error) {
	ret := _m.Called(requestId)

	var r0 *pipelineConfig.DeploymentApprovalRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*pipelineConfig.DeploymentApprovalRequest, error)); ok {
		return rf(requestId)
	}
	if rf, ok := ret.Get(0).(func(int) *pipelineConfig.DeploymentApprovalRequest); ok {
		r0 = rf(requestId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pipelineConfig.DeploymentApprovalRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(requestId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: deploymentApprovalRequest
func (_m *DeploymentApprovalRepository) Save(deploymentApprovalRequest *pipelineConfig.DeploymentApprovalRequest) error {
	ret := _m.Called(deploymentApprovalRequest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pipelineConfig.DeploymentApprovalRequest) error); ok {
		r0 = rf(deploymentApprovalRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveDeploymentUserData provides a mock function with given fields: userData
func (_m *DeploymentApprovalRepository) SaveDeploymentUserData(userData *pipelineConfig.DeploymentApprovalUserData) error {
	ret := _m.Called(userData)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pipelineConfig.DeploymentApprovalUserData) error); ok {
		r0 = rf(userData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: deploymentApprovalRequest
func (_m *DeploymentApprovalRepository) Update(deploymentApprovalRequest *pipelineConfig.DeploymentApprovalRequest) error {
	ret := _m.Called(deploymentApprovalRequest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pipelineConfig.DeploymentApprovalRequest) error); ok {
		r0 = rf(deploymentApprovalRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewDeploymentApprovalRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewDeploymentApprovalRepository creates a new instance of DeploymentApprovalRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDeploymentApprovalRepository(t mockConstructorTestingTNewDeploymentApprovalRepository) *DeploymentApprovalRepository {
	mock := &DeploymentApprovalRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
