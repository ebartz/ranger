FROM registry.suse.com/bci/bci-base:15.4

RUN zypper -n install --no-recommends git-core curl ca-certificates unzip xz gzip sed tar shadow gawk vim-small netcat-openbsd mkisofs && \
    zypper -n clean -a && rm -rf /tmp/* /var/tmp/* /usr/share/doc/packages/* && \
    useradd ranger && \
    mkdir -p /var/lib/ranger /var/lib/cattle /opt/jail /opt/drivers/management-state/bin && \
    chown -R ranger /var/lib/ranger /var/lib/cattle /usr/local/bin

RUN mkdir /root/.kube && \
    ln -s /etc/ranger/k3s/k3s.yaml /root/.kube/k3s.yaml  && \
    ln -s /etc/ranger/k3s/k3s.yaml /root/.kube/config && \
    ln -s $(GOCOVERDIR=rangercoverage /usr/bin/ranger) /usr/bin/reset-password && \
    ln -s $(GOCOVERDIR=rangercoverage /usr/bin/ranger) /usr/bin/ensure-default-admin
WORKDIR /var/lib/ranger

ARG ARCH=amd64
ARG ETCD_UNSUPPORTED_ARCH
ARG IMAGE_REPO=ranger
ARG SYSTEM_CHART_DEFAULT_BRANCH=dev-v2.7
ARG CHART_DEFAULT_BRANCH=dev-v2.7
ARG PARTNER_CHART_DEFAULT_BRANCH=main
ARG RKE2_CHART_DEFAULT_BRANCH=main
# kontainer-driver-metadata branch to be set for specific branch other than dev/master, logic at ranger/ranger/pkg/settings/setting.go
ARG CATTLE_KDM_BRANCH=dev-v2.7

ENV CATTLE_SYSTEM_CHART_DEFAULT_BRANCH=$SYSTEM_CHART_DEFAULT_BRANCH
ENV CATTLE_CHART_DEFAULT_BRANCH=$CHART_DEFAULT_BRANCH
ENV CATTLE_PARTNER_CHART_DEFAULT_BRANCH=$PARTNER_CHART_DEFAULT_BRANCH
ENV CATTLE_RKE2_CHART_DEFAULT_BRANCH=$RKE2_CHART_DEFAULT_BRANCH
ENV CATTLE_HELM_VERSION v2.16.8-ranger2
ENV CATTLE_MACHINE_VERSION v0.15.0-ranger100
ENV CATTLE_K3S_VERSION v1.26.4+k3s1
ENV CATTLE_MACHINE_PROVISION_IMAGE ranger/machine:${CATTLE_MACHINE_VERSION}
ENV CATTLE_ETCD_VERSION v3.5.1
ENV LOGLEVEL_VERSION v0.1.5
ENV TINI_VERSION v0.18.0
ENV TELEMETRY_VERSION v0.5.20
ENV DOCKER_MACHINE_LINODE_VERSION v0.1.8
ENV LINODE_UI_DRIVER_VERSION v0.5.0
# make sure the version number is consistent with the one at Line 100 of pkg/data/management/machinedriver_data.go
ENV DOCKER_MACHINE_HARVESTER_VERSION v0.6.3
ENV CATTLE_KDM_BRANCH ${CATTLE_KDM_BRANCH}
ENV HELM_VERSION v3.11.3
ENV KUSTOMIZE_VERSION v5.0.1
ENV CATTLE_WINS_AGENT_VERSION v0.4.11
ENV CATTLE_WINS_AGENT_INSTALL_SCRIPT https://raw.githubusercontent.com/ranger/wins/${CATTLE_WINS_AGENT_VERSION}/install.ps1
ENV CATTLE_WINS_AGENT_UNINSTALL_SCRIPT https://raw.githubusercontent.com/ranger/wins/${CATTLE_WINS_AGENT_VERSION}/uninstall.ps1
ENV CATTLE_WINS_AGENT_UPGRADE_IMAGE ranger/wins:${CATTLE_WINS_AGENT_VERSION}
ENV CATTLE_CSI_PROXY_AGENT_VERSION v1.1.1
# make sure the CATTLE_SYSTEM_AGENT_VERSION is consistent with tests/v2/codecoverage/package/Dockerfile and pkg/settings/setting.go
ENV CATTLE_SYSTEM_AGENT_VERSION v0.3.3-rc4
ENV CATTLE_SYSTEM_AGENT_INSTALL_SCRIPT https://raw.githubusercontent.com/ranger/system-agent/${CATTLE_SYSTEM_AGENT_VERSION}/install.sh
ENV CATTLE_SYSTEM_AGENT_UNINSTALL_SCRIPT https://raw.githubusercontent.com/ranger/system-agent/${CATTLE_SYSTEM_AGENT_VERSION}/system-agent-uninstall.sh
ENV CATTLE_SYSTEM_AGENT_UPGRADE_IMAGE ranger/system-agent:${CATTLE_SYSTEM_AGENT_VERSION}-suc
ENV CATTLE_SYSTEM_UPGRADE_CONTROLLER_CHART_VERSION 102.1.0+up0.5.0

# System charts minimal version
ENV CATTLE_FLEET_MIN_VERSION=102.1.0+up0.7.0-rc.3
# Deprecated in favor of CATTLE_RANCHER_WEBHOOK_VERSION.
ENV CATTLE_RANCHER_WEBHOOK_MIN_VERSION=''
ENV CATTLE_RANCHER_WEBHOOK_VERSION=2.0.5+up0.3.5-rc6
ENV CATTLE_CSP_ADAPTER_MIN_VERSION=2.0.2+up2.0.2-rc2

RUN mkdir -p /var/lib/ranger-data/local-catalogs/system-library && \
    mkdir -p /var/lib/ranger-data/local-catalogs/library && \
    mkdir -p /var/lib/ranger-data/local-catalogs/helm3-library && \
    mkdir -p /var/lib/ranger-data/local-catalogs/v2 && \
    git clone -b $CATTLE_SYSTEM_CHART_DEFAULT_BRANCH --depth 1 https://github.com/ranger/system-charts /var/lib/ranger-data/local-catalogs/system-library && \
    # Charts need to be copied into the sha256 value of git url computed in https://github.com/ranger/ranger/blob/5ebda9ac23c06e9647b586ec38aa51cc9ff9b031/pkg/catalogv2/git/download.go#L102 to create a unique folder per url
    git clone -b $CATTLE_CHART_DEFAULT_BRANCH --depth 1 https://git.ranger.io/charts /var/lib/ranger-data/local-catalogs/v2/ranger-charts/4b40cac650031b74776e87c1a726b0484d0877c3ec137da0872547ff9b73a721/ && \
    git clone -b $CATTLE_PARTNER_CHART_DEFAULT_BRANCH --depth 1 https://git.ranger.io/partner-charts /var/lib/ranger-data/local-catalogs/v2/ranger-partner-charts/8f17acdce9bffd6e05a58a3798840e408c4ea71783381ecd2e9af30baad65974 && \
    git clone -b $CATTLE_RKE2_CHART_DEFAULT_BRANCH --depth 1 https://git.ranger.io/rke2-charts /var/lib/ranger-data/local-catalogs/v2/ranger-rke2-charts/675f1b63a0a83905972dcab2794479ed599a6f41b86cd6193d69472d0fa889c9 && \
    git clone -b master --depth 1 https://github.com/ranger/charts /var/lib/ranger-data/local-catalogs/library && \
    git clone -b master --depth 1 https://github.com/ranger/helm3-charts /var/lib/ranger-data/local-catalogs/helm3-library

RUN curl -sLf https://github.com/ranger/machine/releases/download/${CATTLE_MACHINE_VERSION}/ranger-machine-${ARCH}.tar.gz | tar xvzf - -C /usr/bin && \
    curl -sLf https://github.com/ranger/loglevel/releases/download/${LOGLEVEL_VERSION}/loglevel-${ARCH}-${LOGLEVEL_VERSION}.tar.gz | tar xvzf - -C /usr/bin && \
    curl -LO https://github.com/linode/docker-machine-driver-linode/releases/download/${DOCKER_MACHINE_LINODE_VERSION}/docker-machine-driver-linode_linux-amd64.zip && \
    unzip docker-machine-driver-linode_linux-amd64.zip -d /opt/drivers/management-state/bin && \
    mkdir -p /usr/share/ranger/ui/assets/ && \
    cp /opt/drivers/management-state/bin/docker-machine-driver-linode /usr/share/ranger/ui/assets/ && \
    rm docker-machine-driver-linode_linux-amd64.zip

RUN curl -LO https://releases.ranger.com/harvester-node-driver/${DOCKER_MACHINE_HARVESTER_VERSION}/docker-machine-driver-harvester-amd64.tar.gz && \
    tar -xf docker-machine-driver-harvester-amd64.tar.gz -C /opt/drivers/management-state/bin && \
    cp /opt/drivers/management-state/bin/docker-machine-driver-harvester /usr/share/ranger/ui/assets/ && \
    rm docker-machine-driver-harvester-amd64.tar.gz

ENV TINI_URL_amd64=https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini \
    TINI_URL_arm64=https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-arm64 \
    TINI_URL_s390x=https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-s390x \
    TINI_URL=TINI_URL_${ARCH}

ENV HELM_URL_V2_amd64=https://github.com/ranger/helm/releases/download/${CATTLE_HELM_VERSION}/ranger-helm \
    HELM_URL_V2_arm64=https://github.com/ranger/helm/releases/download/${CATTLE_HELM_VERSION}/ranger-helm-arm64 \
    HELM_URL_V2_s390x=https://github.com/ranger/helm/releases/download/${CATTLE_HELM_VERSION}/ranger-helm-s390x \
    HELM_URL_V2=HELM_URL_V2_${ARCH} \
    HELM_URL_V3=https://get.helm.sh/helm-${HELM_VERSION}-linux-${ARCH}.tar.gz \
    TILLER_URL_amd64=https://github.com/ranger/helm/releases/download/${CATTLE_HELM_VERSION}/ranger-tiller \
    TILLER_URL_arm64=https://github.com/ranger/helm/releases/download/${CATTLE_HELM_VERSION}/ranger-tiller-arm64 \
    TILLER_URL_s390x=https://github.com/ranger/helm/releases/download/${CATTLE_HELM_VERSION}/ranger-tiller-s390x \
    TILLER_URL=TILLER_URL_${ARCH} \
    ETCD_URL=https://github.com/etcd-io/etcd/releases/download/${CATTLE_ETCD_VERSION}/etcd-${CATTLE_ETCD_VERSION}-linux-${ARCH}.tar.gz \
    KUSTOMIZE_URL=https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_linux_${ARCH}.tar.gz

RUN curl -sLf ${KUSTOMIZE_URL} | tar -xzf - -C /usr/bin

# set up helm 2
RUN curl -sLf ${!HELM_URL_V2} > /usr/bin/ranger-helm && \
    curl -sLf ${!TILLER_URL} > /usr/bin/ranger-tiller && \
    ln -s /usr/bin/ranger-helm /usr/bin/helm && \
    ln -s /usr/bin/ranger-tiller /usr/bin/tiller && \
    chmod +x /usr/bin/ranger-helm /usr/bin/ranger-tiller

# set up helm 3
RUN curl ${HELM_URL_V3} | tar xvzf - --strip-components=1 -C /usr/bin && \
    mv /usr/bin/helm /usr/bin/helm_v3 && \
    chmod +x /usr/bin/kustomize

# Set up K3s: copy the necessary binaries from the K3s image.
COPY --from=ranger/k3s:v1.26.4-k3s1 \
    /bin/blkid \
    /bin/cni \
    /bin/conntrack \
    /bin/containerd \
    /bin/containerd-shim-runc-v2 \
    /bin/ethtool \
    /bin/ip \
    /bin/ipset \
    /bin/k3s \
    /bin/losetup \
    /bin/pigz \
    /bin/runc \
    /bin/which \
    /bin/aux/xtables-legacy-multi \
/usr/bin/

RUN ln -s /usr/bin/cni /usr/bin/bridge && \
    ln -s /usr/bin/cni /usr/bin/flannel && \
    ln -s /usr/bin/cni /usr/bin/host-local && \
    ln -s /usr/bin/cni /usr/bin/loopback && \
    ln -s /usr/bin/cni /usr/bin/portmap && \
    ln -s /usr/bin/k3s /usr/bin/crictl && \
    ln -s /usr/bin/k3s /usr/bin/ctr && \
    ln -s /usr/bin/k3s /usr/bin/k3s-agent && \
    ln -s /usr/bin/k3s /usr/bin/k3s-etcd-snapshot && \
    ln -s /usr/bin/k3s /usr/bin/k3s-server && \
    ln -s /usr/bin/k3s /usr/bin/kubectl && \
    ln -s /usr/bin/pigz /usr/bin/unpigz && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/iptables && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/iptables-save && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/iptables-restore && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/iptables-translate && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/ip6tables && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/ip6tables-save && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/ip6tables-restore && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/ip6tables-translate

RUN curl -sLf ${!TINI_URL} > /usr/bin/tini && \
    mkdir -p /var/lib/ranger/k3s/agent/images/ && \
    curl -sfL ${ETCD_URL} | tar xvzf - --strip-components=1 --no-same-owner -C /usr/bin/ etcd-${CATTLE_ETCD_VERSION}-linux-${ARCH}/etcdctl && \
    curl -sLf https://github.com/ranger/telemetry/releases/download/${TELEMETRY_VERSION}/telemetry-${ARCH} > /usr/bin/telemetry && \
    chmod +x /usr/bin/tini /usr/bin/telemetry && \
    mkdir -p /var/lib/ranger-data/driver-metadata

ENV CATTLE_UI_VERSION 2.7.5-rc2
ENV CATTLE_DASHBOARD_UI_VERSION v2.7.5-rc2
ENV CATTLE_CLI_VERSION v2.7.0

# Base UI brand used as a fallback env setting (not user facing) to indicate this is a non-prime install
ENV CATTLE_BASE_UI_BRAND=

# Please update the api-ui-version in pkg/settings/settings.go when updating the version here.
ENV CATTLE_API_UI_VERSION 1.1.10

RUN mkdir -p /var/log/auditlog
ENV AUDIT_LOG_PATH /var/log/auditlog/ranger-api-audit.log
ENV AUDIT_LOG_MAXAGE 10
ENV AUDIT_LOG_MAXBACKUP 10
ENV AUDIT_LOG_MAXSIZE 100
ENV AUDIT_LEVEL 0

RUN mkdir -p /usr/share/ranger/ui && \
    cd /usr/share/ranger/ui && \
    curl -sL https://releases.ranger.com/ui/${CATTLE_UI_VERSION}.tar.gz | tar xvzf - --strip-components=1 && \
    mkdir -p assets/ranger-ui-driver-linode && \
    cd assets/ranger-ui-driver-linode && \
    curl -O https://linode.github.io/ranger-ui-driver-linode/releases/${LINODE_UI_DRIVER_VERSION}/component.js && \
    curl -O https://linode.github.io/ranger-ui-driver-linode/releases/${LINODE_UI_DRIVER_VERSION}/component.css && \
    curl -O https://linode.github.io/ranger-ui-driver-linode/releases/${LINODE_UI_DRIVER_VERSION}/linode.svg && \
    mkdir -p /usr/share/ranger/ui/api-ui && \
    cd /usr/share/ranger/ui/api-ui && \
    curl -sL https://releases.ranger.com/api-ui/${CATTLE_API_UI_VERSION}.tar.gz | tar xvzf - --strip-components=1 && \
    mkdir -p /usr/share/ranger/ui-dashboard/dashboard && \
    cd /usr/share/ranger/ui-dashboard/dashboard && \
    curl -sL https://releases.ranger.com/dashboard/${CATTLE_DASHBOARD_UI_VERSION}.tar.gz | tar xvzf - --strip-components=2 && \
    ln -s dashboard/index.html ../index.html && \
    cd ../../ui/assets && \
    curl -sfL https://github.com/ranger/system-agent/releases/download/${CATTLE_SYSTEM_AGENT_VERSION}/ranger-system-agent-arm64 -O && \
    curl -sfL https://github.com/ranger/system-agent/releases/download/${CATTLE_SYSTEM_AGENT_VERSION}/ranger-system-agent-amd64 -O && \
    curl -sfL https://github.com/ranger/system-agent/releases/download/${CATTLE_SYSTEM_AGENT_VERSION}/ranger-system-agent-s390x -O && \
    curl -sfL ${CATTLE_SYSTEM_AGENT_INSTALL_SCRIPT} -o system-agent-install.sh && \
    curl -sfL ${CATTLE_SYSTEM_AGENT_UNINSTALL_SCRIPT} -o system-agent-uninstall.sh && \
    curl -sfL https://github.com/ranger/wins/releases/download/${CATTLE_WINS_AGENT_VERSION}/wins.exe -O && \
    curl -sfL https://acs-mirror.azureedge.net/csi-proxy/${CATTLE_CSI_PROXY_AGENT_VERSION}/binaries/csi-proxy-${CATTLE_CSI_PROXY_AGENT_VERSION}.tar.gz -O && \
    curl -sfL ${CATTLE_WINS_AGENT_INSTALL_SCRIPT} -o wins-agent-install.ps1 \
    curl -sfL ${CATTLE_WINS_AGENT_UNINSTALL_SCRIPT} -o wins-agent-uninstall.ps1

ENV CATTLE_CLI_URL_DARWIN  https://releases.ranger.com/cli2/${CATTLE_CLI_VERSION}/ranger-darwin-amd64-${CATTLE_CLI_VERSION}.tar.gz
ENV CATTLE_CLI_URL_LINUX   https://releases.ranger.com/cli2/${CATTLE_CLI_VERSION}/ranger-linux-amd64-${CATTLE_CLI_VERSION}.tar.gz
ENV CATTLE_CLI_URL_WINDOWS https://releases.ranger.com/cli2/${CATTLE_CLI_VERSION}/ranger-windows-386-${CATTLE_CLI_VERSION}.zip

ARG VERSION=dev
ENV CATTLE_SERVER_VERSION ${VERSION}
COPY run_ranger.sh ranger /usr/bin/
COPY kustomize.sh /usr/bin/
COPY jailer.sh /usr/bin/
COPY k3s-airgap-images.tar /var/lib/ranger/k3s/agent/images/
RUN chmod +x /usr/bin/run_ranger.sh
RUN chmod +x /usr/bin/kustomize.sh

COPY data.json /var/lib/ranger-data/driver-metadata/

ENV CATTLE_AGENT_IMAGE ${IMAGE_REPO}/ranger-agent:${VERSION}
ENV CATTLE_SERVER_IMAGE ${IMAGE_REPO}/ranger
ENV ETCDCTL_API=3

ENV SSL_CERT_DIR /etc/ranger/ssl
VOLUME /var/lib/ranger
VOLUME /var/lib/kubelet
VOLUME /var/lib/cni
VOLUME /var/log

ENV ETCD_UNSUPPORTED_ARCH ${ETCD_UNSUPPORTED_ARCH}

ENTRYPOINT ["run_ranger.sh"]
