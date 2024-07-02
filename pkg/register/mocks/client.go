// /*
// Copyright © 2022 - 2024 SUSE LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */
//
//

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rancher/elemental-operator/pkg/register (interfaces: Client)
//
// Generated by this command:
//
//	mockgen -copyright_file=scripts/boilerplate.go.txt -destination=pkg/register/mocks/client.go -package=mocks github.com/rancher/elemental-operator/pkg/register Client
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	v1beta1 "github.com/rancher/elemental-operator/api/v1beta1"
	register "github.com/rancher/elemental-operator/pkg/register"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetNetworkConfig mocks base method.
func (m *MockClient) GetNetworkConfig(arg0 v1beta1.Registration, arg1 []byte, arg2 *register.State) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkConfig", arg0, arg1, arg2)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetworkConfig indicates an expected call of GetNetworkConfig.
func (mr *MockClientMockRecorder) GetNetworkConfig(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkConfig", reflect.TypeOf((*MockClient)(nil).GetNetworkConfig), arg0, arg1, arg2)
}

// Register mocks base method.
func (m *MockClient) Register(arg0 v1beta1.Registration, arg1 []byte, arg2 *register.State) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1, arg2)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockClientMockRecorder) Register(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockClient)(nil).Register), arg0, arg1, arg2)
}
