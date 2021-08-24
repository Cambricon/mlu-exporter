# Changelog

## v1.6.1

- Support MLU365-D2 devices

## v1.6.0

**BREAKING CHANGE**: MLU driver must be equal or above 4.15.2

- Upgrade dependence to cntoolkit 2.2.0
- Get MLU uuid from cndev instead of using MLU sn
- Change default metric and lable names
- Add MLU vf memory usage metric
- **Remove** container_resource_mlu_utilization
- **Remove** container_resource_mlu_memory_utilization
- **Remove** container_resource_mlu_board_power
- Refactor how we deal with errors
- Move metric keys consts to collector package
- Refacor collector function maps
- Refactor MLU vf utilization logic

## v1.5.3

- Watch and reload metrics config dynamically
- Fix MLU220 capacity error
- Fix exporter panics when configured label not applicable

## v1.5.2

- Add host and cnpapi collectors

## v1.5.1

- Add MLU driver, mcu and MLU type labels
- Add mlu_container metric. Use `<metric> * on(boardid) group_right ai_mlu_container` to append k8s container info to a metric.
- **Deprecation:** container_resource_mlu_utilization will be removed in the future
- **Deprecation:** container_resource_mlu_memory_utilization will be removed in the future
- **Deprecation:** container_resource_mlu_board_power will be removed in the future

## v1.5.0

- Open source basic functions.
