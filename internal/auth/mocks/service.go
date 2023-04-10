// Code generated by MockGen. DO NOT EDIT.
// Source: internal/auth/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	auth "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	models "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	api "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models/api"
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

// Authenticate mocks base method.
func (m *MockService) Authenticate(login, hashedPassword string) (models.User, auth.SessionParams, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", login, hashedPassword)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(auth.SessionParams)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockServiceMockRecorder) Authenticate(login, hashedPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockService)(nil).Authenticate), login, hashedPassword)
}

// CheckAuth mocks base method.
func (m *MockService) CheckAuth(userId, sessionId string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuth", userId, sessionId)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAuth indicates an expected call of CheckAuth.
func (mr *MockServiceMockRecorder) CheckAuth(userId, sessionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuth", reflect.TypeOf((*MockService)(nil).CheckAuth), userId, sessionId)
}

// CreateSession mocks base method.
func (m *MockService) CreateSession(userId int) auth.SessionParams {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", userId)
	ret0, _ := ret[0].(auth.SessionParams)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockServiceMockRecorder) CreateSession(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockService)(nil).CreateSession), userId)
}

// DeleteSession mocks base method.
func (m *MockService) DeleteSession(userId, sessionId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", userId, sessionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockServiceMockRecorder) DeleteSession(userId, sessionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockService)(nil).DeleteSession), userId, sessionId)
}

// Register mocks base method.
func (m *MockService) Register(user *api.RegisterParams) (models.User, auth.SessionParams, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(auth.SessionParams)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Register indicates an expected call of Register.
func (mr *MockServiceMockRecorder) Register(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockService)(nil).Register), user)
}

// SetSession mocks base method.
func (m *MockService) SetSession(id string, session *models.Session, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSession", id, session, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSession indicates an expected call of SetSession.
func (mr *MockServiceMockRecorder) SetSession(id, session, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSession", reflect.TypeOf((*MockService)(nil).SetSession), id, session, expiration)
}
