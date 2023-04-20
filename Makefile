# Copyright 2015 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Ensure that 'all' is the default target otherwise it will be the first target from Makefile.common.
all::

ARCH = $(shell go env GOARCH)
OS = $(shell go env GOOS)

# Needs to be defined before including Makefile.common to auto-generate targets
DOCKER_ARCHS ?= amd64

DOCKER_IMAGE_NAME       ?= emqx-exporter:latest

.PHONY: build LOCALBIN
build:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -o .build/${OS}-${ARCH}/emqx-exporter

.PHONY: test
test:
	go test -race --cover -covermode=atomic -coverpkg=./... -coverprofile=cover.out ./...

.PHONY: docker-build
docker-build: build
	docker build -t ${DOCKER_IMAGE_NAME} .

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/.build/${OS}-${ARCH}
$(LOCALBIN):
	mkdir -p $(LOCALBIN)