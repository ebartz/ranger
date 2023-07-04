package librke

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/ranger/norman/types/convert"
	v3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
	"github.com/ranger/rke/cluster"
	"github.com/ranger/rke/pki"
	rketypes "github.com/ranger/rke/types"
)

type rke struct {
}

func (*rke) GenerateRKENodeCerts(ctx context.Context, rkeConfig rketypes.RangerKubernetesEngineConfig, nodeAddress string, certBundle map[string]pki.CertificatePKI) map[string]pki.CertificatePKI {
	return pki.GenerateRKENodeCerts(ctx, rkeConfig, nodeAddress, certBundle)
}

func (*rke) GenerateCerts(config *rketypes.RangerKubernetesEngineConfig) (map[string]pki.CertificatePKI, error) {
	return pki.GenerateRKECerts(context.Background(), *config, "", "")
}

func (*rke) GeneratePlan(ctx context.Context, rkeConfig *rketypes.RangerKubernetesEngineConfig, dockerInfo map[string]types.Info, data map[string]interface{}) (rketypes.RKEPlan, error) {
	return cluster.GeneratePlan(ctx, rkeConfig.DeepCopy(), dockerInfo, data)
}

func GetDockerInfo(node *v3.Node) (map[string]types.Info, error) {
	infos := map[string]types.Info{}
	if node.Status.DockerInfo != nil {
		dockerInfo := types.Info{}
		err := convert.ToObj(node.Status.DockerInfo, &dockerInfo)
		if err != nil {
			return nil, err
		}
		infos[node.Status.NodeConfig.Address] = dockerInfo
	}

	return infos, nil
}
