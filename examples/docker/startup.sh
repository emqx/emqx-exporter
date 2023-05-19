#!/bin/bash

emqxDashboardUsername=admin
emqxDashboardPassword=admin

# which EMQX version your want to deploy?
emqxVersion=$1
case $emqxVersion in
emqx4)
  prometheusConfig=prometheus-emqx4.yaml
  dashboardDefinitions=$(dirname $(dirname $(pwd)))/config/grafana-template/EMQX4
 emqxImage=emqx/emqx:4.4.16
  ;;
emqx4-ee)
  prometheusConfig=prometheus-emqx4.yaml
  dashboardDefinitions=$(dirname $(dirname $(pwd)))/config/grafana-template/EMQX4-enterprise
  emqxImage=emqx/emqx-ee:4.4.16
  ;;
emqx5)
  prometheusConfig=prometheus-emqx5.yaml
  dashboardDefinitions=$(dirname $(dirname $(pwd)))/config/grafana-template/EMQX5
  emqxImage=emqx/emqx:5.0.1
  ;;
*)
  # deploy emqx5 enterprise by default
  prometheusConfig=prometheus-emqx5.yaml
  dashboardDefinitions=$(dirname $(dirname $(pwd)))/config/grafana-template/EMQX5-enterprise
  emqxImage=emqx/emqx-enterprise:5.0.1
  ;;
esac

#docker run -d --name emqx-ee -p 1883:1883 -p 8081:8081 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx-ee:4.4.16
docker run -d --name emqx-demo \
 -v "$(pwd)"/api_secret:/opt/emqx/data/api_secret \
 -e EMQX_DASHBOARD__BOOTSTRAP_USERS_FILE='"/opt/emqx/data/api_secret"' \
 -e EMQX_DASHBOARD__DEFAULT_USER__LOGIN=$emqxDashboardUsername \
 -e EMQX_DASHBOARD__DEFAULT_USER__PASSWORD=$emqxDashboardPassword \
 -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 $emqxImage

# load emqx_prometheus if the EMQX version is 4.4 open-source
if [[ $emqxVersion == "emqx4" ]]; then
while
 plugin=$(docker exec -it emqx-demo ./bin/emqx_ctl plugins list | grep emqx_prometheus | grep true)
 [[ $plugin == "" ]]
 do
   echo "loading plugin emqx_prometheus"
   sleep 5s
   docker exec -it emqx-demo ./bin/emqx_ctl plugins load emqx_prometheus
done
fi

# the username and password is defined in the file api_secret
docker run -d --name exporter-demo -p 8085:8085 emqx/emqx-exporter:latest \
 --emqx.nodes="emqx-demo:18083" --emqx.auth-username=76668f8a2003d597 --emqx.auth-password=CRCDB6lxxzN58e5HoD82llBC0Erg1TVZIAUsdTjPU7N


# use existing config to run prometheus
docker run -d --name prometheus-demo -p 9090:9090 -v $(pwd)/$prometheusConfig:/etc/prometheus/prometheus.yml prom/prometheus

# use provision to run grafana
provisioning=$(dirname $(pwd))/provisioning
docker run -d --name grafana-demo -p 3000:3000 \
 -v "$dashboardDefinitions":/grafana-dashboard-definitions \
 -v "$provisioning"/dashboard.yaml:/etc/grafana/provisioning/dashboards/dashboard.yaml \
 -v "$provisioning"/datasource.yaml:/etc/grafana/provisioning/datasources/datasource.yaml \
 grafana/grafana

res=$(docker network list | grep test)
if [[ $res == "" ]]; then
  docker network create test
fi

docker network connect test emqx-demo
docker network connect test exporter-demo
docker network connect test prometheus-demo
docker network connect test grafana-demo

echo -e -n "\nOpen http://localhost:3000 and sign in Grafana with admin/admin"
echo -e "\nOpen http://localhost:18083 and sign in EMQX dashboard with $emqxDashboardUsername/$emqxDashboardPassword"