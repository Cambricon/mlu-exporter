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

// #cgo LDFLAGS: -ldl -Wl,--unresolved-symbols=ignore-in-object-files
// #include "include/cndev.h"
import "C"

import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

const version = 5

type Cndev interface {
	Init() error
	Release() error

	GetDeviceArmOsMemory(idx uint) (int64, int64, error)
	GetDeviceClusterCount(idx uint) (int, error)
	GetDeviceCoreNum(idx uint) (uint, error)
	GetDeviceCount() (uint, error)
	GetDeviceCPUUtil(idx uint) (uint16, []uint8, error)
	GetDeviceCRCInfo(idx uint) (uint64, uint64, error)
	GetDeviceCurrentPCIeInfo(idx uint) (int, int, error)
	GetDeviceDDRInfo(idx uint) (int, float64, error)
	GetDeviceECCInfo(idx uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error)
	GetDeviceFanSpeed(idx uint) (int, error)
	GetDeviceHealth(idx uint) (int, error)
	GetDeviceImageCodecUtil(idx uint) ([]int, error)
	GetDeviceMemory(idx uint) (int64, int64, int64, int64, error)
	GetDeviceMLULinkCapability(idx, link uint) (uint, uint, error)
	GetDeviceMLULinkCounter(idx, link uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error)
	GetDeviceMLULinkPortMode(idx, link uint) (int, error)
	GetDeviceMLULinkPortNumber(idx uint) int
	GetDeviceMLULinkSpeedInfo(idx, link uint) (float32, int, error)
	GetDeviceMLULinkStatus(idx, link uint) (int, int, error)
	GetDeviceMLULinkVersion(idx, link uint) (uint, uint, uint, error)
	GetDeviceModel(idx uint) string
	GetDeviceNUMANodeID(idx uint) (int, error)
	GetDevicePCIeInfo(idx uint) (int, uint, uint, uint16, uint16, uint, uint, uint, uint, error)
	GetDevicePower(idx uint) (int, error)
	GetDeviceProcessUtil(idx uint) ([]uint32, []uint32, []uint32, []uint32, []uint32, []uint32, error)
	GetDeviceSN(idx uint) (string, error)
	GetDeviceTemperature(idx uint) (int, []int, error)
	GetDeviceTinyCoreUtil(idx uint) ([]int, error)
	GetDeviceUtil(idx uint) (int, []int, error)
	GetDeviceUUID(idx uint) (string, error)
	GetDeviceVersion(idx uint) (uint, uint, uint, uint, uint, uint, error)
	GetDeviceVfState(idx uint) (int, error)
	GetDeviceVideoCodecUtil(idx uint) ([]int, error)
}

type cndev struct{}

func NewCndevClient() Cndev {
	return &cndev{}
}

func (c *cndev) Init() error {
	r := dl.cndevInit()
	if r == C.CNDEV_ERROR_UNINITIALIZED {
		return errors.New("could not load CNDEV library")
	}
	return errorString(r)
}

func (c *cndev) Release() error {
	r := dl.cndevRelease()
	return errorString(r)
}

func (c *cndev) GetDeviceArmOsMemory(idx uint) (int64, int64, error) {
	var cardArmOsMemInfo C.cndevArmOsMemoryInfo_t
	cardArmOsMemInfo.version = C.int(version)
	r := C.cndevGetArmOsMemoryUsage(&cardArmOsMemInfo, C.int(idx))
	return int64(cardArmOsMemInfo.armOsMemoryUsed), int64(cardArmOsMemInfo.armOsMemoryTotal), errorString(r)
}

func (c *cndev) GetDeviceClusterCount(idx uint) (int, error) {
	var cardClusterNum C.cndevCardClusterCount_t
	cardClusterNum.version = C.int(version)
	r := C.cndevGetClusterCount(&cardClusterNum, C.int(idx))
	return int(cardClusterNum.count), errorString(r)
}

func (c *cndev) GetDeviceCoreNum(idx uint) (uint, error) {
	var cardCoreNum C.cndevCardCoreCount_t
	cardCoreNum.version = C.int(version)
	r := C.cndevGetCoreCount(&cardCoreNum, C.int(idx))
	return uint(cardCoreNum.count), errorString(r)
}

func (c *cndev) GetDeviceCount() (uint, error) {
	var cardInfo C.cndevCardInfo_t
	cardInfo.version = C.int(version)
	r := C.cndevGetDeviceCount(&cardInfo)
	return uint(cardInfo.number), errorString(r)
}

func (c *cndev) GetDeviceCPUUtil(idx uint) (uint16, []uint8, error) {
	var cardCPUUtil C.cndevDeviceCPUUtilization_t
	cardCPUUtil.version = C.int(version)
	r := C.cndevGetDeviceCPUUtilization(&cardCPUUtil, C.int(idx))
	if err := errorString(r); err != nil {
		return 0, nil, err
	}
	coreCPUUtil := make([]uint8, uint8(cardCPUUtil.coreNumber))
	for i := uint8(0); i < uint8(cardCPUUtil.coreNumber); i++ {
		coreCPUUtil[i] = uint8(cardCPUUtil.coreUtilization[i])
	}
	return uint16(cardCPUUtil.chipUtilization), coreCPUUtil, nil
}

func (c *cndev) GetDeviceCRCInfo(idx uint) (uint64, uint64, error) {
	var cardCRCInfo C.cndevCRCInfo_t
	cardCRCInfo.version = C.int(version)
	r := C.cndevGetCRCInfo(&cardCRCInfo, C.int(idx))
	return uint64(cardCRCInfo.die2dieCRCError), uint64(cardCRCInfo.die2dieCRCErrorOverflow), errorString(r)
}

func (c *cndev) GetDeviceCurrentPCIeInfo(idx uint) (int, int, error) {
	var currentPCIeInfo C.cndevCurrentPCIInfo_t
	currentPCIeInfo.version = C.int(version)
	r := C.cndevGetCurrentPCIInfo(&currentPCIeInfo, C.int(idx))
	return int(currentPCIeInfo.currentSpeed), int(currentPCIeInfo.currentWidth), errorString(r)
}

func (c *cndev) GetDeviceDDRInfo(idx uint) (int, float64, error) {
	var cardDDRInfo C.cndevDDRInfo_t
	cardDDRInfo.version = C.int(version)
	r := C.cndevGetDDRInfo(&cardDDRInfo, C.int(idx))
	if err := errorString(r); err != nil {
		return 0, 0, err
	}
	bandWidth, err := strconv.ParseFloat(fmt.Sprintf("%d.%d", int(cardDDRInfo.bandWidth), int(cardDDRInfo.bandWidthDecimal)), 64)
	if err != nil {
		return 0, 0, err
	}
	return int(cardDDRInfo.dataWidth), bandWidth, nil
}

func (c *cndev) GetDeviceECCInfo(idx uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error) {
	var cardECCInfo C.cndevECCInfo_t
	cardECCInfo.version = C.int(version)
	r := C.cndevGetECCInfo(&cardECCInfo, C.int(idx))
	return uint64(cardECCInfo.addressForbiddenError), uint64(cardECCInfo.correctedError), uint64(cardECCInfo.multipleError), uint64(cardECCInfo.multipleMultipleError), uint64(cardECCInfo.multipleOneError), uint64(cardECCInfo.oneBitError), uint64(cardECCInfo.totalError), uint64(cardECCInfo.uncorrectedError), errorString(r)
}

func (c *cndev) GetDeviceFanSpeed(idx uint) (int, error) {
	var cardFanSpeed C.cndevFanSpeedInfo_t
	cardFanSpeed.version = C.int(version)
	r := C.cndevGetFanSpeedInfo(&cardFanSpeed, C.int(idx))
	return int(cardFanSpeed.fanSpeed), errorString(r)
}

func (c *cndev) GetDeviceHealth(idx uint) (int, error) {
	var cardHealthState C.cndevCardHealthState_t
	cardHealthState.version = C.int(version)
	r := C.cndevGetCardHealthState(&cardHealthState, C.int(idx))
	return int(cardHealthState.health), errorString(r)
}

func (c *cndev) GetDeviceImageCodecUtil(idx uint) ([]int, error) {
	var cardImageCodecUtil C.cndevImageCodecUtilization_t
	cardImageCodecUtil.version = C.int(version)
	r := C.cndevGetImageCodecUtilization(&cardImageCodecUtil, C.int(idx))
	if err := errorString(r); err != nil {
		return nil, err
	}
	imageCodecUtil := make([]int, uint(cardImageCodecUtil.jpuCount))
	for i := uint(0); i < uint(cardImageCodecUtil.jpuCount); i++ {
		imageCodecUtil[i] = int(cardImageCodecUtil.jpuCodecUtilization[i])
	}
	return imageCodecUtil, nil
}

func (c *cndev) GetDeviceMemory(idx uint) (int64, int64, int64, int64, error) {
	var cardMemInfo C.cndevMemoryInfo_t
	cardMemInfo.version = C.int(version)
	r := C.cndevGetMemoryUsage(&cardMemInfo, C.int(idx))
	return int64(cardMemInfo.physicalMemoryUsed), int64(cardMemInfo.physicalMemoryTotal), int64(cardMemInfo.virtualMemoryUsed), int64(cardMemInfo.virtualMemoryTotal), errorString(r)
}

func (c *cndev) GetDeviceMLULinkCapability(idx, link uint) (uint, uint, error) {
	var cardMLULinkCapability C.cndevMLULinkCapability_t
	cardMLULinkCapability.version = C.int(version)
	r := C.cndevGetMLULinkCapability(&cardMLULinkCapability, C.int(idx), C.int(link))
	return uint(cardMLULinkCapability.p2pTransfer), uint(cardMLULinkCapability.interlakenSerdes), errorString(r)
}

func (c *cndev) GetDeviceMLULinkCounter(idx, link uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error) {
	var cardMLULinkCount C.cndevMLULinkCounter_t
	cardMLULinkCount.version = C.int(version)
	r := C.cndevGetMLULinkCounter(&cardMLULinkCount, C.int(idx), C.int(link))
	return uint64(cardMLULinkCount.cntrReadByte), uint64(cardMLULinkCount.cntrReadPackage), uint64(cardMLULinkCount.cntrWriteByte), uint64(cardMLULinkCount.cntrWritePackage),
		uint64(cardMLULinkCount.errCorrected), uint64(cardMLULinkCount.errCRC24), uint64(cardMLULinkCount.errCRC32), uint64(cardMLULinkCount.errEccDouble),
		uint64(cardMLULinkCount.errFatal), uint64(cardMLULinkCount.errReplay), uint64(cardMLULinkCount.errUncorrected), errorString(r)
}

func (c *cndev) GetDeviceMLULinkPortMode(idx, link uint) (int, error) {
	var cardMLULinkPortMode C.cndevMLULinkPortMode_t
	cardMLULinkPortMode.version = C.int(version)
	r := C.cndevGetMLULinkPortMode(&cardMLULinkPortMode, C.int(idx), C.int(link))
	return int(cardMLULinkPortMode.mode), errorString(r)
}

func (c *cndev) GetDeviceMLULinkPortNumber(idx uint) int {
	return int(C.cndevGetMLULinkPortNumber(C.int(idx)))
}

func (c *cndev) GetDeviceMLULinkSpeedInfo(idx, link uint) (float32, int, error) {
	var cardMLULinkSpeedInfo C.cndevMLULinkSpeed_t
	cardMLULinkSpeedInfo.version = C.int(version)
	r := C.cndevGetMLULinkSpeedInfo(&cardMLULinkSpeedInfo, C.int(idx), C.int(link))
	return float32(cardMLULinkSpeedInfo.speedValue), int(cardMLULinkSpeedInfo.speedFormat), errorString(r)
}

func (c *cndev) GetDeviceMLULinkStatus(idx, link uint) (int, int, error) {
	var cardMLULinkStatus C.cndevMLULinkStatus_t
	cardMLULinkStatus.version = C.int(version)
	r := C.cndevGetMLULinkStatus(&cardMLULinkStatus, C.int(idx), C.int(link))
	return int(cardMLULinkStatus.isActive), int(cardMLULinkStatus.serdesState), errorString(r)
}

func (c *cndev) GetDeviceMLULinkVersion(idx, link uint) (uint, uint, uint, error) {
	var cardMLULinkVersion C.cndevMLULinkVersion_t
	cardMLULinkVersion.version = C.int(version)
	r := C.cndevGetMLULinkVersion(&cardMLULinkVersion, C.int(idx), C.int(link))
	return uint(cardMLULinkVersion.majorVersion), uint(cardMLULinkVersion.minorVersion), uint(cardMLULinkVersion.buildVersion), errorString(r)
}

func (c *cndev) GetDeviceModel(idx uint) string {
	return C.GoString(C.getCardNameStringByDevId(C.int(idx)))
}

func (c *cndev) GetDeviceNUMANodeID(idx uint) (int, error) {
	var cardNUMANodeID C.cndevNUMANodeId_t
	cardNUMANodeID.version = C.int(version)
	r := C.cndevGetNUMANodeIdByDevId(&cardNUMANodeID, C.int(idx))
	return int(cardNUMANodeID.nodeId), errorString(r)
}

func (c *cndev) GetDevicePCIeInfo(idx uint) (int, uint, uint, uint16, uint16, uint, uint, uint, uint, error) {
	var pcieInfo C.cndevPCIeInfo_t
	pcieInfo.version = C.int(version)
	r := C.cndevGetPCIeInfo(&pcieInfo, C.int(idx))
	return int(pcieInfo.slotID), uint(pcieInfo.subsystemId), uint(pcieInfo.deviceId), uint16(pcieInfo.vendor),
		uint16(pcieInfo.subsystemVendor), uint(pcieInfo.domain), uint(pcieInfo.bus), uint(pcieInfo.device),
		uint(pcieInfo.function), errorString(r)
}

func (c *cndev) GetDevicePower(idx uint) (int, error) {
	var cardPower C.cndevPowerInfo_t
	cardPower.version = C.int(version)
	r := C.cndevGetPowerInfo(&cardPower, C.int(idx))
	return int(cardPower.usage), errorString(r)
}

func (c *cndev) GetDeviceProcessUtil(idx uint) ([]uint32, []uint32, []uint32, []uint32, []uint32, []uint32, error) {
	processCount := 10 // maximum number of processes running on an MLU is 10
	var util C.cndevProcessUtilization_t
	utils := (*C.cndevProcessUtilization_t)(C.malloc(C.size_t(processCount) * C.size_t(unsafe.Sizeof(util))))
	defer C.free(unsafe.Pointer(utils))
	utils.version = C.uint(version)
	r := C.cndevGetProcessUtilization((*C.uint)(unsafe.Pointer(&processCount)), utils, C.int(idx))
	if err := errorString(r); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	pids := make([]uint32, processCount)
	ipuUtils := make([]uint32, processCount)
	jpuUtils := make([]uint32, processCount)
	memUtils := make([]uint32, processCount)
	vpuDecUtils := make([]uint32, processCount)
	vpuEncUtils := make([]uint32, processCount)
	array := (*[10]C.cndevProcessUtilization_t)(unsafe.Pointer(utils))
	results := array[:processCount]
	for i := 0; i < processCount; i++ {
		pids[i] = uint32(results[i].pid)
		ipuUtils[i] = uint32(results[i].ipuUtil)
		jpuUtils[i] = uint32(results[i].jpuUtil)
		memUtils[i] = uint32(results[i].memUtil)
		vpuDecUtils[i] = uint32(results[i].vpuDecUtil)
		vpuEncUtils[i] = uint32(results[i].vpuEncUtil)
	}
	return pids, ipuUtils, jpuUtils, memUtils, vpuDecUtils, vpuEncUtils, nil
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

func (c *cndev) GetDeviceTemperature(idx uint) (int, []int, error) {
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
	clusterTemperature := make([]int, cluster)
	for i := 0; i < cluster; i++ {
		clusterTemperature[i] = int(cardTemperature.cluster[i])
	}
	return int(cardTemperature.board), clusterTemperature, nil
}

func (c *cndev) GetDeviceTinyCoreUtil(idx uint) ([]int, error) {
	var cardTinyCoreUtil C.cndevTinyCoreUtilization_t
	cardTinyCoreUtil.version = C.int(version)
	r := C.cndevGetTinyCoreUtilization(&cardTinyCoreUtil, C.int(idx))
	if err := errorString(r); err != nil {
		return nil, err
	}
	tinyCoreUtil := make([]int, uint(cardTinyCoreUtil.tinyCoreCount))
	for i := uint(0); i < uint(cardTinyCoreUtil.tinyCoreCount); i++ {
		tinyCoreUtil[i] = int(cardTinyCoreUtil.tinyCoreUtilization[i])
	}
	return tinyCoreUtil, nil
}

func (c *cndev) GetDeviceUtil(idx uint) (int, []int, error) {
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
	coreUtil := make([]int, core)
	for i := uint(0); i < core; i++ {
		coreUtil[i] = int(cardUtil.coreUtilization[i])
	}
	return int(cardUtil.averageCoreUtilization), coreUtil, nil
}

func (c *cndev) GetDeviceUUID(idx uint) (string, error) {
	var uuidInfo C.cndevUUID_t
	uuidInfo.version = C.int(version)
	r := C.cndevGetUUID(&uuidInfo, C.int(idx))
	if err := errorString(r); err != nil {
		return "", err
	}
	return C.GoString((*C.char)(unsafe.Pointer(&uuidInfo.uuid))), nil
}

func (c *cndev) GetDeviceVersion(idx uint) (uint, uint, uint, uint, uint, uint, error) {
	var versionInfo C.cndevVersionInfo_t
	versionInfo.version = C.int(version)
	r := C.cndevGetVersionInfo(&versionInfo, C.int(idx))
	return uint(versionInfo.mcuMajorVersion), uint(versionInfo.mcuMinorVersion), uint(versionInfo.mcuBuildVersion), uint(versionInfo.driverMajorVersion), uint(versionInfo.driverMinorVersion), uint(versionInfo.driverBuildVersion), errorString(r)
}

// GetDeviceVfState returns cndevCardVfState_t.vfState, which is bitmap of mlu vf state. See cndev user doc for more details.
func (c *cndev) GetDeviceVfState(idx uint) (int, error) {
	var vfstate C.cndevCardVfState_t
	vfstate.version = C.int(version)
	r := C.cndevGetCardVfState(&vfstate, C.int(idx))
	return int(vfstate.vfState), errorString(r)
}

func (c *cndev) GetDeviceVideoCodecUtil(idx uint) ([]int, error) {
	var cardVideoCodecUtil C.cndevVideoCodecUtilization_t
	cardVideoCodecUtil.version = C.int(version)
	r := C.cndevGetVideoCodecUtilization(&cardVideoCodecUtil, C.int(idx))
	if err := errorString(r); err != nil {
		return nil, err
	}
	videoCodecUtil := make([]int, uint(cardVideoCodecUtil.vpuCount))
	for i := uint(0); i < uint(cardVideoCodecUtil.vpuCount); i++ {
		videoCodecUtil[i] = int(cardVideoCodecUtil.vpuCodecUtilization[i])
	}
	return videoCodecUtil, nil
}

func errorString(ret C.cndevRet_t) error {
	if ret == C.CNDEV_SUCCESS {
		return nil
	}
	err := C.GoString(C.cndevGetErrorString(ret))
	return fmt.Errorf("cndev: %v", err)
}
