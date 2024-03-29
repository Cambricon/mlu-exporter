/*
 * Copyright 2022 Cambricon, Inc.
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

#ifndef CNPAPI_H_
#define CNPAPI_H_
#include <stdint.h>
#include "cnpapi_types.h"
#include "callbackapi.h"
#include "activity_api.h"
#include "cnpapi_pmu_api.h"
#ifdef __cplusplus
extern "C" {
#endif
CNPAPI_EXPORT cnpapiResult cnpapiInit();
CNPAPI_EXPORT cnpapiResult cnpapiGetDeviceCount(int* num);
CNPAPI_EXPORT cnpapiResult cnpapiGetDeviceType(int dev_id, cnpapiDeviceType *type);
CNPAPI_EXPORT cnpapiResult cnpapiGetDeviceChipType(int dev_id, cnpapiChipType *chip);
CNPAPI_EXPORT cnpapiResult cnpapiGetResultString(cnpapiResult rst, const char **str);
CNPAPI_EXPORT cnpapiResult cnpapiGetLastError();
CNPAPI_EXPORT uint64_t cnpapiGetTimestamp();
CNPAPI_EXPORT cnpapiResult cnpapiRelease();
CNPAPI_EXPORT cnpapiResult cnpapiGetLibVersion(int *major, int *minor, int *patch);
CNPAPI_EXPORT cnpapiResult cnpapiGetSymbolNameFromCNkernel(void *cnKernel, char **name);
#ifdef __cplusplus
}
#endif

#endif  // CNPAPI_H_
