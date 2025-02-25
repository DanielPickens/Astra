// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/astra/cmdline/cmdline.go

// Package cmdline is a generated GoMock package.
package cmdline

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	kclient "github\.com/danielpickens/astra/pkg/kclient"
)

// MockCmdline is a mock of Cmdline interface.
type MockCmdline struct {
	ctrl     *gomock.Controller
	recorder *MockCmdlineMockRecorder
}

// MockCmdlineMockRecorder is the mock recorder for MockCmdline.
type MockCmdlineMockRecorder struct {
	mock *MockCmdline
}

// NewMockCmdline creates a new mock instance.
func NewMockCmdline(ctrl *gomock.Controller) *MockCmdline {
	mock := &MockCmdline{ctrl: ctrl}
	mock.recorder = &MockCmdlineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCmdline) EXPECT() *MockCmdlineMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockCmdline) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockCmdlineMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockCmdline)(nil).Context))
}

// FlagValue mocks base method.
func (m *MockCmdline) FlagValue(flagName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlagValue", flagName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FlagValue indicates an expected call of FlagValue.
func (mr *MockCmdlineMockRecorder) FlagValue(flagName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlagValue", reflect.TypeOf((*MockCmdline)(nil).FlagValue), flagName)
}

// FlagValueIfSet mocks base method.
func (m *MockCmdline) FlagValueIfSet(flagName string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlagValueIfSet", flagName)
	ret0, _ := ret[0].(string)
	return ret0
}

// FlagValueIfSet indicates an expected call of FlagValueIfSet.
func (mr *MockCmdlineMockRecorder) FlagValueIfSet(flagName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlagValueIfSet", reflect.TypeOf((*MockCmdline)(nil).FlagValueIfSet), flagName)
}

// FlagValues mocks base method.
func (m *MockCmdline) FlagValues(flagName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlagValues", flagName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FlagValues indicates an expected call of FlagValues.
func (mr *MockCmdlineMockRecorder) FlagValues(flagName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlagValues", reflect.TypeOf((*MockCmdline)(nil).FlagValues), flagName)
}

// GetArgsAfterDashes mocks base method.
func (m *MockCmdline) GetArgsAfterDashes(args []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArgsAfterDashes", args)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArgsAfterDashes indicates an expected call of GetArgsAfterDashes.
func (mr *MockCmdlineMockRecorder) GetArgsAfterDashes(args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArgsAfterDashes", reflect.TypeOf((*MockCmdline)(nil).GetArgsAfterDashes), args)
}

// GetFlags mocks base method.
func (m *MockCmdline) GetFlags() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlags")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetFlags indicates an expected call of GetFlags.
func (mr *MockCmdlineMockRecorder) GetFlags() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlags", reflect.TypeOf((*MockCmdline)(nil).GetFlags))
}

// GetKubeClient mocks base method.
func (m *MockCmdline) GetKubeClient() (kclient.ClientInterface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKubeClient")
	ret0, _ := ret[0].(kclient.ClientInterface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKubeClient indicates an expected call of GetKubeClient.
func (mr *MockCmdlineMockRecorder) GetKubeClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKubeClient", reflect.TypeOf((*MockCmdline)(nil).GetKubeClient))
}

// GetName mocks base method.
func (m *MockCmdline) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockCmdlineMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockCmdline)(nil).GetName))
}

// GetParentName mocks base method.
func (m *MockCmdline) GetParentName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParentName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetParentName indicates an expected call of GetParentName.
func (mr *MockCmdlineMockRecorder) GetParentName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParentName", reflect.TypeOf((*MockCmdline)(nil).GetParentName))
}

// GetRootName mocks base method.
func (m *MockCmdline) GetRootName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRootName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetRootName indicates an expected call of GetRootName.
func (mr *MockCmdlineMockRecorder) GetRootName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRootName", reflect.TypeOf((*MockCmdline)(nil).GetRootName))
}

// GetWorkingDirectory mocks base method.
func (m *MockCmdline) GetWorkingDirectory() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkingDirectory")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkingDirectory indicates an expected call of GetWorkingDirectory.
func (mr *MockCmdlineMockRecorder) GetWorkingDirectory() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkingDirectory", reflect.TypeOf((*MockCmdline)(nil).GetWorkingDirectory))
}

// IsFlagSet mocks base method.
func (m *MockCmdline) IsFlagSet(flagName string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFlagSet", flagName)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsFlagSet indicates an expected call of IsFlagSet.
func (mr *MockCmdlineMockRecorder) IsFlagSet(flagName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFlagSet", reflect.TypeOf((*MockCmdline)(nil).IsFlagSet), flagName)
}
