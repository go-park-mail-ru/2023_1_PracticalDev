// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pins/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CheckReadAccess mocks base method.
func (m *MockService) CheckReadAccess(userId, pinId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckReadAccess", userId, pinId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckReadAccess indicates an expected call of CheckReadAccess.
func (mr *MockServiceMockRecorder) CheckReadAccess(userId, pinId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckReadAccess", reflect.TypeOf((*MockService)(nil).CheckReadAccess), userId, pinId)
}

// CheckWriteAccess mocks base method.
func (m *MockService) CheckWriteAccess(userId, pinId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckWriteAccess", userId, pinId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckWriteAccess indicates an expected call of CheckWriteAccess.
func (mr *MockServiceMockRecorder) CheckWriteAccess(userId, pinId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckWriteAccess", reflect.TypeOf((*MockService)(nil).CheckWriteAccess), userId, pinId)
}

// Create mocks base method.
func (m *MockService) Create(params *pins.CreateParams) (models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", params)
	ret0, _ := ret[0].(models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockServiceMockRecorder) Create(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockService)(nil).Create), params)
}

// Delete mocks base method.
func (m *MockService) Delete(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), id)
}

// FullUpdate mocks base method.
func (m *MockService) FullUpdate(params *pins.FullUpdateParams) (models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullUpdate", params)
	ret0, _ := ret[0].(models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FullUpdate indicates an expected call of FullUpdate.
func (mr *MockServiceMockRecorder) FullUpdate(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullUpdate", reflect.TypeOf((*MockService)(nil).FullUpdate), params)
}

// Get mocks base method.
func (m *MockService) Get(id, userId int) (models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id, userId)
	ret0, _ := ret[0].(models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockServiceMockRecorder) Get(id, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockService)(nil).Get), id, userId)
}

// List mocks base method.
func (m *MockService) List(userId, page, limit int) ([]models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", userId, page, limit)
	ret0, _ := ret[0].([]models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockServiceMockRecorder) List(userId, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService)(nil).List), userId, page, limit)
}

// ListByAuthor mocks base method.
func (m *MockService) ListByAuthor(authorId, userId, page, limit int) ([]models.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByAuthor", authorId, userId, page, limit)
	ret0, _ := ret[0].([]models.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByAuthor indicates an expected call of ListByAuthor.
func (mr *MockServiceMockRecorder) ListByAuthor(authorId, userId, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByAuthor", reflect.TypeOf((*MockService)(nil).ListByAuthor), authorId, userId, page, limit)
}

// SetLikedField mocks base method.
func (m *MockService) SetLikedField(pin *models.Pin, userId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLikedField", pin, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLikedField indicates an expected call of SetLikedField.
func (mr *MockServiceMockRecorder) SetLikedField(pin, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLikedField", reflect.TypeOf((*MockService)(nil).SetLikedField), pin, userId)
}
