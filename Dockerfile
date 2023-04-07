ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="EMQX"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/emqx-exporter /bin/emqx-exporter

EXPOSE      8085
USER        nobody
ENTRYPOINT  [ "/bin/emqx-exporter" ]
