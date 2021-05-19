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

ARG BUILDPLATFORM=linux/amd64
FROM --platform=$BUILDPLATFORM golang:1.13 as build
ARG APT_PROXY
ARG GOPROXY
ARG TARGETPLATFORM
RUN set -ex && export http_proxy=$APT_PROXY && \
  apt-get update && \
  apt-get install -y build-essential gcc-aarch64-linux-gnu ca-certificates make
WORKDIR /work/
COPY . .
RUN make build

FROM ubuntu:18.04
ARG TARGETPLATFORM
COPY --from=build /work/mlu-exporter /usr/bin/
COPY libs/$TARGETPLATFORM/*.so /usr/lib/
COPY examples/metrics.yaml /etc/mlu-exporter/metrics.yaml
CMD ["/usr/bin/mlu-exporter"]
