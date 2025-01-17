// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ranger/ranger/pkg/generated/controllers/rke.cattle.io/v1 (interfaces: RKEBootstrapClient,RKEBootstrapCache,RKEControlPlaneController,ETCDSnapshotCache)

// Package mockrkecontrollers is a generated GoMock package.
package mockrkecontrollers

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/ranger/ranger/pkg/apis/rke.cattle.io/v1"
	v10 "github.com/ranger/ranger/pkg/generated/controllers/rke.cattle.io/v1"
	generic "github.com/ranger/wrangler/pkg/generic"
	v11 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// MockRKEBootstrapClient is a mock of RKEBootstrapClient interface.
type MockRKEBootstrapClient struct {
	ctrl     *gomock.Controller
	recorder *MockRKEBootstrapClientMockRecorder
}

// MockRKEBootstrapClientMockRecorder is the mock recorder for MockRKEBootstrapClient.
type MockRKEBootstrapClientMockRecorder struct {
	mock *MockRKEBootstrapClient
}

// NewMockRKEBootstrapClient creates a new mock instance.
func NewMockRKEBootstrapClient(ctrl *gomock.Controller) *MockRKEBootstrapClient {
	mock := &MockRKEBootstrapClient{ctrl: ctrl}
	mock.recorder = &MockRKEBootstrapClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRKEBootstrapClient) EXPECT() *MockRKEBootstrapClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRKEBootstrapClient) Create(arg0 *v1.RKEBootstrap) (*v1.RKEBootstrap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*v1.RKEBootstrap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRKEBootstrapClientMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRKEBootstrapClient)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockRKEBootstrapClient) Delete(arg0, arg1 string, arg2 *v11.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRKEBootstrapClientMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRKEBootstrapClient)(nil).Delete), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockRKEBootstrapClient) Get(arg0, arg1 string, arg2 v11.GetOptions) (*v1.RKEBootstrap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.RKEBootstrap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRKEBootstrapClientMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRKEBootstrapClient)(nil).Get), arg0, arg1, arg2)
}

// List mocks base method.
func (m *MockRKEBootstrapClient) List(arg0 string, arg1 v11.ListOptions) (*v1.RKEBootstrapList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*v1.RKEBootstrapList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRKEBootstrapClientMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRKEBootstrapClient)(nil).List), arg0, arg1)
}

// Patch mocks base method.
func (m *MockRKEBootstrapClient) Patch(arg0, arg1 string, arg2 types.PatchType, arg3 []byte, arg4 ...string) (*v1.RKEBootstrap, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(*v1.RKEBootstrap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Patch indicates an expected call of Patch.
func (mr *MockRKEBootstrapClientMockRecorder) Patch(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockRKEBootstrapClient)(nil).Patch), varargs...)
}

// Update mocks base method.
func (m *MockRKEBootstrapClient) Update(arg0 *v1.RKEBootstrap) (*v1.RKEBootstrap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(*v1.RKEBootstrap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRKEBootstrapClientMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRKEBootstrapClient)(nil).Update), arg0)
}

// UpdateStatus mocks base method.
func (m *MockRKEBootstrapClient) UpdateStatus(arg0 *v1.RKEBootstrap) (*v1.RKEBootstrap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", arg0)
	ret0, _ := ret[0].(*v1.RKEBootstrap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockRKEBootstrapClientMockRecorder) UpdateStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockRKEBootstrapClient)(nil).UpdateStatus), arg0)
}

// Watch mocks base method.
func (m *MockRKEBootstrapClient) Watch(arg0 string, arg1 v11.ListOptions) (watch.Interface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", arg0, arg1)
	ret0, _ := ret[0].(watch.Interface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch.
func (mr *MockRKEBootstrapClientMockRecorder) Watch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockRKEBootstrapClient)(nil).Watch), arg0, arg1)
}

// MockRKEBootstrapCache is a mock of RKEBootstrapCache interface.
type MockRKEBootstrapCache struct {
	ctrl     *gomock.Controller
	recorder *MockRKEBootstrapCacheMockRecorder
}

// MockRKEBootstrapCacheMockRecorder is the mock recorder for MockRKEBootstrapCache.
type MockRKEBootstrapCacheMockRecorder struct {
	mock *MockRKEBootstrapCache
}

// NewMockRKEBootstrapCache creates a new mock instance.
func NewMockRKEBootstrapCache(ctrl *gomock.Controller) *MockRKEBootstrapCache {
	mock := &MockRKEBootstrapCache{ctrl: ctrl}
	mock.recorder = &MockRKEBootstrapCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRKEBootstrapCache) EXPECT() *MockRKEBootstrapCacheMockRecorder {
	return m.recorder
}

// AddIndexer mocks base method.
func (m *MockRKEBootstrapCache) AddIndexer(arg0 string, arg1 v10.RKEBootstrapIndexer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddIndexer", arg0, arg1)
}

// AddIndexer indicates an expected call of AddIndexer.
func (mr *MockRKEBootstrapCacheMockRecorder) AddIndexer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddIndexer", reflect.TypeOf((*MockRKEBootstrapCache)(nil).AddIndexer), arg0, arg1)
}

// Get mocks base method.
func (m *MockRKEBootstrapCache) Get(arg0, arg1 string) (*v1.RKEBootstrap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*v1.RKEBootstrap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRKEBootstrapCacheMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRKEBootstrapCache)(nil).Get), arg0, arg1)
}

// GetByIndex mocks base method.
func (m *MockRKEBootstrapCache) GetByIndex(arg0, arg1 string) ([]*v1.RKEBootstrap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIndex", arg0, arg1)
	ret0, _ := ret[0].([]*v1.RKEBootstrap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIndex indicates an expected call of GetByIndex.
func (mr *MockRKEBootstrapCacheMockRecorder) GetByIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIndex", reflect.TypeOf((*MockRKEBootstrapCache)(nil).GetByIndex), arg0, arg1)
}

// List mocks base method.
func (m *MockRKEBootstrapCache) List(arg0 string, arg1 labels.Selector) ([]*v1.RKEBootstrap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*v1.RKEBootstrap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRKEBootstrapCacheMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRKEBootstrapCache)(nil).List), arg0, arg1)
}

// MockRKEControlPlaneController is a mock of RKEControlPlaneController interface.
type MockRKEControlPlaneController struct {
	ctrl     *gomock.Controller
	recorder *MockRKEControlPlaneControllerMockRecorder
}

// MockRKEControlPlaneControllerMockRecorder is the mock recorder for MockRKEControlPlaneController.
type MockRKEControlPlaneControllerMockRecorder struct {
	mock *MockRKEControlPlaneController
}

// NewMockRKEControlPlaneController creates a new mock instance.
func NewMockRKEControlPlaneController(ctrl *gomock.Controller) *MockRKEControlPlaneController {
	mock := &MockRKEControlPlaneController{ctrl: ctrl}
	mock.recorder = &MockRKEControlPlaneControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRKEControlPlaneController) EXPECT() *MockRKEControlPlaneControllerMockRecorder {
	return m.recorder
}

// AddGenericHandler mocks base method.
func (m *MockRKEControlPlaneController) AddGenericHandler(arg0 context.Context, arg1 string, arg2 generic.Handler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddGenericHandler", arg0, arg1, arg2)
}

// AddGenericHandler indicates an expected call of AddGenericHandler.
func (mr *MockRKEControlPlaneControllerMockRecorder) AddGenericHandler(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddGenericHandler", reflect.TypeOf((*MockRKEControlPlaneController)(nil).AddGenericHandler), arg0, arg1, arg2)
}

// AddGenericRemoveHandler mocks base method.
func (m *MockRKEControlPlaneController) AddGenericRemoveHandler(arg0 context.Context, arg1 string, arg2 generic.Handler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddGenericRemoveHandler", arg0, arg1, arg2)
}

// AddGenericRemoveHandler indicates an expected call of AddGenericRemoveHandler.
func (mr *MockRKEControlPlaneControllerMockRecorder) AddGenericRemoveHandler(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddGenericRemoveHandler", reflect.TypeOf((*MockRKEControlPlaneController)(nil).AddGenericRemoveHandler), arg0, arg1, arg2)
}

// Cache mocks base method.
func (m *MockRKEControlPlaneController) Cache() v10.RKEControlPlaneCache {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cache")
	ret0, _ := ret[0].(v10.RKEControlPlaneCache)
	return ret0
}

// Cache indicates an expected call of Cache.
func (mr *MockRKEControlPlaneControllerMockRecorder) Cache() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cache", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Cache))
}

// Create mocks base method.
func (m *MockRKEControlPlaneController) Create(arg0 *v1.RKEControlPlane) (*v1.RKEControlPlane, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*v1.RKEControlPlane)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRKEControlPlaneControllerMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockRKEControlPlaneController) Delete(arg0, arg1 string, arg2 *v11.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRKEControlPlaneControllerMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Delete), arg0, arg1, arg2)
}

// Enqueue mocks base method.
func (m *MockRKEControlPlaneController) Enqueue(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Enqueue", arg0, arg1)
}

// Enqueue indicates an expected call of Enqueue.
func (mr *MockRKEControlPlaneControllerMockRecorder) Enqueue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enqueue", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Enqueue), arg0, arg1)
}

// EnqueueAfter mocks base method.
func (m *MockRKEControlPlaneController) EnqueueAfter(arg0, arg1 string, arg2 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EnqueueAfter", arg0, arg1, arg2)
}

// EnqueueAfter indicates an expected call of EnqueueAfter.
func (mr *MockRKEControlPlaneControllerMockRecorder) EnqueueAfter(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnqueueAfter", reflect.TypeOf((*MockRKEControlPlaneController)(nil).EnqueueAfter), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockRKEControlPlaneController) Get(arg0, arg1 string, arg2 v11.GetOptions) (*v1.RKEControlPlane, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.RKEControlPlane)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRKEControlPlaneControllerMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Get), arg0, arg1, arg2)
}

// GroupVersionKind mocks base method.
func (m *MockRKEControlPlaneController) GroupVersionKind() schema.GroupVersionKind {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GroupVersionKind")
	ret0, _ := ret[0].(schema.GroupVersionKind)
	return ret0
}

// GroupVersionKind indicates an expected call of GroupVersionKind.
func (mr *MockRKEControlPlaneControllerMockRecorder) GroupVersionKind() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GroupVersionKind", reflect.TypeOf((*MockRKEControlPlaneController)(nil).GroupVersionKind))
}

// Informer mocks base method.
func (m *MockRKEControlPlaneController) Informer() cache.SharedIndexInformer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Informer")
	ret0, _ := ret[0].(cache.SharedIndexInformer)
	return ret0
}

// Informer indicates an expected call of Informer.
func (mr *MockRKEControlPlaneControllerMockRecorder) Informer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Informer", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Informer))
}

// List mocks base method.
func (m *MockRKEControlPlaneController) List(arg0 string, arg1 v11.ListOptions) (*v1.RKEControlPlaneList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*v1.RKEControlPlaneList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRKEControlPlaneControllerMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRKEControlPlaneController)(nil).List), arg0, arg1)
}

// OnChange mocks base method.
func (m *MockRKEControlPlaneController) OnChange(arg0 context.Context, arg1 string, arg2 v10.RKEControlPlaneHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnChange", arg0, arg1, arg2)
}

// OnChange indicates an expected call of OnChange.
func (mr *MockRKEControlPlaneControllerMockRecorder) OnChange(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnChange", reflect.TypeOf((*MockRKEControlPlaneController)(nil).OnChange), arg0, arg1, arg2)
}

// OnRemove mocks base method.
func (m *MockRKEControlPlaneController) OnRemove(arg0 context.Context, arg1 string, arg2 v10.RKEControlPlaneHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnRemove", arg0, arg1, arg2)
}

// OnRemove indicates an expected call of OnRemove.
func (mr *MockRKEControlPlaneControllerMockRecorder) OnRemove(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnRemove", reflect.TypeOf((*MockRKEControlPlaneController)(nil).OnRemove), arg0, arg1, arg2)
}

// Patch mocks base method.
func (m *MockRKEControlPlaneController) Patch(arg0, arg1 string, arg2 types.PatchType, arg3 []byte, arg4 ...string) (*v1.RKEControlPlane, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(*v1.RKEControlPlane)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Patch indicates an expected call of Patch.
func (mr *MockRKEControlPlaneControllerMockRecorder) Patch(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Patch), varargs...)
}

// Update mocks base method.
func (m *MockRKEControlPlaneController) Update(arg0 *v1.RKEControlPlane) (*v1.RKEControlPlane, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(*v1.RKEControlPlane)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRKEControlPlaneControllerMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Update), arg0)
}

// UpdateStatus mocks base method.
func (m *MockRKEControlPlaneController) UpdateStatus(arg0 *v1.RKEControlPlane) (*v1.RKEControlPlane, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", arg0)
	ret0, _ := ret[0].(*v1.RKEControlPlane)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockRKEControlPlaneControllerMockRecorder) UpdateStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockRKEControlPlaneController)(nil).UpdateStatus), arg0)
}

// Updater mocks base method.
func (m *MockRKEControlPlaneController) Updater() generic.Updater {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Updater")
	ret0, _ := ret[0].(generic.Updater)
	return ret0
}

// Updater indicates an expected call of Updater.
func (mr *MockRKEControlPlaneControllerMockRecorder) Updater() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Updater", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Updater))
}

// Watch mocks base method.
func (m *MockRKEControlPlaneController) Watch(arg0 string, arg1 v11.ListOptions) (watch.Interface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", arg0, arg1)
	ret0, _ := ret[0].(watch.Interface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch.
func (mr *MockRKEControlPlaneControllerMockRecorder) Watch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockRKEControlPlaneController)(nil).Watch), arg0, arg1)
}

// MockETCDSnapshotCache is a mock of ETCDSnapshotCache interface.
type MockETCDSnapshotCache struct {
	ctrl     *gomock.Controller
	recorder *MockETCDSnapshotCacheMockRecorder
}

// MockETCDSnapshotCacheMockRecorder is the mock recorder for MockETCDSnapshotCache.
type MockETCDSnapshotCacheMockRecorder struct {
	mock *MockETCDSnapshotCache
}

// NewMockETCDSnapshotCache creates a new mock instance.
func NewMockETCDSnapshotCache(ctrl *gomock.Controller) *MockETCDSnapshotCache {
	mock := &MockETCDSnapshotCache{ctrl: ctrl}
	mock.recorder = &MockETCDSnapshotCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockETCDSnapshotCache) EXPECT() *MockETCDSnapshotCacheMockRecorder {
	return m.recorder
}

// AddIndexer mocks base method.
func (m *MockETCDSnapshotCache) AddIndexer(arg0 string, arg1 v10.ETCDSnapshotIndexer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddIndexer", arg0, arg1)
}

// AddIndexer indicates an expected call of AddIndexer.
func (mr *MockETCDSnapshotCacheMockRecorder) AddIndexer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddIndexer", reflect.TypeOf((*MockETCDSnapshotCache)(nil).AddIndexer), arg0, arg1)
}

// Get mocks base method.
func (m *MockETCDSnapshotCache) Get(arg0, arg1 string) (*v1.ETCDSnapshot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*v1.ETCDSnapshot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockETCDSnapshotCacheMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockETCDSnapshotCache)(nil).Get), arg0, arg1)
}

// GetByIndex mocks base method.
func (m *MockETCDSnapshotCache) GetByIndex(arg0, arg1 string) ([]*v1.ETCDSnapshot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIndex", arg0, arg1)
	ret0, _ := ret[0].([]*v1.ETCDSnapshot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIndex indicates an expected call of GetByIndex.
func (mr *MockETCDSnapshotCacheMockRecorder) GetByIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIndex", reflect.TypeOf((*MockETCDSnapshotCache)(nil).GetByIndex), arg0, arg1)
}

// List mocks base method.
func (m *MockETCDSnapshotCache) List(arg0 string, arg1 labels.Selector) ([]*v1.ETCDSnapshot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*v1.ETCDSnapshot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockETCDSnapshotCacheMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockETCDSnapshotCache)(nil).List), arg0, arg1)
}
