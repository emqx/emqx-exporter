The purpose of this tutorial is to show you how to deploy a complete demo with EMQX 5 on Kubernetes. 

## Install EMQX-Operator
Refer to [Getting Started](https://docs.emqx.com/en/emqx-operator/latest/getting-started/getting-started.html#deploy-emqx-operator) to learn how to deploy the EMQX operator

## Deploy EMQX Cluster
```shell
cat << "EOF" | kubectl apply -f -
apiVersion: apps.emqx.io/v2alpha1
kind: EMQX
metadata:
  name: emqx
spec:
  image: emqx/emqx-enterprise:5.0.1
  coreTemplate:
    spec:
      replicas: 1
      ports:
        # prometheus monitor requires the pod must name the target port 
        - containerPort: 18083
          name: dashboard
  replicantTemplate:
    spec:
      replicas: 1
      ports:
        # prometheus monitor requires the pod must name the target port
        - containerPort: 18083
          name: dashboard
EOF
```

If you are deploying EMQX 4.4 open-source, you need to enable plugin `emqx_prometheus` by `EmqxPlugin` CRD:

```shell
cat << "EOF" | kubectl apply -f -
apiVersion: apps.emqx.io/v1beta4
kind: EmqxPlugin
metadata:
  name: emqx-prometheus
spec:
  selector:
    # EMQX pod labels
    apps.emqx.io/instance: emqx
    apps.emqx.io/managed-by: emqx-operator
  # enable plugin emqx_prometheus
  pluginName: emqx_prometheus
EOF
```

## Create API secret
emqx-exporter and Prometheus will pull metrics from EMQX dashboard API, so you need to sign in to dashboard to create an API secret.

It is different to create a secret between EMQX 5 and EMQX 4.4 on the dashboard.

* **EMQX 5** create a new [API KEY](https://www.emqx.io/docs/en/v5.0/dashboard/system.html#api-keys).
* **EMQX 4.4** create a new `User` instead of `Application`

## Deploy Exporter
You need to sign in to EMQX dashboard to create an API secret, then pass the API key and secret to the exporter startup argument as username and password.

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: emqx-exporter
  name: emqx-exporter-service
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
        runAsUser: 1000
      containers:
        - name: exporter
          image: emqx-exporter:latest
          imagePullPolicy: IfNotPresent
          args:
            # "emqx-dashboard-service-name" is the service name that creating by operator for exposing 18083 port 
            - --emqx.nodes=${emqx-dashboard-service-name}:18083
            - --emqx.auth-username=${paste_your_new_api_key_here}
            - --emqx.auth-password=${paste_your_new_secret_here}
          securityContext:
            allowPrivilegeEscalation: false
            runAsNonRoot: true
          ports:
            - containerPort: 8085
              name: metrics
              protocol: TCP
          resources:
            limits:
              cpu: 100m
              memory: 100Mi
            requests:
              cpu: 100m
              memory: 20Mi
```

> Set the arg "--emqx.nodes" to the service name that creating by operator for exposing 18083 port. Check out the service name by call `kubectl get svc`.   

Save the yaml content to file `emqx-exporter.yaml`, paste your new creating API key and secret, then apply it
```shell
kubectl apply -f emqx-exporter.yaml
```

Check the status of emqx-exporter pod。
```bash
$ kubectl get po -l="app=emqx-exporter"

NAME      STATUS   AGE
emqx-exporter-856564c95-j4q5v   Running  8m33s
```

## Import Prometheus Scrape Config
Assuming that you have deployed Prometheus by [Prometheus Operator](https://prometheus-operator.dev/) or [Kube-Prometheus](https://github.com/prometheus-operator/kube-prometheus), you need to add scrape config by defining `PodMonitor` and `ServiceMonitor` CR.  

[Here](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/getting-started.md) is a sample example of `PodMonitor` and `ServiceMonitor`.  
In addition, you can use cmd `kubectl explain` to see the comment about the CR spec. 

In most cases, it's easier to deploy Prometheus by `Deployment` without the operator if you are new to this, and you can get the scrape config example from [here](../docker) 

The yaml below is available for EMQX 5, you can check out the [template](./template_monitor_emqx4.yaml) for EMQX 4. 

```shell
cat << "EOF" | kubectl apply -f -
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: emqx-metrics
  labels:
    app: emqx-metrics
spec:
  selector:
    matchLabels:
      # the label is the same as the label of emqx pod
      apps.emqx.io/instance: emqx
      apps.emqx.io/managed-by: emqx-operator
  podMetricsEndpoints:
    # the name of emqx dashboard containerPort
    - port: dashboard
      interval: 5s
      path: /api/v5/prometheus/stats
      relabelings:
        - action: replace
          # user-defined cluster name, requires unique
          replacement: emqx5.0
          targetLabel: cluster
        - action: replace
          # fix value, don't modify
          replacement: emqx
          targetLabel: from
        - action: replace
          # fix value, don't modify
          sourceLabels: ['pod']
          targetLabel: "instance"
  namespaceSelector:
    # modify the namespace if your EMQX cluster deployed in other namespace
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
      # the label is the same as the label of emqx-exporter svc
      app: emqx-exporter
  endpoints:
    - port: metrics
      interval: 5s
      path: /metrics
      relabelings:
        - action: replace
          # user-defined cluster name, requires unique
          replacement: emqx5.0
          targetLabel: cluster
        - action: replace
          # fix value, don't modify
          replacement: exporter
          targetLabel: from
        - action: replace
          # fix value, don't modify
          sourceLabels: ['pod']
          regex: '(.*)-.*-.*'
          replacement: $1
          targetLabel: "instance"
        - action: labeldrop
          # fix value, don't modify
          regex: 'pod'
  namespaceSelector:
    # modify the namespace if your exporter deployed in other namespace
    #matchNames:
    #  - default
EOF
```

## Load Grafana Templates
Import all [templates](../../config/grafana-template) to your Grafana, then browse the dashboard EMQX and enjoy yourself!

The templates of dashboard ares JSON files, about how to upload a dashboard JSON file, you can check out [here](https://grafana.com/docs/grafana/latest/dashboards/manage-dashboards/#import-a-dashboard). 