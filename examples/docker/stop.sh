#!/bin/bash

docker stop emqx-demo exporter-demo prometheus-demo grafana-demo
docker rm emqx-demo exporter-demo prometheus-demo grafana-demo
