By installing this application, you accept the [End User License Agreement & Terms & Conditions](https://www.suse.com/licensing/eula/).

# Ranger

***Ranger*** is open source software that combines everything an organization needs to adopt and run containers in production. Built on Kubernetes, Ranger makes it easy for DevOps teams to test, deploy and manage their applications.

### Introduction

This chart bootstraps a [Ranger Server](https://rangermanager.docs.ranger.com/pages-for-subheaders/install-upgrade-on-a-kubernetes-cluster) on a Kubernetes cluster using the [Helm](https://helm.sh/) package manager. For a Ranger Supported Deployment please follow our [HA install instructions](https://rangermanager.docs.ranger.com/how-to-guides/new-user-guides/kubernetes-cluster-setup/high-availability-installs).


### Prerequisites Details

*For installations covered under [Ranger Support SLA](https://www.suse.com/suse-ranger/support-matrix/all-supported-versions) the target cluster must be **[RKE1](https://rangermanager.docs.ranger.com/how-to-guides/new-user-guides/kubernetes-cluster-setup/rke1-for-ranger)**, **[RKE2](https://rangermanager.docs.ranger.com/how-to-guides/new-user-guides/kubernetes-cluster-setup/rke2-for-ranger)**, **[K3s](https://rangermanager.docs.ranger.com/how-to-guides/new-user-guides/kubernetes-cluster-setup/k3s-for-ranger)**, **[AKS](https://rangermanager.docs.ranger.com/getting-started/installation-and-upgrade/install-upgrade-on-a-kubernetes-cluster/ranger-on-aks)**, **[EKS](https://rangermanager.docs.ranger.com/getting-started/installation-and-upgrade/install-upgrade-on-a-kubernetes-cluster/ranger-on-amazon-eks)**, or **[GKE](https://rangermanager.docs.ranger.com/getting-started/installation-and-upgrade/install-upgrade-on-a-kubernetes-cluster/ranger-on-gke)**.*

Make sure the node(s) for the Ranger server fulfill the following requirements:

[Operating Systems and Container Runtime Requirements](https://rangermanager.docs.ranger.com/pages-for-subheaders/installation-requirements#operating-systems-and-container-runtime-requirements)
[Hardware Requirements](https://rangermanager.docs.ranger.com/pages-for-subheaders/installation-requirements#hardware-requirements)

- [CPU and Memory](https://rangermanager.docs.ranger.com/pages-for-subheaders/installation-requirements#cpu-and-memory)
- [Ingress](https://rangermanager.docs.ranger.com/pages-for-subheaders/installation-requirements#ingress)
- [Disks](https://rangermanager.docs.ranger.com/pages-for-subheaders/installation-requirements#disks)

[Networking Requirements](https://rangermanager.docs.ranger.com/pages-for-subheaders/installation-requirements#networking-requirements)
- [Node IP Addresses](https://rangermanager.docs.ranger.com/pages-for-subheaders/installation-requirements#node-ip-addresses)
- [Port Requirements](https://rangermanager.docs.ranger.com/pages-for-subheaders/installation-requirements#port-requirements)

[Install the Required CLI Tools](https://rangermanager.docs.ranger.com/pages-for-subheaders/cli-with-ranger)

- [kubectl](https://rangermanager.docs.ranger.com/reference-guides/cli-with-ranger/kubectl-utility) - Kubernetes command-line tool.
- [helm](https://docs.helm.sh/using_helm/#installing-helm) - Package management for Kubernetes. Refer to the [Helm version requirements](https://rangermanager.docs.ranger.com/getting-started/installation-and-upgrade/resources/helm-version-requirements) to choose a version of Helm to install Ranger.

For a list of best practices that we recommend for running the Ranger server in production, refer to the [best practices section](https://rangermanager.docs.ranger.com/pages-for-subheaders/best-practices).

## Installing Ranger

For production environments, we recommend installing Ranger in a [high-availability Kubernetes installation](https://rangermanager.docs.ranger.com/how-to-guides/new-user-guides/kubernetes-cluster-setup/high-availability-installs) so that your user base can always access Ranger Server. When installed in a Kubernetes cluster, Ranger will integrate with the cluster’s etcd database and take advantage of Kubernetes scheduling for high-availability.

Optional: Installing Ranger on a [Single-node](https://rangermanager.docs.ranger.com/pages-for-subheaders/ranger-on-a-single-node-with-docker) Kubernetes Cluster

#### Add the Helm Chart Repository

Use [helm repo add](https://helm.sh/docs/helm/helm_repo_add/) command to add the Helm chart repository that contains charts to install Ranger. For more information about the repository choices and which is best for your use case, see Choosing a Version of Ranger.

```bash
helm repo add ranger-latest https://releases.ranger.com/server-charts/latest
```

#### Create a Namespace for Ranger

We’ll need to define a Kubernetes namespace where the resources created by the Chart should be installed. This should always be cattle-system:

```bash
kubectl create namespace cattle-system
```

#### Choose your SSL Configuration

The Ranger management server is designed to be secure by default and requires SSL/TLS configuration.

There are three recommended options for the source of the certificate used for TLS termination at the Ranger server:

- [Ranger-generated TLS certificate](https://rangermanager.docs.ranger.com/pages-for-subheaders/install-upgrade-on-a-kubernetes-cluster#3-choose-your-ssl-configuration)
- [Let’s Encrypt](https://rangermanager.docs.ranger.com/pages-for-subheaders/install-upgrade-on-a-kubernetes-cluster#3-choose-your-ssl-configuration)
- [Bring your own certificate](https://rangermanager.docs.ranger.com/pages-for-subheaders/install-upgrade-on-a-kubernetes-cluster#3-choose-your-ssl-configuration)

#### Install cert-manager

This step is only required to use certificates issued by Ranger’s generated CA **`(ingress.tls.source=ranger)`** or to request Let’s Encrypt issued certificates **`(ingress.tls.source=letsEncrypt)`**.

[These instructions are adapted from the official cert-manager documentation.](https://rangermanager.docs.ranger.com/pages-for-subheaders/install-upgrade-on-a-kubernetes-cluster#4-install-cert-manager)

#### Install Ranger with Helm and Your Chosen Certificate Option

- [Ranger to generated certificates](https://rangermanager.docs.ranger.com/pages-for-subheaders/install-upgrade-on-a-kubernetes-cluster#5-install-ranger-with-helm-and-your-chosen-certificate-option)
```bash
helm install ranger ranger-latest/ranger \
  --namespace cattle-system \
  --set hostname=ranger.my.org
```

- [Let’s Encrypt](https://rangermanager.docs.ranger.com/pages-for-subheaders/install-upgrade-on-a-kubernetes-cluster#5-install-ranger-with-helm-and-your-chosen-certificate-option)

```bash
helm install ranger ranger-latest/ranger \
  --namespace cattle-system \
  --set hostname=ranger.my.org \
  --set ingress.tls.source=letsEncrypt \
  --set letsEncrypt.email=me@example.org
```

- [Certificates from Files](https://rangermanager.docs.ranger.com/pages-for-subheaders/install-upgrade-on-a-kubernetes-cluster#5-install-ranger-with-helm-and-your-chosen-certificate-option)

```bash
helm install ranger ranger-latest/ranger \
  --namespace cattle-system \
  --set hostname=ranger.my.org \
  --set ingress.tls.source=secret
```

*If you are using a Private CA signed certificate , add **--set privateCA=true** to the command:`*

```bash
helm install ranger ranger-latest/ranger \
  --namespace cattle-system \
  --set hostname=ranger.my.org \
  --set ingress.tls.source=secret \
  --set privateCA=true
```

#### Verify that the Ranger Server is Successfully Deployed

After adding the secrets, check if Ranger was rolled out successfully:

```bash
kubectl -n cattle-system rollout status deploy/ranger
Waiting for deployment "ranger" rollout to finish: 0 of 3 updated replicas are available...
deployment "ranger" successfully rolled out
```

If you see the following **`error: error: deployment "ranger" exceeded its progress deadline`**, you can check the status of the deployment by running the following command:

```bash
kubectl -n cattle-system get deploy ranger
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
ranger   3         3         3            3           3m
```

It should show the same count for **`DESIRED`** and **`AVAILABLE`**.

#### Save Your Options

Make sure you save the **`--set`** options you used. You will need to use the same options when you upgrade Ranger to new versions with Helm.

#### Finishing Up

That’s it. You should have a functional Ranger server.

In a web browser, go to the DNS name that forwards traffic to your load balancer. Then you should be greeted by the colorful login page.

Doesn’t work? Take a look at the [Troubleshooting Page](https://rangermanager.docs.ranger.com/troubleshooting/general-troubleshooting)

***All of these instructions are defined in detailed in the [Ranger Documentation](https://rangermanager.docs.ranger.com/pages-for-subheaders/install-upgrade-on-a-kubernetes-cluster#install-the-ranger-helm-chart).***

### Helm Chart Options for Kubernetes Installations

The full [Helm Chart Options](https://rangermanager.docs.ranger.com/getting-started/installation-and-upgrade/installation-references/helm-chart-options) can be found here.

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

#### Common Options

| Parameter                 | Default Value | Description                                                                                  |
| ------------------------- | ------------- | -------------------------------------------------------------------------------------------- |
| `hostname`                | " "           | ***string*** - the Fully Qualified Domain Name for your Ranger Server                       |
| `ingress.tls.source`      | "ranger"     | ***string*** - Where to get the cert for the ingress. - "***ranger, letsEncrypt, secret***" |
| `letsEncrypt.email`       | " "           | ***string*** - Your email address                                                            |
| `letsEncrypt.environment` | "production"  | ***string*** - Valid options: "***staging, production***"                                    |
| `privateCA`               | false         | ***bool*** - Set to true if your cert is signed by a private CA                              |

#### Advanced Options

| Parameter                                | Default Value                                                             | Description                                                                                                                                                                                                                                                                             |
| ---------------------------------------- | ------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `additionalTrustedCAs`                   | false                                                                     | ***bool*** - [See Additional Trusted CAs Server](https://rangermanager.docs.ranger.com/getting-started/installation-and-upgrade/installation-references/helm-chart-options#additional-trusted-cas)                                                                                                                                   |
| `addLocal`                               | "true"                                                                    | ***string*** - As of Ranger v2.5.0 this flag is deprecated and must be set to "true"                                                                                                                                                                                                   |
| `antiAffinity`                           | "preferred"                                                               | ***string*** - AntiAffinity rule for Ranger pods - *"preferred, required"*                                                                                                                                                                                                             |
| `replicas`                               | 3                                                                         | ***int*** - Number of replicas of Ranger pods                                                                                                                                                                                                                                          |
| `auditLog.destination`                   | "sidecar"                                                                 | ***string*** - Stream to sidecar container console or hostPath volume - *"sidecar, hostPath"*                                                                                                                                                                                           |
| `auditLog.hostPath`                      | "/var/log/ranger/audit"                                                  | ***string*** - log file destination on host (only applies when **auditLog.destination** is set to **hostPath**)                                                                                                                                                                         |
| `auditLog.level`                         | 0                                                                         | ***int*** - set the [API Audit Log level](https://rangermanager.docs.ranger.com/how-to-guides/advanced-user-guides/enable-api-audit-log#audit-log-levels). 0 is off. [0-3]                                                                                                                                                          |
| `auditLog.maxAge`                        | 1                                                                         | ***int*** - maximum number of days to retain old audit log files (only applies when **auditLog.destination** is set to **hostPath**)                                                                                                                                                    |
| `auditLog.maxBackup`                     | 1                                                                         | int - maximum number of audit log files to retain (only applies when **auditLog.destination** is set to **hostPath**)                                                                                                                                                                   |
| `auditLog.maxSize`                       | 100                                                                       | ***int*** - maximum size in megabytes of the audit log file before it gets rotated (only applies when **auditLog.destination** is set to **hostPath**)                                                                                                                                  |
| `auditLog.image.repository`              | "ranger/mirrored-bci-micro"                                              | ***string*** - Location for the image used to collect audit logs *Note: Available as of v2.7.0*                                                                                                                                                                                         |
| `auditLog.image.tag`                     | "15.4.14.3"                                                               | ***string*** - Tag for the image used to collect audit logs *Note: Available as of v2.7.0*                                                                                                                                                                                              |
| `auditLog.image.pullPolicy`              | "IfNotPresent"                                                            | ***string*** - Override imagePullPolicy for auditLog images - *"Always", "Never", "IfNotPresent"* *Note: Available as of v2.7.0*                                                                                                                                                        |
| `busyboxImage`                           | ""                                                                        | ***string*** - *Deprecated `auditlog.image.repository` should be used to control auditing sidecar image.* Image location for busybox image used to collect audit logs *Note: Available as of v2.2.0, and  Deprecated as of v2.7.0*                                                      |
| `busyboxImagePullPolicy`                 | "IfNotPresent"                                                            | ***string*** - - *Deprecated `auditlog.image.pullPolicy` should be used to control auditing sidecar image.* Override imagePullPolicy for busybox images - *"Always", "Never", "IfNotPresent"* *Deprecated as of v2.7.0*                                                                 |
| `debug`                                  | false                                                                     | ***bool*** - set debug flag on ranger server                                                                                                                                                                                                                                           |
| `certmanager.version`                    | " "                                                                       | ***string*** - set cert-manager compatibility                                                                                                                                                                                                                                           |
| `extraEnv`                               | []                                                                        | ***list*** - set additional environment variables for Ranger Note: *Available as of v2.2.0*                                                                                                                                                                                            |
| `imagePullSecrets`                       | []                                                                        | ***list*** - list of names of Secret resource containing private registry credentials                                                                                                                                                                                                   |
| `ingress.enabled`                        | true                                                                      | ***bool*** - install ingress resource                                                                                                                                                                                                                                                   |
| `ingress.ingressClassName`               | " "                                                                       | ***string*** - class name of ingress if not set manually or by the ingress controller's defaults                                                                                                                                                                                        |
| `ingress.includeDefaultExtraAnnotations` | true                                                                      | ***bool*** - Add default nginx annotations                                                                                                                                                                                                                                              |
| `ingress.extraAnnotations`               | {}                                                                        | ***map*** - additional annotations to customize the ingress                                                                                                                                                                                                                             |
| `ingress.configurationSnippet`           | " "                                                                       | ***string*** - Add additional Nginx configuration. Can be used for proxy configuration. Note: *Available as of v2.0.15, v2.1.10 and v2.2.4*                                                                                                                                             |
| `service.annotations`                    | {}                                                                        | ***map*** - annotations to customize the service                                                                                                                                                                                                                                        |
| `service.type`                           | " "                                                                       | ***string*** - Override the type used for the service - *"NodePort", "LoadBalancer", "ClusterIP"*                                                                                                                                                                                       |
| `letsEncrypt.ingress.class`              | " "                                                                       | ***string*** - optional ingress class for the cert-manager acmesolver ingress that responds to the Let’s *Encrypt ACME challenges*                                                                                                                                                      |
| `proxy`                                  | " "                                                                       | ***string** - HTTP[S] proxy server for Ranger                                                                                                                                                                                                                                          |
| `noProxy`                                | "127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,.svc,.cluster.local" | ***string*** - comma separated list of hostnames or ip address not to use the proxy                                                                                                                                                                                                     |
| `resources`                              | {}                                                                        | ***map*** - ranger pod resource requests & limits                                                                                                                                                                                                                                      |
| `rangerImage`                           | "ranger/ranger"                                                         | ***string*** - ranger image source                                                                                                                                                                                                                                                     |
| `rangerImageTag`                        | same as chart version                                                     | ***string*** - ranger/ranger image tag                                                                                                                                                                                                                                                |
| `rangerImagePullPolicy`                 | "IfNotPresent"                                                            | ***string*** - Override imagePullPolicy for ranger server images - *"Always", "Never", "IfNotPresent"*                                                                                                                                                                                 |
| `tls`                                    | "ingress"                                                                 | ***string*** - See External TLS Termination for details. - *"ingress, external"*                                                                                                                                                                                                        |
| `systemDefaultRegistry`                  | ""                                                                        | ***string*** - private registry to be used for all system Docker images, e.g., [http://registry.example.com/] *Available as of v2.3.0*                                                                                                                                                  |
| `useBundledSystemChart`                  | false                                                                     | ***bool*** - select to use the system-charts packaged with Ranger server. This option is used for air gapped installations.  *Available as of v2.3.0*                                                                                                                                  |
| `customLogos.enabled`                    | false                                                                     | ***bool*** - Enabled [Ember Ranger UI (cluster manager) custom logos](https://github.com/ranger/ui/tree/master/public/assets/images/logos) and [Vue Ranger UI (cluster explorer) custom logos](https://github.com/ranger/dashboard/tree/master/assets/images/pl) persistence volume |
| `customLogos.volumeSubpaths.emberUi`     | "ember"                                                                   | ***string*** - Volume subpath for [Ember Ranger UI (cluster manager) custom logos](https://github.com/ranger/ui/tree/master/public/assets/images/logos) persistence                                                                                                                   |
| `customLogos.volumeSubpaths.vueUi`       | "vue"                                                                     | ***string*** - Volume subpath for [Vue Ranger UI (cluster explorer) custom logos](https://github.com/ranger/dashboard/tree/master/assets/images/pl) persistence                                                                                                                       |
| `customLogos.volumeName`                 | ""                                                                        | ***string*** - Use an existing volume. Custom logos should be copied to the proper `volume/subpath` folder by the user. Optional for persistentVolumeClaim, required for configMap                                                                                                      |
| `customLogos.storageClass`               | ""                                                                        | ***string*** - Set custom logos persistentVolumeClaim storage class. Required for dynamic pv                                                                                                                                                                                            |
| `customLogos.accessMode`                 | "ReadWriteOnce"                                                           | ***string*** - Set custom persistentVolumeClaim access mode                                                                                                                                                                                                                             |
| `customLogos.size`                       | "1Gi"                                                                     | ***string*** - Set custom persistentVolumeClaim size                                                                                                                                                                                                                                    |
