// Code generated by MockGen. DO NOT EDIT.
// Source: internal/contracts.go
//
// Generated by this command:
//
//	mockgen -source=internal/contracts.go -destination=mocks/postgres_mock.go -package=mocks Storage
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	internal "github.com/real-splendid/url-shortener-practicum/internal"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// DeleteUserURLs mocks base method.
func (m *MockStorage) DeleteUserURLs(userID string, shortURLs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserURLs", userID, shortURLs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserURLs indicates an expected call of DeleteUserURLs.
func (mr *MockStorageMockRecorder) DeleteUserURLs(userID, shortURLs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserURLs", reflect.TypeOf((*MockStorage)(nil).DeleteUserURLs), userID, shortURLs)
}

// Get mocks base method.
func (m *MockStorage) Get(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStorageMockRecorder) Get(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStorage)(nil).Get), key)
}

// GetUserURLs mocks base method.
func (m *MockStorage) GetUserURLs(userID string) ([]internal.URLPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", userID)
	ret0, _ := ret[0].([]internal.URLPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockStorageMockRecorder) GetUserURLs(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockStorage)(nil).GetUserURLs), userID)
}

// Set mocks base method.
func (m *MockStorage) Set(key, value, userID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, value, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockStorageMockRecorder) Set(key, value, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStorage)(nil).Set), key, value, userID)
}
