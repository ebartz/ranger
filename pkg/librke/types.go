package librke

import (
	"context"

	rketypes "github.com/ranger/rke/types"

	"github.com/docker/docker/api/types"
	"github.com/ranger/rke/pki"
)

type RKE interface {
	GenerateRKENodeCerts(ctx context.Context, rkeConfig rketypes.RangerKubernetesEngineConfig, nodeAddress string, certBundle map[string]pki.CertificatePKI) map[string]pki.CertificatePKI
	GenerateCerts(config *rketypes.RangerKubernetesEngineConfig) (map[string]pki.CertificatePKI, error)
	GeneratePlan(ctx context.Context, rkeConfig *rketypes.RangerKubernetesEngineConfig, dockerInfo map[string]types.Info, data map[string]interface{}) (rketypes.RKEPlan, error)
}

func New() RKE {
	return (*rke)(nil)
}
