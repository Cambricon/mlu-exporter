# Copyright 2020 Cambricon, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TARGETPLATFORM ?= linux/amd64
export GOOS := $(word 1, $(subst /, ,$(TARGETPLATFORM)))
export GOARCH := $(word 2, $(subst /, ,$(TARGETPLATFORM)))
export CGO_ENABLED := 1
ifeq ($(GOARCH), arm64)
export CC=aarch64-linux-gnu-gcc
endif

generate:
	mockgen -package mock -destination pkg/mock/cndev.go -mock_names=Cndev=Cndev github.com/Cambricon/mlu-exporter/pkg/cndev Cndev
	mockgen -package mock -destination pkg/mock/podrsources.go -mock_names=PodResources=PodResources github.com/Cambricon/mlu-exporter/pkg/podresources PodResources

lint:
	golangci-lint run -v

build:
	go build  -trimpath -ldflags="-s -w" -o mlu-exporter .

test:
	go test -v ./...

addlicense:
	# install with `go get github.com/google/addlicense`
	addlicense -c 'Cambricon, Inc.' -l apache -v .

clean:
	rm -f mlu-exporter
