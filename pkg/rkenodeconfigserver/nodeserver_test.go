package rkenodeconfigserver

import (
	"strings"
	"testing"

	v3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/taints"
	rketypes "github.com/ranger/rke/types"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestAppendKubeletArgs(t *testing.T) {
	type testCase struct {
		name             string
		currentCommand   []string
		taints           []rketypes.RKETaint
		expectedTaintSet map[string]struct{}
	}
	testCases := []testCase{
		testCase{
			name:           "taints args not exists",
			currentCommand: []string{"kubelet", "--register-node"},
			taints: []rketypes.RKETaint{
				rketypes.RKETaint{
					Key:    "test1",
					Value:  "value1",
					Effect: v1.TaintEffectNoSchedule,
				},
				rketypes.RKETaint{
					Key:    "test2",
					Value:  "value2",
					Effect: v1.TaintEffectNoSchedule,
				},
			},
			expectedTaintSet: map[string]struct{}{
				"test1=value1:NoSchedule": struct{}{},
				"test2=value2:NoSchedule": struct{}{},
			},
		},
		testCase{
			name:           "taints args exists",
			currentCommand: []string{"kubelet", "--register-node", "--register-with-taints=node-role.kubernetes.io/controlplane=true:NoSchedule"},
			taints: []rketypes.RKETaint{
				rketypes.RKETaint{
					Key:    "test1",
					Value:  "value1",
					Effect: v1.TaintEffectNoSchedule,
				},
			},
			expectedTaintSet: map[string]struct{}{
				"node-role.kubernetes.io/controlplane=true:NoSchedule": struct{}{},
				"test1=value1:NoSchedule":                              struct{}{},
			},
		},
	}
	for _, tc := range testCases {
		processes := getKubeletProcess(tc.currentCommand)
		afterAppend := AppendTaintsToKubeletArgs(processes, tc.taints)
		appendedCommand := getCommandFromProcesses(afterAppend)
		assert.Equal(t, tc.expectedTaintSet, appendedCommand, "", "")
	}
}

func TestShareMntArgs(t *testing.T) {
	augmentedProcesses := getAugmentedKubeletProcesses()
	args := augmentedProcesses["share-mnt"].Args
	// args are agent call params, by default, arg count is 8 with ca it's 9
	assert.Equal(t, 8, len(args), "default args count for share-mnt should 8")
}

func getKubeletProcess(commands []string) map[string]rketypes.Process {
	return map[string]rketypes.Process{
		"kubelet": rketypes.Process{
			Name:    "kubelet",
			Command: commands,
		},
	}
}

func getAugmentedKubeletProcesses() map[string]rketypes.Process {
	var cluster v3.Cluster
	command := []string{"dummy"}
	binds := []string{"/var/lib/kubelet:/var/lib/kubelet:shared,z", "/var/lib/ranger:/var/lib/ranger:shared,z"}
	processes := map[string]rketypes.Process{
		"kubelet": rketypes.Process{
			Name:    "kubelet",
			Command: command,
			Binds:   binds,
		},
	}

	processes, _ = AugmentProcesses("token", processes, true, "dummynode", &cluster, nil)
	return processes
}

func getCommandFromProcesses(processes map[string]rketypes.Process) map[string]struct{} {
	kubelet, ok := processes["kubelet"]
	if !ok {
		return nil
	}
	rtn := map[string]struct{}{}
	var tmp map[string]int
	for _, command := range kubelet.Command {
		if strings.HasPrefix(command, "--register-with-taints=") {
			tmp = taints.GetTaintSet(taints.GetTaintsFromStrings(strings.Split(strings.TrimPrefix(command, "--register-with-taints="), ",")))
		}
	}
	for key := range tmp {
		rtn[key] = struct{}{}
	}
	return rtn
}
