package hardening

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/pkg/nodes"
	"github.com/sirupsen/logrus"
)

func HardeningNodes(client *ranger.Client, hardened bool, nodes []*nodes.Node, nodeRoles []string) error {
	for key, node := range nodes {
		logrus.Infof("Setting kernel parameters on node %s", node.NodeID)
		_, err := node.ExecuteCommand("sudo bash -c 'echo vm.panic_on_oom=0 >> /etc/sysctl.d/90-kubelet.conf'")
		if err != nil {
			return err
		}
		_, err = node.ExecuteCommand("sudo bash -c 'echo vm.overcommit_memory=1 >> /etc/sysctl.d/90-kubelet.conf'")
		if err != nil {
			return err
		}
		_, err = node.ExecuteCommand("sudo bash -c 'echo kernel.panic=10 >> /etc/sysctl.d/90-kubelet.conf'")
		if err != nil {
			return err
		}
		_, err = node.ExecuteCommand("sudo bash -c 'echo kernel.panic_on_oops=1 >> /etc/sysctl.d/90-kubelet.conf'")
		if err != nil {
			return err
		}
		_, err = node.ExecuteCommand("sudo bash -c 'sysctl -p /etc/sysctl.d/90-kubelet.conf'")
		if err != nil {
			return err
		}
		if strings.Contains(nodeRoles[key], "--etcd") {
			_, err = node.ExecuteCommand("sudo useradd -r -c \"etcd user\" -s /sbin/nologin -M etcd -U")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func PostHardeningConfig(client *ranger.Client, hardened bool, nodes []*nodes.Node, nodeRoles []string) error {
	for key, node := range nodes {
		if strings.Contains(nodeRoles[key], "--controlplane") {
			logrus.Infof("Copying over files to node %s", node.NodeID)
			user, err := user.Current()
			if err != nil {
				return nil
			}

			dirPath := filepath.Join(user.HomeDir, "go/src/github.com/ranger/ranger/tests/framework/extensions/hardening/rke2")
			err = node.SCPFileToNode(dirPath+"/account-update.yaml", "/home/"+node.SSHUser+"/account-update.yaml")
			if err != nil {
				return err
			}
			err = node.SCPFileToNode(dirPath+"/account-update.sh", "/home/"+node.SSHUser+"/account-update.sh")
			if err != nil {
				return err
			}
			_, err = node.ExecuteCommand("sudo bash -c 'mv /home/" + node.SSHUser + "/account-update.yaml /var/lib/ranger/rke2/server/account-update.yaml'")
			if err != nil {
				return err
			}
			_, err = node.ExecuteCommand("sudo bash -c 'mv /home/" + node.SSHUser + "/account-update.sh /var/lib/ranger/rke2/server/account-update.sh'")
			if err != nil {
				return err
			}
			_, err = node.ExecuteCommand("sudo bash -c 'chmod +x /var/lib/ranger/rke2/server/account-update.sh'")
			if err != nil {
				return err
			}
			_, err = node.ExecuteCommand("sudo bash -c 'export KUBECONFIG=/etc/ranger/rke2/rke2.yaml && /var/lib/ranger/rke2/server/account-update.sh'")
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
