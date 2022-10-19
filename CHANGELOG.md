# Changelog

## v1.6.7

- Upgrade dependence to cntoolkit 2.8.2
- Add metric:
  - mlu_process_ipu_utilization
  - mlu_process_jpu_utilization
  - mlu_process_memory_utilization
  - mlu_process_vpu_decode_utilization
  - mlu_process_vpu_encode_utilization
- Set `hostPID` and `hostIPC` to true for exporter daemonset

## v1.6.6

Fix mlu_container does not show up for MLUs not used by pods

## v1.6.5

- Upgrade dependence to cntoolkit 2.7.0
- Add metric:
  - mlu_ecc_address_forbidden_error_total

## v1.6.4

**BREAKING CHANGE**: Rename metric `virtual_memory_utilization` to `virtual_function_memory_utilization`

- Fix virtual function memory utilization formula
- Upgrade dependence to cntoolkit 2.6.0
- Add PCIe, ECC, CRC, MLULink and some utilizations metrics as follows:
  - mlu_pcie_info
  - mlu_chip_cpu_utilization
  - mlu_virtual_memory_total
  - mlu_virtual_memory_used
  - mlu_arm_os_memory_total
  - mlu_arm_os_memory_used
  - mlu_video_codec_utilization
  - mlu_image_codec_utilization
  - mlu_tiny_core_utilization
  - mlu_numa_node_id
  - mlu_ddr_data_width
  - mlu_ddr_band_width
  - mlu_ecc_corrected_error_total
  - mlu_ecc_multiple_error_total
  - mlu_ecc_multiple_multiple_error_total
  - mlu_ecc_multiple_one_error_total
  - mlu_ecc_one_bit_error_total
  - mlu_ecc_error_total
  - mlu_ecc_uncorrected_error_total
  - mlu_mlulink_p2p_transfer_capability
  - mlu_mlulink_interlaken_serdes_capability
  - mlu_mlulink_cntr_read_byte_total
  - mlu_mlulink_cntr_read_package_total
  - mlu_mlulink_cntr_write_byte_total
  - mlu_mlulink_cntr_write_package_total
  - mlu_mlulink_err_corrected_total
  - mlu_mlulink_err_crc24_total
  - mlu_mlulink_err_crc32_total
  - mlu_mlulink_err_ecc_double_total
  - mlu_mlulink_err_fatal_total
  - mlu_mlulink_err_replay_total
  - mlu_mlulink_err_uncorrected_total
  - mlu_mlulink_port_mode
  - mlu_mlulink_speed_format
  - mlu_mlulink_speed
  - mlu_mlulink_status
  - mlu_mlulink_serdes_status
  - mlu_mlulink_version
  - mlu_d2d_crc_error_total
  - mlu_d2d_crc_error_overflow_total

## v1.6.3

- Fix cluster temperature overflow

## v1.6.2

- Fix uuid \x00 suffix
- Fix containers with MLU-sn uuids causing error response

## v1.6.1

- Support new devices

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
