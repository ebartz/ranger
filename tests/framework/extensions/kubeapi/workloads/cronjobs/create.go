package cronjobs

import (
	"context"

	"github.com/ranger/ranger/pkg/api/scheme"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/defaults"
	"github.com/ranger/ranger/tests/framework/extensions/unstructured"
	"github.com/ranger/ranger/tests/framework/pkg/wait"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

// CronJobGroupVersionResource is the required Group Version Resource for accessing cron jobs in a cluster,
// using the dynamic client.
var CronJobGroupVersionResource = schema.GroupVersionResource{
	Group:    "batch",
	Version:  "v1beta1",
	Resource: "cronjobs",
}

// CreateCronJob is a helper function that uses the dynamic client to create a cronjob on a namespace for a specific cluster.
// It registers a delete fuction a wait.WatchWait to ensure the cronjob is deleted cleanly.
func CreateCronJob(client *ranger.Client, clusterName, cronJobName, namespace, schedule string, template corev1.PodTemplateSpec) (*v1beta1.CronJob, error) {
	dynamicClient, err := client.GetDownStreamClusterClient(clusterName)
	if err != nil {
		return nil, err
	}

	template.Spec.RestartPolicy = corev1.RestartPolicyNever
	cronJob := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cronJobName,
			Namespace: namespace,
		},
		Spec: v1beta1.CronJobSpec{
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: template,
				},
			},
			Schedule: schedule,
		},
	}

	cronJobResource := dynamicClient.Resource(CronJobGroupVersionResource).Namespace(namespace)

	unstructuredResp, err := cronJobResource.Create(context.TODO(), unstructured.MustToUnstructured(cronJob), metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	client.Session.RegisterCleanupFunc(func() error {
		err := cronJobResource.Delete(context.TODO(), unstructuredResp.GetName(), metav1.DeleteOptions{})
		if errors.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}

		watchInterface, err := cronJobResource.Watch(context.TODO(), metav1.ListOptions{
			FieldSelector:  "metadata.name=" + unstructuredResp.GetName(),
			TimeoutSeconds: &defaults.WatchTimeoutSeconds,
		})

		if err != nil {
			return err
		}

		return wait.WatchWait(watchInterface, func(event watch.Event) (ready bool, err error) {
			if event.Type == watch.Deleted {
				return true, nil
			}
			return false, nil
		})
	})

	newcronJob := &v1beta1.CronJob{}
	err = scheme.Scheme.Convert(unstructuredResp, newcronJob, unstructuredResp.GroupVersionKind())
	if err != nil {
		return nil, err
	}

	return newcronJob, nil
}
