package bundledclusters

import (
	"github.com/pkg/errors"

	apiv1 "github.com/ranger/ranger/pkg/apis/provisioning.cattle.io/v1"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	v3 "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	v1 "github.com/ranger/ranger/tests/framework/clients/ranger/v1"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
)

// Update is a method of BundledCluster that uses provisioning and management clients
// to update related cluster data depending on cluster version.
func (bc *BundledCluster) Update(client *ranger.Client, cUpdates *BundledCluster) (updatedCluster *BundledCluster, err error) {
	updatedCluster = new(BundledCluster)
	updatedCluster.Meta = bc.Meta

	steveclient := client.Steve.SteveType(clusters.ProvisioningSteveResourceType)
	if err != nil {
		return
	}

	if bc.V1 != nil {
		updatedCluster.V1, err = steveclient.Update(bc.V1, cUpdates.V1)
		if err != nil {
			return
		}
	} else if bc.V3 != nil {
		updatedCluster.V3, err = client.Management.Cluster.Update(bc.V3, cUpdates.V3)
		if err != nil {
			return
		}
	}

	return
}

// UpdateKubernetesVersion is a method of BundledCluster that uses update method of BundledCluster
// depending on cluster's provider information. Returns updated BundledCluster and error if any.
func (bc *BundledCluster) UpdateKubernetesVersion(client *ranger.Client, versionToUpgrade *string) (updatedCluster *BundledCluster, err error) {
	if bc.V1 == nil && bc.V3 == nil {
		return nil, errors.Wrapf(err, "cluster %v doesn't contain related data: v1: {%v}, v3: {%v}", bc.Meta.Name, bc.V1, bc.V3)
	}

	bundledv3 := BundledCluster{V3: new(v3.Cluster)}
	bundledv1 := BundledCluster{V1: new(v1.SteveAPIObject)}

	switch bc.Meta.Provider {
	case clusters.KubernetesProviderRKE:
		bundledv3.V3.Name = bc.Meta.Name
		bundledv3.V3.RangerKubernetesEngineConfig = bc.V3.RangerKubernetesEngineConfig
		bundledv3.V3.RangerKubernetesEngineConfig.Version = *versionToUpgrade

		updatedCluster, err = bc.Update(client, &bundledv3)
		if err != nil {
			return updatedCluster, err
		}
	case clusters.KubernetesProviderRKE2:
		if !bc.Meta.IsImported {
			bundledv3.V3.Name = bc.Meta.Name
			bundledv3.V3.Rke2Config = bc.V3.Rke2Config
			bundledv3.V3.Rke2Config.Version = *versionToUpgrade

			updatedCluster, err = bc.Update(client, &bundledv3)
			if err != nil {
				return updatedCluster, err
			}
		} else if bc.Meta.IsImported {
			bundledv1.V1 = bc.V1
			clusterSpec := &apiv1.ClusterSpec{}
			err = v1.ConvertToK8sType(bundledv1.V1.Spec, clusterSpec)
			if err != nil {
				return updatedCluster, err
			}

			clusterSpec.KubernetesVersion = *versionToUpgrade
			bundledv1.V1.Spec = clusterSpec

			updatedCluster, err = bc.Update(client, &bundledv1)
			if err != nil {
				return updatedCluster, err
			}
		}
	case clusters.KubernetesProviderK3S:
		if !bc.Meta.IsImported {
			bundledv3.V3.Name = bc.Meta.Name
			bundledv3.V3.K3sConfig = bc.V3.K3sConfig
			bundledv3.V3.K3sConfig.Version = *versionToUpgrade

			updatedCluster, err = bc.Update(client, &bundledv3)
			if err != nil {
				return updatedCluster, err
			}
		} else if bc.Meta.IsImported {
			bundledv1.V1 = bc.V1
			clusterSpec := &apiv1.ClusterSpec{}
			err = v1.ConvertToK8sType(bundledv1.V1.Spec, clusterSpec)
			if err != nil {
				return updatedCluster, err
			}

			clusterSpec.KubernetesVersion = *versionToUpgrade
			bundledv1.V1.Spec = clusterSpec

			updatedCluster, err = bc.Update(client, &bundledv1)
			if err != nil {
				return updatedCluster, err
			}
		}
	case clusters.KubernetesProviderGKE:
		bundledv3.V3.Name = bc.Meta.Name
		bundledv3.V3.GKEConfig = bc.V3.GKEConfig
		bundledv3.V3.GKEConfig.KubernetesVersion = versionToUpgrade

		updatedCluster, err = bc.Update(client, &bundledv3)
		if err != nil {
			return updatedCluster, err
		}
	case clusters.KubernetesProviderAKS:
		bundledv3.V3.Name = bc.Meta.Name
		bundledv3.V3.AKSConfig = bc.V3.AKSConfig
		bundledv3.V3.AKSConfig.KubernetesVersion = versionToUpgrade

		updatedCluster, err = bc.Update(client, &bundledv3)
		if err != nil {
			return updatedCluster, err
		}
	case clusters.KubernetesProviderEKS:
		bundledv3.V3.Name = bc.Meta.Name
		bundledv3.V3.EKSConfig = bc.V3.EKSConfig
		bundledv3.V3.EKSConfig.KubernetesVersion = versionToUpgrade

		updatedCluster, err = bc.Update(client, &bundledv3)
		if err != nil {
			return updatedCluster, err
		}
	default:
		return nil, errors.Wrap(err, "kubernetes version upgrade failed")
	}

	return
}

// UpdateNodePoolKubernetesVersions is a method of BundledCluster that uses update method of BundledCluster
// depending on cluster's provider information. Returns updated BundledCluster and error if any.
func (bc *BundledCluster) UpdateNodepoolKubernetesVersions(client *ranger.Client, versionToUpgrade *string) (updatedCluster *BundledCluster, err error) {
	if bc.V3 == nil {
		return nil, errors.Wrapf(err, "cluster %v doesn't contain related data", bc.V3)
	}

	cluster := bc

	switch bc.Meta.Provider {
	case clusters.KubernetesProviderGKE:
		for i := range cluster.V3.GKEConfig.NodePools {
			cluster.V3.GKEConfig.NodePools[i].Version = versionToUpgrade
		}

		updatedCluster, err = bc.Update(client, cluster)
		if err != nil {
			return
		}
	case clusters.KubernetesProviderAKS:
		for i := range cluster.V3.AKSConfig.NodePools {
			cluster.V3.AKSConfig.NodePools[i].OrchestratorVersion = versionToUpgrade
		}

		updatedCluster, err = bc.Update(client, cluster)
		if err != nil {
			return
		}
	case clusters.KubernetesProviderEKS:
		for i := range cluster.V3.EKSConfig.NodeGroups {
			cluster.V3.EKSConfig.NodeGroups[i].Version = versionToUpgrade
		}

		updatedCluster, err = bc.Update(client, cluster)
		if err != nil {
			return
		}
	default:
		return nil, errors.Wrap(err, "node pool kubernetes version upgrade failed")
	}

	return
}
