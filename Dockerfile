# Build the manager binary
FROM golang:1.20.3 as builder

WORKDIR /workspace
COPY . .
RUN go work init \
    && go work use . \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux make build

FROM quay.io/prometheus/busybox:latest

WORKDIR /usr/local/emqx-exporter/bin
COPY --from=builder /workspace/bin/. /usr/local/emqx-exporter/bin

USER nobody:nobody
EXPOSE 8085
ENTRYPOINT [ "/usr/local/emqx-exporter/bin/emqx-exporter" ]
