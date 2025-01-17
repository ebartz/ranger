package assemblers

import (
	"encoding/json"
	"strings"

	apimgmtv3 "github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	v1 "github.com/ranger/ranger/pkg/generated/norman/core/v1"
	"github.com/ranger/ranger/pkg/namespace"

	rketypes "github.com/ranger/rke/types"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	configv1 "k8s.io/apiserver/pkg/apis/config/v1"
	"k8s.io/kubernetes/pkg/credentialprovider"
)

const (
	ClusterType                 = "cluster"
	ClusterTemplateRevisionType = "cluster template revision"
	SecretNamespace             = namespace.GlobalNamespace
	SecretKey                   = "credential"
)

type Assembler func(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error)

// AssemblePrivateRegistryCredential looks up the registry Secret and inserts the keys into the PrivateRegistries list on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssemblePrivateRegistryCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || len(spec.RangerKubernetesEngineConfig.PrivateRegistries) == 0 {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.PrivateRegistries[0].Password != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil
	}
	registrySecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}

	dockerCfg := credentialprovider.DockerConfigJSON{}
	err = json.Unmarshal(registrySecret.Data[corev1.DockerConfigJsonKey], &dockerCfg)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	for i, privateRegistry := range specCopy.RangerKubernetesEngineConfig.PrivateRegistries {
		if reg, ok := dockerCfg.Auths[privateRegistry.URL]; ok {
			specCopy.RangerKubernetesEngineConfig.PrivateRegistries[i].User = reg.Username
			specCopy.RangerKubernetesEngineConfig.PrivateRegistries[i].Password = reg.Password
		}
	}
	return *specCopy, nil
}

// AssembleS3Credential looks up the S3 backup config Secret and inserts the keys into the S3BackupConfig on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleS3Credential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.Services.Etcd.BackupConfig == nil || spec.RangerKubernetesEngineConfig.Services.Etcd.BackupConfig.S3BackupConfig == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.Services.Etcd.BackupConfig.S3BackupConfig.SecretKey != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil
	}
	s3Cred, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	specCopy.RangerKubernetesEngineConfig.Services.Etcd.BackupConfig.S3BackupConfig.SecretKey = string(s3Cred.Data[SecretKey])
	return *specCopy, nil
}

// AssembleWeaveCredential looks up the weave Secret and inserts the keys into the network provider config on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleWeaveCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.Network.WeaveNetworkProvider == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.Network.WeaveNetworkProvider.Password != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil
	}
	weaveSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	specCopy.RangerKubernetesEngineConfig.Network.WeaveNetworkProvider.Password = string(weaveSecret.Data[SecretKey])
	return *specCopy, nil
}

// AssembleVsphereGlobalCredential looks up the vsphere global Secret and inserts the keys into the cloud provider config on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleVsphereGlobalCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.CloudProvider.VsphereCloudProvider == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.CloudProvider.VsphereCloudProvider.Global.Password != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil
	}
	vsphereSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	specCopy.RangerKubernetesEngineConfig.CloudProvider.VsphereCloudProvider.Global.Password = string(vsphereSecret.Data[SecretKey])
	return *specCopy, nil
}

// AssembleVsphereVirtualCenterCredential looks up the vsphere virtualcenter Secret and inserts the keys into the cloud provider config on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleVsphereVirtualCenterCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.CloudProvider.VsphereCloudProvider == nil {
		return spec, nil
	}
	if secretRef == "" {
		for _, v := range spec.RangerKubernetesEngineConfig.CloudProvider.VsphereCloudProvider.VirtualCenter {
			if v.Password != "" {
				logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
				break
			}
		}
		return spec, nil

	}
	vcenterSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	for k, v := range vcenterSecret.Data {
		vCenter := specCopy.RangerKubernetesEngineConfig.CloudProvider.VsphereCloudProvider.VirtualCenter[k]
		vCenter.Password = string(v)
		specCopy.RangerKubernetesEngineConfig.CloudProvider.VsphereCloudProvider.VirtualCenter[k] = vCenter
	}
	return *specCopy, nil
}

// AssembleOpenStackCredential looks up the OpenStack Secret and inserts the keys into the cloud provider config on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleOpenStackCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.CloudProvider.OpenstackCloudProvider == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.CloudProvider.OpenstackCloudProvider.Global.Password != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil

	}
	openStackSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	specCopy.RangerKubernetesEngineConfig.CloudProvider.OpenstackCloudProvider.Global.Password = string(openStackSecret.Data[SecretKey])
	return *specCopy, nil
}

// AssembleAADClientSecretCredential looks up the AAD client secret Secret and inserts the keys into the cloud provider config on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleAADClientSecretCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.CloudProvider.AzureCloudProvider == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.CloudProvider.AzureCloudProvider.AADClientSecret != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil

	}
	aadClientSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	specCopy.RangerKubernetesEngineConfig.CloudProvider.AzureCloudProvider.AADClientSecret = string(aadClientSecret.Data[SecretKey])
	return *specCopy, nil
}

// AssembleAADCertCredential looks up the AAD client cert password Secret and inserts the keys into the cloud provider config on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleAADCertCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.CloudProvider.AzureCloudProvider == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.CloudProvider.AzureCloudProvider.AADClientCertPassword != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil
	}
	aadCertSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	specCopy.RangerKubernetesEngineConfig.CloudProvider.AzureCloudProvider.AADClientCertPassword = string(aadCertSecret.Data[SecretKey])
	return *specCopy, nil
}

// AssembleACIAPICUserKeyCredential looks up the aci apic user key Secret and inserts the keys into the AciNetworkProvider on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleACIAPICUserKeyCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.Network.AciNetworkProvider == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.Network.AciNetworkProvider.ApicUserKey != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil

	}
	aciUserKeySecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	specCopy.RangerKubernetesEngineConfig.Network.AciNetworkProvider.ApicUserKey = string(aciUserKeySecret.Data[SecretKey])
	return *specCopy, nil
}

// AssembleACITokenCredential looks up the aci token Secret and inserts the keys into the AciNetworkProvider on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleACITokenCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.Network.AciNetworkProvider == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.Network.AciNetworkProvider.Token != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil

	}
	aciTokenSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	specCopy.RangerKubernetesEngineConfig.Network.AciNetworkProvider.Token = string(aciTokenSecret.Data[SecretKey])
	return *specCopy, nil
}

// AssembleACIKafkaClientKeyCredential looks up the aci kafka client key Secret and inserts the keys into the AciNetworkProvider on the Cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleACIKafkaClientKeyCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil || spec.RangerKubernetesEngineConfig.Network.AciNetworkProvider == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.Network.AciNetworkProvider.KafkaClientKey != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil

	}
	aciKafkaClientKeySecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	specCopy := spec.DeepCopy()
	specCopy.RangerKubernetesEngineConfig.Network.AciNetworkProvider.KafkaClientKey = string(aciKafkaClientKeySecret.Data[SecretKey])
	return *specCopy, nil
}

// AssembleSecretsEncryptionProvidersSecretCredential looks up the rke KubeAPI secrets encryption configuration and
// inserts it back into the cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleSecretsEncryptionProvidersSecretCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil ||
		spec.RangerKubernetesEngineConfig.Services.KubeAPI.SecretsEncryptionConfig == nil ||
		spec.RangerKubernetesEngineConfig.Services.KubeAPI.SecretsEncryptionConfig.CustomConfig == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.Services.KubeAPI.SecretsEncryptionConfig.CustomConfig.Resources != nil {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil
	}
	secretsEncryptionProvidersSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	var resource []configv1.ResourceConfiguration
	err = json.Unmarshal(secretsEncryptionProvidersSecret.Data[SecretKey], &resource)
	if err != nil {
		return spec, err
	}
	spec.RangerKubernetesEngineConfig.Services.KubeAPI.SecretsEncryptionConfig.CustomConfig.Resources = resource
	return spec, nil
}

// AssembleBastionHostSSHKeyCredential looks up bastion host ssh key and inserts it back into the cluster spec.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleBastionHostSSHKeyCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil {
		return spec, nil
	}
	if secretRef == "" {
		if spec.RangerKubernetesEngineConfig.BastionHost.SSHKey != "" {
			logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
		}
		return spec, nil
	}
	bastionHostSSHKeySecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	spec.RangerKubernetesEngineConfig.BastionHost.SSHKey = string(bastionHostSSHKeySecret.Data[SecretKey])
	return spec, nil
}

// AssembleKubeletExtraEnvCredential looks up the AWS_SECRET_ACCESS_KEY extraEnv for the kubelet if it exists.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssembleKubeletExtraEnvCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil {
		return spec, nil
	}
	if secretRef == "" {
		for _, e := range spec.RangerKubernetesEngineConfig.Services.Kubelet.ExtraEnv {
			if strings.Contains(e, "AWS_SECRET_ACCESS_KEY") {
				logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
				break
			}
		}
		return spec, nil
	}
	kubeletExtraEnvSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}
	env := "AWS_SECRET_ACCESS_KEY=" + string(kubeletExtraEnvSecret.Data[SecretKey])
	spec.RangerKubernetesEngineConfig.Services.Kubelet.ExtraEnv = append(spec.RangerKubernetesEngineConfig.Services.Kubelet.ExtraEnv, env)
	return spec, nil
}

// AssemblePrivateRegistryECRCredential looks up Private Registry's ECR credential auth info, if it exists.
// It returns a new copy of the spec without modifying the original. The Cluster is never updated.
func AssemblePrivateRegistryECRCredential(secretRef, objType, objName string, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	if spec.RangerKubernetesEngineConfig == nil ||
		len(spec.RangerKubernetesEngineConfig.PrivateRegistries) == 0 {
		return spec, nil
	}
	if secretRef == "" {
		for _, r := range spec.RangerKubernetesEngineConfig.PrivateRegistries {
			if ecr := r.ECRCredentialPlugin; ecr != nil && (ecr.AwsSecretAccessKey != "" || ecr.AwsSessionToken != "") {
				logrus.Warnf("[secretmigrator] secrets for %s %s are not finished migrating", objType, objName)
				break
			}
		}
		return spec, nil
	}
	privateRegistryECRSecret, err := secretLister.Get(SecretNamespace, secretRef)
	if err != nil {
		return spec, err
	}

	data, ok := privateRegistryECRSecret.Data[SecretKey]
	if !ok {
		return spec, nil
	}
	var registries map[string]string
	err = json.Unmarshal(data, &registries)
	if err != nil {
		return spec, err
	}

	for i, reg := range spec.RangerKubernetesEngineConfig.PrivateRegistries {
		if ecrData, ok := registries[reg.URL]; ok {
			var ecr rketypes.ECRCredentialPlugin
			err := json.Unmarshal([]byte(ecrData), &ecr)
			if err != nil {
				return spec, err
			}
			spec.RangerKubernetesEngineConfig.PrivateRegistries[i].ECRCredentialPlugin.AwsSecretAccessKey = ecr.AwsSecretAccessKey
			spec.RangerKubernetesEngineConfig.PrivateRegistries[i].ECRCredentialPlugin.AwsSessionToken = ecr.AwsSessionToken
		}
	}

	return spec, nil
}

// AssembleRKEConfigSpec is a wrapper assembler for assembling configs on Clusters.
func AssembleRKEConfigSpec(cluster *apimgmtv3.Cluster, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	spec, err := AssembleS3Credential(cluster.GetSecret("S3CredentialSecret"), ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssemblePrivateRegistryCredential(cluster.GetSecret("PrivateRegistrySecret"), ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleWeaveCredential(cluster.GetSecret("WeavePasswordSecret"), ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleVsphereGlobalCredential(cluster.GetSecret("VsphereSecret"), ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleVsphereVirtualCenterCredential(cluster.GetSecret("VirtualCenterSecret"), ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleOpenStackCredential(cluster.GetSecret("OpenStackSecret"), ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleAADClientSecretCredential(cluster.GetSecret("AADClientSecret"), ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleAADCertCredential(cluster.GetSecret("AADClientCertSecret"), ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleACIAPICUserKeyCredential(cluster.Spec.ClusterSecrets.ACIAPICUserKeySecret, ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleACITokenCredential(cluster.Spec.ClusterSecrets.ACITokenSecret, ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleACIKafkaClientKeyCredential(cluster.Spec.ClusterSecrets.ACIKafkaClientKeySecret, ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleSecretsEncryptionProvidersSecretCredential(cluster.Spec.ClusterSecrets.SecretsEncryptionProvidersSecret, ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleBastionHostSSHKeyCredential(cluster.Spec.ClusterSecrets.BastionHostSSHKeySecret, ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleKubeletExtraEnvCredential(cluster.Spec.ClusterSecrets.KubeletExtraEnvSecret, ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssemblePrivateRegistryECRCredential(cluster.Spec.ClusterSecrets.PrivateRegistryECRSecret, ClusterType, cluster.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	return spec, nil
}

// AssembleRKEConfigTemplateSpec is a wrapper assembler for assembling configs on ClusterTemplateRevisions. It returns a ClusterSpec.
func AssembleRKEConfigTemplateSpec(template *apimgmtv3.ClusterTemplateRevision, spec apimgmtv3.ClusterSpec, secretLister v1.SecretLister) (apimgmtv3.ClusterSpec, error) {
	spec, err := AssembleS3Credential(template.Status.S3CredentialSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssemblePrivateRegistryCredential(template.Status.PrivateRegistrySecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleWeaveCredential(template.Status.WeavePasswordSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleVsphereGlobalCredential(template.Status.VsphereSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleVsphereVirtualCenterCredential(template.Status.VirtualCenterSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleOpenStackCredential(template.Status.OpenStackSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleAADClientSecretCredential(template.Status.AADClientSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleAADCertCredential(template.Status.AADClientCertSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleACIAPICUserKeyCredential(template.Status.ACIAPICUserKeySecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleACITokenCredential(template.Status.ACITokenSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleACIKafkaClientKeyCredential(template.Status.ACIKafkaClientKeySecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleSecretsEncryptionProvidersSecretCredential(template.Status.SecretsEncryptionProvidersSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleBastionHostSSHKeyCredential(template.Status.BastionHostSSHKeySecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssembleKubeletExtraEnvCredential(template.Status.KubeletExtraEnvSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	spec, err = AssemblePrivateRegistryECRCredential(template.Status.PrivateRegistryECRSecret, ClusterTemplateRevisionType, template.Name, spec, secretLister)
	if err != nil {
		return spec, err
	}
	return spec, nil
}

// AssembleSMTPCredential looks up the SMTP Secret and inserts the keys into the Notifier.
// It returns a new copy of the Notifier without modifying the original. The Notifier is never updated.
func AssembleSMTPCredential(notifier *apimgmtv3.Notifier, secretLister v1.SecretLister) (*apimgmtv3.NotifierSpec, error) {
	if notifier.Spec.SMTPConfig == nil {
		return &notifier.Spec, nil
	}
	if notifier.Status.SMTPCredentialSecret == "" {
		if notifier.Spec.SMTPConfig.Password != "" {
			logrus.Warnf("[secretmigrator] secrets for notifier %s are not finished migrating", notifier.Name)
		}
		return &notifier.Spec, nil
	}
	smtpSecret, err := secretLister.Get(SecretNamespace, notifier.Status.SMTPCredentialSecret)
	if err != nil {
		return &notifier.Spec, err
	}
	spec := notifier.Spec.DeepCopy()
	spec.SMTPConfig.Password = string(smtpSecret.Data[SecretKey])
	return spec, nil
}

// AssembleWechatCredential looks up the Wechat Secret and inserts the keys into the Notifier.
// It returns a new copy of the Notifier without modifying the original. The Notifier is never updated.
func AssembleWechatCredential(notifier *apimgmtv3.Notifier, secretLister v1.SecretLister) (*apimgmtv3.NotifierSpec, error) {
	if notifier.Spec.WechatConfig == nil {
		return &notifier.Spec, nil
	}
	if notifier.Status.WechatCredentialSecret == "" {
		if notifier.Spec.WechatConfig.Secret != "" {
			logrus.Warnf("[secretmigrator] secrets for notifier %s are not finished migrating", notifier.Name)
		}
		return &notifier.Spec, nil
	}
	wechatSecret, err := secretLister.Get(SecretNamespace, notifier.Status.WechatCredentialSecret)
	if err != nil {
		return &notifier.Spec, err
	}
	spec := notifier.Spec.DeepCopy()
	spec.WechatConfig.Secret = string(wechatSecret.Data[SecretKey])
	return spec, nil
}

// AssembleDingtalkCredential looks up the Dingtalk Secret and inserts the keys into the Notifier.
// It returns a new copy of the Notifier without modifying the original. The Notifier is never updated.
func AssembleDingtalkCredential(notifier *apimgmtv3.Notifier, secretLister v1.SecretLister) (*apimgmtv3.NotifierSpec, error) {
	if notifier.Spec.DingtalkConfig == nil {
		return &notifier.Spec, nil
	}
	if notifier.Status.DingtalkCredentialSecret == "" {
		if notifier.Spec.DingtalkConfig.Secret != "" {
			logrus.Warnf("[secretmigrator] secrets for notifier %s are not finished migrating", notifier.Name)
		}
		return &notifier.Spec, nil
	}
	secret, err := secretLister.Get(SecretNamespace, notifier.Status.DingtalkCredentialSecret)
	if err != nil {
		return &notifier.Spec, err
	}
	spec := notifier.Spec.DeepCopy()
	spec.DingtalkConfig.Secret = string(secret.Data[SecretKey])
	return spec, nil
}
