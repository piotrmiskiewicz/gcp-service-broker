// Copyright 2017 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Automatically generated by MockGen. DO NOT EDIT!
// Source: google.golang.org/genproto/googleapis/devtools/cloudprofiler/v2 (interfaces: ProfilerServiceClient)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	context "golang.org/x/net/context"
	v2 "google.golang.org/genproto/googleapis/devtools/cloudprofiler/v2"
	grpc "google.golang.org/grpc"
)

// Mock of ProfilerServiceClient interface
type MockProfilerServiceClient struct {
	ctrl     *gomock.Controller
	recorder *_MockProfilerServiceClientRecorder
}

// Recorder for MockProfilerServiceClient (not exported)
type _MockProfilerServiceClientRecorder struct {
	mock *MockProfilerServiceClient
}

func NewMockProfilerServiceClient(ctrl *gomock.Controller) *MockProfilerServiceClient {
	mock := &MockProfilerServiceClient{ctrl: ctrl}
	mock.recorder = &_MockProfilerServiceClientRecorder{mock}
	return mock
}

func (_m *MockProfilerServiceClient) EXPECT() *_MockProfilerServiceClientRecorder {
	return _m.recorder
}

func (_m *MockProfilerServiceClient) CreateProfile(_param0 context.Context, _param1 *v2.CreateProfileRequest, _param2 ...grpc.CallOption) (*v2.Profile, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "CreateProfile", _s...)
	ret0, _ := ret[0].(*v2.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockProfilerServiceClientRecorder) CreateProfile(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateProfile", _s...)
}

func (_m *MockProfilerServiceClient) UpdateProfile(_param0 context.Context, _param1 *v2.UpdateProfileRequest, _param2 ...grpc.CallOption) (*v2.Profile, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "UpdateProfile", _s...)
	ret0, _ := ret[0].(*v2.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockProfilerServiceClientRecorder) UpdateProfile(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateProfile", _s...)
}