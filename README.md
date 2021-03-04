# Cambricon MLU Exporter

Prometheus exporter for Cambricon MLU metrics, written in Go with pluggable metric collectors.

## Installation and Usage

### Preparing your MLU Nodes

It is assumed that the Cambricon driver and neuware are installed on your MLU Nodes.

### Download and build

```shell
git clone https://github.com/Cambricon/mlu-exporter.git
cd mlu-exporter
```

Set the following environment variables if you need.

| env       | description                                                                   |
| --------- | ----------------------------------------------------------------------------- |
| APT_PROXY | apt proxy address                                                             |
| GOPROXY   | golang proxy address                                                          |
| ARCH      | target platform architecture, amd64 or arm64, amd64 by default                |
| LIBCNDEV  | absolute path of the libcndev.so binary, neuware installation path by default |
| LIBCNPAPI  | absolute path of the libcnpapi.so binary, neuware installation path by default |

> If you want to cross build, make sure docker version >= 19.03.

for amd64:

```shell
./build_image.sh
```

for arm64:

```shell
export ARCH=arm64
./build_image.sh
```

Please make sure Cambricon neuware is installed in your compiling environment.
It uses **libcndev.so** binary in your compiling machine and generates Cambricon MLU exporter docker images under folder ./image.

### Usage

#### Docker

Load the image you built in ./image folder on your node.

Use the following command to start the exporter.

```shell
docker run -d \
-p 30108:30108 \
--privileged=true \
cambricon-mlu-exporter:v1.5.2
```

Then use the following command to get the metrics.

```shell
curl localhost:30108/metrics
```

If you want to configure command args by yourself , see the following example:

```shell
docker run -d \
-p 30108:30108 \
-v examples/metrics.yaml:/var/lib/mlu-exporter/metrics.yaml \
--privileged=true \
cambricon-mlu-exporter:v1.5.2 \
mlu-exporter \
--metrics-config=/var/lib/mlu-exporter/metrics.yaml \
--metrics-path=/metrics \
--port=30108 \
--hostname=hostname \
--metrics-prefix=ai \
--collector=cndev \
--collector=host
```

Command Args Description

| arg            | description                                |
| -------------- | ------------------------------------------ |
| metrics-config | configuration file of MLU exporter metrics |
| metrics-path   | metrics path of the exporter service       |
| hostname       | machine hostname                           |
| port           | exporter service port                      |
| metrics-prefix | prefix of all metric names                 |
| collector      | collector names, cndev by default          |

available collectors:

- cndev: collects basic MLU metrics
- podresources: collects MLU usage metrics in containers managed by Kubernetes
- cnpapi: collects cnpapi pmu api metrics. **Please note that cnpapi can only be used by one single process on a machine. Not recommended for production scenarios.**
- host: collects host machine metrics

And set the metrics configuration file passed by your metrics-config arg as you like, see examples/metrics.yaml for an example.

#### Kubernetes And Prometheus

Load the image you built in ./image folder on your node.

```shell
kubectl apply -f examples/cambricon-mlu-exporter-cm.yaml
kubectl apply -f examples/cambricon-mlu-exporter-ds.yaml
kubectl apply -f examples/cambricon-mlu-exporter-svc.yaml
```

use the command to get the endpoint

```shell
kubectl get ep -n kube-system -l app=exporter-mlu-monitoring
```

curl  http://\<endpoint\>/metrics to get the metrics.

And if you want to create a Prometheus service monitor

```shell
kubectl apply -f examples/cambricon-mlu-exporter-sm.yaml
```

Then checkout your Prometheus to get the MLU metrics.

You can also set the command args described above in the MLU exporter daemonset spec. And set the metrics configuration in the MLU exporter configMap.

## Upgrade Notice

**Please see [changelog](CHANGELOG.md) for deprecation and breaking changes.**
