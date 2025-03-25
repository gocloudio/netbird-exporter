# NetBird Exporter Helm Chart

A Helm chart for deploying the NetBird Prometheus exporter to Kubernetes as a DaemonSet.

## Introduction

This chart deploys the NetBird exporter as a DaemonSet, which exports metrics from NetBird peer connections for Prometheus monitoring. The exporter runs on each node where NetBird is installed.

## Prerequisites

- Kubernetes 1.16+
- Helm 3.0+
- NetBird daemon running on the cluster nodes

## Installing the Chart

### From Helm Repository

Add the repository:

```bash
helm repo add netbird-exporter https://gocloudio.github.io/netbird-exporter
helm repo update
```

Install the chart:

```bash
helm install netbird-exporter netbird-exporter/netbird-exporter
```

### From OCI Registry (Recommended)

This chart is also available as an OCI artifact. You can install it directly from the GitHub Container Registry:

```bash
# Helm 3.8.0 or later
helm install netbird-exporter oci://ghcr.io/gocloudio/netbird-exporter/charts/netbird-exporter --version 0.1.0
```

You can also pull the chart first:

```bash
helm pull oci://ghcr.io/gocloudio/netbird-exporter/charts/netbird-exporter --version 0.1.0
helm install netbird-exporter ./netbird-exporter-0.1.0.tgz
```

## Configuration

The following table lists the configurable parameters of the NetBird exporter chart and their default values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Image repository | `ghcr.io/gocloudio/netbird-exporter` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `image.tag` | Image tag | `latest` |
| `serviceAccount.create` | Create service account | `true` |
| `serviceAccount.annotations` | Service account annotations | `{}` |
| `serviceAccount.name` | Service account name | `""` |
| `podAnnotations` | Pod annotations | `{}` |
| `podSecurityContext` | Pod security context | `{}` |
| `securityContext` | Container security context | `{}` |
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `9101` |
| `hostNetwork` | Use host network | `false` |
| `resources` | Pod resources | `{}` |
| `nodeSelector` | Node labels for pod assignment | `{}` |
| `tolerations` | Tolerations for pod assignment | `[]` |
| `affinity` | Node/pod affinities | `{}` |
| `netbird.socket` | NetBird socket path | `unix:///var/run/netbird.sock` |
| `netbird.ignoreFQDNDomain` | Domain suffix to remove from FQDNs | `""` |
| `hostPathSocketDir` | Host path for NetBird socket | `/var/run` |
| `serviceMonitor.enabled` | Create ServiceMonitor | `false` |
| `serviceMonitor.additionalLabels` | ServiceMonitor labels | `{}` |
| `serviceMonitor.interval` | ServiceMonitor scrape interval | `30s` |
| `serviceMonitor.scrapeTimeout` | ServiceMonitor scrape timeout | `10s` |
| `serviceMonitor.endpointRelabelings` | ServiceMonitor relabelings | `[{"sourceLabels": ["__meta_kubernetes_pod_node_name"], "targetLabel": "node"}]` |
| `serviceMonitor.implementation` | ServiceMonitor implementation type | `prometheus-operator` |

## Node Affinity and Tolerations

Since the NetBird exporter needs to run only on nodes with NetBird installed, you may want to use node selectors or tolerations to target specific nodes.

Example for targeting nodes with the `netbird=true` label:

```yaml
nodeSelector:
  netbird: "true"
```

If you need to run on all nodes, including master/control-plane nodes:

```yaml
tolerations:
  - key: "node-role.kubernetes.io/master"
    operator: "Exists"
    effect: "NoSchedule"
  - key: "node-role.kubernetes.io/control-plane"
    operator: "Exists"
    effect: "NoSchedule"
```

## Prometheus Integration

### Standard Prometheus

Add a scrape config to your Prometheus:

```yaml
scrape_configs:
  - job_name: 'netbird'
    kubernetes_sd_configs:
      - role: endpoints
        namespaces:
          names:
            - <namespace>
    relabel_configs:
      - source_labels: [__meta_kubernetes_service_label_app_kubernetes_io_name]
        regex: netbird-exporter
        action: keep
      - source_labels: [__meta_kubernetes_pod_node_name]
        target_label: node
```

### Prometheus Operator

Enable the ServiceMonitor:

```yaml
serviceMonitor:
  enabled: true
  additionalLabels:
    release: prometheus-stack
```

### VictoriaMetrics Operator

The chart also supports VictoriaMetrics Operator for monitoring:

```yaml
serviceMonitor:
  enabled: true
  implementation: "victoriametrics"
  additionalLabels:
    release: victoria-metrics-stack
```

## Accessing Metrics

Since this is deployed as a DaemonSet, each node will have its own exporter pod. To access metrics from a specific pod:

```bash
# Get pod name
kubectl get pods -l "app.kubernetes.io/name=netbird-exporter" -o wide

# Forward port from a specific pod
kubectl port-forward <pod-name> 9101:9101

# Access metrics
curl http://localhost:9101/metrics
``` 