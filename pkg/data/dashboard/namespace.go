package dashboard

import (
	"context"

	"github.com/ranger/ranger/pkg/features"
	fleetconst "github.com/ranger/ranger/pkg/fleet"
	"github.com/ranger/ranger/pkg/namespace"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func addCattleGlobalNamespaces(ctx context.Context, k8s kubernetes.Interface) error {
	if features.Fleet.Enabled() {
		_, err := k8s.CoreV1().Namespaces().Get(ctx, fleetconst.ClustersLocalNamespace, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			_, err = k8s.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: fleetconst.ClustersLocalNamespace,
				},
			}, metav1.CreateOptions{})
		}
		if err != nil {
			return err
		}
	}
	_, err := k8s.CoreV1().Namespaces().Get(ctx, namespace.System, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = k8s.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace.System,
			},
		}, metav1.CreateOptions{})
	}
	return err
}
