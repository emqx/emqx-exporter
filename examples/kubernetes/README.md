The purpose of this tutorial is to show you how to deploy a complete demo with EMQX 5 on Kubernetes. 

## Requirements

+ EMQX Operator
  
  Refer to [Getting Started](https://docs.emqx.com/en/emqx-operator/latest/getting-started/getting-started.html#deploy-emqx-operator) to learn how to deploy the EMQX operator

+ CRDs for prometheus stack

  ```shell
  git clone https://github.com/prometheus-operator/kube-prometheus.git
  cd kube-prometheus
  kubectl apply --server-side -f manifests/setup
  kubectl wait \
  	--for condition=Established \
  	--all CustomResourceDefinition \
  	--namespace=monitoring
  ```

## Deploy example

  ```shell
  kubectl apply -k examples/kubernetes --server-side
  ```
 