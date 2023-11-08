<h1 align="center" style="border-bottom: none">
    EMQX Exporter
</h1>

<div align="center">

[![GitHub Release](https://img.shields.io/github/release/emqx/emqx-exporter?color=brightgreen)](https://github.com/emqx/emqx-exporter/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/emqx/emqx-exporter)](https://hub.docker.com/r/emqx/emqx-exporter)
[![codecov](https://codecov.io/gh/emqx/emqx-exporter/graph/badge.svg?token=XXXYUFHOQR)](https://codecov.io/gh/emqx/emqx-exporter)

</div>

<p align="center">The `emqx-exporter` is designed to expose partial metrics that are not included in the EMQX Prometheus API.</p>
<p align="center">It is compatible with EMQX 4.4 and EMQX 5, both open-source and enterprise.</p>

![Dashboard](https://assets.emqx.com/images/d0529c5355782a6d027de58cce0de69d.png)

## Structure
![Structure](https://assets.emqx.com/images/58adbe48aefb3388f6868854812b18ed.png)

## Metrics
See the documentation [Instruction](grafana-dashboard/template/README.md) for an explanation of the metrics on the dashboard

## Building and running
The `emqx-exporter` listens on HTTP port 8085 by default. See the `--help` output for more options.

### Required
EMQX exporter requires access to the EMQX dashboard API with basic auth, so you need to sign in to the dashboard to create an API secret
Note that it is different to create a secret between EMQX 5 and EMQX 4.4 on the dashboard.

* **EMQX 5**

  * Create a new [API KEY](https://www.emqx.io/docs/en/v5.0/dashboard/system.html#api-keys).

* **EMQX 4.4**

  * Create a new `User` instead of `Application`

  + Make sure the `emqx_prometheus` plugin has been started on all nodes, check it one by one on the dashboard <http://your_cluster_addr:18083/#/plugins>.

### Build

    make build

### Running

    ./bin/emqx-exporter <flags>

### Docker Compose

Refer to the [example](examples/docker-compose) to deploy a complete demo by docker compose.

### Kubernetes

Refer to the [example](examples/kubernetes/README.md) to learn how to deploy `emqx-exporter` on the Kubernetes.

## Configuration

Sample config file like this

```
metrics:
  target: 127.0.0.1:18083
  api_key: "some_api_key"
  api_secret: "some_api_secret"
probes:
  - target: 127.0.0.1:1883
```

The `metrics` and the `probes` are not required configuration items, if not set `metrics`, the metrics feature will disable, and if not set `probes`, the probe feature will disable.

## Prometheus Config

The scrape config below is available for EMQX 5

```yaml
scrape_configs:
- job_name: 'emqx-self-metrics'
  metrics_path: /api/v5/prometheus/stats
  scrape_interval: 5s
  honor_labels: true
  static_configs:
    # a list of addresses of all EMQX nodes
    - targets: [${your_emqx_addr}:18083]
      labels:
        # label the cluster name of where the metrics data from
        cluster: ${your_emqx_addr}
        # fix value, don't modify
        from: emqx
- job_name: 'exporter-metrics'
  metrics_path: /metrics
  scrape_interval: 5s
  static_configs:
    - targets: [${your_exporter_addr}:8085]
      labels:
        # label the cluster name of where the metrics data from
        cluster: ${your_cluster_name}
        # fix value, don't modify
        from: exporter
- job_name: 'exporter-probe'
  metrics_path: /probes
  params:
    target:
      # must equal the `probes[$index].taget` in config file
      - "127.0.0.1:1883"
  scrape_interval: 5s
  static_configs:
    - targets: [${your_exporter_addr}:8085]
      labels:
        # label the cluster name of where the metrics data from
        cluster: ${your_cluster_name}
        # fix value, don't modify
        from: exporter
```

## Grafana Dashboard
Import all [templates](./grafana-dashboard/template) to your Grafana, then browse the dashboard `EMQX` and enjoy yourself!

The templates of dashboard ares JSON files, about how to upload a dashboard JSON file, you can check out [here](https://grafana.com/docs/grafana/latest/dashboards/manage-dashboards/#import-a-dashboard).

## TLS endpoint

**EXPERIMENTAL**

The exporter supports TLS via a new web configuration file.

```console
./emqx-exporter --web.config.file=web-config.yml
```

See the [exporter-toolkit https package](https://github.com/prometheus/exporter-toolkit/blob/v0.1.0/https/README.md) for more details.
