# Build the manager binary
FROM golang:1.20.3 as builder

ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace
COPY . .
RUN go work init \
    && go work use . \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -o emqx-exporter

FROM quay.io/prometheus/busybox:latest
LABEL maintainer="EMQX"

COPY --from=builder /workspace/emqx-exporter /bin/emqx-exporter

EXPOSE      8085
USER        nobody
ENTRYPOINT  [ "/bin/emqx-exporter" ]
