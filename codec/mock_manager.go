// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/codec (interfaces: Manager)

// Package codec is a generated GoMock package.
package codec

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Marshal mocks base method.
func (m *MockManager) Marshal(arg0 uint16, arg1 interface{}) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Marshal", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Marshal indicates an expected call of Marshal.
func (mr *MockManagerMockRecorder) Marshal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Marshal", reflect.TypeOf((*MockManager)(nil).Marshal), arg0, arg1)
}

// RegisterCodec mocks base method.
func (m *MockManager) RegisterCodec(arg0 uint16, arg1 Codec) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterCodec", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterCodec indicates an expected call of RegisterCodec.
func (mr *MockManagerMockRecorder) RegisterCodec(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCodec", reflect.TypeOf((*MockManager)(nil).RegisterCodec), arg0, arg1)
}

// SetMaxSize mocks base method.
func (m *MockManager) SetMaxSize(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxSize", arg0)
}

// SetMaxSize indicates an expected call of SetMaxSize.
func (mr *MockManagerMockRecorder) SetMaxSize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxSize", reflect.TypeOf((*MockManager)(nil).SetMaxSize), arg0)
}

// Size mocks base method.
func (m *MockManager) Size(arg0 uint16, arg1 interface{}) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Size indicates an expected call of Size.
func (mr *MockManagerMockRecorder) Size(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockManager)(nil).Size), arg0, arg1)
}

// Unmarshal mocks base method.
func (m *MockManager) Unmarshal(arg0 []byte, arg1 interface{}) (uint16, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", arg0, arg1)
	ret0, _ := ret[0].(uint16)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unmarshal indicates an expected call of Unmarshal.
func (mr *MockManagerMockRecorder) Unmarshal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockManager)(nil).Unmarshal), arg0, arg1)
}
