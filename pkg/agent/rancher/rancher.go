package ranger

import (
	"context"
	"io/ioutil"
	"os"
	"sync"

	"github.com/ranger/ranger/pkg/agent/cluster"
	"github.com/ranger/ranger/pkg/features"
	"github.com/ranger/ranger/pkg/namespace"
	"github.com/ranger/ranger/pkg/ranger"
	"github.com/ranger/wrangler/pkg/apply"
	corefactory "github.com/ranger/wrangler/pkg/generated/controllers/core"
	corecontrollers "github.com/ranger/wrangler/pkg/generated/controllers/core/v1"
	"github.com/ranger/wrangler/pkg/kubeconfig"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var (
	started bool
)

func Run(ctx context.Context) error {
	if err := setupSteveAggregation(); err != nil {
		return err
	}

	if started {
		return nil
	}

	if !features.MCMAgent.Enabled() {
		return nil
	}

	cfg, err := kubeconfig.GetNonInteractiveClientConfig("").ClientConfig()
	if err != nil {
		return err
	}

	core, err := corefactory.NewFactoryFromConfig(cfg)
	if err != nil {
		return err
	}

	h := handler{
		ctx:          ctx,
		serviceCache: core.Core().V1().Service().Cache(),
	}

	core.Core().V1().Service().OnChange(ctx, "ranger-installed", h.OnChange)
	if err := core.Start(ctx, 1); err != nil {
		return err
	}
	started = true
	return nil
}

type handler struct {
	lock            sync.Mutex
	ctx             context.Context
	rangerNotFound *bool
	serviceCache    corecontrollers.ServiceCache
}

func (h *handler) startRanger() {
	clientConfig := kubeconfig.GetNonInteractiveClientConfig("")
	server, err := ranger.New(h.ctx, clientConfig, &ranger.Options{
		HTTPListenPort:  80,
		HTTPSListenPort: 443,
		Features:        os.Getenv("CATTLE_FEATURES"),
		AddLocal:        "true",
		ClusterRegistry: os.Getenv("CATTLE_CLUSTER_REGISTRY"),
	})
	if err != nil {
		logrus.Fatalf("Embedded ranger failed to initialize: %v", err)
	}
	go func() {
		err = server.ListenAndServe(h.ctx)
		logrus.Fatalf("Embedded ranger failed to start: %v", err)
	}()
}

func (h *handler) OnChange(key string, service *corev1.Service) (*corev1.Service, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.rangerNotFound == nil {
		_, err := h.serviceCache.Get(namespace.System, "ranger")
		if notFound := apierror.IsNotFound(err); notFound {
			h.rangerNotFound = &notFound
			h.startRanger()
		} else if err != nil {
			return nil, err
		} else {
			h.rangerNotFound = &notFound
		}
	}

	if service == nil {
		if key == namespace.System+"/ranger" {
			logrus.Info("Ranger has been uninstalled, restarting")
			os.Exit(0)
		}
	} else if service.Namespace == namespace.System && service.Name == "ranger" && *h.rangerNotFound {
		logrus.Info("Ranger has been installed, restarting")
		os.Exit(0)
	}

	return service, nil
}

func setupSteveAggregation() error {
	c, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	apply, err := apply.NewForConfig(c)
	if err != nil {
		return err
	}

	token, url, err := cluster.TokenAndURL()
	if err != nil {
		return err
	}

	data := map[string][]byte{
		"CATTLE_SERVER":      []byte(url),
		"CATTLE_TOKEN":       []byte(token),
		"CATTLE_CA_CHECKSUM": []byte(cluster.CAChecksum()),
		"url":                []byte(url + "/v3/connect"),
		"token":              []byte("stv-cluster-" + token),
	}

	ca, err := ioutil.ReadFile("/etc/kubernetes/ssl/certs/serverca")
	if os.IsNotExist(err) {
	} else if err != nil {
		return err
	} else {
		data["ca.crt"] = ca
	}

	return apply.
		WithDynamicLookup().
		WithSetID("ranger-stv-aggregation").
		WithListerNamespace(namespace.System).
		ApplyObjects(&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace.System,
				Name:      "stv-aggregation",
			},
			Data: data,
		})
}
