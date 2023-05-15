// Code generated by MockGen. DO NOT EDIT.
// Source: internal/followings/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	followings "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
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
func (m *MockRepository) Create(followerId, followeeId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", followerId, followeeId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(followerId, followeeId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), followerId, followeeId)
}

// Delete mocks base method.
func (m *MockRepository) Delete(followerId, followeeId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", followerId, followeeId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(followerId, followeeId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), followerId, followeeId)
}

// FollowingExists mocks base method.
func (m *MockRepository) FollowingExists(followerId, followeeId int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FollowingExists", followerId, followeeId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FollowingExists indicates an expected call of FollowingExists.
func (mr *MockRepositoryMockRecorder) FollowingExists(followerId, followeeId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FollowingExists", reflect.TypeOf((*MockRepository)(nil).FollowingExists), followerId, followeeId)
}

// GetFollowees mocks base method.
func (m *MockRepository) GetFollowees(userId int) ([]followings.Followee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFollowees", userId)
	ret0, _ := ret[0].([]followings.Followee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFollowees indicates an expected call of GetFollowees.
func (mr *MockRepositoryMockRecorder) GetFollowees(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFollowees", reflect.TypeOf((*MockRepository)(nil).GetFollowees), userId)
}

// GetFollowers mocks base method.
func (m *MockRepository) GetFollowers(userId int) ([]followings.Follower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFollowers", userId)
	ret0, _ := ret[0].([]followings.Follower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFollowers indicates an expected call of GetFollowers.
func (mr *MockRepositoryMockRecorder) GetFollowers(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFollowers", reflect.TypeOf((*MockRepository)(nil).GetFollowers), userId)
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