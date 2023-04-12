// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pins/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	gomock "github.com/golang/mock/gomock"
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

// CheckReadAccess mocks base method.
func (m *MockRepository) CheckReadAccess(userId, pinId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckReadAccess", userId, pinId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckReadAccess indicates an expected call of CheckReadAccess.
func (mr *MockRepositoryMockRecorder) CheckReadAccess(userId, pinId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckReadAccess", reflect.TypeOf((*MockRepository)(nil).CheckReadAccess), userId, pinId)
}

// CheckWriteAccess mocks base method.
func (m *MockRepository) CheckWriteAccess(userId, pinId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckWriteAccess", userId, pinId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckWriteAccess indicates an expected call of CheckWriteAccess.
func (mr *MockRepositoryMockRecorder) CheckWriteAccess(userId, pinId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckWriteAccess", reflect.TypeOf((*MockRepository)(nil).CheckWriteAccess), userId, pinId)
}

// Create mocks base method.
func (m *MockRepository) Create(params *pins.CreateParams) (models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", params)
	ret0, _ := ret[0].(models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), params)
}

// Delete mocks base method.
func (m *MockRepository) Delete(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), id)
}

// FullUpdate mocks base method.
func (m *MockRepository) FullUpdate(params *pins.FullUpdateParams) (models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullUpdate", params)
	ret0, _ := ret[0].(models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FullUpdate indicates an expected call of FullUpdate.
func (mr *MockRepositoryMockRecorder) FullUpdate(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullUpdate", reflect.TypeOf((*MockRepository)(nil).FullUpdate), params)
}

// Get mocks base method.
func (m *MockRepository) Get(id int) (models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), id)
}

// List mocks base method.
func (m *MockRepository) List(page, limit int) ([]models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", page, limit)
	ret0, _ := ret[0].([]models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRepositoryMockRecorder) List(page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepository)(nil).List), page, limit)
}

// ListByUser mocks base method.
func (m *MockRepository) ListByUser(userId, page, limit int) ([]models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByUser", userId, page, limit)
	ret0, _ := ret[0].([]models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByUser indicates an expected call of ListByUser.
func (mr *MockRepositoryMockRecorder) ListByUser(userId, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByUser", reflect.TypeOf((*MockRepository)(nil).ListByUser), userId, page, limit)
}
