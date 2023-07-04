package client

const (
	ComposeSpecType                = "composeSpec"
	ComposeSpecFieldRangerCompose = "rangerCompose"
)

type ComposeSpec struct {
	RangerCompose string `json:"rangerCompose,omitempty" yaml:"rangerCompose,omitempty"`
}
