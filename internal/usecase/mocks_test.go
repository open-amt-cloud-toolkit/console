// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package usecase_test is a generated GoMock package.
package usecase_test

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/open-amt-cloud-toolkit/console/internal/entity"
)

// MockDomain is a mock of Domain interface.
type MockDomain struct {
	ctrl     *gomock.Controller
	recorder *MockDomainMockRecorder
}

// MockDomainMockRecorder is the mock recorder for MockDomain.
type MockDomainMockRecorder struct {
	mock *MockDomain
}

// NewMockDomain creates a new mock instance.
func NewMockDomain(ctrl *gomock.Controller) *MockDomain {
	mock := &MockDomain{ctrl: ctrl}
	mock.recorder = &MockDomainMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDomain) EXPECT() *MockDomainMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDomain) Delete(ctx context.Context, domainName, tenantId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, domainName, tenantId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockDomainMockRecorder) Delete(ctx, domainName, tenantId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDomain)(nil).Delete), ctx, domainName, tenantId)
}

// Get mocks base method.
func (m *MockDomain) Get(ctx context.Context, top, skip int, tenantId string) ([]entity.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, top, skip, tenantId)
	ret0, _ := ret[0].([]entity.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDomainMockRecorder) Get(ctx, top, skip, tenantId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDomain)(nil).Get), ctx, top, skip, tenantId)
}

// GetByName mocks base method.
func (m *MockDomain) GetByName(ctx context.Context, domainName, tenantId string) (*entity.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, domainName, tenantId)
	ret0, _ := ret[0].(*entity.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockDomainMockRecorder) GetByName(ctx, domainName, tenantId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockDomain)(nil).GetByName), ctx, domainName, tenantId)
}

// GetCount mocks base method.
func (m *MockDomain) GetCount(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockDomainMockRecorder) GetCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockDomain)(nil).GetCount), arg0, arg1)
}

// GetDomainByDomainSuffix mocks base method.
func (m *MockDomain) GetDomainByDomainSuffix(ctx context.Context, domainSuffix, tenantId string) (*entity.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainByDomainSuffix", ctx, domainSuffix, tenantId)
	ret0, _ := ret[0].(*entity.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainByDomainSuffix indicates an expected call of GetDomainByDomainSuffix.
func (mr *MockDomainMockRecorder) GetDomainByDomainSuffix(ctx, domainSuffix, tenantId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainByDomainSuffix", reflect.TypeOf((*MockDomain)(nil).GetDomainByDomainSuffix), ctx, domainSuffix, tenantId)
}

// Insert mocks base method.
func (m *MockDomain) Insert(ctx context.Context, d *entity.Domain) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, d)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockDomainMockRecorder) Insert(ctx, d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockDomain)(nil).Insert), ctx, d)
}

// Update mocks base method.
func (m *MockDomain) Update(ctx context.Context, d *entity.Domain) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, d)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockDomainMockRecorder) Update(ctx, d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDomain)(nil).Update), ctx, d)
}
