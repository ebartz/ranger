package charts

import (
	"time"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	v1 "github.com/ranger/ranger/tests/framework/clients/ranger/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	// Project that charts are installed in
	gatekeeperProjectName = "gatekeeper-project"
	// namespace that is created without a label
	RangerDisallowedNamespace  = "no-label"
	ConstraintResourceSteveType = "constraints.gatekeeper.sh.k8srequiredlabels"
)

type ConstraintStatus struct {
	AuditTimestamp  string
	ByPod           interface{}
	TotalViolations int64
	Violations      []interface{}
}

func getAuditTimestamp(client *ranger.Client, project *management.Project) error {
	// wait until the first audit finishes running.
	// AuditTimestamp will be empty string until first audit finishes
	steveClient, err := client.Steve.ProxyDownstream(project.ClusterID)
	if err != nil {
		return err
	}
	return wait.Poll(1*time.Second, 5*time.Minute, func() (done bool, err error) {

		// get list of constraints
		auditList, err := steveClient.SteveType(ConstraintResourceSteveType).List(nil)
		if err != nil {
			return false, nil
		}

		constraintsStatusType := &ConstraintStatus{}
		constraintStatus := auditList.Data[0].Status
		err = v1.ConvertToK8sType(constraintStatus, constraintsStatusType)
		if err != nil {
			return false, nil
		}

		// extract the timestamp of the last constraint audit
		auditTime := constraintsStatusType.AuditTimestamp
		if auditTime == "" {
			return false, nil
		}
		return true, nil
	})

}
