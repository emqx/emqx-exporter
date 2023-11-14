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

VSN ?= $(shell git describe --tags --always)
OS ?= $(shell go env GOOS)

all: build

.PHONY: build
build:
	GOOS=$(OS) go build -o $(LOCALBIN)/$(PROJECT_NAME)
	@cp $(PROJECT_DIR)/config/example/config.yaml $(LOCALBIN)/config.yaml
	@tar -zcvf emqx-exporter-$(OS)-$(VSN).tgz bin

.PHONY: test
test:
	go test -v -race --cover -covermode=atomic -coverpkg=./... -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o=cover.html

.PHONY: docker
docker:
	docker build -t $(PROJECT_NAME) .

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
PROJECT_NAME := emqx-exporter
LOCALBIN ?= $(PROJECT_DIR)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
