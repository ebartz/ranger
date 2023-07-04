package rangerversion

const (
	ConfigurationFileKey = "prime"
)

type Config struct {
	Brand          string `json:"brand" yaml:"brand"`
	GitCommit      string `json:"gitCommit" yaml:"gitCommit"`
	IsPrime        bool   `json:"isPrime" yaml:"isPrime" default:false`
	RangerVersion string `json:"rangerVersion" yaml:"rangerVersion"`
	Registry       string `json:"registry" yaml:"registry"`
}
