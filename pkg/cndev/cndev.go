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
	"unsafe"
)

const version = 5

type Cndev interface {
	Init() error
	Release() error
	GetDeviceCount() (uint, error)
	GetDeviceSN(idx uint) (string, error)
	GetDeviceUUID(idx uint) (string, error)
	GetDeviceModel(idx uint) string
	GetDeviceClusterCount(idx uint) (uint, error)
	GetDeviceTemperature(idx uint) (uint, []uint, error)
	GetDeviceHealth(idx uint) (uint, error)
	GetDeviceMemory(idx uint) (uint, uint, uint, uint, error)
	GetDevicePower(idx uint) (uint, error)
	GetDeviceCoreNum(idx uint) (uint, error)
	GetDeviceUtil(idx uint) (uint, []uint, error)
	GetDeviceFanSpeed(idx uint) (uint, error)
	GetDeviceVersion(idx uint) (string, string, error)
	GetDeviceVfState(idx uint) (int, error)
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
	return uint(cardInfos.number), errorString(r)
}

func (c *cndev) GetDeviceSN(idx uint) (string, error) {
	var cardSN C.cndevCardSN_t
	cardSN.version = C.int(version)
	r := C.cndevGetCardSN(&cardSN, C.int(idx))
	if err := errorString(r); err != nil {
		return "", err
	}
	sn := fmt.Sprintf("%x", int64(cardSN.sn))
	return sn, nil
}

func (c *cndev) GetDeviceUUID(idx uint) (string, error) {
	var uuidInfo C.cndevUUID_t
	uuidInfo.version = C.int(version)
	r := C.cndevGetUUID(&uuidInfo, C.int(idx))
	if err := errorString(r); err != nil {
		return "", err
	}
	uuid := *(*[C.UUID_SIZE]C.uchar)(unsafe.Pointer(&uuidInfo.uuid))
	return fmt.Sprintf("%s", uuid), nil
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
		clusterTemperature[i] = uint(cardTemperature.cluster[i])
	}
	return uint(cardTemperature.board), clusterTemperature, nil
}

func (c *cndev) GetDeviceHealth(idx uint) (uint, error) {
	var cardHealthState C.cndevCardHealthState_t
	cardHealthState.version = C.int(version)
	r := C.cndevGetCardHealthState(&cardHealthState, C.int(idx))
	return uint(cardHealthState.health), errorString(r)
}

func (c *cndev) GetDeviceMemory(idx uint) (uint, uint, uint, uint, error) {
	var cardMemInfo C.cndevMemoryInfo_t
	cardMemInfo.version = C.int(version)
	r := C.cndevGetMemoryUsage(&cardMemInfo, C.int(idx))
	return uint(cardMemInfo.physicalMemoryUsed), uint(cardMemInfo.physicalMemoryTotal), uint(cardMemInfo.virtualMemoryUsed), uint(cardMemInfo.virtualMemoryTotal), errorString(r)
}

func (c *cndev) GetDevicePower(idx uint) (uint, error) {
	var cardPower C.cndevPowerInfo_t
	cardPower.version = C.int(version)
	r := C.cndevGetPowerInfo(&cardPower, C.int(idx))
	return uint(cardPower.usage), errorString(r)
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
		coreUtil[i] = uint(cardUtil.coreUtilization[i])
	}
	return uint(cardUtil.averageCoreUtilization), coreUtil, nil
}

func (c *cndev) GetDeviceFanSpeed(idx uint) (uint, error) {
	var cardFanSpeed C.cndevFanSpeedInfo_t
	cardFanSpeed.version = C.int(version)
	r := C.cndevGetFanSpeedInfo(&cardFanSpeed, C.int(idx))
	return uint(cardFanSpeed.fanSpeed), errorString(r)
}

// GetDeviceVfState returns cndevCardVfState_t.vfState, which is bitmap of mlu vf state. See cndev user doc for more details.
func (c *cndev) GetDeviceVfState(idx uint) (int, error) {
	var vfstate C.cndevCardVfState_t
	vfstate.version = C.int(version)
	r := C.cndevGetCardVfState(&vfstate, C.int(idx))
	return int(vfstate.vfState), errorString(r)
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
	return getVersion(uint(versionInfo.mcuMajorVersion), uint(versionInfo.mcuMinorVersion), uint(versionInfo.mcuBuildVersion)), getVersion(uint(versionInfo.driverMajorVersion), uint(versionInfo.driverMinorVersion), uint(versionInfo.driverBuildVersion)), nil
}
