// Code generated by MockGen. DO NOT EDIT.
// Source: internal/boards/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	boards "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	models "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
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
func (m *MockService) CheckReadAccess(userId, boardId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckReadAccess", userId, boardId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckReadAccess indicates an expected call of CheckReadAccess.
func (mr *MockServiceMockRecorder) CheckReadAccess(userId, boardId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckReadAccess", reflect.TypeOf((*MockService)(nil).CheckReadAccess), userId, boardId)
}

// CheckWriteAccess mocks base method.
func (m *MockService) CheckWriteAccess(userId, boardId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckWriteAccess", userId, boardId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckWriteAccess indicates an expected call of CheckWriteAccess.
func (mr *MockServiceMockRecorder) CheckWriteAccess(userId, boardId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckWriteAccess", reflect.TypeOf((*MockService)(nil).CheckWriteAccess), userId, boardId)
}

// Create mocks base method.
func (m *MockService) Create(params *boards.CreateParams) (models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", params)
	ret0, _ := ret[0].(models.Board)
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
func (m *MockService) FullUpdate(params *boards.FullUpdateParams) (models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullUpdate", params)
	ret0, _ := ret[0].(models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FullUpdate indicates an expected call of FullUpdate.
func (mr *MockServiceMockRecorder) FullUpdate(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullUpdate", reflect.TypeOf((*MockService)(nil).FullUpdate), params)
}

// Get mocks base method.
func (m *MockService) Get(id int) (models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockServiceMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockService)(nil).Get), id)
}

// List mocks base method.
func (m *MockService) List(userId int) ([]models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", userId)
	ret0, _ := ret[0].([]models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockServiceMockRecorder) List(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService)(nil).List), userId)
}

// PartialUpdate mocks base method.
func (m *MockService) PartialUpdate(params *boards.PartialUpdateParams) (models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PartialUpdate", params)
	ret0, _ := ret[0].(models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PartialUpdate indicates an expected call of PartialUpdate.
func (mr *MockServiceMockRecorder) PartialUpdate(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PartialUpdate", reflect.TypeOf((*MockService)(nil).PartialUpdate), params)
}