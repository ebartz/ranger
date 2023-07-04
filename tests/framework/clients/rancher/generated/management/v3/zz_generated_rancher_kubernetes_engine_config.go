package client

const (
	RangerKubernetesEngineConfigType                     = "rangerKubernetesEngineConfig"
	RangerKubernetesEngineConfigFieldAddonJobTimeout     = "addonJobTimeout"
	RangerKubernetesEngineConfigFieldAddons              = "addons"
	RangerKubernetesEngineConfigFieldAddonsInclude       = "addonsInclude"
	RangerKubernetesEngineConfigFieldAuthentication      = "authentication"
	RangerKubernetesEngineConfigFieldAuthorization       = "authorization"
	RangerKubernetesEngineConfigFieldBastionHost         = "bastionHost"
	RangerKubernetesEngineConfigFieldCloudProvider       = "cloudProvider"
	RangerKubernetesEngineConfigFieldClusterName         = "clusterName"
	RangerKubernetesEngineConfigFieldDNS                 = "dns"
	RangerKubernetesEngineConfigFieldEnableCRIDockerd    = "enableCriDockerd"
	RangerKubernetesEngineConfigFieldIgnoreDockerVersion = "ignoreDockerVersion"
	RangerKubernetesEngineConfigFieldIngress             = "ingress"
	RangerKubernetesEngineConfigFieldMonitoring          = "monitoring"
	RangerKubernetesEngineConfigFieldNetwork             = "network"
	RangerKubernetesEngineConfigFieldNodes               = "nodes"
	RangerKubernetesEngineConfigFieldPrefixPath          = "prefixPath"
	RangerKubernetesEngineConfigFieldPrivateRegistries   = "privateRegistries"
	RangerKubernetesEngineConfigFieldRestore             = "restore"
	RangerKubernetesEngineConfigFieldRotateCertificates  = "rotateCertificates"
	RangerKubernetesEngineConfigFieldRotateEncryptionKey = "rotateEncryptionKey"
	RangerKubernetesEngineConfigFieldSSHAgentAuth        = "sshAgentAuth"
	RangerKubernetesEngineConfigFieldSSHCertPath         = "sshCertPath"
	RangerKubernetesEngineConfigFieldSSHKeyPath          = "sshKeyPath"
	RangerKubernetesEngineConfigFieldServices            = "services"
	RangerKubernetesEngineConfigFieldUpgradeStrategy     = "upgradeStrategy"
	RangerKubernetesEngineConfigFieldVersion             = "kubernetesVersion"
	RangerKubernetesEngineConfigFieldWindowsPrefixPath   = "winPrefixPath"
)

type RangerKubernetesEngineConfig struct {
	AddonJobTimeout     int64                `json:"addonJobTimeout,omitempty" yaml:"addonJobTimeout,omitempty"`
	Addons              string               `json:"addons,omitempty" yaml:"addons,omitempty"`
	AddonsInclude       []string             `json:"addonsInclude,omitempty" yaml:"addonsInclude,omitempty"`
	Authentication      *AuthnConfig         `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	Authorization       *AuthzConfig         `json:"authorization,omitempty" yaml:"authorization,omitempty"`
	BastionHost         *BastionHost         `json:"bastionHost,omitempty" yaml:"bastionHost,omitempty"`
	CloudProvider       *CloudProvider       `json:"cloudProvider,omitempty" yaml:"cloudProvider,omitempty"`
	ClusterName         string               `json:"clusterName,omitempty" yaml:"clusterName,omitempty"`
	DNS                 *DNSConfig           `json:"dns,omitempty" yaml:"dns,omitempty"`
	EnableCRIDockerd    *bool                `json:"enableCriDockerd,omitempty" yaml:"enableCriDockerd,omitempty"`
	IgnoreDockerVersion *bool                `json:"ignoreDockerVersion,omitempty" yaml:"ignoreDockerVersion,omitempty"`
	Ingress             *IngressConfig       `json:"ingress,omitempty" yaml:"ingress,omitempty"`
	Monitoring          *MonitoringConfig    `json:"monitoring,omitempty" yaml:"monitoring,omitempty"`
	Network             *NetworkConfig       `json:"network,omitempty" yaml:"network,omitempty"`
	Nodes               []RKEConfigNode      `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	PrefixPath          string               `json:"prefixPath,omitempty" yaml:"prefixPath,omitempty"`
	PrivateRegistries   []PrivateRegistry    `json:"privateRegistries,omitempty" yaml:"privateRegistries,omitempty"`
	Restore             *RestoreConfig       `json:"restore,omitempty" yaml:"restore,omitempty"`
	RotateCertificates  *RotateCertificates  `json:"rotateCertificates,omitempty" yaml:"rotateCertificates,omitempty"`
	RotateEncryptionKey bool                 `json:"rotateEncryptionKey,omitempty" yaml:"rotateEncryptionKey,omitempty"`
	SSHAgentAuth        bool                 `json:"sshAgentAuth,omitempty" yaml:"sshAgentAuth,omitempty"`
	SSHCertPath         string               `json:"sshCertPath,omitempty" yaml:"sshCertPath,omitempty"`
	SSHKeyPath          string               `json:"sshKeyPath,omitempty" yaml:"sshKeyPath,omitempty"`
	Services            *RKEConfigServices   `json:"services,omitempty" yaml:"services,omitempty"`
	UpgradeStrategy     *NodeUpgradeStrategy `json:"upgradeStrategy,omitempty" yaml:"upgradeStrategy,omitempty"`
	Version             string               `json:"kubernetesVersion,omitempty" yaml:"kubernetesVersion,omitempty"`
	WindowsPrefixPath   string               `json:"winPrefixPath,omitempty" yaml:"winPrefixPath,omitempty"`
}
