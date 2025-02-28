# Cambricon MLU Exporter

Prometheus exporter for Cambricon MLU metrics, written in Go with pluggable metric collectors.

## Prerequisites

The prerequisites for running Cambricon MLU Exporter:

- MLU270, MLU270-X5K, MLU220, MLU290, MLU370 devices
- For MLU 2xx needs driver >= 4.9.13
- For MLU 3xx needs driver >= 4.20.9
- For MLU 2xx、3xx needs cntoolkit >= 2.8.2 on your building machine

For MLU driver version before 4.9.13, please use [release v1.5.3].

## Installation and Usage

### Preparing your MLU Nodes

It is assumed that the Cambricon driver and neuware are installed on your MLU Nodes.

### Download and build

```shell
git clone https://github.com/Cambricon/mlu-exporter.git
cd mlu-exporter
```

Set the following environment variables if you need.

| env        | description                                                                   |
| ---------- | ----------------------------------------------------------------------------- |
| APT_PROXY  | apt proxy address                                                             |
| GOPROXY    | golang proxy address                                                          |
| ARCH       | target platform architecture, amd64 or arm64, amd64 by default                |
| LIBCNDEV   | absolute path of the libcndev.so binary, neuware installation path by default |
| BASE_IMAGE | mlu exporter base image                                                       |

Docker should be >= 17.05.0 on your building machine. If you want to cross build, make sure docker version >= 19.03.

for amd64:

```shell
GOPROXY=https://goproxy.cn ./build_image.sh
```

for arm64:

```shell
ARCH=arm64 GOPROXY=https://goproxy.cn ./build_image.sh
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
--pid=host \
-e ENV_NODE_NAME={nodeName} \
cambricon-mlu-exporter:v2.0.17
```

Then use the following command to get the metrics.

```shell
curl localhost:30108/metrics
```

If you want to configure command args by yourself , see the following example:

```shell
docker run -d \
-p 30108:30108 \
-v examples/metrics.yaml:/etc/mlu-exporter/metrics.yaml \
--privileged=true \
--pid=host \
cambricon-mlu-exporter:v2.0.17 \
mlu-exporter \
--metrics-config=/etc/mlu-exporter/metrics.yaml \
--metrics-path=/metrics \
--log-level=info \
--port=30108 \
--hostname=hostname \
--metrics-prefix=mlu \
--collector=cndev \
--collector=host
```

Command Args Description

| arg                            | description                                                                                                                                          |
| ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| collector                      | collector names, cndev by default                                                                                                                    |
| env-share-num                  | vf numbers under env share mode, should set virtual-mode to env-share                                                                                |
| hostname                       | machine hostname, or env:"ENV_NODE_NAME"                                                                                                             |
| log-level                      | set log level: trace/debug/info/warn/error/fatal/panic" default:"info"                                                                               |
| metrics-config                 | configuration file of MLU exporter metrics                                                                                                           |
| metrics-path                   | metrics path of the exporter service                                                                                                                 |
| metrics-prefix                 | prefix of all metric names                                                                                                                           |
| port                           | exporter service port                                                                                                                                |
| virtual-mode                   | virtual mode, default "", support env-share                                                                                                          |
| push-gateway-url               | If set, metrics with push enabled will push to this server via prometheus push gateway protocol                                                      |
| push-interval-ms               | numbers of metrics push interval in milliseconds, minimum 100, default 500                                                                           |
| push-job-name                  | metrics push job name, default mlu-push-monitoring                                                                                                   |
| cluster-name                   | cluster name, add cluster label for metrics push                                                                                                     |
| xid-error-metric-name          | xid error metric name, push to push gateway immediately when xid event occurs                                                                        |
| xid-error-retry-times          | retry times for push xid error metric failed                                                                                                         |
| log-file-for-xid-metric-failed | log file to record xid information when push failed, should set related volume mount in `examples/cambricon-mlu-exporter-ds.yaml` for log collection |

available collectors:

- cndev: collects basic MLU metrics
- podresources: collects MLU usage metrics in containers managed by Kubernetes. For Kubernetes lower than 1.15, make sure `KubeletPodResources` [feature gate] is enabled by setting the `feature-gates` [kubelet option] in your kubelet configuration.
- host: collects host machine metrics

And set the metrics configuration file passed by your metrics-config arg as you like, see examples/metrics.yaml for an example.

#### Kubernetes And Prometheus

Load the image you built in ./image folder on your node.

```shell
kubectl apply -f examples/cambricon-mlu-exporter-cm.yaml
kubectl apply -f examples/cambricon-mlu-exporter-ds.yaml
kubectl apply -f examples/cambricon-mlu-exporter-svc.yaml
```

You can also set the command args described above as you like in the MLU exporter daemonset spec.

use the command to get the endpoint

```shell
kubectl get ep -n kube-system -l app=exporter-mlu-monitoring
```

curl <http://{endpoints}/metrics> to get the metrics.

And if you want to create a Prometheus service monitor

```shell
kubectl apply -f examples/cambricon-mlu-exporter-sm.yaml
```

if you want to create a Prometheus [additional scrape configs]

```shell
kubectl create secret generic additional-scrape-configs --from-file=examples/cambricon-mlu-exporter-additional.yaml
```

##### Install with helm

This step is optional, you can also install mlu exporter via helm

1. Download and install helm firstly, refer <https://helm.sh/docs/intro/install/>

2. Clone this repo, `cd deploys/helm` and run `helm install mlu-exporter-release-x.x.x mlu-exporter/`

Then checkout your Prometheus to get the MLU metrics.

##### Change log level

Support to change log level and validate without restarting

change to debug

```shell
curl -i http://{{mlu-exporter-pod-ip}}:30108/logLevel?level=debug
```

change to info

```shell
curl -i http://{{mlu-exporter-pod-ip}}:30108/logLevel?level=info
```

##### Group Metrics

In a Kubernetes cluster, if the podresources collector is enabled, you can get namespace/pod/container info of an allocated MLU by the metric `mlu_container`.

To attach namespace/pod/container info to another MLU metric, use `mlu_container` as the following example in your Prometheus:

```text
mlu_power_usage * on(uuid) group_right mlu_container
```

And for env-share VFs:

```text
mlu_utilization * on(uuid,vf) group_right mlu_container
```

#### Metrics and Labels

For installation by docker, you can modify the configuration file you passed by metrics-config arg to change metric and label names.

For Kubernetes, you can modify the MLU exporter configMap to change metric and label names.

For MLU370, mlu_temperature does not support cluster temperature, all cluster temperature metrics are set to 0.

#### Push metrics

If you want to push metrics to Prometheus Pushgateway or VictoriaMetrics, you can run exporter with setting `push-gateway-url` flag,
and config `push: true` on metrics witch need pushing in metrics configuration.

See examples/metrics-push.yaml for pushing some metrics.

```yaml
utilization:
  help: The utilization of Cambricon MLU, unit is '%'
  labels:
    driver: driver
    mcu: mcu
    mlu: mlu
    model: model
    node: node
    sn: sn
    type: type
    uuid: uuid
    vf: vf
  name: utilization
  # enable push
  push: true
```

The metrics with `push: true` will be pushing to push-gateway-url, and metrics without this flag will export on metrics path.

Warning, the metrics with `push: true` should not export on metrics path such as `/metrics`.

## Upgrade Notice

**Please see [changelog](CHANGELOG.md) for deprecation and breaking changes.**

[release v1.5.3]: https://github.com/Cambricon/mlu-exporter/releases/tag/v1.5.3
[feature gate]: https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
[kubelet option]: https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/#options
[additional scrape configs]: https://github.com/prometheus-operator/prometheus-operator/tree/main/example/additional-scrape-configs
