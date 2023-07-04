#!/bin/sh
set -e -x

# This script publishes docker manifests and images required for multi-architecture support that we (Ranger) don't have automated publishing for.
# We don't have automate publishing for these images because the upstream providers already build images for multiple architectures.
# We need to publish these images under the ranger organization for two reasons:
# 1. To support deploying Ranger in places where quay.io gcr.io are not available
# 2. To support multiple architectures seemlessly through docker manifests. 
# TODO Automate this so that we don't need to run this manually


# quay.io/coreos/etcd:v3.2.24
docker pull quay.io/coreos/etcd:v3.2.24
docker tag quay.io/coreos/etcd:v3.2.24 ranger/coreos-etcd:v3.2.24-amd64
docker pull quay.io/coreos/etcd:v3.2.24-arm64
docker tag quay.io/coreos/etcd:v3.2.24-arm64 ranger/coreos-etcd:v3.2.24-arm64
docker push ranger/coreos-etcd:v3.2.24-amd64
docker push ranger/coreos-etcd:v3.2.24-arm64
docker manifest create ranger/coreos-etcd:v3.2.24-ranger1 ranger/coreos-etcd:v3.2.24-amd64 ranger/coreos-etcd:v3.2.24-arm64
docker manifest annotate ranger/coreos-etcd:v3.2.24-ranger1 ranger/coreos-etcd:v3.2.24-amd64 --arch amd64
docker manifest annotate ranger/coreos-etcd:v3.2.24-ranger1 ranger/coreos-etcd:v3.2.24-arm64 --arch arm64
docker manifest push -p ranger/coreos-etcd:v3.2.24-ranger1

# quay.io/coreos/flannel:v0.10.0
docker pull quay.io/coreos/flannel:v0.10.0
docker tag quay.io/coreos/flannel:v0.10.0 ranger/coreos-flannel:v0.10.0-amd64
docker pull quay.io/coreos/flannel:v0.10.0-arm64
docker tag quay.io/coreos/flannel:v0.10.0-arm64 ranger/coreos-flannel:v0.10.0-arm64
docker push ranger/coreos-flannel:v0.10.0-amd64
docker push ranger/coreos-flannel:v0.10.0-arm64
docker manifest create ranger/coreos-flannel:v0.10.0-ranger1 ranger/coreos-flannel:v0.10.0-amd64 ranger/coreos-flannel:v0.10.0-arm64
docker manifest annotate ranger/coreos-flannel:v0.10.0-ranger1 ranger/coreos-flannel:v0.10.0-amd64 --arch amd64
docker manifest annotate ranger/coreos-flannel:v0.10.0-ranger1 ranger/coreos-flannel:v0.10.0-arm64 --arch arm64
docker manifest push -p ranger/coreos-flannel:v0.10.0-ranger1

# k8s.gcr.io/defaultbackend-amd64:1.4
docker pull k8s.gcr.io/defaultbackend-amd64:1.4
docker tag k8s.gcr.io/defaultbackend-amd64:1.4 ranger/nginx-ingress-controller-defaultbackend:1.4-amd64
docker pull k8s.gcr.io/defaultbackend-arm64:1.4
docker tag k8s.gcr.io/defaultbackend-arm64:1.4 ranger/nginx-ingress-controller-defaultbackend:1.4-arm64
docker push ranger/nginx-ingress-controller-defaultbackend:1.4-amd64
docker push ranger/nginx-ingress-controller-defaultbackend:1.4-arm64
docker manifest create ranger/nginx-ingress-controller-defaultbackend:1.4-ranger1 ranger/nginx-ingress-controller-defaultbackend:1.4-amd64 ranger/nginx-ingress-controller-defaultbackend:1.4-arm64
docker manifest annotate ranger/nginx-ingress-controller-defaultbackend:1.4-ranger1 ranger/nginx-ingress-controller-defaultbackend:1.4-amd64 --arch amd64
docker manifest annotate ranger/nginx-ingress-controller-defaultbackend:1.4-ranger1 ranger/nginx-ingress-controller-defaultbackend:1.4-arm64 --arch arm64
docker manifest push -p ranger/nginx-ingress-controller-defaultbackend:1.4-ranger1

# gcr.io/google_containers/k8s-dns-kube-dns:1.14.13
docker pull gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.13
docker tag gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.13 ranger/k8s-dns-kube-dns:1.14.13-amd64
docker pull gcr.io/google_containers/k8s-dns-kube-dns-arm64:1.14.13
docker tag gcr.io/google_containers/k8s-dns-kube-dns-arm64:1.14.13 ranger/k8s-dns-kube-dns:1.14.13-arm64
docker push ranger/k8s-dns-kube-dns:1.14.13-amd64
docker push ranger/k8s-dns-kube-dns:1.14.13-arm64
docker manifest create ranger/k8s-dns-kube-dns:1.14.13 ranger/k8s-dns-kube-dns:1.14.13-amd64 ranger/k8s-dns-kube-dns:1.14.13-arm64
docker manifest annotate ranger/k8s-dns-kube-dns:1.14.13 ranger/k8s-dns-kube-dns:1.14.13-amd64 --arch amd64
docker manifest annotate ranger/k8s-dns-kube-dns:1.14.13 ranger/k8s-dns-kube-dns:1.14.13-arm64 --arch arm64
docker manifest push -p ranger/k8s-dns-kube-dns:1.14.13

#gcr.io/google_containers/k8s-dns-dnsmasq-nanny:1.14.13
docker pull gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.13
docker tag gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.13 ranger/k8s-dns-dnsmasq-nanny:1.14.13-amd64
docker pull gcr.io/google_containers/k8s-dns-dnsmasq-nanny-arm64:1.14.13
docker tag gcr.io/google_containers/k8s-dns-dnsmasq-nanny-arm64:1.14.13 ranger/k8s-dns-dnsmasq-nanny:1.14.13-arm64
docker push ranger/k8s-dns-dnsmasq-nanny:1.14.13-amd64
docker push ranger/k8s-dns-dnsmasq-nanny:1.14.13-arm64
docker manifest create ranger/k8s-dns-dnsmasq-nanny:1.14.13 ranger/k8s-dns-dnsmasq-nanny:1.14.13-amd64 ranger/k8s-dns-dnsmasq-nanny:1.14.13-arm64
docker manifest annotate ranger/k8s-dns-dnsmasq-nanny:1.14.13 ranger/k8s-dns-dnsmasq-nanny:1.14.13-amd64 --arch amd64
docker manifest annotate ranger/k8s-dns-dnsmasq-nanny:1.14.13 ranger/k8s-dns-dnsmasq-nanny:1.14.13-arm64 --arch arm64
docker manifest push -p ranger/k8s-dns-dnsmasq-nanny:1.14.13

#gcr.io/google_containers/k8s-dns-sidecar:1.14.13
docker pull gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.13
docker tag gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.13 ranger/k8s-dns-sidecar:1.14.13-amd64
docker pull gcr.io/google_containers/k8s-dns-sidecar-arm64:1.14.13
docker tag gcr.io/google_containers/k8s-dns-sidecar-arm64:1.14.13 ranger/k8s-dns-sidecar:1.14.13-arm64
docker push ranger/k8s-dns-sidecar:1.14.13-amd64
docker push ranger/k8s-dns-sidecar:1.14.13-arm64
docker manifest create ranger/k8s-dns-sidecar:1.14.13 ranger/k8s-dns-sidecar:1.14.13-amd64 ranger/k8s-dns-sidecar:1.14.13-arm64
docker manifest annotate ranger/k8s-dns-sidecar:1.14.13 ranger/k8s-dns-sidecar:1.14.13-amd64 --arch amd64
docker manifest annotate ranger/k8s-dns-sidecar:1.14.13 ranger/k8s-dns-sidecar:1.14.13-arm64 --arch arm64
docker manifest push -p ranger/k8s-dns-sidecar:1.14.13

#gcr.io/google_containers/cluster-proportional-autoscaler:1.0.0
docker pull gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0
docker tag gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0 ranger/cluster-proportional-autoscaler:1.0.0-amd64
docker pull gcr.io/google_containers/cluster-proportional-autoscaler-arm64:1.0.0
docker tag gcr.io/google_containers/cluster-proportional-autoscaler-arm64:1.0.0 ranger/cluster-proportional-autoscaler:1.0.0-arm64
docker push ranger/cluster-proportional-autoscaler:1.0.0-amd64
docker push ranger/cluster-proportional-autoscaler:1.0.0-arm64
docker manifest create ranger/cluster-proportional-autoscaler:1.0.0 ranger/cluster-proportional-autoscaler:1.0.0-amd64 ranger/cluster-proportional-autoscaler:1.0.0-arm64
docker manifest annotate ranger/cluster-proportional-autoscaler:1.0.0 ranger/cluster-proportional-autoscaler:1.0.0-amd64 --arch amd64
docker manifest annotate ranger/cluster-proportional-autoscaler:1.0.0 ranger/cluster-proportional-autoscaler:1.0.0-arm64 --arch arm64
docker manifest push -p ranger/cluster-proportional-autoscaler:1.0.0

#gcr.io/google_containers/pause:3.1
docker pull gcr.io/google_containers/pause-amd64:3.1
docker tag gcr.io/google_containers/pause-amd64:3.1 ranger/pause:3.1-amd64
docker pull gcr.io/google_containers/pause-arm64:3.1
docker tag gcr.io/google_containers/pause-arm64:3.1 ranger/pause:3.1-arm64
docker pull gcr.io/google_containers/pause-arm:3.1
docker tag gcr.io/google_containers/pause-arm:3.1 ranger/pause:3.1-arm
docker push ranger/pause:3.1-amd64
docker push ranger/pause:3.1-arm64
docker push ranger/pause:3.1-arm
docker manifest create ranger/pause:3.1 ranger/pause:3.1-amd64 ranger/pause:3.1-arm64 ranger/pause:3.1-arm
docker manifest annotate ranger/pause:3.1 ranger/pause:3.1-amd64 --arch amd64
docker manifest annotate ranger/pause:3.1 ranger/pause:3.1-arm64 --arch arm64
docker manifest annotate ranger/pause:3.1 ranger/pause:3.1-arm --arch arm
docker manifest push -p ranger/pause:3.1

#gcr.io/google_containers/metrics-server:v0.3.1
docker pull gcr.io/google_containers/metrics-server-amd64:v0.3.1
docker tag gcr.io/google_containers/metrics-server-amd64:v0.3.1 ranger/metrics-server:v0.3.1-amd64
docker pull gcr.io/google_containers/metrics-server-arm64:v0.3.1
docker tag gcr.io/google_containers/metrics-server-arm64:v0.3.1 ranger/metrics-server:v0.3.1-arm64
docker push ranger/metrics-server:v0.3.1-amd64
docker push ranger/metrics-server:v0.3.1-arm64
docker manifest create ranger/metrics-server:v0.3.1 ranger/metrics-server:v0.3.1-amd64 ranger/metrics-server:v0.3.1-arm64
docker manifest annotate ranger/metrics-server:v0.3.1 ranger/metrics-server:v0.3.1-amd64 --arch amd64
docker manifest annotate ranger/metrics-server:v0.3.1 ranger/metrics-server:v0.3.1-arm64 --arch arm64
docker manifest push -p ranger/metrics-server:v0.3.1

#gcr.io/google_containers/k8s-dns-kube-dns:1.15.2
docker pull gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.15.2
docker tag gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.15.2 ranger/k8s-dns-kube-dns:1.15.2-amd64
docker pull gcr.io/google_containers/k8s-dns-kube-dns-arm64:1.15.2
docker tag gcr.io/google_containers/k8s-dns-kube-dns-arm64:1.15.2 ranger/k8s-dns-kube-dns:1.15.2-arm64
docker push ranger/k8s-dns-kube-dns:1.15.2-amd64
docker push ranger/k8s-dns-kube-dns:1.15.2-arm64
docker manifest create ranger/k8s-dns-kube-dns:1.15.2 ranger/k8s-dns-kube-dns:1.15.2-amd64 ranger/k8s-dns-kube-dns:1.15.2-arm64
docker manifest annotate ranger/k8s-dns-kube-dns:1.15.2 ranger/k8s-dns-kube-dns:1.15.2-amd64 --arch amd64
docker manifest annotate ranger/k8s-dns-kube-dns:1.15.2 ranger/k8s-dns-kube-dns:1.15.2-arm64 --arch arm64
docker manifest push -p ranger/k8s-dns-kube-dns:1.15.2

#gcr.io/google_containers/k8s-dns-dnsmasq-nanny:1.15.2
docker pull gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.15.2
docker tag gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.15.2 ranger/k8s-dns-dnsmasq-nanny:1.15.2-amd64
docker pull gcr.io/google_containers/k8s-dns-dnsmasq-nanny-arm64:1.15.2
docker tag gcr.io/google_containers/k8s-dns-dnsmasq-nanny-arm64:1.15.2 ranger/k8s-dns-dnsmasq-nanny:1.15.2-arm64
docker push ranger/k8s-dns-dnsmasq-nanny:1.15.2-amd64
docker push ranger/k8s-dns-dnsmasq-nanny:1.15.2-arm64
docker manifest create ranger/k8s-dns-dnsmasq-nanny:1.15.2 ranger/k8s-dns-dnsmasq-nanny:1.15.2-amd64 ranger/k8s-dns-dnsmasq-nanny:1.15.2-arm64
docker manifest annotate ranger/k8s-dns-dnsmasq-nanny:1.15.2 ranger/k8s-dns-dnsmasq-nanny:1.15.2-amd64 --arch amd64
docker manifest annotate ranger/k8s-dns-dnsmasq-nanny:1.15.2 ranger/k8s-dns-dnsmasq-nanny:1.15.2-arm64 --arch arm64
docker manifest push -p ranger/k8s-dns-dnsmasq-nanny:1.15.2

#gcr.io/google_containers/k8s-dns-sidecar:1.15.2
docker pull gcr.io/google_containers/k8s-dns-sidecar-amd64:1.15.2
docker tag gcr.io/google_containers/k8s-dns-sidecar-amd64:1.15.2 ranger/k8s-dns-sidecar:1.15.2-amd64
docker pull gcr.io/google_containers/k8s-dns-sidecar-arm64:1.15.2
docker tag gcr.io/google_containers/k8s-dns-sidecar-arm64:1.15.2 ranger/k8s-dns-sidecar:1.15.2-arm64
docker push ranger/k8s-dns-sidecar:1.15.2-amd64
docker push ranger/k8s-dns-sidecar:1.15.2-arm64
docker manifest create ranger/k8s-dns-sidecar:1.15.2 ranger/k8s-dns-sidecar:1.15.2-amd64 ranger/k8s-dns-sidecar:1.15.2-arm64
docker manifest annotate ranger/k8s-dns-sidecar:1.15.2 ranger/k8s-dns-sidecar:1.15.2-amd64 --arch amd64
docker manifest annotate ranger/k8s-dns-sidecar:1.15.2 ranger/k8s-dns-sidecar:1.15.2-arm64 --arch arm64
docker manifest push -p ranger/k8s-dns-sidecar:1.15.2
