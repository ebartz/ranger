package machinenodelookup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ranger/lasso/pkg/dynamic"
	rkev1 "github.com/ranger/ranger/pkg/apis/rke.cattle.io/v1"
	"github.com/ranger/ranger/pkg/capr"
	capicontrollers "github.com/ranger/ranger/pkg/generated/controllers/cluster.x-k8s.io/v1beta1"
	rangercontrollers "github.com/ranger/ranger/pkg/generated/controllers/provisioning.cattle.io/v1"
	rkecontroller "github.com/ranger/ranger/pkg/generated/controllers/rke.cattle.io/v1"
	"github.com/ranger/ranger/pkg/provisioningv2/kubeconfig"
	"github.com/ranger/ranger/pkg/wrangler"
	"github.com/ranger/wrangler/pkg/condition"
	"github.com/ranger/wrangler/pkg/data"
	"github.com/ranger/wrangler/pkg/generic"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	nodeErrorEnqueueTime = 15 * time.Second
)

var (
	bootstrapAPIVersion = fmt.Sprintf("%s/%s", rkev1.SchemeGroupVersion.Group, rkev1.SchemeGroupVersion.Version)
)

type handler struct {
	rangerClusterCache rangercontrollers.ClusterCache
	machineCache        capicontrollers.MachineCache
	machines            capicontrollers.MachineController
	rkeBootstrap        rkecontroller.RKEBootstrapController
	kubeconfigManager   *kubeconfig.Manager
	dynamic             *dynamic.Controller
}

func Register(ctx context.Context, clients *wrangler.Context, kubeconfigManager *kubeconfig.Manager) {
	h := &handler{
		rangerClusterCache: clients.Provisioning.Cluster().Cache(),
		machines:            clients.CAPI.Machine(),
		machineCache:        clients.CAPI.Machine().Cache(),
		rkeBootstrap:        clients.RKE.RKEBootstrap(),
		kubeconfigManager:   kubeconfigManager,
		dynamic:             clients.Dynamic,
	}

	clients.RKE.RKEBootstrap().OnChange(ctx, "machine-node-lookup", h.associateMachineWithNode)
}

func (h *handler) associateMachineWithNode(_ string, bootstrap *rkev1.RKEBootstrap) (*rkev1.RKEBootstrap, error) {
	if bootstrap == nil || bootstrap.DeletionTimestamp != nil {
		return bootstrap, nil
	}

	if !bootstrap.Status.Ready || bootstrap.Status.DataSecretName == nil || *bootstrap.Status.DataSecretName == "" {
		return bootstrap, nil
	}

	machine, err := capr.GetMachineByOwner(h.machineCache, bootstrap)
	if err != nil {
		if errors.Is(err, capr.ErrNoMachineOwnerRef) {
			return bootstrap, generic.ErrSkip
		}
		return bootstrap, err
	}

	if machine.Spec.ProviderID != nil && *machine.Spec.ProviderID != "" {
		// If the machine already has its provider ID set, then we do not need to continue
		return bootstrap, nil
	}

	rangerCluster, err := h.rangerClusterCache.Get(machine.Namespace, machine.Spec.ClusterName)
	if err != nil {
		return bootstrap, err
	}

	config, err := h.kubeconfigManager.GetRESTConfig(rangerCluster, rangerCluster.Status)
	if err != nil {
		return bootstrap, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return bootstrap, err
	}

	nodeLabelSelector := metav1.LabelSelector{MatchLabels: map[string]string{capr.MachineUIDLabel: string(machine.GetUID())}}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(nodeLabelSelector.MatchLabels).String()})
	if err != nil || len(nodes.Items) == 0 || nodes.Items[0].Spec.ProviderID == "" || !condition.Cond("Ready").IsTrue(nodes.Items[0]) {
		logrus.Debugf("Searching for providerID for selector %s in cluster %s/%s, machine %s: %v",
			labels.Set(nodeLabelSelector.MatchLabels), rangerCluster.Namespace, rangerCluster.Name, machine.Name, err)
		h.rkeBootstrap.EnqueueAfter(bootstrap.Namespace, bootstrap.Name, nodeErrorEnqueueTime)
		return bootstrap, nil
	}

	return bootstrap, h.updateMachine(&nodes.Items[0], machine)
}

func (h *handler) updateMachine(node *corev1.Node, machine *capi.Machine) error {
	gvk := schema.FromAPIVersionAndKind(machine.Spec.InfrastructureRef.APIVersion, machine.Spec.InfrastructureRef.Kind)
	infra, err := h.dynamic.Get(gvk, machine.Namespace, machine.Spec.InfrastructureRef.Name)
	if apierror.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	d, err := data.Convert(infra)
	if err != nil {
		return err
	}

	if d.String("spec", "providerID") != node.Spec.ProviderID {
		obj, err := data.Convert(infra.DeepCopyObject())
		if err != nil {
			return err
		}

		obj.SetNested(node.Status.Addresses, "status", "addresses")
		newObj, err := h.dynamic.UpdateStatus(&unstructured.Unstructured{
			Object: obj,
		})
		if err != nil {
			return err
		}

		obj, err = data.Convert(newObj)
		if err != nil {
			return err
		}

		obj.SetNested(node.Spec.ProviderID, "spec", "providerID")
		_, err = h.dynamic.Update(&unstructured.Unstructured{
			Object: obj,
		})
		return err
	}

	return nil
}
