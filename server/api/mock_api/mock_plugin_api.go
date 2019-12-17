// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-solar-lottery/server/api (interfaces: PluginAPI)

// Package mock_api is a generated GoMock package.
package mock_api

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/mattermost/mattermost-server/v5/model"
	reflect "reflect"
)

// MockPluginAPI is a mock of PluginAPI interface
type MockPluginAPI struct {
	ctrl     *gomock.Controller
	recorder *MockPluginAPIMockRecorder
}

// MockPluginAPIMockRecorder is the mock recorder for MockPluginAPI
type MockPluginAPIMockRecorder struct {
	mock *MockPluginAPI
}

// NewMockPluginAPI creates a new mock instance
func NewMockPluginAPI(ctrl *gomock.Controller) *MockPluginAPI {
	mock := &MockPluginAPI{ctrl: ctrl}
	mock.recorder = &MockPluginAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPluginAPI) EXPECT() *MockPluginAPIMockRecorder {
	return m.recorder
}

// GetMattermostUser mocks base method
func (m *MockPluginAPI) GetMattermostUser(arg0 string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMattermostUser", arg0)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMattermostUser indicates an expected call of GetMattermostUser
func (mr *MockPluginAPIMockRecorder) GetMattermostUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMattermostUser", reflect.TypeOf((*MockPluginAPI)(nil).GetMattermostUser), arg0)
}

// GetMattermostUserByUsername mocks base method
func (m *MockPluginAPI) GetMattermostUserByUsername(arg0 string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMattermostUserByUsername", arg0)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMattermostUserByUsername indicates an expected call of GetMattermostUserByUsername
func (mr *MockPluginAPIMockRecorder) GetMattermostUserByUsername(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMattermostUserByUsername", reflect.TypeOf((*MockPluginAPI)(nil).GetMattermostUserByUsername), arg0)
}

// IsPluginAdmin mocks base method
func (m *MockPluginAPI) IsPluginAdmin(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPluginAdmin", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPluginAdmin indicates an expected call of IsPluginAdmin
func (mr *MockPluginAPIMockRecorder) IsPluginAdmin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPluginAdmin", reflect.TypeOf((*MockPluginAPI)(nil).IsPluginAdmin), arg0)
}
