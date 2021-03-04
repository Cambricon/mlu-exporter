/*
 * Copyright 2021 Cambricon, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifndef CNPAPI_PMU_API_H_
#define CNPAPI_PMU_API_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include "cnpapi_types.h"  // NOLINT

typedef enum {
  PMU_AUTO_FLUSH,
  PMU_EXPLICIT_FLUSH
} pmuFlushMode_t;

/**
 * @brief: used to get supported counters of a certain type mlu card.
 *
 * @param: dev_type[in]: the card type.
 * @param: counter_dst[out]: counter id array pointer. the supported counter id will be copyed to this dst.
 * @param: size[in/out]: the size of given counter_dst, and the supported counter size will be set to this addr.
 * @ret: *CNPAPI_SUCCESS*                   on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*     cnpapi not initialized.
 *       *CNPAPI_ERROR_INSUFFICIENT_MEMORY* size less than needed.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*    size or counter_dst is NULL.
 *       *CNPAPI_ERROR_INVALID_DEVICE_TYPE* unknown device type.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuGetCounterSupported(cnpapiDeviceType_t dev_type,
                                                        uint64_t *counter_dst, uint64_t *size);

/**
 * @brief: used to get counter name by counter id.
 *
 * @param: counter_id[in]: counter id.
 * @param: name[out]: name dst.
 * @ret: *CNPAPI_SUCCESS*                      on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*        cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*       name is NULL.
 *       *CNPAPI_ERROR_INVALID_PMU_COUNTER_ID* unknown counter id.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuGetCounterName(uint64_t counter_id, const char** name);

/**
 * @brief: used to get conter id by counter name.
 *
 * @param: name[in]: counter name.
 * @param: counter_id[out]: counter id addr.
 * @ret: *CNPAPI_SUCCESS*                       on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*         cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*        counter_id is NULL or name is not a valid counter name.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuGetCounterIdByName(const char* name, uint64_t *counter_id);

/**
 * @brief: used to enable conter by counter id.
 *
 * @param: dev_id[in]: device id.
 * @param: counter_id[in]: counter id.
 * @param: is_enable[in]: true to enable counter and false to disable counter.
 * @ret: *CNPAPI_SUCCESS*                           on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*             cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_PMU_COUNTER_ID*      invalid counter id.
 *       *CNPAPI_ERROR_INVALID_DEVICE_ID*           invalid device.
 *       *CNPAPI_ERROR_PMU_COUNTER_NOT_ENABLED*     counter not enabled.
 *       *CNPAPI_ERROR_DRIVER_COMMUNICATION_FAILED* driver communication failed.
 *       *CNPAPI_ERROR_DEVICE_BUSY*                 device busy.
 *       *CNPAPI_ERROR_ALREADY_IN_USE*              pmu module already in use.
 *       *CNPAPI_ERROR_UNKNOWN*                     internal error.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuEnableCounter(int dev_id, uint64_t counter_id, bool is_enable);

/**
 * @brief: used to get conter value by counter id.
 *
 * @param: dev_id[in]: device id.
 * @param: counter_id[in]: counter id.
 * @param: value[out]: addr to get counter value.
 * @ret: *CNPAPI_SUCCESS*                           on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*             cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*            value is NULL.
 *       *CNPAPI_ERROR_INVALID_PMU_COUNTER_ID*      invalid counter id.
 *       *CNPAPI_ERROR_INVALID_DEVICE_ID*           invalid device.
 *       *CNPAPI_ERROR_PMU_COUNTER_NOT_ENABLED*     counter not enabled.
 *       *CNPAPI_ERROR_DRIVER_COMMUNICATION_FAILED* driver communication failed.
 *       *CNPAPI_ERROR_DEVICE_BUSY*                 device busy.
 *       *CNPAPI_ERROR_ALREADY_IN_USE*              pmu module already in use.
 *       *CNPAPI_ERROR_UNKNOWN*                     internal error.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuGetCounterValue(int dev_id, uint64_t counter_id, uint64_t *value);

/**
 * @brief: used to flush conter data of a device.
 *
 * @param: dev_id[in]: device id.
 * @ret: *CNPAPI_SUCCESS*                           on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*             cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_DEVICE_ID*           invalid device.
 *       *CNPAPI_ERROR_DRIVER_COMMUNICATION_FAILED* driver communication failed.
 *       *CNPAPI_ERROR_DEVICE_BUSY*                 device busy.
 *       *CNPAPI_ERROR_ALREADY_IN_USE*              pmu api already in use on other execution.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuFlushData(int dev_id);

/**
 * @brief: used to set data flush mode, the default mode is PMU_AUTO_FLUSH.
 *
 * @param: mode[in]: flush mode to set
 *                   PMU_AUTO_FLUSH will automatic flush data when calling cnpapiPmuGetCounterValue,
 *                   PMU_EXPLICIT_FLUSH will flush data only when calling cnpapiPmuFlushData.
 * @ret: *CNPAPI_SUCCESS*                           on success.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*            invalid mode.
 */
/* set data flush mode, the default mode is PMU_AUTO_FLUSH. */
CNPAPI_EXPORT cnpapiResult cnpapiPmuSetFlushMode(pmuFlushMode_t mode);

#ifdef __cplusplus
}
#endif
#endif  // CNPAPI_PMU_API_H_
