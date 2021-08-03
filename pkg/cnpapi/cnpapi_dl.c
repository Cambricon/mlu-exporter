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

#include <stddef.h>
#include <dlfcn.h>
#include <stdlib.h>
#include "cnpapi_dl.h"
#include "./include/cnpapi.h"
#include "./include/cnpapi_types.h"
#include "./include/cnpapi_pmu_api.h"

void *cnpapiHandle;

cnpapiResult (*cnpapiInitFunc)(void);

cnpapiResult (*cnpapiGetResultStringFunc)(cnpapiResult rst, const char **str);
cnpapiResult cnpapiGetResultString(cnpapiResult rst, const char **str)
{
    if (cnpapiGetResultStringFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    return cnpapiGetResultStringFunc(rst, str);
}

cnpapiResult (*cnpapiGetDeviceCountFunc)(int *num);
cnpapiResult cnpapiGetDeviceCount(int *num)
{
    if (cnpapiGetDeviceCountFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    return cnpapiGetDeviceCountFunc(num);
}

cnpapiResult (*cnpapiPmuGetCounterIdByNameFunc)(const char *name, uint64_t *counter_id);
cnpapiResult cnpapiPmuGetCounterIdByName(const char *name, uint64_t *counter_id)
{
    if (cnpapiPmuGetCounterIdByNameFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    return cnpapiPmuGetCounterIdByNameFunc(name, counter_id);
}

cnpapiResult (*cnpapiPmuEnableCounterFunc)(int dev_id, uint64_t counter_id, bool is_enable);
cnpapiResult cnpapiPmuEnableCounter(int dev_id, uint64_t counter_id, bool is_enable)
{
    if (cnpapiPmuEnableCounterFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    return cnpapiPmuEnableCounterFunc(dev_id, counter_id, is_enable);
}

cnpapiResult (*cnpapiPmuGetCounterValueFunc)(int dev_id, uint64_t counter_id, uint64_t *value);
cnpapiResult cnpapiPmuGetCounterValue(int dev_id, uint64_t counter_id, uint64_t *value)
{
    if (cnpapiPmuGetCounterValueFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    return cnpapiPmuGetCounterValueFunc(dev_id, counter_id, value);
}

cnpapiResult (*cnpapiPmuSetFlushModeFunc)(pmuFlushMode_t mode);
cnpapiResult cnpapiPmuSetFlushMode(pmuFlushMode_t mode)
{
    if (cnpapiPmuSetFlushModeFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    return cnpapiPmuSetFlushModeFunc(mode);
}

cnpapiResult (*cnpapiPmuFlushDataFunc)(int dev_id);
cnpapiResult cnpapiPmuFlushData(int dev_id)
{
    if (cnpapiPmuFlushDataFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    return cnpapiPmuFlushDataFunc(dev_id);
}

cnpapiResult (*cnpapiGetDeviceTypeFunc)(int dev_id, cnpapiDeviceType_t *type);
cnpapiResult cnpapiGetDeviceType(int dev_id, cnpapiDeviceType_t *type)
{
    if (cnpapiGetDeviceTypeFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    return cnpapiGetDeviceTypeFunc(dev_id, type);
}

cnpapiResult CNPAPI_DL(cnpapiInit)(void)
{
    cnpapiHandle = dlopen("libcnpapi.so", RTLD_LAZY);
    if (cnpapiHandle == NULL)
    {
        return CNPAPI_ERROR_NOT_INITIALIZED;
    }
    cnpapiInitFunc = dlsym(cnpapiHandle, "cnpapiInit");
    if (cnpapiInitFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    cnpapiGetDeviceCountFunc = dlsym(cnpapiHandle, "cnpapiGetDeviceCount");
    if (cnpapiGetDeviceCountFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    cnpapiGetResultStringFunc = dlsym(cnpapiHandle, "cnpapiGetResultString");
    if (cnpapiGetResultStringFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    cnpapiPmuGetCounterIdByNameFunc = dlsym(cnpapiHandle, "cnpapiPmuGetCounterIdByName");
    if (cnpapiPmuGetCounterIdByNameFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    cnpapiPmuEnableCounterFunc = dlsym(cnpapiHandle, "cnpapiPmuEnableCounter");
    if (cnpapiPmuEnableCounterFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    cnpapiPmuGetCounterValueFunc = dlsym(cnpapiHandle, "cnpapiPmuGetCounterValue");
    if (cnpapiPmuGetCounterValueFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    cnpapiPmuSetFlushModeFunc = dlsym(cnpapiHandle, "cnpapiPmuSetFlushMode");
    if (cnpapiPmuSetFlushModeFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    cnpapiPmuFlushDataFunc = dlsym(cnpapiHandle, "cnpapiPmuFlushData");
    if (cnpapiPmuFlushDataFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    cnpapiGetDeviceTypeFunc = dlsym(cnpapiHandle, "cnpapiGetDeviceType");
    if (cnpapiGetDeviceTypeFunc == NULL)
    {
        return CNPAPI_ERROR_UNKNOWN;
    }
    cnpapiResult result = cnpapiInitFunc();
    if (result != CNPAPI_SUCCESS)
    {
        dlclose(cnpapiHandle);
        cnpapiHandle = NULL;
        return result;
    }
    return CNPAPI_SUCCESS;
}
