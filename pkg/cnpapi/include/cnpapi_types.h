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

#ifndef CNPAPI_TYPES_H_
#define CNPAPI_TYPES_H_
#include <stdint.h>
#ifdef __cplusplus
extern "C" {
#endif

#define CNPAPI_EXPORT __attribute__((visibility("default")))
#ifdef NDEBUG
#define CNPAPI_DEBUG_EXPORT
#else
#define CNPAPI_DEBUG_EXPORT __attribute__((visibility("default")))
#endif
#define CNPAPI_DISABLE_EXPORT __attribute__((visibility("hidden")))

typedef enum {
  CNPAPI_SUCCESS = 0,
  CNPAPI_ERROR_NOT_INITIALIZED = 1,
  CNPAPI_ERROR_INVALID_DEVICE = 2,
  CNPAPI_ERROR_INVALID_DEVICE_ID = 2,
  CNPAPI_ERROR_INVALID_ARGUMENT = 3,
  CNPAPI_ERROR_EVENT_GROUP_ENABLED = 4,  // when manipulate with event group, should disable it first
  CNPAPI_ERROR_INSUFFICIENT_MEMORY = 5,
  CNPAPI_ERROR_NO_DRIVER = 6,
  CNPAPI_ERROR_RESERVED0 = 7,
  CNPAPI_ERROR_UNKNOWN = 8,
  CNPAPI_ERROR_MAX_LIMIT_REACHED = 9,
  CNPAPI_ERROR_DRIVER_COMMUNICATION_FAILED = 10,
  CNPAPI_ERROR_DEVICE_BUSY = 11,
  CNPAPI_ERROR_ALREADY_IN_USE = 12,
  CNPAPI_ERROR_ACTIVITY_CALLBACK_NOT_REGISTERED = 13,
  CNPAPI_ERROR_ACTIVITY_CALLBACK_ALREADY_REGISTERED = 14,
  CNPAPI_ERROR_PMU_COUNTER_NOT_ENABLED = 15,
  CNPAPI_ERROR_INVALID_PMU_COUNTER_ID = 16,
  CNPAPI_ERROR_INVALID_PMU_COUNTER_NAME = 17,
  CNPAPI_ERROR_INVALID_DEVICE_TYPE = 18,
  CNPAPI_ERROR_UNSUPPORTED_DRIVER_VERSION = 19,
  CNPAPI_ERROR_INVALID_CHIP_TYPE = 20,
} cnpapiResult;

typedef enum {
  CNPAPI_DEVICE_TYPE_UNKNOWN = -1,
  CNPAPI_DEVICE_TYPE_MLU220 = 0,
  CNPAPI_DEVICE_TYPE_MLU270,
  CNPAPI_DEVICE_TYPE_MLU290,
  CNPAPI_DEVICE_TYPE_MLU370,
  CNPAPI_DEVICE_TYPE_SIZE,
  CNPAPI_MLU220 = CNPAPI_DEVICE_TYPE_MLU220,
  CNPAPI_MLU270 = CNPAPI_DEVICE_TYPE_MLU270,
  CNPAPI_MLU290 = CNPAPI_DEVICE_TYPE_MLU290
} cnpapiDeviceType_t;

typedef enum {
  CNPAPI_CHIP_TYPE_UNKNOWN = -1,
  CNPAPI_CHIP_TYPE_C20E = 0,
  CNPAPI_CHIP_TYPE_C20L,
  CNPAPI_CHIP_TYPE_C20,
  CNPAPI_CHIP_TYPE_C30S,
  CNPAPI_CHIP_TYPE_C30D,
  CNPAPI_CHIP_TYPE_SIZE
} cnpapiChipType_t;

#ifdef __cplusplus
}
#endif

#endif  // CNPAPI_TYPES_H_
