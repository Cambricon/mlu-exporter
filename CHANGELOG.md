# Changelog

## v1.5.3

+ Watch and reload metrics config dynamically
+ Fix MLU220 capacity error
+ Fix exporter panics when configured label not applicable

## v1.5.2

+ Add host and cnpapi collectors

## v1.5.1

+ Add MLU driver, mcu and mlu type labels
+ Add mlu_container metric. Use `<metric> * on(boardid) group_right ai_mlu_container` to append k8s container info to a metric.
+ **Deprecation:** container_resource_mlu_utilization will be removed in the future
+ **Deprecation:** container_resource_mlu_memory_utilization will be removed in the future
+ **Deprecation:** container_resource_mlu_board_power will be removed in the future

## v1.5.0

+ Open source basic functions.
