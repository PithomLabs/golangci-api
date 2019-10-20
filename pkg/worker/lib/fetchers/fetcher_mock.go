// Code generated by MockGen. DO NOT EDIT.
// Source: fetcher.go

// Package fetchers is a generated GoMock package.
package fetchers

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	result "github.com/golangci/golangci-api/pkg/goenvbuild/result"
	executors "github.com/golangci/golangci-api/pkg/worker/lib/executors"
	reflect "reflect"
)

// MockFetcher is a mock of Fetcher interface
type MockFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockFetcherMockRecorder
}

// MockFetcherMockRecorder is the mock recorder for MockFetcher
type MockFetcherMockRecorder struct {
	mock *MockFetcher
}

// NewMockFetcher creates a new mock instance
func NewMockFetcher(ctrl *gomock.Controller) *MockFetcher {
	mock := &MockFetcher{ctrl: ctrl}
	mock.recorder = &MockFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFetcher) EXPECT() *MockFetcherMockRecorder {
	return m.recorder
}

// Fetch mocks base method
func (m *MockFetcher) Fetch(ctx context.Context, sg *result.StepGroup, repo *Repo, exec executors.Executor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", ctx, sg, repo, exec)
	ret0, _ := ret[0].(error)
	return ret0
}

// Fetch indicates an expected call of Fetch
func (mr *MockFetcherMockRecorder) Fetch(ctx, sg, repo, exec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockFetcher)(nil).Fetch), ctx, sg, repo, exec)
}
