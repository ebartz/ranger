// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ranger/wrangler/pkg/generated/controllers/core/v1 (interfaces: SecretClient,SecretCache,ConfigMapCache)

// Package mockcorecontrollers is a generated GoMock package.
package mockcorecontrollers

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/ranger/wrangler/pkg/generated/controllers/core/v1"
	v10 "k8s.io/api/core/v1"
	v11 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
)

// MockSecretClient is a mock of SecretClient interface.
type MockSecretClient struct {
	ctrl     *gomock.Controller
	recorder *MockSecretClientMockRecorder
}

// MockSecretClientMockRecorder is the mock recorder for MockSecretClient.
type MockSecretClientMockRecorder struct {
	mock *MockSecretClient
}

// NewMockSecretClient creates a new mock instance.
func NewMockSecretClient(ctrl *gomock.Controller) *MockSecretClient {
	mock := &MockSecretClient{ctrl: ctrl}
	mock.recorder = &MockSecretClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretClient) EXPECT() *MockSecretClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSecretClient) Create(arg0 *v10.Secret) (*v10.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*v10.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSecretClientMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSecretClient)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockSecretClient) Delete(arg0, arg1 string, arg2 *v11.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSecretClientMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSecretClient)(nil).Delete), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockSecretClient) Get(arg0, arg1 string, arg2 v11.GetOptions) (*v10.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v10.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSecretClientMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSecretClient)(nil).Get), arg0, arg1, arg2)
}

// List mocks base method.
func (m *MockSecretClient) List(arg0 string, arg1 v11.ListOptions) (*v10.SecretList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*v10.SecretList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockSecretClientMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSecretClient)(nil).List), arg0, arg1)
}

// Patch mocks base method.
func (m *MockSecretClient) Patch(arg0, arg1 string, arg2 types.PatchType, arg3 []byte, arg4 ...string) (*v10.Secret, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(*v10.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Patch indicates an expected call of Patch.
func (mr *MockSecretClientMockRecorder) Patch(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockSecretClient)(nil).Patch), varargs...)
}

// Update mocks base method.
func (m *MockSecretClient) Update(arg0 *v10.Secret) (*v10.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(*v10.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockSecretClientMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSecretClient)(nil).Update), arg0)
}

// Watch mocks base method.
func (m *MockSecretClient) Watch(arg0 string, arg1 v11.ListOptions) (watch.Interface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", arg0, arg1)
	ret0, _ := ret[0].(watch.Interface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch.
func (mr *MockSecretClientMockRecorder) Watch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockSecretClient)(nil).Watch), arg0, arg1)
}

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
