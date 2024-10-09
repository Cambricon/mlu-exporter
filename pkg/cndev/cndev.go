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

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
)

type MimInfo struct {
	InstanceID     int
	Name           string
	PlacementSize  int
	PlacementStart int
	UUID           string
}

type SmluInfo struct {
	InstanceID int
	IpuUtil    int
	MemUsed    int64
	MemTotal   int64
	Name       string
	UUID       string
}

type DeviceInfo struct {
	Device            uint
	ComputeInstanceID int32
	MLUInstanceID     int32
}

type XIDInfo struct {
	XID int64
	DeviceInfo
}

type XIDInfoWithTimestamp struct {
	Timestamp string
	XIDInfo
}

type Cndev interface {
	Init(healthCheck bool) error
	Release() error

	DeviceMimModeEnabled(idx uint) (bool, error)
	DeviceSmluModeEnabled(idx uint) (bool, error)
	GetAllMLUInstanceInfo(idx uint) ([]MimInfo, error)
	GetAllSMluInfo(idx uint) ([]SmluInfo, error)
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
	GetDeviceHeartbeatCount(idx uint) (uint32, error)
	GetDeviceImageCodecUtil(idx uint) ([]int, error)
	GetDeviceMemEccCounter(idx uint) (uint32, uint32, uint32, uint32, uint32, error)
	GetDeviceMemory(idx uint) (int64, int64, int64, int64, error)
	GetDeviceMemoryDieCount(idx uint) (int, error)
	GetDeviceMLULinkCapability(idx, link uint) (uint, uint, error)
	GetDeviceMLULinkCounter(idx, link uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error)
	GetDeviceMLULinkPortMode(idx, link uint) (int, error)
	GetDeviceMLULinkPortNumber(idx uint) int
	GetDeviceMLULinkSpeedInfo(idx, link uint) (float32, int, error)
	GetDeviceMLULinkStatus(idx, link uint) (int, int, error)
	GetDeviceMLULinkVersion(idx, link uint) (uint, uint, uint, error)
	GetDeviceMLUInstanceByID(idx uint, instanceID int) (uint, error)
	GetDeviceModel(idx uint) string
	GetDeviceNUMANodeID(idx uint) (int, error)
	GetDeviceParityError(idx uint) (int, error)
	GetDevicePCIeInfo(idx uint) (int, uint, uint, uint16, uint16, uint, uint, uint, uint, error)
	GetDevicePCIeReplayCount(idx uint) (uint32, error)
	GetDevicePCIeThroughput(idx uint) (int64, int64, error)
	GetDevicePower(idx uint) (int, error)
	GetDeviceProcessUtil(idx uint) ([]uint32, []uint32, []uint32, []uint32, []uint32, []uint32, error)
	GetDeviceRemappedRows(idx uint) (uint32, uint32, uint32, uint32, error)
	GetDeviceRetiredPageInfo(idx uint) (int, uint32, error)
	GetDeviceSN(idx uint) (string, error)
	GetDeviceTemperature(idx uint) (int, int, int, []int, []int, error)
	GetDeviceTinyCoreUtil(idx uint) ([]int, error)
	GetDeviceUtil(idx uint) (int, []int, error)
	GetDeviceFrequency(idx uint) (int, int, error)
	GetDeviceUUID(idx uint) (string, error)
	GetDeviceVersion(idx uint) (uint, uint, uint, uint, uint, uint, error)
	GetDeviceVfState(idx uint) (int, error)
	GetDeviceVideoCodecUtil(idx uint) ([]int, error)
	RegisterEventsHandleAndWait(slots []int, ch chan XIDInfoWithTimestamp) error
	GetSupportedEventTypes(idx uint) error
}

type cndev struct {
	cndevHandleMap map[uint]C.cndevDevice_t
}

func NewCndevClient() Cndev {
	return &cndev{
		cndevHandleMap: map[uint]C.cndevDevice_t{},
	}
}

func (c *cndev) Init(healthCheck bool) error {
	r := dl.cndevInit()
	if r == C.CNDEV_ERROR_UNINITIALIZED {
		return errors.New("could not load CNDEV library")
	}
	if healthCheck {
		return errorString(r)
	}
	return c.generateDeviceHandleMap()
}

func (c *cndev) Release() error {
	r := dl.cndevRelease()
	return errorString(r)
}

func (c *cndev) DeviceMimModeEnabled(idx uint) (bool, error) {
	var mode C.cndevMimMode_t
	mode.version = C.CNDEV_VERSION_5
	r := C.cndevGetMimMode(&mode, c.cndevHandleMap[idx])
	return mode.mimMode == C.CNDEV_FEATURE_ENABLED, errorString(r)
}

func (c *cndev) DeviceSmluModeEnabled(idx uint) (bool, error) {
	var mode C.cndevSMLUMode_t
	mode.version = C.CNDEV_VERSION_5
	r := C.cndevGetSMLUMode(&mode, c.cndevHandleMap[idx])
	return mode.smluMode == C.CNDEV_FEATURE_ENABLED, errorString(r)
}

func (c *cndev) GetAllMLUInstanceInfo(idx uint) ([]MimInfo, error) {
	miCount := C.int(1 << 4)
	var miInfos []MimInfo
	var miInfo C.cndevMluInstanceInfo_t
	infs := (*C.cndevMluInstanceInfo_t)(C.malloc(C.size_t(miCount) * C.size_t(unsafe.Sizeof(miInfo))))
	defer func() {
		if infs != nil {
			C.free(unsafe.Pointer(infs))
		}
	}()

	infs.version = C.CNDEV_VERSION_6
	r := C.cndevGetAllMluInstanceInfo(&miCount, infs, c.cndevHandleMap[idx])
	if errorString(r) != nil && r != C.CNDEV_ERROR_INSUFFICIENT_SPACE {
		return miInfos, errorString(r)
	}

	// handle the case when the initial count is insufficient,
	// after cndevGetAllMluInstanceInfo miCount will be set to real count,
	// so just use it to reallocate and cndevGetAllMluInstanceInfo again.
	if r == C.CNDEV_ERROR_INSUFFICIENT_SPACE {
		log.Debugf("cndevGetAllMluInstanceInfo with insufficient space with real counts %d, with slot: %d, will try with the real counts", miCount, idx)
		newInfs := (*C.cndevMluInstanceInfo_t)(C.realloc(unsafe.Pointer(infs), C.size_t(miCount)*C.size_t(unsafe.Sizeof(miInfo))))
		if newInfs == nil {
			return miInfos, fmt.Errorf("realloc failed for cndevGetAllMluInstanceInfo")
		}
		infs = newInfs
		r := C.cndevGetAllMluInstanceInfo(&miCount, infs, c.cndevHandleMap[idx])
		if errorString(r) != nil {
			return miInfos, errorString(r)
		}
	}

	infos := (*[1 << 16]C.cndevMluInstanceInfo_t)(unsafe.Pointer(infs))[:miCount]
	for i := 0; i < int(miCount); i++ {
		info := infos[i]
		miInfos = append(miInfos,
			MimInfo{
				Name:           C.GoString((*C.char)(unsafe.Pointer(&info.profileName))),
				UUID:           C.GoString((*C.char)(unsafe.Pointer(&info.uuid))),
				InstanceID:     int(info.instanceId),
				PlacementStart: int(info.placement.start),
				PlacementSize:  int(info.placement.size),
			})
	}
	log.Debugf("mim infos for device %d are %+v", idx, miInfos)

	return miInfos, nil
}

func (c *cndev) GetAllSMluInfo(idx uint) ([]SmluInfo, error) {
	smluCount := C.int(1 << 7)
	var smluInfos []SmluInfo
	var smluInfo C.cndevSMluInfo_t
	infs := (*C.cndevSMluInfo_t)(C.malloc(C.size_t(smluCount) * C.size_t(unsafe.Sizeof(smluInfo))))
	defer func() {
		if infs != nil {
			C.free(unsafe.Pointer(infs))
		}
	}()

	infs.version = C.CNDEV_VERSION_6
	r := C.cndevGetAllSMluInstanceInfo(&smluCount, infs, c.cndevHandleMap[idx])
	if errorString(r) != nil && r != C.CNDEV_ERROR_INSUFFICIENT_SPACE {
		return smluInfos, errorString(r)
	}

	// handle the case when the initial count is insufficient,
	// after cndevGetAllSMluInstanceInfo smluCount will be set to real count,
	// so just use it to reallocate and cndevGetAllSMluInstanceInfo again.
	if r == C.CNDEV_ERROR_INSUFFICIENT_SPACE {
		log.Debugf("cndevGetAllSMluInstanceInfo with insufficient space with real counts %d, with slot: %d, will try with the real counts", smluCount, idx)
		newInfs := (*C.cndevSMluInfo_t)(C.realloc(unsafe.Pointer(infs), C.size_t(smluCount)*C.size_t(unsafe.Sizeof(smluInfo))))
		if newInfs == nil {
			return smluInfos, fmt.Errorf("realloc failed for cndevGetAllSMluInstanceInfo")
		}
		infs = newInfs
		r := C.cndevGetAllSMluInstanceInfo(&smluCount, infs, c.cndevHandleMap[idx])
		if errorString(r) != nil {
			return smluInfos, errorString(r)
		}
	}

	infos := (*[1 << 16]C.cndevSMluInfo_t)(unsafe.Pointer(infs))[:smluCount]
	for i := 0; i < int(smluCount); i++ {
		info := infos[i]
		smluInfos = append(smluInfos,
			SmluInfo{
				InstanceID: int(info.instanceId),
				Name:       C.GoString((*C.char)(unsafe.Pointer(&info.profileName))),
				IpuUtil:    int(info.mluQuota[C.CNDEV_SMLU_USAGE]),
				MemUsed:    int64(info.memorySize[C.CNDEV_SMLU_USAGE]),
				MemTotal:   int64(info.memorySize[C.CNDEV_SMLU_MAX]),
				UUID:       C.GoString((*C.char)(unsafe.Pointer(&info.uuid))),
			})
	}
	log.Debugf("smlu infos for device %d are %+v", idx, smluInfos)

	return smluInfos, nil
}

func (c *cndev) GetDeviceArmOsMemory(idx uint) (int64, int64, error) {
	type armOsMem struct {
		version     int
		memoryTotal int64
		memoryUsed  int64
	}
	var armOsMemInfo armOsMem
	memInfo := (*C.cndevArmOsMemoryInfo_t)(C.malloc(C.size_t(unsafe.Sizeof(armOsMemInfo))))
	defer C.free(unsafe.Pointer(memInfo))
	memInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetArmOsMemoryUsage(memInfo, c.cndevHandleMap[idx])
	armOsMemoryTotal := *(*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(memInfo)) + unsafe.Sizeof(armOsMemInfo.version)))
	armOsMemoryUsed := *(*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(memInfo)) + unsafe.Sizeof(armOsMemInfo.version) + unsafe.Sizeof(armOsMemInfo.memoryTotal)))
	return int64(armOsMemoryUsed), int64(armOsMemoryTotal), errorString(r)
}

func (c *cndev) GetDeviceClusterCount(idx uint) (int, error) {
	var cardClusterNum C.cndevCardClusterCount_t
	cardClusterNum.version = C.CNDEV_VERSION_5
	r := C.cndevGetClusterCount(&cardClusterNum, c.cndevHandleMap[idx])
	return int(cardClusterNum.count), errorString(r)
}

func (c *cndev) GetDeviceCoreNum(idx uint) (uint, error) {
	var cardCoreNum C.cndevCardCoreCount_t
	cardCoreNum.version = C.CNDEV_VERSION_5
	r := C.cndevGetCoreCount(&cardCoreNum, c.cndevHandleMap[idx])
	return uint(cardCoreNum.count), errorString(r)
}

func (c *cndev) GetDeviceCount() (uint, error) {
	var cardInfo C.cndevCardInfo_t
	cardInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetDeviceCount(&cardInfo)
	return uint(cardInfo.number), errorString(r)
}

func (c *cndev) GetDeviceCPUUtil(idx uint) (uint16, []uint8, error) {
	var cardCPUUtil C.cndevDeviceCPUUtilization_t
	cardCPUUtil.version = C.CNDEV_VERSION_5
	r := C.cndevGetDeviceCPUUtilization(&cardCPUUtil, c.cndevHandleMap[idx])
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
	cardCRCInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetCRCInfo(&cardCRCInfo, c.cndevHandleMap[idx])
	return uint64(cardCRCInfo.die2dieCRCError), uint64(cardCRCInfo.die2dieCRCErrorOverflow), errorString(r)
}

func (c *cndev) GetDeviceCurrentPCIeInfo(idx uint) (int, int, error) {
	var currentPCIeInfo C.cndevCurrentPCIInfo_t
	currentPCIeInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetCurrentPCIInfo(&currentPCIeInfo, c.cndevHandleMap[idx])
	return int(currentPCIeInfo.currentSpeed), int(currentPCIeInfo.currentWidth), errorString(r)
}

func (c *cndev) GetDeviceDDRInfo(idx uint) (int, float64, error) {
	var cardDDRInfo C.cndevDDRInfo_t
	cardDDRInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetDDRInfo(&cardDDRInfo, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return 0, 0, err
	}
	bandWidth, err := strconv.ParseFloat(fmt.Sprintf("%d.%d", int(cardDDRInfo.bandWidth), int(cardDDRInfo.bandWidthDecimal)), 64)
	if err != nil {
		return 0, 0, err
	}
	return int(cardDDRInfo.dataWidth), bandWidth, nil
}

func (c *cndev) GetDeviceDriverVersion(idx uint) (uint, uint, uint, error) {
	var cardVersionInfo C.cndevVersionInfo_t
	cardVersionInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetVersionInfo(&cardVersionInfo, c.cndevHandleMap[idx])
	return uint(cardVersionInfo.driverMajorVersion), uint(cardVersionInfo.driverMinorVersion), uint(cardVersionInfo.driverBuildVersion), errorString(r)
}

func (c *cndev) GetDeviceECCInfo(idx uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error) {
	var cardECCInfo C.cndevECCInfo_t
	cardECCInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetECCInfo(&cardECCInfo, c.cndevHandleMap[idx])
	return uint64(cardECCInfo.addressForbiddenError), uint64(cardECCInfo.correctedError), uint64(cardECCInfo.multipleError), uint64(cardECCInfo.multipleMultipleError), uint64(cardECCInfo.multipleOneError), uint64(cardECCInfo.oneBitError), uint64(cardECCInfo.totalError), uint64(cardECCInfo.uncorrectedError), errorString(r)
}

func (c *cndev) GetDeviceFanSpeed(idx uint) (int, error) {
	var cardFanSpeed C.cndevFanSpeedInfo_t
	cardFanSpeed.version = C.CNDEV_VERSION_5
	r := C.cndevGetFanSpeedInfo(&cardFanSpeed, c.cndevHandleMap[idx])
	return int(cardFanSpeed.fanSpeed), errorString(r)
}

func (c *cndev) GetDeviceFrequency(idx uint) (int, int, error) {
	var cardFrequency C.cndevFrequencyInfo_t
	cardFrequency.version = C.CNDEV_VERSION_5
	r := C.cndevGetFrequencyInfo(&cardFrequency, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return 0, 0, err
	}
	return int(cardFrequency.boardFreq), int(cardFrequency.ddrFreq), nil
}

func (c *cndev) GetDeviceHealth(idx uint) (int, error) {
	var cardHealthState C.cndevCardHealthState_t
	cardHealthState.version = C.CNDEV_VERSION_5
	r := C.cndevGetCardHealthState(&cardHealthState, c.cndevHandleMap[idx])
	return int(cardHealthState.health), errorString(r)
}

func (c *cndev) GetDeviceHeartbeatCount(idx uint) (uint32, error) {
	var heartbeatCount C.cndevCardHeartbeatCount_t
	heartbeatCount.version = C.CNDEV_VERSION_5
	index := C.int(idx)
	if id, ok := c.cndevHandleMap[idx]; ok {
		index = id
	}
	r := C.cndevGetCardHeartbeatCount(&heartbeatCount, index)
	if err := errorString(r); err != nil {
		return 0, err
	}
	return uint32(heartbeatCount.heartbeatCount), nil
}

func (c *cndev) GetDeviceImageCodecUtil(idx uint) ([]int, error) {
	var cardImageCodecUtil C.cndevImageCodecUtilization_t
	cardImageCodecUtil.version = C.CNDEV_VERSION_5
	r := C.cndevGetImageCodecUtilization(&cardImageCodecUtil, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return nil, err
	}
	imageCodecUtil := make([]int, uint(cardImageCodecUtil.jpuCount))
	for i := uint(0); i < uint(cardImageCodecUtil.jpuCount); i++ {
		imageCodecUtil[i] = int(cardImageCodecUtil.jpuCodecUtilization[i])
	}
	return imageCodecUtil, nil
}

func (c *cndev) GetDeviceMemEccCounter(idx uint) (uint32, uint32, uint32, uint32, uint32, error) {
	var memEccCounter C.cndevMemEccCounter_t
	memEccCounter.version = C.CNDEV_VERSION_5
	index := C.int(idx)
	if id, ok := c.cndevHandleMap[idx]; ok {
		index = id
	}
	r := C.cndevGetMemEccCounter(&memEccCounter, index)
	if err := errorString(r); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	return uint32(memEccCounter.volatileSramEccSbeCounter), uint32(memEccCounter.volatileSramEccDbeCounter), uint32(memEccCounter.volatileSramEccParityCounter), uint32(memEccCounter.volatileDramEccSbeCounter), uint32(memEccCounter.volatileDramEccDbeCounter), nil
}

func (c *cndev) GetDeviceMemory(idx uint) (int64, int64, int64, int64, error) {
	var cardMemInfo C.cndevMemoryInfo_t
	cardMemInfo.version = C.CNDEV_VERSION_5
	index := C.int(idx)
	if id, ok := c.cndevHandleMap[idx]; ok {
		index = id
	}
	r := C.cndevGetMemoryUsage(&cardMemInfo, index)
	return int64(cardMemInfo.physicalMemoryUsed), int64(cardMemInfo.physicalMemoryTotal), int64(cardMemInfo.virtualMemoryUsed), int64(cardMemInfo.virtualMemoryTotal), errorString(r)
}

func (c *cndev) GetDeviceMemoryDieCount(idx uint) (int, error) {
	var cardMemoryDieNum C.cndevCardMemoryDieCount_t
	cardMemoryDieNum.version = C.CNDEV_VERSION_5
	r := C.cndevGetMemoryDieCount(&cardMemoryDieNum, c.cndevHandleMap[idx])
	return int(cardMemoryDieNum.count), errorString(r)
}

func (c *cndev) GetDeviceMLULinkCapability(idx, link uint) (uint, uint, error) {
	var cardMLULinkCapability C.cndevMLULinkCapability_t
	cardMLULinkCapability.version = C.CNDEV_VERSION_5
	r := C.cndevGetMLULinkCapability(&cardMLULinkCapability, c.cndevHandleMap[idx], C.int(link))
	return uint(cardMLULinkCapability.p2pTransfer), uint(cardMLULinkCapability.interlakenSerdes), errorString(r)
}

func (c *cndev) GetDeviceMLULinkCounter(idx, link uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error) {
	var cardMLULinkCount C.cndevMLULinkCounter_t
	cardMLULinkCount.version = C.CNDEV_VERSION_5
	r := C.cndevGetMLULinkCounter(&cardMLULinkCount, c.cndevHandleMap[idx], C.int(link))
	return uint64(cardMLULinkCount.cntrReadByte), uint64(cardMLULinkCount.cntrReadPackage), uint64(cardMLULinkCount.cntrWriteByte), uint64(cardMLULinkCount.cntrWritePackage),
		uint64(cardMLULinkCount.errCorrected), uint64(cardMLULinkCount.errCRC24), uint64(cardMLULinkCount.errCRC32), uint64(cardMLULinkCount.errEccDouble),
		uint64(cardMLULinkCount.errFatal), uint64(cardMLULinkCount.errReplay), uint64(cardMLULinkCount.errUncorrected), uint64(cardMLULinkCount.cntrCnpPackage),
		uint64(cardMLULinkCount.cntrPfcPackage), errorString(r)
}

func (c *cndev) GetDeviceMLULinkPortMode(idx, link uint) (int, error) {
	var cardMLULinkPortMode C.cndevMLULinkPortMode_t
	cardMLULinkPortMode.version = C.CNDEV_VERSION_5
	r := C.cndevGetMLULinkPortMode(&cardMLULinkPortMode, c.cndevHandleMap[idx], C.int(link))
	return int(cardMLULinkPortMode.mode), errorString(r)
}

func (c *cndev) GetDeviceMLULinkPortNumber(idx uint) int {
	return int(C.cndevGetMLULinkPortNumber(c.cndevHandleMap[idx]))
}

func (c *cndev) GetDeviceMLULinkSpeedInfo(idx, link uint) (float32, int, error) {
	var cardMLULinkSpeedInfo C.cndevMLULinkSpeed_t
	cardMLULinkSpeedInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetMLULinkSpeedInfo(&cardMLULinkSpeedInfo, c.cndevHandleMap[idx], C.int(link))
	return float32(cardMLULinkSpeedInfo.speedValue), int(cardMLULinkSpeedInfo.speedFormat), errorString(r)
}

func (c *cndev) GetDeviceMLULinkStatus(idx, link uint) (int, int, error) {
	var cardMLULinkStatus C.cndevMLULinkStatus_t
	cardMLULinkStatus.version = C.CNDEV_VERSION_5
	r := C.cndevGetMLULinkStatus(&cardMLULinkStatus, c.cndevHandleMap[idx], C.int(link))
	return int(cardMLULinkStatus.isActive), int(cardMLULinkStatus.serdesState), errorString(r)
}

func (c *cndev) GetDeviceMLULinkVersion(idx, link uint) (uint, uint, uint, error) {
	var cardMLULinkVersion C.cndevMLULinkVersion_t
	cardMLULinkVersion.version = C.CNDEV_VERSION_5
	r := C.cndevGetMLULinkVersion(&cardMLULinkVersion, c.cndevHandleMap[idx], C.int(link))
	return uint(cardMLULinkVersion.majorVersion), uint(cardMLULinkVersion.minorVersion), uint(cardMLULinkVersion.buildVersion), errorString(r)
}

func (c *cndev) GetDeviceMLUInstanceByID(idx uint, instanceID int) (uint, error) {
	var handle C.cndevMluInstance_t
	r := C.cndevGetMluInstanceById(&handle, C.int(instanceID), c.cndevHandleMap[idx])
	return uint(handle), errorString(r)
}

func (c *cndev) GetDeviceModel(idx uint) string {
	return C.GoString(C.cndevGetCardNameStringByDevId(c.cndevHandleMap[idx]))
}

func (c *cndev) GetDeviceNUMANodeID(idx uint) (int, error) {
	var cardNUMANodeID C.cndevNUMANodeId_t
	cardNUMANodeID.version = C.CNDEV_VERSION_5
	r := C.cndevGetNUMANodeIdByDevId(&cardNUMANodeID, c.cndevHandleMap[idx])
	return int(cardNUMANodeID.nodeId), errorString(r)
}

func (c *cndev) GetDeviceParityError(idx uint) (int, error) {
	var parityError C.cndevParityError_t
	parityError.version = C.CNDEV_VERSION_5
	r := C.cndevGetParityError(&parityError, c.cndevHandleMap[idx])
	return int(parityError.counter), errorString(r)
}

func (c *cndev) GetDevicePCIeInfo(idx uint) (int, uint, uint, uint16, uint16, uint, uint, uint, uint, error) {
	var pcieInfo C.cndevPCIeInfo_t
	pcieInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetPCIeInfo(&pcieInfo, c.cndevHandleMap[idx])
	return int(pcieInfo.slotID), uint(pcieInfo.subsystemId), uint(pcieInfo.deviceId), uint16(pcieInfo.vendor),
		uint16(pcieInfo.subsystemVendor), uint(pcieInfo.domain), uint(pcieInfo.bus), uint(pcieInfo.device),
		uint(pcieInfo.function), errorString(r)
}

func (c *cndev) GetDevicePCIeReplayCount(idx uint) (uint32, error) {
	var pcieReplay C.cndevPcieReplayCounter_t
	pcieReplay.version = C.CNDEV_VERSION_5
	r := C.cndevGetPcieReplayCounter(&pcieReplay, c.cndevHandleMap[idx])
	return uint32(pcieReplay.counter), errorString(r)
}

func (c *cndev) GetDevicePCIeThroughput(idx uint) (int64, int64, error) {
	var pcieThroughput C.cndevPCIethroughput_t
	pcieThroughput.version = C.CNDEV_VERSION_5
	r := C.cndevGetPCIethroughput(&pcieThroughput, c.cndevHandleMap[idx])
	return int64(pcieThroughput.pcieRead), int64(pcieThroughput.pcieWrite), errorString(r)
}

func (c *cndev) GetDevicePower(idx uint) (int, error) {
	var cardPower C.cndevPowerInfo_t
	cardPower.version = C.CNDEV_VERSION_5
	index := C.int(idx)
	if id, ok := c.cndevHandleMap[idx]; ok {
		index = id
	}
	r := C.cndevGetPowerInfo(&cardPower, index)
	return int(cardPower.usage), errorString(r)
}

func (c *cndev) GetDeviceProcessUtil(idx uint) ([]uint32, []uint32, []uint32, []uint32, []uint32, []uint32, error) {
	processCount := 1 << 4

	var util C.cndevProcessUtilization_t
	utils := (*C.cndevProcessUtilization_t)(C.malloc(C.size_t(processCount) * C.size_t(unsafe.Sizeof(util))))
	defer func() {
		if utils != nil {
			C.free(unsafe.Pointer(utils))
		}
	}()

	utils.version = C.CNDEV_VERSION_5
	r := C.cndevGetProcessUtilization((*C.uint)(unsafe.Pointer(&processCount)), utils, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil && r != C.CNDEV_ERROR_INSUFFICIENT_SPACE {
		return nil, nil, nil, nil, nil, nil, err
	}

	// handle the case when the initial count is insufficient,
	// after cndevGetProcessUtilization processCount will be set to real count,
	// so just use it to reallocate and cndevGetProcessUtilization again.
	if r == C.CNDEV_ERROR_INSUFFICIENT_SPACE {
		log.Debugf("cndevGetProcessUtilization with insufficient space with real counts %d, with slot: %d, will try with the real counts", processCount, idx)
		newUtils := (*C.cndevProcessUtilization_t)(C.realloc(unsafe.Pointer(utils), C.size_t(processCount)*C.size_t(unsafe.Sizeof(util))))
		if newUtils == nil {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("realloc failed for cndevGetProcessUtilization")
		}
		utils = newUtils
		r = C.cndevGetProcessUtilization((*C.uint)(unsafe.Pointer(&processCount)), utils, c.cndevHandleMap[idx])
		if err := errorString(r); err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}
	}

	pids := make([]uint32, processCount)
	ipuUtils := make([]uint32, processCount)
	jpuUtils := make([]uint32, processCount)
	memUtils := make([]uint32, processCount)
	vpuDecUtils := make([]uint32, processCount)
	vpuEncUtils := make([]uint32, processCount)
	results := (*[1 << 16]C.cndevProcessUtilization_t)(unsafe.Pointer(utils))[:processCount]
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

func (c *cndev) GetDeviceRemappedRows(idx uint) (uint32, uint32, uint32, uint32, error) {
	major, minor, build, err := c.GetDeviceDriverVersion(idx)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	v1, err := semver.NewVersion(fmt.Sprintf("%d.%d.%d", major, minor, build))
	if err != nil {
		return 0, 0, 0, 0, err
	}
	// cndevGetRemappedRowsV2 is only available in driver version later than 6.0.2
	v2 := semver.MustParse("6.0.2")

	if v1.GreaterThan(v2) {
		var remappedRowsV2 C.cndevRemappedRowV2_t
		remappedRowsV2.version = C.CNDEV_VERSION_5
		r := C.cndevGetRemappedRowsV2(&remappedRowsV2, c.cndevHandleMap[idx])
		if err := errorString(r); err != nil {
			return 0, 0, 0, 0, err
		}
		var repairStatus C.cndevRepairStatus_t
		repairStatus.version = C.CNDEV_VERSION_5
		r = C.cndevGetRepairStatus(&repairStatus, c.cndevHandleMap[idx])
		if err := errorString(r); err != nil {
			return 0, 0, 0, 0, err
		}
		var pendingRows, failedRows uint32
		if bool(repairStatus.isPending) {
			pendingRows = 1
		} else {
			pendingRows = 0
		}
		if bool(repairStatus.isFailure) {
			failedRows = 1
		} else {
			failedRows = 0
		}
		return uint32(remappedRowsV2.correctCounts), uint32(remappedRowsV2.uncorrectCounts),
			uint32(pendingRows), uint32(failedRows), nil
	}

	var remappedRows C.cndevRemappedRow_t
	remappedRows.version = C.CNDEV_VERSION_5
	r := C.cndevGetRemappedRows(&remappedRows, c.cndevHandleMap[idx])
	return uint32(remappedRows.correctRows), uint32(remappedRows.uncorrectRows),
		uint32(remappedRows.pendingRows), uint32(remappedRows.failedRows), errorString(r)
}

func (c *cndev) GetDeviceRetiredPageInfo(idx uint) (int, uint32, error) {
	var retiredPage C.cndevRetiredPageInfo_t
	retiredPage.version = C.CNDEV_VERSION_5
	r := C.cndevGetRetiredPages(&retiredPage, c.cndevHandleMap[idx])
	return int(retiredPage.cause), uint32(retiredPage.pageCount), errorString(r)
}

func (c *cndev) GetDeviceSN(idx uint) (string, error) {
	var cardSN C.cndevCardSN_t
	cardSN.version = C.CNDEV_VERSION_5
	r := C.cndevGetCardSN(&cardSN, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return "", err
	}
	sn := fmt.Sprintf("%x", uint64(cardSN.sn))
	return sn, nil
}

func (c *cndev) GetDeviceTemperature(idx uint) (int, int, int, []int, []int, error) {
	var cardTemperature C.cndevTemperatureInfo_t
	cardTemperature.version = C.CNDEV_VERSION_5
	cluster, err := c.GetDeviceClusterCount(idx)
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}
	r := C.cndevGetTemperatureInfo(&cardTemperature, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return 0, 0, 0, nil, nil, err
	}
	clusterTemperature := make([]int, cluster)
	for i := 0; i < cluster; i++ {
		clusterTemperature[i] = int(cardTemperature.cluster[i])
	}
	memDie, err := c.GetDeviceMemoryDieCount(idx)
	if err != nil { // to comply with old version
		memDie = 0
	}
	memDieTemperature := make([]int, memDie)
	for i := 0; i < memDie; i++ {
		memDieTemperature[i] = int(cardTemperature.memoryDie[i])
	}

	return int(cardTemperature.board), int(cardTemperature.memory), int(cardTemperature.chip), clusterTemperature, memDieTemperature, nil
}

func (c *cndev) GetDeviceTinyCoreUtil(idx uint) ([]int, error) {
	var cardTinyCoreUtil C.cndevTinyCoreUtilization_t
	cardTinyCoreUtil.version = C.CNDEV_VERSION_5
	r := C.cndevGetTinyCoreUtilization(&cardTinyCoreUtil, c.cndevHandleMap[idx])
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
	cardUtil.version = C.CNDEV_VERSION_5
	core, err := c.GetDeviceCoreNum(idx)
	if err != nil {
		return 0, nil, err
	}
	r := C.cndevGetDeviceUtilizationInfo(&cardUtil, c.cndevHandleMap[idx])
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
	uuidInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetUUID(&uuidInfo, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return "", err
	}
	return C.GoString((*C.char)(unsafe.Pointer(&uuidInfo.uuid))), nil
}

func (c *cndev) GetDeviceVersion(idx uint) (uint, uint, uint, uint, uint, uint, error) {
	var versionInfo C.cndevVersionInfo_t
	versionInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetVersionInfo(&versionInfo, c.cndevHandleMap[idx])
	return uint(versionInfo.mcuMajorVersion), uint(versionInfo.mcuMinorVersion), uint(versionInfo.mcuBuildVersion), uint(versionInfo.driverMajorVersion), uint(versionInfo.driverMinorVersion), uint(versionInfo.driverBuildVersion), errorString(r)
}

// GetDeviceVfState returns cndevCardVfState_t.vfState, which is bitmap of mlu vf state. See cndev user doc for more details.
func (c *cndev) GetDeviceVfState(idx uint) (int, error) {
	var vfstate C.cndevCardVfState_t
	vfstate.version = C.CNDEV_VERSION_5
	r := C.cndevGetCardVfState(&vfstate, c.cndevHandleMap[idx])
	log.Debugf("slot:%d, vfstate is:%d", idx, int(vfstate.vfState))
	return int(vfstate.vfState), errorString(r)
}

func (c *cndev) GetDeviceVideoCodecUtil(idx uint) ([]int, error) {
	var cardVideoCodecUtil C.cndevVideoCodecUtilization_t
	cardVideoCodecUtil.version = C.CNDEV_VERSION_5
	r := C.cndevGetVideoCodecUtilization(&cardVideoCodecUtil, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return nil, err
	}
	videoCodecUtil := make([]int, uint(cardVideoCodecUtil.vpuCount))
	for i := uint(0); i < uint(cardVideoCodecUtil.vpuCount); i++ {
		videoCodecUtil[i] = int(cardVideoCodecUtil.vpuCodecUtilization[i])
	}
	return videoCodecUtil, nil
}

func (c *cndev) GetSupportedEventTypes(idx uint) error {
	types := C.ulonglong(C.cndevEventTypeAll)
	r := C.cndevGetSupportedEventTypes(&types, c.cndevHandleMap[idx])

	return errorString(r)
}

func (c *cndev) RegisterEventsHandleAndWait(slots []int, ch chan XIDInfoWithTimestamp) error {
	var handle C.cndevEventHandle
	r := C.cndevEventHandleCreate(&handle)
	if err := errorString(r); err != nil {
		return err
	}

	for _, idx := range slots {
		r := C.cndevRegisterEvents(handle, C.cndevEventTypeAll, c.cndevHandleMap[uint(idx)])
		if err := errorString(r); err != nil {
			log.Errorf("register event error: %v for slot %d", err, idx)
		}
	}

	go waitEvents(handle, ch)
	return nil
}

func waitEvents(handle C.cndevEventHandle, ch chan XIDInfoWithTimestamp) {
	var ret C.cndevRet_t
	var eventData C.cndevEventData_t
	for {
		ret = C.cndevEventWait(handle, &eventData, -1)
		if err := errorString(ret); err != nil {
			log.Errorf("wait event failed: %v", err)
			continue
		}
		if eventData.eventType != 1 {
			log.Debugf("not xid event: %v", eventData)
			continue
		}

		info := XIDInfoWithTimestamp{
			XIDInfo: XIDInfo{
				XID: int64(eventData.eventData),
				DeviceInfo: DeviceInfo{
					Device:            uint(eventData.device),
					ComputeInstanceID: int32(eventData.computeInstanceId),
					MLUInstanceID:     int32(eventData.mluInstanceId),
				},
			},

			Timestamp: fmt.Sprintf("%d", eventData.timestamp),
		}

		select {
		case ch <- info:
		default:
			// channel is full, delete the oldest data
			<-ch
			ch <- info
		}
	}
}

func errorString(ret C.cndevRet_t) error {
	if ret == C.CNDEV_SUCCESS {
		return nil
	}
	err := C.GoString(C.cndevGetErrorString(ret))
	return fmt.Errorf("cndev: %v", err)
}

func (c *cndev) generateDeviceHandleMap() error {
	num, err := c.GetDeviceCount()
	if err != nil {
		return err
	}
	for i := uint(0); i < num; i++ {
		var handle C.cndevDevice_t
		r := C.cndevGetDeviceHandleByIndex(C.int(i), &handle)
		if errorString(r) != nil {
			return errorString(r)
		}
		c.cndevHandleMap[i] = handle
	}
	return nil
}
