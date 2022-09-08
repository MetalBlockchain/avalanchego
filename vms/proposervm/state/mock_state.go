// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/vms/proposervm/state (interfaces: State)

// Package state is a generated GoMock package.
package state

import (
	reflect "reflect"

	versiondb "github.com/MetalBlockchain/metalgo/database/versiondb"
	ids "github.com/MetalBlockchain/metalgo/ids"
	choices "github.com/MetalBlockchain/metalgo/snow/choices"
	logging "github.com/MetalBlockchain/metalgo/utils/logging"
	block "github.com/MetalBlockchain/metalgo/vms/proposervm/block"
	gomock "github.com/golang/mock/gomock"
)

// MockState is a mock of State interface.
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState.
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance.
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockState) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockStateMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockState)(nil).Commit))
}

// DeleteCheckpoint mocks base method.
func (m *MockState) DeleteCheckpoint() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCheckpoint")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCheckpoint indicates an expected call of DeleteCheckpoint.
func (mr *MockStateMockRecorder) DeleteCheckpoint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCheckpoint", reflect.TypeOf((*MockState)(nil).DeleteCheckpoint))
}

// DeleteLastAccepted mocks base method.
func (m *MockState) DeleteLastAccepted() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastAccepted")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastAccepted indicates an expected call of DeleteLastAccepted.
func (mr *MockStateMockRecorder) DeleteLastAccepted() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastAccepted", reflect.TypeOf((*MockState)(nil).DeleteLastAccepted))
}

// GetBlock mocks base method.
func (m *MockState) GetBlock(arg0 ids.ID) (block.Block, choices.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", arg0)
	ret0, _ := ret[0].(block.Block)
	ret1, _ := ret[1].(choices.Status)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBlock indicates an expected call of GetBlock.
func (mr *MockStateMockRecorder) GetBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockState)(nil).GetBlock), arg0)
}

// GetBlockIDAtHeight mocks base method.
func (m *MockState) GetBlockIDAtHeight(arg0 uint64) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockIDAtHeight", arg0)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockIDAtHeight indicates an expected call of GetBlockIDAtHeight.
func (mr *MockStateMockRecorder) GetBlockIDAtHeight(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockIDAtHeight", reflect.TypeOf((*MockState)(nil).GetBlockIDAtHeight), arg0)
}

// GetCheckpoint mocks base method.
func (m *MockState) GetCheckpoint() (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCheckpoint")
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCheckpoint indicates an expected call of GetCheckpoint.
func (mr *MockStateMockRecorder) GetCheckpoint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCheckpoint", reflect.TypeOf((*MockState)(nil).GetCheckpoint))
}

// GetForkHeight mocks base method.
func (m *MockState) GetForkHeight() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForkHeight")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForkHeight indicates an expected call of GetForkHeight.
func (mr *MockStateMockRecorder) GetForkHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForkHeight", reflect.TypeOf((*MockState)(nil).GetForkHeight))
}

// GetLastAccepted mocks base method.
func (m *MockState) GetLastAccepted() (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAccepted")
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastAccepted indicates an expected call of GetLastAccepted.
func (mr *MockStateMockRecorder) GetLastAccepted() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAccepted", reflect.TypeOf((*MockState)(nil).GetLastAccepted))
}

// HasIndexReset mocks base method.
func (m *MockState) HasIndexReset() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasIndexReset")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasIndexReset indicates an expected call of HasIndexReset.
func (mr *MockStateMockRecorder) HasIndexReset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasIndexReset", reflect.TypeOf((*MockState)(nil).HasIndexReset))
}

// IsIndexEmpty mocks base method.
func (m *MockState) IsIndexEmpty() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsIndexEmpty")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsIndexEmpty indicates an expected call of IsIndexEmpty.
func (mr *MockStateMockRecorder) IsIndexEmpty() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsIndexEmpty", reflect.TypeOf((*MockState)(nil).IsIndexEmpty))
}

// PutBlock mocks base method.
func (m *MockState) PutBlock(arg0 block.Block, arg1 choices.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutBlock", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutBlock indicates an expected call of PutBlock.
func (mr *MockStateMockRecorder) PutBlock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutBlock", reflect.TypeOf((*MockState)(nil).PutBlock), arg0, arg1)
}

// ResetHeightIndex mocks base method.
func (m *MockState) ResetHeightIndex(arg0 logging.Logger, arg1 versiondb.Commitable) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetHeightIndex", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetHeightIndex indicates an expected call of ResetHeightIndex.
func (mr *MockStateMockRecorder) ResetHeightIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetHeightIndex", reflect.TypeOf((*MockState)(nil).ResetHeightIndex), arg0, arg1)
}

// SetBlockIDAtHeight mocks base method.
func (m *MockState) SetBlockIDAtHeight(arg0 uint64, arg1 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBlockIDAtHeight", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetBlockIDAtHeight indicates an expected call of SetBlockIDAtHeight.
func (mr *MockStateMockRecorder) SetBlockIDAtHeight(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBlockIDAtHeight", reflect.TypeOf((*MockState)(nil).SetBlockIDAtHeight), arg0, arg1)
}

// SetCheckpoint mocks base method.
func (m *MockState) SetCheckpoint(arg0 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCheckpoint", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCheckpoint indicates an expected call of SetCheckpoint.
func (mr *MockStateMockRecorder) SetCheckpoint(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCheckpoint", reflect.TypeOf((*MockState)(nil).SetCheckpoint), arg0)
}

// SetForkHeight mocks base method.
func (m *MockState) SetForkHeight(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetForkHeight", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetForkHeight indicates an expected call of SetForkHeight.
func (mr *MockStateMockRecorder) SetForkHeight(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetForkHeight", reflect.TypeOf((*MockState)(nil).SetForkHeight), arg0)
}

// SetIndexHasReset mocks base method.
func (m *MockState) SetIndexHasReset() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetIndexHasReset")
	ret0, _ := ret[0].(error)
	return ret0
}

// SetIndexHasReset indicates an expected call of SetIndexHasReset.
func (mr *MockStateMockRecorder) SetIndexHasReset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetIndexHasReset", reflect.TypeOf((*MockState)(nil).SetIndexHasReset))
}

// SetLastAccepted mocks base method.
func (m *MockState) SetLastAccepted(arg0 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLastAccepted", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLastAccepted indicates an expected call of SetLastAccepted.
func (mr *MockStateMockRecorder) SetLastAccepted(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLastAccepted", reflect.TypeOf((*MockState)(nil).SetLastAccepted), arg0)
}
