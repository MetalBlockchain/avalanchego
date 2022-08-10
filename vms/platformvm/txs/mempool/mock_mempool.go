// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/avalanchego/vms/platformvm/txs/mempool (interfaces: Mempool)

// Package mempool is a generated GoMock package.
package mempool

import (
	ids "github.com/MetalBlockchain/avalanchego/ids"
	txs "github.com/MetalBlockchain/avalanchego/vms/platformvm/txs"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMempool is a mock of Mempool interface
type MockMempool struct {
	ctrl     *gomock.Controller
	recorder *MockMempoolMockRecorder
}

// MockMempoolMockRecorder is the mock recorder for MockMempool
type MockMempoolMockRecorder struct {
	mock *MockMempool
}

// NewMockMempool creates a new mock instance
func NewMockMempool(ctrl *gomock.Controller) *MockMempool {
	mock := &MockMempool{ctrl: ctrl}
	mock.recorder = &MockMempoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMempool) EXPECT() *MockMempoolMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockMempool) Add(arg0 *txs.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockMempoolMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockMempool)(nil).Add), arg0)
}

// AddDecisionTx mocks base method
func (m *MockMempool) AddDecisionTx(arg0 *txs.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddDecisionTx", arg0)
}

// AddDecisionTx indicates an expected call of AddDecisionTx
func (mr *MockMempoolMockRecorder) AddDecisionTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecisionTx", reflect.TypeOf((*MockMempool)(nil).AddDecisionTx), arg0)
}

// AddProposalTx mocks base method
func (m *MockMempool) AddProposalTx(arg0 *txs.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddProposalTx", arg0)
}

// AddProposalTx indicates an expected call of AddProposalTx
func (mr *MockMempoolMockRecorder) AddProposalTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProposalTx", reflect.TypeOf((*MockMempool)(nil).AddProposalTx), arg0)
}

// DisableAdding mocks base method
func (m *MockMempool) DisableAdding() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DisableAdding")
}

// DisableAdding indicates an expected call of DisableAdding
func (mr *MockMempoolMockRecorder) DisableAdding() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableAdding", reflect.TypeOf((*MockMempool)(nil).DisableAdding))
}

// EnableAdding mocks base method
func (m *MockMempool) EnableAdding() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EnableAdding")
}

// EnableAdding indicates an expected call of EnableAdding
func (mr *MockMempoolMockRecorder) EnableAdding() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableAdding", reflect.TypeOf((*MockMempool)(nil).EnableAdding))
}

// Get mocks base method
func (m *MockMempool) Get(arg0 ids.ID) *txs.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*txs.Tx)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockMempoolMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMempool)(nil).Get), arg0)
}

// GetDropReason mocks base method
func (m *MockMempool) GetDropReason(arg0 ids.ID) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDropReason", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetDropReason indicates an expected call of GetDropReason
func (mr *MockMempoolMockRecorder) GetDropReason(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDropReason", reflect.TypeOf((*MockMempool)(nil).GetDropReason), arg0)
}

// Has mocks base method
func (m *MockMempool) Has(arg0 ids.ID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Has indicates an expected call of Has
func (mr *MockMempoolMockRecorder) Has(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockMempool)(nil).Has), arg0)
}

// HasDecisionTxs mocks base method
func (m *MockMempool) HasDecisionTxs() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasDecisionTxs")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasDecisionTxs indicates an expected call of HasDecisionTxs
func (mr *MockMempoolMockRecorder) HasDecisionTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasDecisionTxs", reflect.TypeOf((*MockMempool)(nil).HasDecisionTxs))
}

// HasProposalTx mocks base method
func (m *MockMempool) HasProposalTx() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasProposalTx")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasProposalTx indicates an expected call of HasProposalTx
func (mr *MockMempoolMockRecorder) HasProposalTx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasProposalTx", reflect.TypeOf((*MockMempool)(nil).HasProposalTx))
}

// MarkDropped mocks base method
func (m *MockMempool) MarkDropped(arg0 ids.ID, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkDropped", arg0, arg1)
}

// MarkDropped indicates an expected call of MarkDropped
func (mr *MockMempoolMockRecorder) MarkDropped(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkDropped", reflect.TypeOf((*MockMempool)(nil).MarkDropped), arg0, arg1)
}

// PopDecisionTxs mocks base method
func (m *MockMempool) PopDecisionTxs(arg0 int) []*txs.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PopDecisionTxs", arg0)
	ret0, _ := ret[0].([]*txs.Tx)
	return ret0
}

// PopDecisionTxs indicates an expected call of PopDecisionTxs
func (mr *MockMempoolMockRecorder) PopDecisionTxs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopDecisionTxs", reflect.TypeOf((*MockMempool)(nil).PopDecisionTxs), arg0)
}

// PopProposalTx mocks base method
func (m *MockMempool) PopProposalTx() *txs.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PopProposalTx")
	ret0, _ := ret[0].(*txs.Tx)
	return ret0
}

// PopProposalTx indicates an expected call of PopProposalTx
func (mr *MockMempoolMockRecorder) PopProposalTx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopProposalTx", reflect.TypeOf((*MockMempool)(nil).PopProposalTx))
}

// RemoveDecisionTxs mocks base method
func (m *MockMempool) RemoveDecisionTxs(arg0 []*txs.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveDecisionTxs", arg0)
}

// RemoveDecisionTxs indicates an expected call of RemoveDecisionTxs
func (mr *MockMempoolMockRecorder) RemoveDecisionTxs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDecisionTxs", reflect.TypeOf((*MockMempool)(nil).RemoveDecisionTxs), arg0)
}

// RemoveProposalTx mocks base method
func (m *MockMempool) RemoveProposalTx(arg0 *txs.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveProposalTx", arg0)
}

// RemoveProposalTx indicates an expected call of RemoveProposalTx
func (mr *MockMempoolMockRecorder) RemoveProposalTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveProposalTx", reflect.TypeOf((*MockMempool)(nil).RemoveProposalTx), arg0)
}
