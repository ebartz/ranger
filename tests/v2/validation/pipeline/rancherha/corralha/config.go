package corralha

// The json/yaml config key for the corral package to be build ..
const (
	CorralRangerHAConfigConfigurationFileKey = "corralRangerHA"
)

// CorralPackages is a struct that has the path to the packages
type CorralRangerHA struct {
	Name string `json:"name" yaml:"name"`
}
