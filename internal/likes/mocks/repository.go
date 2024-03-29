// Code generated by MockGen. DO NOT EDIT.
// Source: internal/likes/repository.go

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
func (m *MockRepository) Create(pinId, authorId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", pinId, authorId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(pinId, authorId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), pinId, authorId)
}

// Delete mocks base method.
func (m *MockRepository) Delete(pinId, authorId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", pinId, authorId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(pinId, authorId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), pinId, authorId)
}

// LikeExists mocks base method.
func (m *MockRepository) LikeExists(pinId, authorId int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikeExists", pinId, authorId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LikeExists indicates an expected call of LikeExists.
func (mr *MockRepositoryMockRecorder) LikeExists(pinId, authorId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikeExists", reflect.TypeOf((*MockRepository)(nil).LikeExists), pinId, authorId)
}

// ListByAuthor mocks base method.
func (m *MockRepository) ListByAuthor(authorId int) ([]models.Like, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByAuthor", authorId)
	ret0, _ := ret[0].([]models.Like)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByAuthor indicates an expected call of ListByAuthor.
func (mr *MockRepositoryMockRecorder) ListByAuthor(authorId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByAuthor", reflect.TypeOf((*MockRepository)(nil).ListByAuthor), authorId)
}

// ListByPin mocks base method.
func (m *MockRepository) ListByPin(pinId int) ([]models.Like, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByPin", pinId)
	ret0, _ := ret[0].([]models.Like)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByPin indicates an expected call of ListByPin.
func (mr *MockRepositoryMockRecorder) ListByPin(pinId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByPin", reflect.TypeOf((*MockRepository)(nil).ListByPin), pinId)
}

// PinExists mocks base method.
func (m *MockRepository) PinExists(pinId int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PinExists", pinId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PinExists indicates an expected call of PinExists.
func (mr *MockRepositoryMockRecorder) PinExists(pinId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PinExists", reflect.TypeOf((*MockRepository)(nil).PinExists), pinId)
}

// UserExists mocks base method.
func (m *MockRepository) UserExists(userId int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserExists", userId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserExists indicates an expected call of UserExists.
func (mr *MockRepositoryMockRecorder) UserExists(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserExists", reflect.TypeOf((*MockRepository)(nil).UserExists), userId)
}
