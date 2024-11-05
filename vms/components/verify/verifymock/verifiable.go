// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/vms/components/verify (interfaces: Verifiable)
//
// Generated by this command:
//
//	mockgen -package=verifymock -destination=vms/components/verify/verifymock/verifiable.go -mock_names=Verifiable=Verifiable github.com/MetalBlockchain/metalgo/vms/components/verify Verifiable
//

// Package verifymock is a generated GoMock package.
package verifymock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Verifiable is a mock of Verifiable interface.
type Verifiable struct {
	ctrl     *gomock.Controller
	recorder *VerifiableMockRecorder
}

// VerifiableMockRecorder is the mock recorder for Verifiable.
type VerifiableMockRecorder struct {
	mock *Verifiable
}

// NewVerifiable creates a new mock instance.
func NewVerifiable(ctrl *gomock.Controller) *Verifiable {
	mock := &Verifiable{ctrl: ctrl}
	mock.recorder = &VerifiableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Verifiable) EXPECT() *VerifiableMockRecorder {
	return m.recorder
}

// Verify mocks base method.
func (m *Verifiable) Verify() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify")
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *VerifiableMockRecorder) Verify() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*Verifiable)(nil).Verify))
}
