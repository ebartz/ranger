package nodes

import (
	"net/url"
	"time"

	"github.com/ranger/norman/types"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	v1 "github.com/ranger/ranger/tests/framework/clients/ranger/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	active                   = "active"
	machineSteveResourceType = "cluster.x-k8s.io.machine"
	etcdLabel                = "rke.cattle.io/etcd-role"
	clusterLabel             = "cluster.x-k8s.io/cluster-name"
	PollInterval             = time.Duration(5 * time.Second)
	PollTimeout              = time.Duration(15 * time.Minute)
)

// IsNodeReady is a helper method that will loop and check if the node is ready in the RKE1 cluster.
// It will return an error if the node is not ready after set amount of time.
func IsNodeReady(client *ranger.Client, ClusterID string) error {
	err := wait.Poll(500*time.Millisecond, 30*time.Minute, func() (bool, error) {
		nodes, err := client.Management.Node.ListAll(&types.ListOpts{
			Filters: map[string]interface{}{
				"clusterId": ClusterID,
			},
		})
		if err != nil {
			return false, err
		}

		for _, node := range nodes.Data {
			node, err := client.Management.Node.ByID(node.ID)
			if err != nil {
				return false, nil
			}

			if node.State != active {
				return false, nil
			}
		}
		logrus.Infof("All nodes in the cluster are in an active state!")
		return true, nil
	})

	return err
}

func IsRKE1EtcdNodeReplaced(client *ranger.Client, etcdNodeToDelete management.Node, clusterResp *management.Cluster, numOfEtcdNodesBeforeDeletion int) (bool, error) {
	numOfEtcdNodesAfterDeletion := 0

	err := wait.Poll(PollInterval, PollTimeout, func() (done bool, err error) {
		machines, err := client.Management.Node.List(&types.ListOpts{Filters: map[string]interface{}{
			"clusterId": clusterResp.ID,
		}})
		if err != nil {
			return false, err
		}
		numOfEtcdNodesAfterDeletion = 0
		for _, machine := range machines.Data {
			if machine.Etcd {
				if machine.ID == etcdNodeToDelete.ID {
					return false, nil
				}
				numOfEtcdNodesAfterDeletion++
			}
		}
		logrus.Info("new etcd node : ")
		for _, machine := range machines.Data {
			if machine.Etcd {
				logrus.Info(machine.NodeName)
			}
		}
		return true, nil
	})
	return numOfEtcdNodesBeforeDeletion == numOfEtcdNodesAfterDeletion, err
}

func IsRKE2K3SEtcdNodeReplaced(client *ranger.Client, query url.Values, clusterName string, etcdNodeToDelete v1.SteveAPIObject, numOfEtcdNodesBeforeDeletion int) (bool, error) {
	numOfEtcdNodesAfterDeletion := 0

	err := wait.Poll(PollInterval, PollTimeout, func() (done bool, err error) {
		machines, err := client.Steve.SteveType(machineSteveResourceType).List(query)
		if err != nil {
			return false, err
		}

		numOfEtcdNodesAfterDeletion = 0
		for _, machine := range machines.Data {
			if machine.Labels[etcdLabel] == "true" && machine.Labels[clusterLabel] == clusterName {
				if machine.Name == etcdNodeToDelete.Name {
					return false, nil
				}
				numOfEtcdNodesAfterDeletion++
			}
		}
		logrus.Info("new etcd node : ")
		for _, machine := range machines.Data {
			if machine.Labels[etcdLabel] == "true" && machine.Labels[clusterLabel] == clusterName {
				logrus.Info(machine.Name)
			}
		}
		return true, nil
	})
	return numOfEtcdNodesBeforeDeletion == numOfEtcdNodesAfterDeletion, err
}
