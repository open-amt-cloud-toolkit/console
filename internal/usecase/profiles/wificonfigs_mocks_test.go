// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecase/wificonfigs/interfaces.go
//
// Generated by this command:
//
//	mockgen -source ./internal/usecase/wificonfigs/interfaces.go -package profiles_test -mock_names Repository=MockwificonfigsRepository,Feature=MockwificonfigsFeature
//

// Package profiles_test is a generated GoMock package.
package profiles_test

import (
	context "context"
	reflect "reflect"

	entity "github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	gomock "go.uber.org/mock/gomock"
)

// MockwificonfigsRepository is a mock of Repository interface.
type MockwificonfigsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockwificonfigsRepositoryMockRecorder
}

// MockwificonfigsRepositoryMockRecorder is the mock recorder for MockwificonfigsRepository.
type MockwificonfigsRepositoryMockRecorder struct {
	mock *MockwificonfigsRepository
}

// NewMockwificonfigsRepository creates a new mock instance.
func NewMockwificonfigsRepository(ctrl *gomock.Controller) *MockwificonfigsRepository {
	mock := &MockwificonfigsRepository{ctrl: ctrl}
	mock.recorder = &MockwificonfigsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockwificonfigsRepository) EXPECT() *MockwificonfigsRepositoryMockRecorder {
	return m.recorder
}

// CheckProfileExists mocks base method.
func (m *MockwificonfigsRepository) CheckProfileExists(ctx context.Context, profileName, tenantID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckProfileExists", ctx, profileName, tenantID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckProfileExists indicates an expected call of CheckProfileExists.
func (mr *MockwificonfigsRepositoryMockRecorder) CheckProfileExists(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckProfileExists", reflect.TypeOf((*MockwificonfigsRepository)(nil).CheckProfileExists), ctx, profileName, tenantID)
}

// Delete mocks base method.
func (m *MockwificonfigsRepository) Delete(ctx context.Context, profileName, tenantID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, profileName, tenantID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockwificonfigsRepositoryMockRecorder) Delete(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockwificonfigsRepository)(nil).Delete), ctx, profileName, tenantID)
}

// Get mocks base method.
func (m *MockwificonfigsRepository) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.WirelessConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, top, skip, tenantID)
	ret0, _ := ret[0].([]entity.WirelessConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockwificonfigsRepositoryMockRecorder) Get(ctx, top, skip, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockwificonfigsRepository)(nil).Get), ctx, top, skip, tenantID)
}

// GetByName mocks base method.
func (m *MockwificonfigsRepository) GetByName(ctx context.Context, guid, tenantID string) (*entity.WirelessConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, guid, tenantID)
	ret0, _ := ret[0].(*entity.WirelessConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockwificonfigsRepositoryMockRecorder) GetByName(ctx, guid, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockwificonfigsRepository)(nil).GetByName), ctx, guid, tenantID)
}

// GetCount mocks base method.
func (m *MockwificonfigsRepository) GetCount(ctx context.Context, tenantID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", ctx, tenantID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockwificonfigsRepositoryMockRecorder) GetCount(ctx, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockwificonfigsRepository)(nil).GetCount), ctx, tenantID)
}

// Insert mocks base method.
func (m *MockwificonfigsRepository) Insert(ctx context.Context, p *entity.WirelessConfig) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, p)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockwificonfigsRepositoryMockRecorder) Insert(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockwificonfigsRepository)(nil).Insert), ctx, p)
}

// Update mocks base method.
func (m *MockwificonfigsRepository) Update(ctx context.Context, p *entity.WirelessConfig) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, p)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockwificonfigsRepositoryMockRecorder) Update(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockwificonfigsRepository)(nil).Update), ctx, p)
}

// MockwificonfigsFeature is a mock of Feature interface.
type MockwificonfigsFeature struct {
	ctrl     *gomock.Controller
	recorder *MockwificonfigsFeatureMockRecorder
}

// MockwificonfigsFeatureMockRecorder is the mock recorder for MockwificonfigsFeature.
type MockwificonfigsFeatureMockRecorder struct {
	mock *MockwificonfigsFeature
}

// NewMockwificonfigsFeature creates a new mock instance.
func NewMockwificonfigsFeature(ctrl *gomock.Controller) *MockwificonfigsFeature {
	mock := &MockwificonfigsFeature{ctrl: ctrl}
	mock.recorder = &MockwificonfigsFeatureMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockwificonfigsFeature) EXPECT() *MockwificonfigsFeatureMockRecorder {
	return m.recorder
}

// CheckProfileExists mocks base method.
func (m *MockwificonfigsFeature) CheckProfileExists(ctx context.Context, profileName, tenantID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckProfileExists", ctx, profileName, tenantID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckProfileExists indicates an expected call of CheckProfileExists.
func (mr *MockwificonfigsFeatureMockRecorder) CheckProfileExists(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckProfileExists", reflect.TypeOf((*MockwificonfigsFeature)(nil).CheckProfileExists), ctx, profileName, tenantID)
}

// Delete mocks base method.
func (m *MockwificonfigsFeature) Delete(ctx context.Context, profileName, tenantID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, profileName, tenantID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockwificonfigsFeatureMockRecorder) Delete(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockwificonfigsFeature)(nil).Delete), ctx, profileName, tenantID)
}

// Get mocks base method.
func (m *MockwificonfigsFeature) Get(ctx context.Context, top, skip int, tenantID string) ([]dtov1.WirelessConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, top, skip, tenantID)
	ret0, _ := ret[0].([]dtov1.WirelessConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockwificonfigsFeatureMockRecorder) Get(ctx, top, skip, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockwificonfigsFeature)(nil).Get), ctx, top, skip, tenantID)
}

// GetByName mocks base method.
func (m *MockwificonfigsFeature) GetByName(ctx context.Context, guid, tenantID string) (*dtov1.WirelessConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, guid, tenantID)
	ret0, _ := ret[0].(*dtov1.WirelessConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockwificonfigsFeatureMockRecorder) GetByName(ctx, guid, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockwificonfigsFeature)(nil).GetByName), ctx, guid, tenantID)
}

// GetCount mocks base method.
func (m *MockwificonfigsFeature) GetCount(ctx context.Context, tenantID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", ctx, tenantID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockwificonfigsFeatureMockRecorder) GetCount(ctx, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockwificonfigsFeature)(nil).GetCount), ctx, tenantID)
}

// Insert mocks base method.
func (m *MockwificonfigsFeature) Insert(ctx context.Context, p *dtov1.WirelessConfig) (*dtov1.WirelessConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, p)
	ret0, _ := ret[0].(*dtov1.WirelessConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockwificonfigsFeatureMockRecorder) Insert(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockwificonfigsFeature)(nil).Insert), ctx, p)
}

// Update mocks base method.
func (m *MockwificonfigsFeature) Update(ctx context.Context, p *dtov1.WirelessConfig) (*dtov1.WirelessConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, p)
	ret0, _ := ret[0].(*dtov1.WirelessConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockwificonfigsFeatureMockRecorder) Update(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockwificonfigsFeature)(nil).Update), ctx, p)
}
