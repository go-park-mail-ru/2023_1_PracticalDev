// Code generated by MockGen. DO NOT EDIT.
// Source: internal/notifications/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
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

// Create mocks base method.
func (m *MockRepository) Create(userID int, notificationType string, data interface{}) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userID, notificationType, data)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(userID, notificationType, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), userID, notificationType, data)
}

// Get mocks base method.
func (m *MockRepository) Get(notificationID int) (*models.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", notificationID)
	ret0, _ := ret[0].(*models.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(notificationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), notificationID)
}

// ListUnreadByUser mocks base method.
func (m *MockRepository) ListUnreadByUser(userID int) ([]models.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUnreadByUser", userID)
	ret0, _ := ret[0].([]models.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUnreadByUser indicates an expected call of ListUnreadByUser.
func (mr *MockRepositoryMockRecorder) ListUnreadByUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUnreadByUser", reflect.TypeOf((*MockRepository)(nil).ListUnreadByUser), userID)
}

// MarkAsRead mocks base method.
func (m *MockRepository) MarkAsRead(notificationID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkAsRead", notificationID)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkAsRead indicates an expected call of MarkAsRead.
func (mr *MockRepositoryMockRecorder) MarkAsRead(notificationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkAsRead", reflect.TypeOf((*MockRepository)(nil).MarkAsRead), notificationID)
}