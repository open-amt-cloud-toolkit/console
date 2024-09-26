// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecase/domains/interfaces.go
//
// Generated by this command:
//
//	mockgen -source ./internal/usecase/domains/interfaces.go -package domains_test
//

// Package domains_test is a generated GoMock package.
package domains_test

import (
	context "context"
	reflect "reflect"

	entity "github.com/open-amt-cloud-toolkit/console/internal/entity"
	dto "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
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
func (m *MockRepository) Delete(ctx context.Context, name, tenantID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, name, tenantID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, name, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, name, tenantID)
}

// Get mocks base method.
func (m *MockRepository) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, top, skip, tenantID)
	ret0, _ := ret[0].([]entity.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(ctx, top, skip, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), ctx, top, skip, tenantID)
}

// GetByName mocks base method.
func (m *MockRepository) GetByName(ctx context.Context, name, tenantID string) (*entity.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, name, tenantID)
	ret0, _ := ret[0].(*entity.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockRepositoryMockRecorder) GetByName(ctx, name, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockRepository)(nil).GetByName), ctx, name, tenantID)
}

// GetCount mocks base method.
func (m *MockRepository) GetCount(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockRepositoryMockRecorder) GetCount(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockRepository)(nil).GetCount), arg0, arg1)
}

// GetDomainByDomainSuffix mocks base method.
func (m *MockRepository) GetDomainByDomainSuffix(ctx context.Context, domainSuffix, tenantID string) (*entity.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainByDomainSuffix", ctx, domainSuffix, tenantID)
	ret0, _ := ret[0].(*entity.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainByDomainSuffix indicates an expected call of GetDomainByDomainSuffix.
func (mr *MockRepositoryMockRecorder) GetDomainByDomainSuffix(ctx, domainSuffix, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainByDomainSuffix", reflect.TypeOf((*MockRepository)(nil).GetDomainByDomainSuffix), ctx, domainSuffix, tenantID)
}

// Insert mocks base method.
func (m *MockRepository) Insert(ctx context.Context, d *entity.Domain) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, d)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockRepositoryMockRecorder) Insert(ctx, d any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockRepository)(nil).Insert), ctx, d)
}

// Update mocks base method.
func (m *MockRepository) Update(ctx context.Context, d *entity.Domain) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, d)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(ctx, d any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), ctx, d)
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
func (m *MockFeature) Delete(ctx context.Context, name, tenantID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, name, tenantID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockFeatureMockRecorder) Delete(ctx, name, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFeature)(nil).Delete), ctx, name, tenantID)
}

// Get mocks base method.
func (m *MockFeature) Get(ctx context.Context, top, skip int, tenantID string) ([]dto.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, top, skip, tenantID)
	ret0, _ := ret[0].([]dto.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockFeatureMockRecorder) Get(ctx, top, skip, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockFeature)(nil).Get), ctx, top, skip, tenantID)
}

// GetByName mocks base method.
func (m *MockFeature) GetByName(ctx context.Context, name, tenantID string) (*dto.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, name, tenantID)
	ret0, _ := ret[0].(*dto.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockFeatureMockRecorder) GetByName(ctx, name, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockFeature)(nil).GetByName), ctx, name, tenantID)
}

// GetCount mocks base method.
func (m *MockFeature) GetCount(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockFeatureMockRecorder) GetCount(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockFeature)(nil).GetCount), arg0, arg1)
}

// GetDomainByDomainSuffix mocks base method.
func (m *MockFeature) GetDomainByDomainSuffix(ctx context.Context, domainSuffix, tenantID string) (*dto.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainByDomainSuffix", ctx, domainSuffix, tenantID)
	ret0, _ := ret[0].(*dto.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainByDomainSuffix indicates an expected call of GetDomainByDomainSuffix.
func (mr *MockFeatureMockRecorder) GetDomainByDomainSuffix(ctx, domainSuffix, tenantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainByDomainSuffix", reflect.TypeOf((*MockFeature)(nil).GetDomainByDomainSuffix), ctx, domainSuffix, tenantID)
}

// Insert mocks base method.
func (m *MockFeature) Insert(ctx context.Context, d *dto.Domain) (*dto.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, d)
	ret0, _ := ret[0].(*dto.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockFeatureMockRecorder) Insert(ctx, d any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockFeature)(nil).Insert), ctx, d)
}

// Update mocks base method.
func (m *MockFeature) Update(ctx context.Context, d *dto.Domain) (*dto.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, d)
	ret0, _ := ret[0].(*dto.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockFeatureMockRecorder) Update(ctx, d any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockFeature)(nil).Update), ctx, d)
}
