package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type MockUrlSaver struct {
	ctrl     *gomock.Controller
	recorder *MockUrlSaverMockRecorder
}

type MockUrlSaverMockRecorder struct {
	mock *MockUrlSaver
}

func NewMockUrlSaver(ctrl *gomock.Controller) *MockUrlSaver {
	mock := &MockUrlSaver{ctrl: ctrl}
	mock.recorder = &MockUrlSaverMockRecorder{mock}
	return mock
}

func (m *MockUrlSaver) EXPECT() *MockUrlSaverMockRecorder {
	return m.recorder
}

func (m *MockUrlSaver) SaveUrl(ctx context.Context, urlToSave, alias string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUrl", ctx, urlToSave, alias)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUrlSaverMockRecorder) SaveUrl(ctx, urlToSave, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUrl", reflect.TypeOf((*MockUrlSaver)(nil).SaveUrl), ctx, urlToSave, alias)
}

type MockUrlRedirecter struct {
	ctrl     *gomock.Controller
	recorder *MockUrlRedirecterMockRecorder
}

type MockUrlRedirecterMockRecorder struct {
	mock *MockUrlRedirecter
}

func NewMockUrlRedirecter(ctrl *gomock.Controller) *MockUrlRedirecter {
	mock := &MockUrlRedirecter{ctrl: ctrl}
	mock.recorder = &MockUrlRedirecterMockRecorder{mock}
	return mock
}

func (m *MockUrlRedirecter) EXPECT() *MockUrlRedirecterMockRecorder {
	return m.recorder
}

func (m *MockUrlRedirecter) GetUrl(ctx context.Context, alias string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUrl", ctx, alias)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUrlRedirecterMockRecorder) GetUrl(ctx, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUrl", reflect.TypeOf((*MockUrlRedirecter)(nil).GetUrl), ctx, alias)
}
