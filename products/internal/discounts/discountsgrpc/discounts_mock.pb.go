// Code generated by MockGen. DO NOT EDIT.
// Source: ../products/internal/discounts/discountsgrpc/discounts_grpc.pb.go

// Package discountsgrpc is a generated GoMock package.
package discountsgrpc

import (
	context "context"
	reflect "reflect"

	domainsgrpc "github.com/dalmarcogd/mobstore/products/internal/domains/domainsgrpc"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockDiscountsClient is a mock of DiscountsClient interface.
type MockDiscountsClient struct {
	ctrl     *gomock.Controller
	recorder *MockDiscountsClientMockRecorder
}

// MockDiscountsClientMockRecorder is the mock recorder for MockDiscountsClient.
type MockDiscountsClientMockRecorder struct {
	mock *MockDiscountsClient
}

// NewMockDiscountsClient creates a new mock instance.
func NewMockDiscountsClient(ctrl *gomock.Controller) *MockDiscountsClient {
	mock := &MockDiscountsClient{ctrl: ctrl}
	mock.recorder = &MockDiscountsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDiscountsClient) EXPECT() *MockDiscountsClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockDiscountsClient) Get(ctx context.Context, in *domainsgrpc.DiscountRequest, opts ...grpc.CallOption) (*domainsgrpc.DiscountResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*domainsgrpc.DiscountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDiscountsClientMockRecorder) Get(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDiscountsClient)(nil).Get), varargs...)
}
