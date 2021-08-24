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
#ifdef __cplusplus
}
#endif

#endif  // CNPAPI_H_
