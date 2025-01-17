package fleetcluster_test

import (
	"testing"
	"time"

	fleetv1api "github.com/ranger/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	mgmt "github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"github.com/ranger/ranger/tests/v2prov/clients"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	waitFor = 5 * time.Minute
	tick    = 2 * time.Second
)

var (
	builtinAffinity = corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
				{
					Weight: 1,
					Preference: corev1.NodeSelectorTerm{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "fleet.cattle.io/agent",
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{"true"},
							},
						},
					},
				},
			},
		},
	}
	linuxAffinity = corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "kubernetes.io/os",
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{"linux"},
							},
						},
					},
				},
			},
		},
	}
	resourceReq = &corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse("1"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
	}
	tolerations = []corev1.Toleration{
		{
			Key:      "key",
			Operator: corev1.TolerationOpEqual,
			Value:    "value",
		},
	}
)

func TestFleetCluster(t *testing.T) {
	require := require.New(t)
	clients, err := clients.New()
	if err != nil {
		t.Fatal(err)
	}
	defer clients.Close()

	cluster := &fleetv1api.Cluster{}
	// wait for fleet local cluster with default affinity
	require.Eventually(func() bool {
		cluster, err = clients.Fleet.Cluster().Get("fleet-local", "local", metav1.GetOptions{})
		return err == nil && cluster != nil && cluster.Status.Summary.Ready > 0
	}, waitFor, tick)
	require.Equal(cluster.Spec.AgentAffinity, &builtinAffinity)
	require.Nil(cluster.Spec.AgentResources)
	require.Empty(cluster.Spec.AgentTolerations)

	// fleet-agent deployment has affinity
	agent, err := clients.Apps.Deployment().Get(cluster.Status.Agent.Namespace, "fleet-agent", metav1.GetOptions{})
	require.NoError(err)
	require.Equal(agent.Spec.Template.Spec.Affinity, &builtinAffinity)
	require.Len(agent.Spec.Template.Spec.Containers, 1)
	require.Empty(agent.Spec.Template.Spec.Containers[0].Resources)
	require.NotEmpty(agent.Spec.Template.Spec.Tolerations) // Fleet has built-in tolerations

	// change settings on management cluster, results should show up in fleet-agent deployment
	mc, err := clients.Mgmt.Cluster().Get("local", metav1.GetOptions{})
	require.NoError(err)

	mc.Spec.FleetAgentDeploymentCustomization = &mgmt.AgentDeploymentCustomization{
		OverrideAffinity:             &linuxAffinity,
		OverrideResourceRequirements: resourceReq,
		AppendTolerations:            tolerations,
	}

	_, err = clients.Mgmt.Cluster().Update(mc)
	require.NoError(err)

	// changes propagate to fleet cluster
	require.Eventually(func() bool {
		cluster, err = clients.Fleet.Cluster().Get("fleet-local", "local", metav1.GetOptions{})
		if err == nil && cluster != nil && cluster.Status.Summary.Ready > 0 {
			assert.Equal(t, cluster.Spec.AgentAffinity, &linuxAffinity)
		}
		return false
	}, waitFor, tick)

	require.Equal(cluster.Spec.AgentAffinity, &linuxAffinity)
	require.Equal(cluster.Spec.AgentResources, resourceReq)
	require.Contains(cluster.Spec.AgentTolerations, tolerations[0])

	// changes are present in deployment
	agent, err = clients.Apps.Deployment().Get(cluster.Status.Agent.Namespace, "fleet-agent", metav1.GetOptions{})
	require.NoError(err)
	require.Equal(agent.Spec.Template.Spec.Affinity, &linuxAffinity)
	require.Len(agent.Spec.Template.Spec.Containers, 1)
	require.Equal(agent.Spec.Template.Spec.Containers[0].Resources.Limits, resourceReq.Limits)
	require.Equal(agent.Spec.Template.Spec.Containers[0].Resources.Requests, resourceReq.Requests)
	require.Contains(agent.Spec.Template.Spec.Tolerations, tolerations[0])
}
