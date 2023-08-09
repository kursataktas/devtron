// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	pg "github.com/go-pg/pg"
	mock "github.com/stretchr/testify/mock"

	pipelineConfig "github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
)

// CiPipelineMaterialRepository is an autogenerated mock type for the CiPipelineMaterialRepository type
type CiPipelineMaterialRepository struct {
	mock.Mock
}

// CheckRegexExistsForMaterial provides a mock function with given fields: id
func (_m *CiPipelineMaterialRepository) CheckRegexExistsForMaterial(id int) bool {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindByCiPipelineIdsIn provides a mock function with given fields: ids
func (_m *CiPipelineMaterialRepository) FindByCiPipelineIdsIn(ids []int) ([]*pipelineConfig.CiPipelineMaterial, error) {
	ret := _m.Called(ids)

	var r0 []*pipelineConfig.CiPipelineMaterial
	var r1 error
	if rf, ok := ret.Get(0).(func([]int) ([]*pipelineConfig.CiPipelineMaterial, error)); ok {
		return rf(ids)
	}
	if rf, ok := ret.Get(0).(func([]int) []*pipelineConfig.CiPipelineMaterial); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.CiPipelineMaterial)
		}
	}

	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllExceptUnsetRegexBranch provides a mock function with given fields:
func (_m *CiPipelineMaterialRepository) GetAllExceptUnsetRegexBranch() ([]*pipelineConfig.CiPipelineMaterial, error) {
	ret := _m.Called()

	var r0 []*pipelineConfig.CiPipelineMaterial
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*pipelineConfig.CiPipelineMaterial, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*pipelineConfig.CiPipelineMaterial); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.CiPipelineMaterial)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByCiPipelineIdsExceptUnsetRegexBranch provides a mock function with given fields: ids
func (_m *CiPipelineMaterialRepository) GetByCiPipelineIdsExceptUnsetRegexBranch(ids []int) ([]*pipelineConfig.CiPipelineMaterial, error) {
	ret := _m.Called(ids)

	var r0 []*pipelineConfig.CiPipelineMaterial
	var r1 error
	if rf, ok := ret.Get(0).(func([]int) ([]*pipelineConfig.CiPipelineMaterial, error)); ok {
		return rf(ids)
	}
	if rf, ok := ret.Get(0).(func([]int) []*pipelineConfig.CiPipelineMaterial); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.CiPipelineMaterial)
		}
	}

	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id
func (_m *CiPipelineMaterialRepository) GetById(id int) (*pipelineConfig.CiPipelineMaterial, error) {
	ret := _m.Called(id)

	var r0 *pipelineConfig.CiPipelineMaterial
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*pipelineConfig.CiPipelineMaterial, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *pipelineConfig.CiPipelineMaterial); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pipelineConfig.CiPipelineMaterial)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByPipelineId provides a mock function with given fields: id
func (_m *CiPipelineMaterialRepository) GetByPipelineId(id int) ([]*pipelineConfig.CiPipelineMaterial, error) {
	ret := _m.Called(id)

	var r0 []*pipelineConfig.CiPipelineMaterial
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]*pipelineConfig.CiPipelineMaterial, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) []*pipelineConfig.CiPipelineMaterial); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.CiPipelineMaterial)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByPipelineIdAndGitMaterialId provides a mock function with given fields: id, gitMaterialId
func (_m *CiPipelineMaterialRepository) GetByPipelineIdAndGitMaterialId(id int, gitMaterialId int) ([]*pipelineConfig.CiPipelineMaterial, error) {
	ret := _m.Called(id, gitMaterialId)

	var r0 []*pipelineConfig.CiPipelineMaterial
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) ([]*pipelineConfig.CiPipelineMaterial, error)); ok {
		return rf(id, gitMaterialId)
	}
	if rf, ok := ret.Get(0).(func(int, int) []*pipelineConfig.CiPipelineMaterial); ok {
		r0 = rf(id, gitMaterialId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.CiPipelineMaterial)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(id, gitMaterialId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByPipelineIdForRegexAndFixed provides a mock function with given fields: id
func (_m *CiPipelineMaterialRepository) GetByPipelineIdForRegexAndFixed(id int) ([]*pipelineConfig.CiPipelineMaterial, error) {
	ret := _m.Called(id)

	var r0 []*pipelineConfig.CiPipelineMaterial
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]*pipelineConfig.CiPipelineMaterial, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) []*pipelineConfig.CiPipelineMaterial); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.CiPipelineMaterial)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCheckoutPath provides a mock function with given fields: gitMaterialId
func (_m *CiPipelineMaterialRepository) GetCheckoutPath(gitMaterialId int) (string, error) {
	ret := _m.Called(gitMaterialId)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (string, error)); ok {
		return rf(gitMaterialId)
	}
	if rf, ok := ret.Get(0).(func(int) string); ok {
		r0 = rf(gitMaterialId)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(gitMaterialId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegexByPipelineId provides a mock function with given fields: id
func (_m *CiPipelineMaterialRepository) GetRegexByPipelineId(id int) ([]*pipelineConfig.CiPipelineMaterial, error) {
	ret := _m.Called(id)

	var r0 []*pipelineConfig.CiPipelineMaterial
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]*pipelineConfig.CiPipelineMaterial, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) []*pipelineConfig.CiPipelineMaterial); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pipelineConfig.CiPipelineMaterial)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: tx, pipeline
func (_m *CiPipelineMaterialRepository) Save(tx *pg.Tx, pipeline ...*pipelineConfig.CiPipelineMaterial) error {
	_va := make([]interface{}, len(pipeline))
	for _i := range pipeline {
		_va[_i] = pipeline[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, tx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pg.Tx, ...*pipelineConfig.CiPipelineMaterial) error); ok {
		r0 = rf(tx, pipeline...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: tx, material
func (_m *CiPipelineMaterialRepository) Update(tx *pg.Tx, material ...*pipelineConfig.CiPipelineMaterial) error {
	_va := make([]interface{}, len(material))
	for _i := range material {
		_va[_i] = material[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, tx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pg.Tx, ...*pipelineConfig.CiPipelineMaterial) error); ok {
		r0 = rf(tx, material...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateNotNull provides a mock function with given fields: tx, material
func (_m *CiPipelineMaterialRepository) UpdateNotNull(tx *pg.Tx, material ...*pipelineConfig.CiPipelineMaterial) error {
	_va := make([]interface{}, len(material))
	for _i := range material {
		_va[_i] = material[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, tx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pg.Tx, ...*pipelineConfig.CiPipelineMaterial) error); ok {
		r0 = rf(tx, material...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCiPipelineMaterialRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewCiPipelineMaterialRepository creates a new instance of CiPipelineMaterialRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCiPipelineMaterialRepository(t mockConstructorTestingTNewCiPipelineMaterialRepository) *CiPipelineMaterialRepository {
	mock := &CiPipelineMaterialRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
