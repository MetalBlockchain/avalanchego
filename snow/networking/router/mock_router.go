// Code generated by MockGen. DO NOT EDIT.
// Source: snow/networking/router/router.go
//
// Generated by this command:
//
//	mockgen -source=snow/networking/router/router.go -destination=snow/networking/router/mock_router.go -package=router -exclude_interfaces=InternalHandler
//

// Package router is a generated GoMock package.
package router

import (
	context "context"
	reflect "reflect"
	time "time"

	ids "github.com/MetalBlockchain/metalgo/ids"
	message "github.com/MetalBlockchain/metalgo/message"
	p2p "github.com/MetalBlockchain/metalgo/proto/pb/p2p"
	handler "github.com/MetalBlockchain/metalgo/snow/networking/handler"
	timeout "github.com/MetalBlockchain/metalgo/snow/networking/timeout"
	logging "github.com/MetalBlockchain/metalgo/utils/logging"
	set "github.com/MetalBlockchain/metalgo/utils/set"
	version "github.com/MetalBlockchain/metalgo/version"
	prometheus "github.com/prometheus/client_golang/prometheus"
	gomock "go.uber.org/mock/gomock"
)

// MockRouter is a mock of Router interface.
type MockRouter struct {
	ctrl     *gomock.Controller
	recorder *MockRouterMockRecorder
}

// MockRouterMockRecorder is the mock recorder for MockRouter.
type MockRouterMockRecorder struct {
	mock *MockRouter
}

// NewMockRouter creates a new mock instance.
func NewMockRouter(ctrl *gomock.Controller) *MockRouter {
	mock := &MockRouter{ctrl: ctrl}
	mock.recorder = &MockRouterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRouter) EXPECT() *MockRouterMockRecorder {
	return m.recorder
}

// AddChain mocks base method.
func (m *MockRouter) AddChain(ctx context.Context, chain handler.Handler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddChain", ctx, chain)
}

// AddChain indicates an expected call of AddChain.
func (mr *MockRouterMockRecorder) AddChain(ctx, chain any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddChain", reflect.TypeOf((*MockRouter)(nil).AddChain), ctx, chain)
}

// Benched mocks base method.
func (m *MockRouter) Benched(chainID ids.ID, validatorID ids.NodeID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Benched", chainID, validatorID)
}

// Benched indicates an expected call of Benched.
func (mr *MockRouterMockRecorder) Benched(chainID, validatorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Benched", reflect.TypeOf((*MockRouter)(nil).Benched), chainID, validatorID)
}

// Connected mocks base method.
func (m *MockRouter) Connected(nodeID ids.NodeID, nodeVersion *version.Application, subnetID ids.ID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Connected", nodeID, nodeVersion, subnetID)
}

// Connected indicates an expected call of Connected.
func (mr *MockRouterMockRecorder) Connected(nodeID, nodeVersion, subnetID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connected", reflect.TypeOf((*MockRouter)(nil).Connected), nodeID, nodeVersion, subnetID)
}

// Disconnected mocks base method.
func (m *MockRouter) Disconnected(nodeID ids.NodeID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Disconnected", nodeID)
}

// Disconnected indicates an expected call of Disconnected.
func (mr *MockRouterMockRecorder) Disconnected(nodeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnected", reflect.TypeOf((*MockRouter)(nil).Disconnected), nodeID)
}

// HandleInbound mocks base method.
func (m *MockRouter) HandleInbound(arg0 context.Context, arg1 message.InboundMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleInbound", arg0, arg1)
}

// HandleInbound indicates an expected call of HandleInbound.
func (mr *MockRouterMockRecorder) HandleInbound(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleInbound", reflect.TypeOf((*MockRouter)(nil).HandleInbound), arg0, arg1)
}

// HealthCheck mocks base method.
func (m *MockRouter) HealthCheck(arg0 context.Context) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HealthCheck indicates an expected call of HealthCheck.
func (mr *MockRouterMockRecorder) HealthCheck(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockRouter)(nil).HealthCheck), arg0)
}

// Initialize mocks base method.
func (m *MockRouter) Initialize(nodeID ids.NodeID, log logging.Logger, timeouts timeout.Manager, shutdownTimeout time.Duration, criticalChains set.Set[ids.ID], sybilProtectionEnabled bool, trackedSubnets set.Set[ids.ID], onFatal func(int), healthConfig HealthConfig, reg prometheus.Registerer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", nodeID, log, timeouts, shutdownTimeout, criticalChains, sybilProtectionEnabled, trackedSubnets, onFatal, healthConfig, reg)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockRouterMockRecorder) Initialize(nodeID, log, timeouts, shutdownTimeout, criticalChains, sybilProtectionEnabled, trackedSubnets, onFatal, healthConfig, reg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockRouter)(nil).Initialize), nodeID, log, timeouts, shutdownTimeout, criticalChains, sybilProtectionEnabled, trackedSubnets, onFatal, healthConfig, reg)
}

// RegisterRequest mocks base method.
func (m *MockRouter) RegisterRequest(ctx context.Context, nodeID ids.NodeID, sourceChainID, destinationChainID ids.ID, requestID uint32, op message.Op, failedMsg message.InboundMessage, engineType p2p.EngineType) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterRequest", ctx, nodeID, sourceChainID, destinationChainID, requestID, op, failedMsg, engineType)
}

// RegisterRequest indicates an expected call of RegisterRequest.
func (mr *MockRouterMockRecorder) RegisterRequest(ctx, nodeID, sourceChainID, destinationChainID, requestID, op, failedMsg, engineType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterRequest", reflect.TypeOf((*MockRouter)(nil).RegisterRequest), ctx, nodeID, sourceChainID, destinationChainID, requestID, op, failedMsg, engineType)
}

// Shutdown mocks base method.
func (m *MockRouter) Shutdown(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown", arg0)
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockRouterMockRecorder) Shutdown(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockRouter)(nil).Shutdown), arg0)
}

// Unbenched mocks base method.
func (m *MockRouter) Unbenched(chainID ids.ID, validatorID ids.NodeID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unbenched", chainID, validatorID)
}

// Unbenched indicates an expected call of Unbenched.
func (mr *MockRouterMockRecorder) Unbenched(chainID, validatorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unbenched", reflect.TypeOf((*MockRouter)(nil).Unbenched), chainID, validatorID)
}
