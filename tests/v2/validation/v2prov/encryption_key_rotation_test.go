package v2prov

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	apiv1 "github.com/ranger/ranger/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/ranger/ranger/pkg/apis/rke.cattle.io/v1"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	v1 "github.com/ranger/ranger/tests/framework/clients/ranger/v1"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
	"github.com/ranger/ranger/tests/framework/extensions/kubeapi"
	"github.com/ranger/ranger/tests/framework/extensions/kubeapi/secrets"
	namegen "github.com/ranger/ranger/tests/framework/pkg/namegenerator"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/ranger/ranger/tests/framework/pkg/wait"
	"github.com/ranger/ranger/tests/integration/pkg/defaults"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kwait "k8s.io/apimachinery/pkg/util/wait"
)

type V2ProvEncryptionKeyRotationTestSuite struct {
	suite.Suite
	session     *session.Session
	client      *ranger.Client
	clusterName string
}

const (
	namespace    = "fleet-default"
	totalSecrets = 10000
)

var phases = []rkev1.RotateEncryptionKeysPhase{
	rkev1.RotateEncryptionKeysPhasePrepare,
	rkev1.RotateEncryptionKeysPhasePostPrepareRestart,
	rkev1.RotateEncryptionKeysPhaseRotate,
	rkev1.RotateEncryptionKeysPhasePostRotateRestart,
	rkev1.RotateEncryptionKeysPhaseReencrypt,
	rkev1.RotateEncryptionKeysPhasePostReencryptRestart,
	rkev1.RotateEncryptionKeysPhaseDone,
}

func (r *V2ProvEncryptionKeyRotationTestSuite) TearDownSuite() {
	r.session.Cleanup()
}

func (r *V2ProvEncryptionKeyRotationTestSuite) SetupSuite() {
	testSession := session.NewSession()
	r.session = testSession

	client, err := ranger.NewClient("", testSession)
	require.NoError(r.T(), err)

	r.client = client

	r.clusterName = r.client.RangerConfig.ClusterName
}

func rotateEncryptionKeys(t *testing.T, client *ranger.Client, steveID string, generation int64, timeout time.Duration) {
	t.Logf("Applying encryption key rotation generation %d for cluster %s", generation, steveID)

	kubeProvisioningClient, err := client.GetKubeAPIProvisioningClient()
	require.NoError(t, err)

	cluster, err := client.Steve.SteveType(clusters.ProvisioningSteveResourceType).ByID(steveID)
	require.NoError(t, err)

	clusterSpec := &apiv1.ClusterSpec{}
	err = v1.ConvertToK8sType(cluster.Spec, clusterSpec)
	require.NoError(t, err)

	updatedCluster := *cluster

	clusterSpec.RKEConfig.RotateEncryptionKeys = &rkev1.RotateEncryptionKeys{
		Generation: generation,
	}

	updatedCluster.Spec = *clusterSpec

	cluster, err = client.Steve.SteveType(clusters.ProvisioningSteveResourceType).Update(cluster, updatedCluster)
	require.NoError(t, err)

	for _, phase := range phases {
		err = kwait.Poll(10*time.Second, timeout, IsAtLeast(t, client, namespace, cluster.ObjectMeta.Name, phase))
		require.NoError(t, err)
	}

	clusterWait, err := kubeProvisioningClient.Clusters(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + cluster.ObjectMeta.Name,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	})
	require.NoError(t, err)

	err = wait.WatchWait(clusterWait, clusters.IsProvisioningClusterReady)
	require.NoError(t, err)

	t.Logf("Successfully completed encryption key rotation for %s", cluster.ObjectMeta.Name)
}

func createSecretsForCluster(t *testing.T, client *ranger.Client, steveID string, scale int) {
	t.Logf("Creating %d secrets in namespace default for encryption key rotation", scale)

	_, clusterName, found := strings.Cut(steveID, "/")
	require.True(t, found)

	clusterID, err := clusters.GetClusterIDByName(client, clusterName)
	require.NoError(t, err)
	secretResource, err := kubeapi.ResourceForClient(client, clusterID, "default", secrets.SecretGroupVersionResource)
	require.NoError(t, err)

	for i := 0; i < scale; i++ {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: fmt.Sprintf("encryption-key-rotation-test-%d-", i),
			},
			Data: map[string][]byte{
				"key": []byte(namegen.RandStringLower(5)),
			},
		}
		_, err = secrets.CreateSecret(secretResource, secret)
		require.NoError(t, err)
	}
}

func (r *V2ProvEncryptionKeyRotationTestSuite) TestEncryptionKeyRotation() {
	subSession := r.session.NewSession()
	defer subSession.Cleanup()

	id, err := clusters.GetClusterIDByName(r.client, r.clusterName)
	require.NoError(r.T(), err)

	prefix := "encryption-key-rotation-"
	r.Run(prefix+"new-cluster", func() {
		rotateEncryptionKeys(r.T(), r.client, id, 1, 10*time.Minute)
	})

	// create 10k secrets for stress test, takes ~30 minutes
	createSecretsForCluster(r.T(), r.client, id, totalSecrets)

	r.Run(prefix+"stress-test", func() {
		rotateEncryptionKeys(r.T(), r.client, id, 2, 1*time.Hour) // takes ~45 minutes for HA
	})
}

func IsAtLeast(t *testing.T, client *ranger.Client, namespace, name string, phase rkev1.RotateEncryptionKeysPhase) kwait.ConditionFunc {
	return func() (ready bool, err error) {
		kubeRKEClient, err := client.GetKubeAPIRKEClient()
		require.NoError(t, err)

		controlPlane, err := kubeRKEClient.RKEControlPlanes(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		require.NoError(t, err)

		if controlPlane.Status.RotateEncryptionKeysPhase == rkev1.RotateEncryptionKeysPhaseFailed {
			t.Errorf("Encryption key rotation failed waiting to reach %s", phase)
			return ready, fmt.Errorf("encryption key rotation failed")
		}

		desiredPhase := -1
		currentPhase := -1

		for i, v := range phases {
			if v == phase {
				desiredPhase = i
			}
			if v == controlPlane.Status.RotateEncryptionKeysPhase {
				currentPhase = i
			}
			if desiredPhase != -1 && currentPhase != -1 {
				break
			}
		}

		if currentPhase < desiredPhase {
			return false, nil
		}

		t.Logf("Encryption key rotation successfully entered %s", phase)

		return true, nil
	}
}

func TestEncryptionKeyRotation(t *testing.T) {
	suite.Run(t, new(V2ProvEncryptionKeyRotationTestSuite))
}
