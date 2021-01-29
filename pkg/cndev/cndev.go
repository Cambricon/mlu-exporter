// Copyright 2020 Cambricon, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cndev

// #cgo LDFLAGS: -ldl
// #include "include/cndev.h"
// #include "cndev_dl.h"
import "C"

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const version = 3

type Cndev interface {
	Init() error
	Release() error
	GetDeviceCount() (uint, error)
	GetDeviceSN(idx uint) (string, error)
	GetDeviceModel(idx uint) string
	GetDeviceClusterCount(idx uint) (uint, error)
	GetDeviceTemperature(idx uint) (uint, []uint, error)
	GetDeviceHealth(idx uint) (uint, error)
	GetDeviceMemory(idx uint) (uint, uint, error)
	GetDevicePower(idx uint) (uint, error)
	GetDeviceCoreNum(idx uint) (uint, error)
	GetDeviceUtil(idx uint) (uint, []uint, error)
	GetDeviceFanSpeed(idx uint) (uint, error)
	GetDevicePCIeID(idx uint) (string, error)
	GetDeviceVersion(idx uint) (string, string, error)
}

type cndev struct{}

func NewCndevClient() Cndev {
	return &cndev{}
}

func errorString(ret C.cndevRet_t) error {
	if ret == C.CNDEV_SUCCESS {
		return nil
	}
	err := C.GoString(C.cndevGetErrorString(ret))
	return fmt.Errorf("cndev: %v", err)
}

func (c *cndev) Init() error {
	r := C.cndevInit_dl()
	if r == C.CNDEV_ERROR_UNINITIALIZED {
		return errors.New("could not load CNDEV library")
	}
	return errorString(r)
}

func (c *cndev) Release() error {
	r := C.cndevRelease_dl()
	return errorString(r)
}

func (c *cndev) GetDeviceCount() (uint, error) {
	var cardInfos C.cndevCardInfo_t
	cardInfos.version = C.int(version)
	r := C.cndevGetDeviceCount(&cardInfos)
	return uint(cardInfos.Number), errorString(r)
}

func (c *cndev) GetDeviceSN(idx uint) (string, error) {
	var cardSN C.cndevCardSN_t
	cardSN.version = C.int(version)
	r := C.cndevGetCardSN(&cardSN, C.int(idx))
	if err := errorString(r); err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x", int64(cardSN.sn))
	return uuid, nil
}

func (c *cndev) GetDeviceModel(idx uint) string {
	return C.GoString(C.getCardNameStringByDevId(C.int(idx)))
}

func (c *cndev) GetDeviceClusterCount(idx uint) (uint, error) {
	var cardClusterNum C.cndevCardClusterCount_t
	cardClusterNum.version = C.int(version)
	r := C.cndevGetClusterCount(&cardClusterNum, C.int(idx))
	return uint(cardClusterNum.count), errorString(r)
}

func (c *cndev) GetDeviceTemperature(idx uint) (uint, []uint, error) {
	var cardTemperature C.cndevTemperatureInfo_t
	cardTemperature.version = C.int(version)
	cluster, err := c.GetDeviceClusterCount(idx)
	if err != nil {
		return 0, nil, err
	}
	r := C.cndevGetTemperatureInfo(&cardTemperature, C.int(idx))
	if err := errorString(r); err != nil {
		return 0, nil, err
	}
	clusterTemperature := make([]uint, cluster)
	for i := uint(0); i < cluster; i++ {
		clusterTemperature[i] = uint(cardTemperature.Cluster[i])
	}
	return uint(cardTemperature.Board), clusterTemperature, nil
}

func (c *cndev) GetDeviceHealth(idx uint) (uint, error) {
	var cardHealthState C.cndevCardHealthState_t
	cardHealthState.version = C.int(version)
	r := C.cndevGetCardHealthState(&cardHealthState, C.int(idx))
	return uint(cardHealthState.health), errorString(r)
}

func (c *cndev) GetDeviceMemory(idx uint) (uint, uint, error) {
	var cardMemInfo C.cndevMemoryInfo_t
	cardMemInfo.version = C.int(version)
	r := C.cndevGetMemoryUsage(&cardMemInfo, C.int(idx))
	return uint(cardMemInfo.PhysicalMemoryUsed), uint(cardMemInfo.PhysicalMemoryTotal), errorString(r)
}

func (c *cndev) GetDevicePower(idx uint) (uint, error) {
	var cardPower C.cndevPowerInfo_t
	cardPower.version = C.int(version)
	r := C.cndevGetPowerInfo(&cardPower, C.int(idx))
	return uint(cardPower.Usage), errorString(r)
}

func (c *cndev) GetDeviceCoreNum(idx uint) (uint, error) {
	var cardCoreNum C.cndevCardCoreCount_t
	cardCoreNum.version = C.int(version)
	r := C.cndevGetCoreCount(&cardCoreNum, C.int(idx))
	return uint(cardCoreNum.count), errorString(r)
}

func (c *cndev) GetDeviceUtil(idx uint) (uint, []uint, error) {
	var cardUtil C.cndevUtilizationInfo_t
	cardUtil.version = C.int(version)
	core, err := c.GetDeviceCoreNum(idx)
	if err != nil {
		return 0, nil, err
	}
	r := C.cndevGetDeviceUtilizationInfo(&cardUtil, C.int(idx))
	if err := errorString(r); err != nil {
		return 0, nil, err
	}
	coreUtil := make([]uint, core)
	for i := uint(0); i < core; i++ {
		coreUtil[i] = uint(cardUtil.CoreUtilization[i])
	}
	return uint(cardUtil.BoardUtilization), coreUtil, nil
}

func (c *cndev) GetDeviceFanSpeed(idx uint) (uint, error) {
	var cardFanSpeed C.cndevFanSpeedInfo_t
	cardFanSpeed.version = C.int(version)
	r := C.cndevGetFanSpeedInfo(&cardFanSpeed, C.int(idx))
	return uint(cardFanSpeed.FanSpeed), errorString(r)
}

func getPCIeID(domain int64, bus int64, device int64, function int64) string {
	domainStr := strconv.FormatInt(domain, 16)
	domainStr = strings.Repeat("0", 4-len([]byte(domainStr))) + domainStr
	busStr := strconv.FormatInt(bus, 16)
	if bus < 16 {
		busStr = "0" + busStr
	}
	deviceStr := strconv.FormatInt(device, 16)
	if device < 16 {
		deviceStr = "0" + deviceStr
	}
	functionStr := strconv.FormatInt(function, 16)
	return domainStr + ":" + busStr + ":" + deviceStr + "." + functionStr
}

func (c *cndev) GetDevicePCIeID(idx uint) (string, error) {
	var pcieInfo C.cndevPCIeInfo_t
	pcieInfo.version = C.int(version)
	r := C.cndevGetPCIeInfo(&pcieInfo, C.int(idx))
	if err := errorString(r); err != nil {
		return "", err
	}
	return getPCIeID(int64(pcieInfo.domain), int64(pcieInfo.bus), int64(pcieInfo.device), int64(pcieInfo.function)), nil
}

func getVersion(major uint, minor uint, build uint) string {
	return fmt.Sprintf("v%d.%d.%d", major, minor, build)
}

func (c *cndev) GetDeviceVersion(idx uint) (string, string, error) {
	var versionInfo C.cndevVersionInfo_t
	versionInfo.version = C.int(version)
	r := C.cndevGetVersionInfo(&versionInfo, C.int(idx))
	if err := errorString(r); err != nil {
		return "", "", err
	}
	return getVersion(uint(versionInfo.MCUMajorVersion), uint(versionInfo.MCUMinorVersion), uint(versionInfo.MCUBuildVersion)), getVersion(uint(versionInfo.DriverMajorVersion), uint(versionInfo.DriverMinorVersion), uint(versionInfo.DriverBuildVersion)), nil
}
