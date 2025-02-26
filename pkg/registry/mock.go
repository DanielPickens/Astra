// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/registry/interface.go

// Package registry is a generated GoMock package.
package registry

import (
	context "context"
	reflect "reflect"

	v1alpha2 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	util "github.com/devfile/library/v2/pkg/util"
	library "github.com/devfile/registry-support/registry-library/library"
	gomock "github.com/golang/mock/gomock"
	api "github\.com/danielpickens/astra/pkg/api"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// DownloadFileInMemory mocks base method.
func (m *MockClient) DownloadFileInMemory(params util.HTTPRequestParams) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFileInMemory", params)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadFileInMemory indicates an expected call of DownloadFileInMemory.
func (mr *MockClientMockRecorder) DownloadFileInMemory(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFileInMemory", reflect.TypeOf((*MockClient)(nil).DownloadFileInMemory), params)
}

// DownloadStarterProject mocks base method.
func (m *MockClient) DownloadStarterProject(starterProject *v1alpha2.StarterProject, decryptedToken, contextDir string, verbose bool) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadStarterProject", starterProject, decryptedToken, contextDir, verbose)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadStarterProject indicates an expected call of DownloadStarterProject.
func (mr *MockClientMockRecorder) DownloadStarterProject(starterProject, decryptedToken, contextDir, verbose interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadStarterProject", reflect.TypeOf((*MockClient)(nil).DownloadStarterProject), starterProject, decryptedToken, contextDir, verbose)
}

// GetDevfileRegistries mocks base method.
func (m *MockClient) GetDevfileRegistries(registryName string) ([]api.Registry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevfileRegistries", registryName)
	ret0, _ := ret[0].([]api.Registry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDevfileRegistries indicates an expected call of GetDevfileRegistries.
func (mr *MockClientMockRecorder) GetDevfileRegistries(registryName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevfileRegistries", reflect.TypeOf((*MockClient)(nil).GetDevfileRegistries), registryName)
}

// ListDevfileStacks mocks base method.
func (m *MockClient) ListDevfileStacks(ctx context.Context, registryName, devfileFlag, filterFlag string, detailsFlag, withDevfileContent bool) (DevfileStackList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDevfileStacks", ctx, registryName, devfileFlag, filterFlag, detailsFlag, withDevfileContent)
	ret0, _ := ret[0].(DevfileStackList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDevfileStacks indicates an expected call of ListDevfileStacks.
func (mr *MockClientMockRecorder) ListDevfileStacks(ctx, registryName, devfileFlag, filterFlag, detailsFlag, withDevfileContent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDevfileStacks", reflect.TypeOf((*MockClient)(nil).ListDevfileStacks), ctx, registryName, devfileFlag, filterFlag, detailsFlag, withDevfileContent)
}

// PullStackFromRegistry mocks base method.
func (m *MockClient) PullStackFromRegistry(registry, stack, destDir string, options library.RegistryOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullStackFromRegistry", registry, stack, destDir, options)
	ret0, _ := ret[0].(error)
	return ret0
}

// PullStackFromRegistry indicates an expected call of PullStackFromRegistry.
func (mr *MockClientMockRecorder) PullStackFromRegistry(registry, stack, destDir, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullStackFromRegistry", reflect.TypeOf((*MockClient)(nil).PullStackFromRegistry), registry, stack, destDir, options)
}
