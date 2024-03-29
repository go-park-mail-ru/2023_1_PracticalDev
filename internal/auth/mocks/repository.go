// Code generated by MockGen. DO NOT EDIT.
// Source: internal/auth/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

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

// Authenticate mocks base method.
func (m *MockRepository) Authenticate(email, password string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", email, password)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockRepositoryMockRecorder) Authenticate(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockRepository)(nil).Authenticate), email, password)
}

// CheckAuth mocks base method.
func (m *MockRepository) CheckAuth(userId, sessionId string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuth", userId, sessionId)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAuth indicates an expected call of CheckAuth.
func (mr *MockRepositoryMockRecorder) CheckAuth(userId, sessionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuth", reflect.TypeOf((*MockRepository)(nil).CheckAuth), userId, sessionId)
}

// DeleteSession mocks base method.
func (m *MockRepository) DeleteSession(userId, sessionId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", userId, sessionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockRepositoryMockRecorder) DeleteSession(userId, sessionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockRepository)(nil).DeleteSession), userId, sessionId)
}

// Register mocks base method.
func (m *MockRepository) Register(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockRepositoryMockRecorder) Register(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockRepository)(nil).Register), user)
}

// SetSession mocks base method.
func (m *MockRepository) SetSession(id string, session *models.Session, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSession", id, session, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSession indicates an expected call of SetSession.
func (mr *MockRepositoryMockRecorder) SetSession(id, session, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSession", reflect.TypeOf((*MockRepository)(nil).SetSession), id, session, expiration)
}
