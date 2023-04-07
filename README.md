# EMQX exporter

Prometheus exporter for EMQX cluster metrics exposed

## Installation and Usage
The `emqx-exporter` listens on HTTP port 8085 by default. See the `--help` output for more options.

### API Secret
The `emqx-exporter` is designed to export partial metrics that doesn't include in the EMQX prometheus API.
It requires access to the dashboard API with basic auth, so you need to sign in dashboard to create an API secret,
then pass the API key and Secret Key to startup argument as username and password.

### Docker
```bash
docker run -d \
  -p 8085:8085 \
  emqx-exporter:latest \
  --emqx.nodes="your_cluster_addr:18083"  \
  --emqx.auth-username="apiKey" \
  --emqx.auth-password="secretKey"
```

The arg `emqx.nodes` is a host list, the exporter will choose one to establish connection.

For excluding metrics about exporter itself, add a flag `--web.disable-exporter-metrics`.

```yaml
version: '3.8'

services:
  emqx-exporter:
    image: emqx-exporter:latest
    container_name: emqx-exporter
    command:
      - '--emqx.nodes=your_cluster_addr:18083'
      - '--emqx.auth-username=apiKey'
      - '--emqx.auth-password=secretKey'
    restart: unless-stopped
```

### Kubernetes
```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: emqx-exporter
  name: emqx-exporter-metrics-service
spec:
  ports:
  - name: metrics
    port: 8085
    targetPort: metrics
  selector:
    app: emqx-exporter
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: emqx-exporter
  labels:
    app: emqx-exporter
spec:
  selector:
    matchLabels:
      app: emqx-exporter
  replicas: 1
  template:
    metadata:
      labels:
        app: emqx-exporter
    spec:
      securityContext:
        runAsNonRoot: true
      containers:
      - name: exporter
        image: emqx-exporter:latest
        imagePullPolicy: IfNotPresent
        args:
        - --emqx.nodes=your_cluster_addr1:18083,your_cluster_addr2:18083
        - --emqx.auth-username=apiKey
        - --emqx.auth-password=secretKey        
        securityContext:
          allowPrivilegeEscalation: false
        ports:
        - containerPort: 8085
          name: metrics
          protocol: HTTP
        resources:
          limits:
            cpu: 100m
            memory: 100Mi
          requests:
            cpu: 100m
            memory: 20Mi
```

## Prometheus Config
For EMQX5 version, make sure the EMQX cluster has exposed metrics by prometheus, check it in dashboard(http://your_cluster_addr:18083/#/monitoring/integration).  

__Note that disable the prometheus push mode(Pushgateway)__  

Scrape Config:
```yaml
scrape_configs:
- job_name: 'emqx'
  metrics_path: /api/v5/prometheus/stats
  scrape_interval: 5s
  honor_labels: true
  static_configs:
    # EMQX IP address and port
    - targets: [emqx-enterprise:18083,emqx-enterprise2:18083]
      labels:
        # label the cluster name of where the metrics data from
        cluster: emqx_name
        # fix value
        from: emqx
  relabel_configs:
    - source_labels: ["__address__"]
      target_label: "instance"
      regex: (.*):.*
      replacement: $1
- job_name: 'exporter'
  metrics_path: /metrics
  scrape_interval: 5s
  static_configs:
    - targets: [192.168.1.1:8085]
      labels:
        # label the cluster name of where the metrics data from
        cluster: emqx_name
        # fix value
        from: exporter
```

For EMQX4 version, the EMQX cluster has exposed metrics by prometheus by default.

Scrape Config:
```yaml
scrape_configs:
- job_name: 'emqx'
  metrics_path: /api/v4/emqx_prometheus?type=prometheus
  scrape_interval: 5s
  honor_labels: true
  static_configs:
    # EMQX IP address and port
    - targets: [emqx-enterprise:18083,emqx-enterprise2:18083]
      labels:
        # label the cluster name of where the metrics data from
        cluster: emqx_name
        # fix value
        from: emqx
  relabel_configs:
    - source_labels: ["__address__"]
      target_label: "instance"
      regex: (.*):.*
      replacement: $1
- job_name: 'exporter'
  metrics_path: /metrics
  scrape_interval: 5s
  static_configs:
    - targets: [192.168.1.1:8085]
      labels:
        # label the cluster name of where the metrics data from
        cluster: emqx_name
        # fix value
        from: exporter
```

If you deployed prometheus by [operator](https://prometheus-operator.dev/), then you need to create two service monitor for
exporter to add a scrape job to prometheus config.

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: emqx-dashboard
  labels:
    app: emqx-dashboard
spec:
  selector:
    matchLabels:
      # the label in emqx dashboard svc
      app: emqx-dashboard
  endpoints:
    - port: http
      interval: 15s
      relabelings:
        - action: replace
          replacement: emqx_name
          targetLabel: cluster
        - action: replace
          replacement: emqx
          targetLabel: from
  namespaceSelector:
  #matchNames:
  #  - default
  
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: emqx-exporter
  labels:
    app: emqx-exporter
spec:
  selector:
    matchLabels:
      app: emqx-exporter
  endpoints:
    - port: http
      interval: 15s
      relabelings:
      - action: replace
        replacement: emqx_name
        targetLabel: cluster
      - action: replace
        replacement: exporter
        targetLabel: from
  namespaceSelector:
    #matchNames:
    #  - default
```

## Grafana Dashboard

## Development building and running

Prerequisites:

* [Go compiler](https://golang.org/dl/)
* RHEL/CentOS: `glibc-static` package.

Building:

    git clone https://github.com/emqx/emqx-exporter.git
    cd emqx-exporter
    make build
    # cd to output folder `./build/$OS_$ARCH/`
    ./emqx-exporter <flags>

To see all available configuration flags:

    ./emqx-exporter -h

## TLS endpoint

** EXPERIMENTAL **

The exporter supports TLS via a new web configuration file.

```console
./emqx-exporter --web.config.file=web-config.yml
```

See the [exporter-toolkit https package](https://github.com/prometheus/exporter-toolkit/blob/v0.1.0/https/README.md) for more details.
