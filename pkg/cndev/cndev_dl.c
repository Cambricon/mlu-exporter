/*
 * Copyright 2020 Cambricon, Inc.
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
#include "cndev_dl.h"
#include "./include/cndev.h"

void *cndevHandle;

cndevRet_t (*cndevInitFunc)(int);

cndevRet_t (*cndevReleaseFunc)(void);

const char *(*cndevGetErrorStringFunc)(cndevRet_t errorId);
const char *cndevGetErrorString(cndevRet_t errorId)
{
    if (cndevGetErrorStringFunc == NULL)
    {
        return "cndevGetErrorString Function not found";
    }
    return cndevGetErrorStringFunc(errorId);
}

cndevRet_t (*cndevGetDeviceCountFunc)(cndevCardInfo_t *cardNum);
cndevRet_t cndevGetDeviceCount(cndevCardInfo_t *cardNum)
{
    if (cndevGetDeviceCountFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetDeviceCountFunc(cardNum);
}

cndevRet_t (*cndevGetCardSNFunc)(cndevCardSN_t *cardSN, int devId);
cndevRet_t cndevGetCardSN(cndevCardSN_t *cardSN, int devId)
{
    if (cndevGetCardSNFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetCardSNFunc(cardSN, devId);
}

const char *(*getCardNameStringByDevIdFunc)(int devId);
const char *getCardNameStringByDevId(int devId)
{
    if (getCardNameStringByDevIdFunc == NULL)
    {
        return "getCardNameStringByDevId Function not found";
    }
    return getCardNameStringByDevIdFunc(devId);
}

cndevRet_t (*cndevGetClusterCountFunc)(cndevCardClusterCount_t *clusCount, int devId);
cndevRet_t cndevGetClusterCount(cndevCardClusterCount_t *clusCount, int devId)
{
    if (cndevGetClusterCountFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetClusterCountFunc(clusCount, devId);
}

cndevRet_t (*cndevGetTemperatureInfoFunc)(cndevTemperatureInfo_t *tempInfo, int devId);
cndevRet_t cndevGetTemperatureInfo(cndevTemperatureInfo_t *tempInfo, int devId)
{
    if (cndevGetTemperatureInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetTemperatureInfoFunc(tempInfo, devId);
}

cndevRet_t (*cndevGetCardHealthStateFunc)(cndevCardHealthState_t *cardHealthState, int devId);
cndevRet_t cndevGetCardHealthState(cndevCardHealthState_t *cardHealthState, int devId)
{
    if (cndevGetCardHealthStateFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetCardHealthStateFunc(cardHealthState, devId);
}

cndevRet_t (*cndevGetMemoryUsageFunc)(cndevMemoryInfo_t *memInfo, int devId);
cndevRet_t cndevGetMemoryUsage(cndevMemoryInfo_t *memInfo, int devId)
{
    if (cndevGetMemoryUsageFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetMemoryUsageFunc(memInfo, devId);
}

cndevRet_t (*cndevGetPowerInfoFunc)(cndevPowerInfo_t *powerInfo, int devId);
cndevRet_t cndevGetPowerInfo(cndevPowerInfo_t *powerInfo, int devId)
{
    if (cndevGetPowerInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetPowerInfoFunc(powerInfo, devId);
}

cndevRet_t (*cndevGetCoreCountFunc)(cndevCardCoreCount_t *cardCoreCount, int devId);
cndevRet_t cndevGetCoreCount(cndevCardCoreCount_t *cardCoreCount, int devId)
{
    if (cndevGetCoreCountFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetCoreCountFunc(cardCoreCount, devId);
}

cndevRet_t (*cndevGetDeviceUtilizationInfoFunc)(cndevUtilizationInfo_t *utilInfo, int devId);
cndevRet_t cndevGetDeviceUtilizationInfo(cndevUtilizationInfo_t *utilInfo, int devId)
{
    if (cndevGetDeviceUtilizationInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetDeviceUtilizationInfoFunc(utilInfo, devId);
}

cndevRet_t (*cndevGetFanSpeedInfoFunc)(cndevFanSpeedInfo_t *fanInfo, int devId);
cndevRet_t cndevGetFanSpeedInfo(cndevFanSpeedInfo_t *fanInfo, int devId)
{
    if (cndevGetFanSpeedInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetFanSpeedInfoFunc(fanInfo, devId);
}

cndevRet_t (*cndevGetPCIeInfoFunc)(cndevPCIeInfo_t *deviceInfo, int devId);
cndevRet_t cndevGetPCIeInfo(cndevPCIeInfo_t *deviceInfo, int devId)
{
    if (cndevGetPCIeInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetPCIeInfoFunc(deviceInfo, devId);
}

cndevRet_t (*cndevGetVersionInfoFunc)(cndevVersionInfo_t *versInfo, int devId);
cndevRet_t cndevGetVersionInfo(cndevVersionInfo_t *versInfo, int devId)
{
    if (cndevGetVersionInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetVersionInfoFunc(versInfo, devId);
}

cndevRet_t (*cndevGetUUIDFunc)(cndevUUID_t *uuidInfo, int devId);
cndevRet_t cndevGetUUID(cndevUUID_t *uuidInfo, int devId)
{
    if (cndevGetUUIDFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetUUIDFunc(uuidInfo, devId);
}

cndevRet_t (*cndevGetCardVfStateFunc)(cndevCardVfState_t *vfstate, int devId);
cndevRet_t cndevGetCardVfState(cndevCardVfState_t *vfstate, int devId)
{
    if (cndevGetCardVfStateFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    return cndevGetCardVfStateFunc(vfstate, devId);
}
cndevRet_t CNDEV_DL(cndevInit)(void)
{
    cndevHandle = dlopen("libcndev.so", RTLD_LAZY);
    if (cndevHandle == NULL)
    {
        return CNDEV_ERROR_UNINITIALIZED;
    }
    cndevInitFunc = dlsym(cndevHandle, "cndevInit");
    if (cndevInitFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevReleaseFunc = dlsym(cndevHandle, "cndevRelease");
    if (cndevReleaseFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetErrorStringFunc = dlsym(cndevHandle, "cndevGetErrorString");
    if (cndevGetErrorStringFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetDeviceCountFunc = dlsym(cndevHandle, "cndevGetDeviceCount");
    if (cndevGetDeviceCountFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetCardSNFunc = dlsym(cndevHandle, "cndevGetCardSN");
    if (cndevGetCardSNFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetUUIDFunc = dlsym(cndevHandle, "cndevGetUUID");
    if (cndevGetUUIDFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetCardVfStateFunc = dlsym(cndevHandle, "cndevGetCardVfState");
    if (cndevGetCardVfStateFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    getCardNameStringByDevIdFunc = dlsym(cndevHandle, "getCardNameStringByDevId");
    if (getCardNameStringByDevIdFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetClusterCountFunc = dlsym(cndevHandle, "cndevGetClusterCount");
    if (cndevGetClusterCountFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetTemperatureInfoFunc = dlsym(cndevHandle, "cndevGetTemperatureInfo");
    if (cndevGetTemperatureInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetCardHealthStateFunc = dlsym(cndevHandle, "cndevGetCardHealthState");
    if (cndevGetCardHealthStateFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetMemoryUsageFunc = dlsym(cndevHandle, "cndevGetMemoryUsage");
    if (cndevGetMemoryUsageFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetPowerInfoFunc = dlsym(cndevHandle, "cndevGetPowerInfo");
    if (cndevGetPowerInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetCoreCountFunc = dlsym(cndevHandle, "cndevGetCoreCount");
    if (cndevGetCoreCountFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetDeviceUtilizationInfoFunc = dlsym(cndevHandle, "cndevGetDeviceUtilizationInfo");
    if (cndevGetDeviceUtilizationInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetFanSpeedInfoFunc = dlsym(cndevHandle, "cndevGetFanSpeedInfo");
    if (cndevGetFanSpeedInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetPCIeInfoFunc = dlsym(cndevHandle, "cndevGetPCIeInfo");
    if (cndevGetPCIeInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevGetVersionInfoFunc = dlsym(cndevHandle, "cndevGetVersionInfo");
    if (cndevGetVersionInfoFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevRet_t result = cndevInitFunc(0);
    if (result != CNDEV_SUCCESS)
    {
        dlclose(cndevHandle);
        cndevHandle = NULL;
        return result;
    }
    return CNDEV_SUCCESS;
}

cndevRet_t CNDEV_DL(cndevRelease)(void)
{
    if (cndevHandle == NULL)
    {
        return CNDEV_SUCCESS;
    }
    if (cndevReleaseFunc == NULL)
    {
        return CNDEV_ERROR_UNKNOWN;
    }
    cndevRet_t r = cndevReleaseFunc();
    if (r != CNDEV_SUCCESS)
    {
        return r;
    }
    return (dlclose(cndevHandle) ? CNDEV_ERROR_UNKNOWN : CNDEV_SUCCESS);
}
