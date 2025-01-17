package apiservice

import (
	"fmt"

	v3 "github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/features"
	"github.com/ranger/ranger/pkg/namespace"
	"github.com/ranger/ranger/pkg/settings"
	corev1 "k8s.io/api/core/v1"
)

const RangerServiceName = "ranger"

func (h *handler) SetupInternalServerURL(key string, setting *v3.Setting) (*v3.Setting, error) {
	if key != settings.ServerURL.Name {
		return setting, nil
	}

	internalURL, internalCA, err := h.getInternalServerAndURL()
	if err != nil {
		return nil, err
	}

	// purposely update CA before URL, because we only wait for internalURL != "" when checking if it's initialized
	if settings.InternalCACerts.Get() != internalCA {
		if err := settings.InternalCACerts.Set(internalCA); err != nil {
			return setting, err
		}
	}

	if settings.InternalServerURL.Get() != internalURL {
		if err := settings.InternalServerURL.Set(internalURL); err != nil {
			return setting, err
		}
	}

	return setting, nil
}

func (h *handler) getClusterIP() (string, error) {
	serviceName := RangerServiceName
	if features.MCMAgent.Enabled() {
		serviceName = "cattle-cluster-agent"
	}
	service, err := h.services.Get(namespace.System, serviceName)
	if err != nil {
		return "", err
	}
	if service.Spec.ClusterIP == "" {
		return "", fmt.Errorf("waiting on service %s/%s to be assigned a ClusterIP", namespace.System, serviceName)
	}
	return service.Spec.ClusterIP, nil
}

func (h *handler) getInternalServerAndURL() (string, string, error) {
	serverURL := settings.ServerURL.Get()
	ca := settings.CACerts.Get()

	tlsSecret, err := h.secretsCache.Get(namespace.System, "tls-ranger-internal-ca")
	if err != nil {
		return "", "", err
	}
	internalCA := string(tlsSecret.Data[corev1.TLSCertKey])

	clusterIPService := false
	if features.MCMAgent.Enabled() {
		if _, err := h.deploymentCache.Get(namespace.System, "cattle-cluster-agent"); err == nil {
			clusterIPService = true
		}
	} else {
		if dp, err := h.deploymentCache.Get(namespace.System, RangerServiceName); err == nil && dp.Spec.Replicas != nil && *dp.Spec.Replicas != 0 {
			clusterIPService = true
		} else if _, err := h.daemonSetCache.Get(namespace.System, RangerServiceName); err == nil {
			clusterIPService = true
		}
	}

	if h.embedded {
		clusterIPService = true
	}

	if clusterIPService {
		clusterIP, err := h.getClusterIP()
		return fmt.Sprintf("https://%s", clusterIP), internalCA, err
	}

	return serverURL, ca, nil
}
