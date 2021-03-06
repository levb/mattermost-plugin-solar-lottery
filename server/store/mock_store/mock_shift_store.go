// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-solar-lottery/server/store (interfaces: ShiftStore)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	gomock "github.com/golang/mock/gomock"
	store "github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	reflect "reflect"
)

// MockShiftStore is a mock of ShiftStore interface
type MockShiftStore struct {
	ctrl     *gomock.Controller
	recorder *MockShiftStoreMockRecorder
}

// MockShiftStoreMockRecorder is the mock recorder for MockShiftStore
type MockShiftStoreMockRecorder struct {
	mock *MockShiftStore
}

// NewMockShiftStore creates a new mock instance
func NewMockShiftStore(ctrl *gomock.Controller) *MockShiftStore {
	mock := &MockShiftStore{ctrl: ctrl}
	mock.recorder = &MockShiftStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockShiftStore) EXPECT() *MockShiftStoreMockRecorder {
	return m.recorder
}

// DeleteShift mocks base method
func (m *MockShiftStore) DeleteShift(arg0 string, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteShift", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteShift indicates an expected call of DeleteShift
func (mr *MockShiftStoreMockRecorder) DeleteShift(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteShift", reflect.TypeOf((*MockShiftStore)(nil).DeleteShift), arg0, arg1)
}

// LoadShift mocks base method
func (m *MockShiftStore) LoadShift(arg0 string, arg1 int) (*store.Shift, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadShift", arg0, arg1)
	ret0, _ := ret[0].(*store.Shift)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadShift indicates an expected call of LoadShift
func (mr *MockShiftStoreMockRecorder) LoadShift(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadShift", reflect.TypeOf((*MockShiftStore)(nil).LoadShift), arg0, arg1)
}

// StoreShift mocks base method
func (m *MockShiftStore) StoreShift(arg0 string, arg1 int, arg2 *store.Shift) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreShift", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreShift indicates an expected call of StoreShift
func (mr *MockShiftStoreMockRecorder) StoreShift(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreShift", reflect.TypeOf((*MockShiftStore)(nil).StoreShift), arg0, arg1, arg2)
}
