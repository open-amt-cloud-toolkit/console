// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecase/ciraconfigs/interfaces.go
//
// Generated by this command:
//
//	mockgen -source ./internal/usecase/ciraconfigs/interfaces.go -package ciraconfigs_test
//

// Package ciraconfigs_test is a generated GoMock package.
package ciraconfigs_test

import (
	context "context"
	reflect "reflect"

	entity "github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockRepository) Delete(ctx context.Context, profileName, tenantID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, profileName, tenantID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, profileName, tenantID)
}

// Get mocks base method.
func (m *MockRepository) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.CIRAConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, top, skip, tenantID)
	ret0, _ := ret[0].([]entity.CIRAConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(ctx, top, skip, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), ctx, top, skip, tenantID)
}

// GetByName mocks base method.
func (m *MockRepository) GetByName(ctx context.Context, configName, tenantID string) (*entity.CIRAConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, configName, tenantID)
	ret0, _ := ret[0].(*entity.CIRAConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockRepositoryMockRecorder) GetByName(ctx, configName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockRepository)(nil).GetByName), ctx, configName, tenantID)
}

// GetCount mocks base method.
func (m *MockRepository) GetCount(ctx context.Context, tenantID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", ctx, tenantID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockRepositoryMockRecorder) GetCount(ctx, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockRepository)(nil).GetCount), ctx, tenantID)
}

// Insert mocks base method.
func (m *MockRepository) Insert(ctx context.Context, p *entity.CIRAConfig) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, p)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockRepositoryMockRecorder) Insert(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockRepository)(nil).Insert), ctx, p)
}

// Update mocks base method.
func (m *MockRepository) Update(ctx context.Context, p *entity.CIRAConfig) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, p)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), ctx, p)
}

// MockFeature is a mock of Feature interface.
type MockFeature struct {
	ctrl     *gomock.Controller
	recorder *MockFeatureMockRecorder
}

// MockFeatureMockRecorder is the mock recorder for MockFeature.
type MockFeatureMockRecorder struct {
	mock *MockFeature
}

// NewMockFeature creates a new mock instance.
func NewMockFeature(ctrl *gomock.Controller) *MockFeature {
	mock := &MockFeature{ctrl: ctrl}
	mock.recorder = &MockFeatureMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFeature) EXPECT() *MockFeatureMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockFeature) Delete(ctx context.Context, profileName, tenantID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, profileName, tenantID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockFeatureMockRecorder) Delete(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFeature)(nil).Delete), ctx, profileName, tenantID)
}

// Get mocks base method.
func (m *MockFeature) Get(ctx context.Context, top, skip int, tenantID string) ([]dtov1.CIRAConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, top, skip, tenantID)
	ret0, _ := ret[0].([]dtov1.CIRAConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockFeatureMockRecorder) Get(ctx, top, skip, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockFeature)(nil).Get), ctx, top, skip, tenantID)
}

// GetByName mocks base method.
func (m *MockFeature) GetByName(ctx context.Context, configName, tenantID string) (*dtov1.CIRAConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, configName, tenantID)
	ret0, _ := ret[0].(*dtov1.CIRAConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockFeatureMockRecorder) GetByName(ctx, configName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockFeature)(nil).GetByName), ctx, configName, tenantID)
}

// GetCount mocks base method.
func (m *MockFeature) GetCount(ctx context.Context, tenantID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", ctx, tenantID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockFeatureMockRecorder) GetCount(ctx, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockFeature)(nil).GetCount), ctx, tenantID)
}

// Insert mocks base method.
func (m *MockFeature) Insert(ctx context.Context, p *dtov1.CIRAConfig) (*dtov1.CIRAConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, p)
	ret0, _ := ret[0].(*dtov1.CIRAConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockFeatureMockRecorder) Insert(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockFeature)(nil).Insert), ctx, p)
}

// Update mocks base method.
func (m *MockFeature) Update(ctx context.Context, p *dtov1.CIRAConfig) (*dtov1.CIRAConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, p)
	ret0, _ := ret[0].(*dtov1.CIRAConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockFeatureMockRecorder) Update(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockFeature)(nil).Update), ctx, p)
}
