package systemcharts

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	catalog "github.com/ranger/ranger/pkg/apis/catalog.cattle.io/v1"
	"github.com/ranger/ranger/pkg/controllers/dashboard/chart"
	"github.com/ranger/ranger/pkg/controllers/dashboard/chart/fake"
	"github.com/ranger/ranger/pkg/features"
	"github.com/ranger/ranger/pkg/namespace"
	"github.com/ranger/ranger/pkg/settings"
	corev1 "github.com/ranger/wrangler/pkg/generated/controllers/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	errUnimplemented  = fmt.Errorf("unimplemented")
	errNotFound       = fmt.Errorf("not found")
	priorityClassName = "ranger-critical"
	operatorNamespace = "ranger-operator-system"
)

// Test_ChartInstallation test that all expected charts are installed or uninstalled with expected configuration
func Test_ChartInstallation(t *testing.T) {
	repo := &catalog.ClusterRepo{
		ObjectMeta: metav1.ObjectMeta{
			Name: repoName,
		},
	}
	h := &handler{
		chartsConfig: chart.RangerConfigGetter{ConfigCache: &mockCache{}},
	}
	tests := []struct {
		name             string
		setup            func(*gomock.Controller) chart.Manager
		registryOverride string
		wantErr          bool
	}{
		{
			name: "normal installation",
			setup: func(ctrl *gomock.Controller) chart.Manager {
				settings.ConfigMapName.Set("pass")
				settings.RangerWebhookVersion.Set("2.0.0")
				manager := fake.NewMockManager(ctrl)
				expectedValues := map[string]interface{}{
					"priorityClassName": priorityClassName,
					"capi": map[string]interface{}{
						"enabled": features.EmbeddedClusterAPI.Enabled(),
					},
					"mcm": map[string]interface{}{
						"enabled": features.MCM.Enabled(),
					},
					"global": map[string]interface{}{
						"cattle": map[string]interface{}{
							"systemDefaultRegistry": settings.SystemDefaultRegistry.Get(),
						},
					},
				}
				var b bool
				manager.EXPECT().Ensure(
					namespace.System,
					"ranger-webhook",
					"",
					"2.0.0",
					expectedValues,
					gomock.AssignableToTypeOf(b),
					"",
				).Return(nil)

				manager.EXPECT().Uninstall(operatorNamespace, "ranger-operator").Return(nil)
				return manager
			},
		},
		{
			name: "installation without webhook priority class",
			setup: func(ctrl *gomock.Controller) chart.Manager {
				settings.ConfigMapName.Set("fail")
				settings.RangerWebhookVersion.Set("2.0.0")
				manager := fake.NewMockManager(ctrl)
				expectedValues := map[string]interface{}{
					"capi": map[string]interface{}{
						"enabled": features.EmbeddedClusterAPI.Enabled(),
					},
					"mcm": map[string]interface{}{
						"enabled": features.MCM.Enabled(),
					},
					"global": map[string]interface{}{
						"cattle": map[string]interface{}{
							"systemDefaultRegistry": settings.SystemDefaultRegistry.Get(),
						},
					},
				}
				var b bool
				manager.EXPECT().Ensure(
					namespace.System,
					"ranger-webhook",
					"",
					"2.0.0",
					expectedValues,
					gomock.AssignableToTypeOf(b),
					"",
				).Return(nil)

				manager.EXPECT().Uninstall(operatorNamespace, "ranger-operator").Return(nil)
				return manager

			},
		},
		{
			name: "installation with image override",
			setup: func(ctrl *gomock.Controller) chart.Manager {
				settings.ConfigMapName.Set("fail")
				settings.RangerWebhookVersion.Set("2.0.1")
				manager := fake.NewMockManager(ctrl)
				expectedValues := map[string]interface{}{
					"capi": map[string]interface{}{
						"enabled": features.EmbeddedClusterAPI.Enabled(),
					},
					"mcm": map[string]interface{}{
						"enabled": features.MCM.Enabled(),
					},
					"global": map[string]interface{}{
						"cattle": map[string]interface{}{
							"systemDefaultRegistry": "",
						},
					},
					"image": map[string]interface{}{
						"repository": "ranger-test.io/ranger/ranger-webhook",
					},
				}
				var b bool
				manager.EXPECT().Ensure(
					namespace.System,
					"ranger-webhook",
					"",
					"2.0.1",
					expectedValues,
					gomock.AssignableToTypeOf(b),
					"ranger-test.io/"+settings.ShellImage.Get(),
				).Return(nil)

				manager.EXPECT().Uninstall(operatorNamespace, "ranger-operator").Return(nil)
				return manager
			},
			registryOverride: "ranger-test.io",
		},
		{
			name: "installation with min version override",
			setup: func(ctrl *gomock.Controller) chart.Manager {
				settings.ConfigMapName.Set("fail")
				settings.RangerWebhookMinVersion.Set("2.0.1")
				settings.RangerWebhookVersion.Set("2.0.4")
				manager := fake.NewMockManager(ctrl)
				expectedValues := map[string]interface{}{
					"capi": map[string]interface{}{
						"enabled": features.EmbeddedClusterAPI.Enabled(),
					},
					"mcm": map[string]interface{}{
						"enabled": features.MCM.Enabled(),
					},
					"global": map[string]interface{}{
						"cattle": map[string]interface{}{
							"systemDefaultRegistry": "",
						},
					},
					"image": map[string]interface{}{
						"repository": "ranger-test.io/ranger/ranger-webhook",
					},
				}
				var b bool
				manager.EXPECT().Ensure(
					namespace.System,
					"ranger-webhook",
					"2.0.1",
					"",
					expectedValues,
					gomock.AssignableToTypeOf(b),
					"ranger-test.io/"+settings.ShellImage.Get(),
				).Return(nil)

				manager.EXPECT().Uninstall(operatorNamespace, "ranger-operator").Return(nil)
				return manager
			},
			registryOverride: "ranger-test.io",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			namespaceCtrl := NewMockNamespaceController(ctrl)
			namespaceCtrl.EXPECT().Delete(operatorNamespace, nil).Return(nil)
			h.manager = tt.setup(ctrl)
			h.namespaces = namespaceCtrl
			h.registryOverride = tt.registryOverride
			_, err := h.onRepo("", repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.onRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type mockCache struct {
	Maps []*v1.ConfigMap
}

func (m *mockCache) Get(namespace string, name string) (*v1.ConfigMap, error) {
	if name != "pass" {
		return nil, errNotFound
	}
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{"priorityClassName": priorityClassName},
	}, nil
}

func (m *mockCache) List(namespace string, selector labels.Selector) ([]*v1.ConfigMap, error) {
	return nil, errUnimplemented
}

func (m *mockCache) AddIndexer(indexName string, indexer corev1.ConfigMapIndexer) {}

func (m *mockCache) GetByIndex(indexName, key string) ([]*v1.ConfigMap, error) {
	return nil, errUnimplemented
}
