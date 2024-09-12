// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecase/profiles/interfaces.go
//
// Generated by this command:
//
//	mockgen -source ./internal/usecase/profiles/interfaces.go -package v1 -mock_names Repository=MockProfilesRepository,Feature=MockProfilesFeature
//

// Package v1 is a generated GoMock package.
package v1

import (
	context "context"
	reflect "reflect"

	entity "github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	gomock "go.uber.org/mock/gomock"
)

// MockProfilesRepository is a mock of Repository interface.
type MockProfilesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProfilesRepositoryMockRecorder
}

// MockProfilesRepositoryMockRecorder is the mock recorder for MockProfilesRepository.
type MockProfilesRepositoryMockRecorder struct {
	mock *MockProfilesRepository
}

// NewMockProfilesRepository creates a new mock instance.
func NewMockProfilesRepository(ctrl *gomock.Controller) *MockProfilesRepository {
	mock := &MockProfilesRepository{ctrl: ctrl}
	mock.recorder = &MockProfilesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfilesRepository) EXPECT() *MockProfilesRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockProfilesRepository) Delete(ctx context.Context, profileName, tenantID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, profileName, tenantID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockProfilesRepositoryMockRecorder) Delete(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProfilesRepository)(nil).Delete), ctx, profileName, tenantID)
}

// Get mocks base method.
func (m *MockProfilesRepository) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, top, skip, tenantID)
	ret0, _ := ret[0].([]entity.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockProfilesRepositoryMockRecorder) Get(ctx, top, skip, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockProfilesRepository)(nil).Get), ctx, top, skip, tenantID)
}

// GetByName mocks base method.
func (m *MockProfilesRepository) GetByName(ctx context.Context, profileName, tenantID string) (*entity.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, profileName, tenantID)
	ret0, _ := ret[0].(*entity.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockProfilesRepositoryMockRecorder) GetByName(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockProfilesRepository)(nil).GetByName), ctx, profileName, tenantID)
}

// GetCount mocks base method.
func (m *MockProfilesRepository) GetCount(ctx context.Context, tenantID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", ctx, tenantID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockProfilesRepositoryMockRecorder) GetCount(ctx, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockProfilesRepository)(nil).GetCount), ctx, tenantID)
}

// Insert mocks base method.
func (m *MockProfilesRepository) Insert(ctx context.Context, p *entity.Profile) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, p)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockProfilesRepositoryMockRecorder) Insert(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockProfilesRepository)(nil).Insert), ctx, p)
}

// Update mocks base method.
func (m *MockProfilesRepository) Update(ctx context.Context, p *entity.Profile) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, p)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockProfilesRepositoryMockRecorder) Update(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockProfilesRepository)(nil).Update), ctx, p)
}

// MockProfilesFeature is a mock of Feature interface.
type MockProfilesFeature struct {
	ctrl     *gomock.Controller
	recorder *MockProfilesFeatureMockRecorder
}

// MockProfilesFeatureMockRecorder is the mock recorder for MockProfilesFeature.
type MockProfilesFeatureMockRecorder struct {
	mock *MockProfilesFeature
}

// NewMockProfilesFeature creates a new mock instance.
func NewMockProfilesFeature(ctrl *gomock.Controller) *MockProfilesFeature {
	mock := &MockProfilesFeature{ctrl: ctrl}
	mock.recorder = &MockProfilesFeatureMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfilesFeature) EXPECT() *MockProfilesFeatureMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockProfilesFeature) Delete(ctx context.Context, profileName, tenantID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, profileName, tenantID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockProfilesFeatureMockRecorder) Delete(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProfilesFeature)(nil).Delete), ctx, profileName, tenantID)
}

// Get mocks base method.
func (m *MockProfilesFeature) Get(ctx context.Context, top, skip int, tenantID string) ([]dtov1.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, top, skip, tenantID)
	ret0, _ := ret[0].([]dtov1.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockProfilesFeatureMockRecorder) Get(ctx, top, skip, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockProfilesFeature)(nil).Get), ctx, top, skip, tenantID)
}

// GetByName mocks base method.
func (m *MockProfilesFeature) GetByName(ctx context.Context, profileName, tenantID string) (*dtov1.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, profileName, tenantID)
	ret0, _ := ret[0].(*dtov1.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockProfilesFeatureMockRecorder) GetByName(ctx, profileName, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockProfilesFeature)(nil).GetByName), ctx, profileName, tenantID)
}

// GetCount mocks base method.
func (m *MockProfilesFeature) GetCount(ctx context.Context, tenantID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", ctx, tenantID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockProfilesFeatureMockRecorder) GetCount(ctx, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockProfilesFeature)(nil).GetCount), ctx, tenantID)
}

// Insert mocks base method.
func (m *MockProfilesFeature) Insert(ctx context.Context, p *dtov1.Profile) (*dtov1.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, p)
	ret0, _ := ret[0].(*dtov1.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockProfilesFeatureMockRecorder) Insert(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockProfilesFeature)(nil).Insert), ctx, p)
}

// Update mocks base method.
func (m *MockProfilesFeature) Update(ctx context.Context, p *dtov1.Profile) (*dtov1.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, p)
	ret0, _ := ret[0].(*dtov1.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockProfilesFeatureMockRecorder) Update(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockProfilesFeature)(nil).Update), ctx, p)
}
