package ranger

// The json/yaml config key for the ranger config
const ConfigurationFileKey = "ranger"

// Config is configuration need to test against a ranger instance
type Config struct {
	Host          string `yaml:"host"`
	AdminToken    string `yaml:"adminToken"`
	AdminPassword string `yaml:"adminPassword"`
	Insecure      *bool  `yaml:"insecure" default:"true"`
	Cleanup       *bool  `yaml:"cleanup" default:"true"`
	CAFile        string `yaml:"caFile" default:""`
	CACerts       string `yaml:"caCerts" default:""`
	ClusterName   string `yaml:"clusterName" default:""`
}
