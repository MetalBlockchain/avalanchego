// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/database (interfaces: Iterator)
//
// Generated by this command:
//
//	mockgen -package=databasemock -destination=database/databasemock/iterator.go -mock_names=Iterator=Iterator github.com/MetalBlockchain/metalgo/database Iterator
//

// Package databasemock is a generated GoMock package.
package databasemock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Iterator is a mock of Iterator interface.
type Iterator struct {
	ctrl     *gomock.Controller
	recorder *IteratorMockRecorder
}

// IteratorMockRecorder is the mock recorder for Iterator.
type IteratorMockRecorder struct {
	mock *Iterator
}

// NewIterator creates a new mock instance.
func NewIterator(ctrl *gomock.Controller) *Iterator {
	mock := &Iterator{ctrl: ctrl}
	mock.recorder = &IteratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Iterator) EXPECT() *IteratorMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *Iterator) Error() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error")
	ret0, _ := ret[0].(error)
	return ret0
}

// Error indicates an expected call of Error.
func (mr *IteratorMockRecorder) Error() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Iterator)(nil).Error))
}

// Key mocks base method.
func (m *Iterator) Key() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Key")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Key indicates an expected call of Key.
func (mr *IteratorMockRecorder) Key() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Key", reflect.TypeOf((*Iterator)(nil).Key))
}

// Next mocks base method.
func (m *Iterator) Next() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next.
func (mr *IteratorMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*Iterator)(nil).Next))
}

// Release mocks base method.
func (m *Iterator) Release() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Release")
}

// Release indicates an expected call of Release.
func (mr *IteratorMockRecorder) Release() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*Iterator)(nil).Release))
}

// Value mocks base method.
func (m *Iterator) Value() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Value")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Value indicates an expected call of Value.
func (mr *IteratorMockRecorder) Value() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Value", reflect.TypeOf((*Iterator)(nil).Value))
}