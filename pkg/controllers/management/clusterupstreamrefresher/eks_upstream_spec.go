package clusterupstreamrefresher

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	ekscontroller "github.com/ranger/eks-operator/controller"
	eksv1 "github.com/ranger/eks-operator/pkg/apis/eks.cattle.io/v1"
	mgmtv3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/settings"
	"github.com/ranger/ranger/pkg/wrangler"
	wranglerv1 "github.com/ranger/wrangler/pkg/generated/controllers/core/v1"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	eksUpstreamRefresh       = "eks-refresh"
	eksRefreshCronDeprecated = "eks-refresh-cron"
	eksRefreshCronAnnotation = "settings.management.cattle.io/migrated"
)

func BuildEKSUpstreamSpec(secretsCache wranglerv1.SecretCache, cluster *mgmtv3.Cluster) (*eksv1.EKSClusterConfigSpec, error) {
	sess, eksService, err := ekscontroller.StartAWSSessions(secretsCache, *cluster.Spec.EKSConfig)
	if err != nil {
		return nil, err
	}

	clusterState, err := eksService.DescribeCluster(
		&eks.DescribeClusterInput{
			Name: aws.String(cluster.Spec.EKSConfig.DisplayName),
		})
	if err != nil {
		return nil, err
	}

	ngs, err := eksService.ListNodegroups(
		&eks.ListNodegroupsInput{
			ClusterName: aws.String(cluster.Spec.EKSConfig.DisplayName),
		})
	if err != nil {
		return nil, err
	}

	// gather upstream node groups states
	var nodeGroupStates []*eks.DescribeNodegroupOutput
	var errs []string
	for _, ngName := range ngs.Nodegroups {
		ng, err := eksService.DescribeNodegroup(
			&eks.DescribeNodegroupInput{
				ClusterName:   aws.String(cluster.Spec.EKSConfig.DisplayName),
				NodegroupName: ngName,
			})
		if err != nil {
			return nil, err
		}

		nodeGroupStates = append(nodeGroupStates, ng)
		var nodeGroupMustBeDeleted string
		if len(ng.Nodegroup.Health.Issues) != 0 {
			var issueMessages []string
			for _, issue := range ng.Nodegroup.Health.Issues {
				issueMessages = append(issueMessages, aws.StringValue(issue.Message))
				if !ekscontroller.NodeGroupIssueIsUpdatable(aws.StringValue(issue.Code)) {
					nodeGroupMustBeDeleted = ": node group cannot be updated, must be deleted and recreated"
				}
			}
			errs = append(errs, fmt.Sprintf("health error for node group [%s] in cluster [%s]: %s%s",
				aws.StringValue(ng.Nodegroup.NodegroupName),
				cluster.Name,
				strings.Join(issueMessages, "; "),
				nodeGroupMustBeDeleted,
			))
		}
	}

	upstreamSpec, _, err := ekscontroller.BuildUpstreamClusterState(cluster.Spec.DisplayName, cluster.Status.EKSStatus.ManagedLaunchTemplateID, clusterState, nodeGroupStates, ec2.New(sess), false)
	if err != nil {
		// If we get an error here, then syncing is broken
		return nil, err
	}

	upstreamSpec.DisplayName = cluster.Spec.EKSConfig.DisplayName
	upstreamSpec.Region = cluster.Spec.EKSConfig.Region
	upstreamSpec.AmazonCredentialSecret = cluster.Spec.EKSConfig.AmazonCredentialSecret
	upstreamSpec.Imported = cluster.Spec.EKSConfig.Imported
	upstreamSpec.Subnets = cluster.Spec.EKSConfig.Subnets
	upstreamSpec.SecurityGroups = cluster.Spec.EKSConfig.SecurityGroups
	upstreamSpec.ServiceRole = cluster.Spec.EKSConfig.ServiceRole

	if len(errs) != 0 {
		// If there are errors here, we can still sync, but there are problems with the nodegroups that should be reported
		err = fmt.Errorf("error for cluster [%s]: %s",
			cluster.Name,
			strings.Join(errs, "\n"))
	}

	return upstreamSpec, err
}

// MigrateEksRefreshCronSetting migrates the deprecated eks-refresh-cron setting to new
// setting only if default setting was changed
// This function will be run only once during startup by pkg/multiclustermanager/app.go
func MigrateEksRefreshCronSetting(wContext *wrangler.Context) {
	settingsClient := wContext.Mgmt.Setting()
	eksCronSetting, err := settingsClient.Get(eksRefreshCronDeprecated, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return
	} else if err != nil {
		logrus.Errorf("Unable to complete EKS cron migration, will attempt at next ranger startup. "+
			"Error getting %s setting: %v", eksRefreshCronDeprecated, err)
		return
	}
	if eksCronSetting.Annotations != nil && eksCronSetting.Annotations[eksRefreshCronAnnotation] == "true" {
		return
	}

	eksCronAnnotate := make(map[string]string)
	if eksCronSetting.Annotations != nil {
		eksCronAnnotate = eksCronSetting.Annotations
	}
	eksCronAnnotate[eksRefreshCronAnnotation] = "true"

	settingsClientCache := wContext.Mgmt.Setting().Cache()
	eksRefreshSetting, err := settingsClientCache.Get(eksUpstreamRefresh)
	if errors.IsNotFound(err) {
		return
	} else if err != nil {
		logrus.Errorf("Unable to complete EKS cron migration, will attempt at next ranger startup. "+
			"Error getting %s setting: %v", eksUpstreamRefresh, err)
		return
	}

	if eksRefreshSetting.Value != "" || eksCronSetting.Value == "" {
		eksCronSetting.SetAnnotations(eksCronAnnotate)
		if _, err = settingsClient.Update(eksCronSetting); err != nil {
			logrus.Errorf("Unable to complete EKS cron migration, will attempt at next ranger startup. "+
				"Error annotating eks-refresh-cron setting: %v", err)
		}
		return
	}

	eksSchedule, err := cron.ParseStandard(eksCronSetting.Value)
	if err != nil {
		logrus.Errorf("Unable to complete EKS cron migration, will attempt at next ranger startup. "+
			"Error parsing cron schedule %s setting: %v", eksRefreshCronDeprecated, err)
		return
	}

	next := eksSchedule.Next(time.Now())
	refreshTime := int(eksSchedule.Next(next).Sub(next) / time.Second)

	err = settings.EKSUpstreamRefresh.Set(fmt.Sprint(refreshTime))
	if err != nil {
		logrus.Errorf("Unable to complete EKS cron migration, will attempt at next ranger startup. "+
			"Error updating eks-refresh setting: %v", err)
	}
	eksCronSetting.SetAnnotations(eksCronAnnotate)
	if _, err = settingsClient.Update(eksCronSetting); err != nil {
		logrus.Errorf("Unable to complete EKS cron migration, will attempt at next ranger startup. "+
			"Error annotating eks-refresh-cron setting: %v", err)
	}
}
