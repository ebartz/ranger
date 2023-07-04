// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ranger/wrangler/pkg/generated/controllers/core/v1 (interfaces: SecretCache,ConfigMapCache)

// Package hostedcluster is a generated GoMock package.
package hostedcluster

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/ranger/wrangler/pkg/generated/controllers/core/v1"
	v10 "k8s.io/api/core/v1"
	labels "k8s.io/apimachinery/pkg/labels"
)

// MockSecretCache is a mock of SecretCache interface.
type MockSecretCache struct {
	ctrl     *gomock.Controller
	recorder *MockSecretCacheMockRecorder
}

// MockSecretCacheMockRecorder is the mock recorder for MockSecretCache.
type MockSecretCacheMockRecorder struct {
	mock *MockSecretCache
}

// NewMockSecretCache creates a new mock instance.
func NewMockSecretCache(ctrl *gomock.Controller) *MockSecretCache {
	mock := &MockSecretCache{ctrl: ctrl}
	mock.recorder = &MockSecretCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretCache) EXPECT() *MockSecretCacheMockRecorder {
	return m.recorder
}

// AddIndexer mocks base method.
func (m *MockSecretCache) AddIndexer(arg0 string, arg1 v1.SecretIndexer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddIndexer", arg0, arg1)
}

// AddIndexer indicates an expected call of AddIndexer.
func (mr *MockSecretCacheMockRecorder) AddIndexer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddIndexer", reflect.TypeOf((*MockSecretCache)(nil).AddIndexer), arg0, arg1)
}

// Get mocks base method.
func (m *MockSecretCache) Get(arg0, arg1 string) (*v10.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*v10.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSecretCacheMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSecretCache)(nil).Get), arg0, arg1)
}

// GetByIndex mocks base method.
func (m *MockSecretCache) GetByIndex(arg0, arg1 string) ([]*v10.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIndex", arg0, arg1)
	ret0, _ := ret[0].([]*v10.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIndex indicates an expected call of GetByIndex.
func (mr *MockSecretCacheMockRecorder) GetByIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIndex", reflect.TypeOf((*MockSecretCache)(nil).GetByIndex), arg0, arg1)
}

// List mocks base method.
func (m *MockSecretCache) List(arg0 string, arg1 labels.Selector) ([]*v10.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*v10.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockSecretCacheMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSecretCache)(nil).List), arg0, arg1)
}

// MockConfigMapCache is a mock of ConfigMapCache interface.
type MockConfigMapCache struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMapCacheMockRecorder
}

// MockConfigMapCacheMockRecorder is the mock recorder for MockConfigMapCache.
type MockConfigMapCacheMockRecorder struct {
	mock *MockConfigMapCache
}

// NewMockConfigMapCache creates a new mock instance.
func NewMockConfigMapCache(ctrl *gomock.Controller) *MockConfigMapCache {
	mock := &MockConfigMapCache{ctrl: ctrl}
	mock.recorder = &MockConfigMapCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigMapCache) EXPECT() *MockConfigMapCacheMockRecorder {
	return m.recorder
}

// AddIndexer mocks base method.
func (m *MockConfigMapCache) AddIndexer(arg0 string, arg1 v1.ConfigMapIndexer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddIndexer", arg0, arg1)
}

// AddIndexer indicates an expected call of AddIndexer.
func (mr *MockConfigMapCacheMockRecorder) AddIndexer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddIndexer", reflect.TypeOf((*MockConfigMapCache)(nil).AddIndexer), arg0, arg1)
}

// Get mocks base method.
func (m *MockConfigMapCache) Get(arg0, arg1 string) (*v10.ConfigMap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*v10.ConfigMap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockConfigMapCacheMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConfigMapCache)(nil).Get), arg0, arg1)
}

// GetByIndex mocks base method.
func (m *MockConfigMapCache) GetByIndex(arg0, arg1 string) ([]*v10.ConfigMap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIndex", arg0, arg1)
	ret0, _ := ret[0].([]*v10.ConfigMap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIndex indicates an expected call of GetByIndex.
func (mr *MockConfigMapCacheMockRecorder) GetByIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIndex", reflect.TypeOf((*MockConfigMapCache)(nil).GetByIndex), arg0, arg1)
}

// List mocks base method.
func (m *MockConfigMapCache) List(arg0 string, arg1 labels.Selector) ([]*v10.ConfigMap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*v10.ConfigMap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockConfigMapCacheMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockConfigMapCache)(nil).List), arg0, arg1)
}
