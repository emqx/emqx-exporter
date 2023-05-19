#!/bin/bash
echo "stop containers"
docker stop emqx-demo exporter-demo prometheus-demo grafana-demo

sleep 1s
echo "rm containers"
docker rm emqx-demo exporter-demo prometheus-demo grafana-demo
