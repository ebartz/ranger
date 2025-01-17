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
		_, err := node.ExecuteCommand("sudo setenforce 1")
		if err != nil {
			return err
		}
		_, err = node.ExecuteCommand("sudo bash -c 'echo vm.panic_on_oom=0 >> /etc/sysctl.conf'")
		if err != nil {
			return err
		}
		_, err = node.ExecuteCommand("sudo bash -c 'echo kernel.panic=10 >> /etc/sysctl.conf'")
		if err != nil {
			return err
		}
		_, err = node.ExecuteCommand("sudo bash -c 'echo kernel.panic_on_oops=1 >> /etc/sysctl.conf'")
		if err != nil {
			return err
		}
		_, err = node.ExecuteCommand("sudo bash -c 'echo kernel.keys.root_maxbytes=25000000 >> /etc/sysctl.conf'")
		if err != nil {
			return err
		}
		_, err = node.ExecuteCommand("sudo bash -c 'sysctl -p /etc/sysctl.conf'")
		if err != nil {
			return err
		}

		if strings.Contains(nodeRoles[key], "--controlplane") {
			logrus.Infof("Copying over files to node %s", node.NodeID)
			user, err := user.Current()
			if err != nil {
				return nil
			}

			dirPath := filepath.Join(user.HomeDir, "go/src/github.com/ranger/ranger/tests/framework/extensions/hardening/k3s")
			err = node.SCPFileToNode(dirPath+"/audit.yaml", "/home/"+node.SSHUser+"/audit.yaml")
			if err != nil {
				return err
			}
			err = node.SCPFileToNode(dirPath+"/psp.yaml", "/home/"+node.SSHUser+"/psp.yaml")
			if err != nil {
				return err
			}
			err = node.SCPFileToNode(dirPath+"/system-policy.yaml", "/home/"+node.SSHUser+"/system-policy.yaml")
			if err != nil {
				return err
			}

			logrus.Infof("Applying hardened YAML files to node: %s", node.NodeID)
			_, err = node.ExecuteCommand("sudo bash -c 'mv /home/" + node.SSHUser + "/audit.yaml /var/lib/ranger/k3s/server/audit.yaml'")
			if err != nil {
				return err
			}
			_, err = node.ExecuteCommand("sudo bash -c 'mv /home/" + node.SSHUser + "/psp.yaml /var/lib/ranger/k3s/psp.yaml'")
			if err != nil {
				return err
			}
			_, err = node.ExecuteCommand("sudo bash -c 'mv /home/" + node.SSHUser + "/system-policy.yaml /var/lib/ranger/k3s/system-policy.yaml'")
			if err != nil {
				return err
			}

			_, err = node.ExecuteCommand("sudo bash -c 'export KUBECONFIG=/etc/ranger/k3s/k3s.yaml && kubectl apply -f /var/lib/ranger/k3s/psp.yaml'")
			if err != nil {
				return err
			}
			_, err = node.ExecuteCommand("sudo bash -c 'export KUBECONFIG=/etc/ranger/k3s/k3s.yaml && kubectl apply -f /var/lib/ranger/k3s/system-policy.yaml'")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
