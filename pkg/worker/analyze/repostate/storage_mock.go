// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package repostate is a generated GoMock package.
package repostate

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// UpdateState mocks base method
func (m *MockStorage) UpdateState(ctx context.Context, owner, name, analysisID string, state *State) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateState", ctx, owner, name, analysisID, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateState indicates an expected call of UpdateState
func (mr *MockStorageMockRecorder) UpdateState(ctx, owner, name, analysisID, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateState", reflect.TypeOf((*MockStorage)(nil).UpdateState), ctx, owner, name, analysisID, state)
}

// GetState mocks base method
func (m *MockStorage) GetState(ctx context.Context, owner, name, analysisID string) (*State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", ctx, owner, name, analysisID)
	ret0, _ := ret[0].(*State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState
func (mr *MockStorageMockRecorder) GetState(ctx, owner, name, analysisID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockStorage)(nil).GetState), ctx, owner, name, analysisID)
}
