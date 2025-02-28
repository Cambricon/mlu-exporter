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

// #cgo LDFLAGS: -ldl -Wl,--export-dynamic -Wl,--unresolved-symbols=ignore-in-object-files
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

type MimInstanceInfo struct {
	InstanceName string
	InstanceID   int
	UUID         string
}

type MimProfileInfo struct {
	GDMACount    int
	JPUCount     int
	MemorySize   int
	MLUCoreCount int
	ProfileID    int
	ProfileName  string
	VPUCount     int
}

type MimInfo struct {
	InstanceInfo   MimInstanceInfo
	ProfileInfo    MimProfileInfo
	PlacementStart int
	PlacementSize  int
}

type SmluInstanceInfo struct {
	InstanceID   int
	InstanceName string
	UUID         string
}

type SmluProfileInfo struct {
	IpuTotal    uint32
	MemTotal    uint64
	ProfileID   int
	ProfileName string
}

type SmluInfo struct {
	InstanceInfo SmluInstanceInfo
	ProfileInfo  SmluProfileInfo
	IpuUtil      uint32
	MemUsed      uint64
}

type DeviceInfo struct {
	Device            uint
	ComputeInstanceID int32
	MLUInstanceID     int32
}

type XIDInfo struct {
	XIDBase10 int64
	XID       string
	DeviceInfo
}

type XIDInfoWithTimestamp struct {
	Timestamp int64
	XIDInfo
}

type Cndev interface {
	Init(healthCheck bool) error
	Release() error

	DeviceMimModeEnabled(idx uint) (bool, error)
	DeviceSmluModeEnabled(idx uint) (bool, error)
	GetAllMLUInstanceInfo(idx uint) ([]MimInfo, error)
	GetAllSMluInfo(idx uint) ([]SmluInfo, error)
	GetDeviceAddressSwaps(idx uint) (uint32, uint32, uint32, uint32, uint32, error)
	GetDeviceBAR4MemoryInfo(idx uint) (uint64, uint64, uint64, error)
	GetDeviceClusterCount(idx uint) (int, error)
	GetDeviceCndevVersion() (uint, uint, uint, error)
	GetDeviceComputeCapability(idx uint) (uint, uint, error)
	GetDeviceCoreNum(idx uint) (uint, error)
	GetDeviceCount() (uint, error)
	GetDeviceCPUUtil(idx uint) (uint16, []uint8, []uint8, []uint8, []uint8, []uint8, []uint8, error)
	GetDeviceCRCInfo(idx uint) (uint64, uint64, error)
	GetDeviceCurrentInfo(idx uint) (int, int, int, error)
	GetDeviceCurrentPCIeInfo(idx uint) (int, int, error)
	GetDeviceDDRInfo(idx uint) (int, float64, error)
	GetDeviceECCInfo(idx uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error)
	GetDeviceEccMode(idx uint) (int, int, error)
	GetDeviceFanSpeed(idx uint) (int, error)
	GetDeviceHealth(idx uint) (int, error)
	GetDeviceHeartbeatCount(idx uint) (uint32, error)
	GetDeviceImageCodecUtil(idx uint) ([]int, error)
	GetDeviceMaxPCIeInfo(idx uint) (int, int, error)
	GetDeviceMemEccCounter(idx uint) ([]uint32, []uint64, error)
	GetDeviceMemory(idx uint) (int64, int64, int64, int64, error)
	GetDeviceMemoryDieCount(idx uint) (int, error)
	GetDeviceMimProfileInfo(idx uint) ([]MimProfileInfo, error)
	GetDeviceMimProfileMaxInstanceCount(idx, profile uint) (int, error)
	GetDeviceMLULinkCapability(idx, link uint) (uint, uint, error)
	GetDeviceMLULinkCounter(idx, link uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error)
	GetDeviceMLULinkErrorCounter(idx, link uint) (uint64, error)
	GetDeviceMLULinkEventCounter(idx, link uint) (uint64, error)
	GetDeviceMLULinkPortMode(idx, link uint) (int, error)
	GetDeviceMLULinkPortNumber(idx uint) int
	GetDeviceMLULinkPPI(idx, link uint) (string, error)
	GetDeviceMLULinkRemoteInfo(idx, link uint) (uint64, uint64, uint32, uint32, int32, uint64, string, string, string, string, error)
	GetDeviceMLULinkSpeedInfo(idx, link uint) (float32, int, error)
	GetDeviceMLULinkStatus(idx, link uint) (int, int, int, error)
	GetDeviceMLULinkVersion(idx, link uint) (uint, uint, uint, error)
	GetDeviceMLUInstanceByID(idx uint, instanceID int) (uint, error)
	GetDeviceModel(idx uint) string
	GetDeviceNUMANodeID(idx uint) (int, error)
	GetDeviceOpticalInfo(idx, link uint) (uint8, float32, float32, []float32, []float32, error)
	GetDeviceOsMemory(idx uint) (int64, int64, error)
	GetDeviceOverTemperatureInfo(idx uint) (uint32, uint32, error)
	GetDeviceOverTemperatureShutdownThreshold(idx uint) (int, error)
	GetDeviceOverTemperatureSlowdownThreshold(idx uint) (int, error)
	GetDeviceParityError(idx uint) (int, error)
	GetDevicePCIeInfo(idx uint) (int, uint, uint, uint16, uint16, uint, uint, uint, uint, error)
	GetDevicePCIeReplayCount(idx uint) (uint32, error)
	GetDevicePCIeThroughput(idx uint) (int64, int64, error)
	GetDevicePerformanceThrottleReason(idx uint) (bool, error)
	GetDevicePower(idx uint) (int, error)
	GetDevicePowerManagementDefaultLimitation(idx uint) (uint16, error)
	GetDevicePowerManagementLimitation(idx uint) (uint16, error)
	GetDevicePowerManagementLimitRange(idx uint) (uint16, uint16, error)
	GetDeviceProcessUtil(idx uint) ([]uint32, []uint32, []uint32, []uint32, []uint32, []uint32, error)
	GetDeviceRemappedRows(idx uint) (uint32, uint32, uint32, uint32, uint32, error)
	GetDeviceRetiredPageInfo(idx uint) (uint32, uint32, error)
	GetDeviceRetiredPagesOperation(idx uint) (int, error)
	GetDeviceSMluProfileIDInfo(idx uint) ([]int, error)
	GetDeviceSMluProfileInfo(idx, profile uint) (SmluProfileInfo, uint32, uint32, error)
	GetDeviceSN(idx uint) (string, error)
	GetDeviceTemperature(idx uint) (int, int, int, []int, []int, error)
	GetDeviceTensorUtil(idx uint) (int, error)
	GetDeviceTinyCoreUtil(idx uint) ([]int, error)
	GetDeviceUtil(idx uint) (int, []int, error)
	GetDeviceFrequency(idx uint) (int, int, error)
	GetDeviceUUID(idx uint) (string, error)
	GetDeviceVersion(idx uint) (uint, uint, uint, uint, uint, uint, error)
	GetDeviceVfState(idx uint) (int, error)
	GetDeviceVideoCodecUtil(idx uint) ([]int, []int, error)
	GetDeviceVoltageInfo(idx uint) (int, int, int, error)
	RegisterEventsHandleAndWait(slots []int, ch chan XIDInfoWithTimestamp) error
	GetTopologyRelationship(domain1, bus1, device1, function1, domain2, bus2, device2, function2 uint) (int, error)
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
	if infs == nil {
		return miInfos, fmt.Errorf("malloc failed for GetAllMLUInstanceInfo")
	}
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
	// while new process may be produced between two cndevGetAllMluInstanceInfo calls,
	// so need to loop until this corner case vanishes
	for {
		err := errorString(r)
		if err == nil {
			break
		}
		if r == C.CNDEV_ERROR_INSUFFICIENT_SPACE {
			log.Debugf("cndevGetAllMluInstanceInfo with insufficient space with real counts %d, with slot: %d, will try with the real counts", miCount, idx)
			newInfs := (*C.cndevMluInstanceInfo_t)(C.realloc(unsafe.Pointer(infs), C.size_t(miCount)*C.size_t(unsafe.Sizeof(miInfo))))
			if newInfs == nil {
				return miInfos, fmt.Errorf("realloc failed for cndevGetAllMluInstanceInfo")
			}
			infs = newInfs
			r = C.cndevGetAllMluInstanceInfo(&miCount, infs, c.cndevHandleMap[idx])
			continue
		}
		return miInfos, err
	}

	infos := (*[1 << 16]C.cndevMluInstanceInfo_t)(unsafe.Pointer(infs))[:miCount]
	for i := 0; i < int(miCount); i++ {
		info := infos[i]
		miInfos = append(miInfos,
			MimInfo{
				InstanceInfo: MimInstanceInfo{
					InstanceName: C.GoString((*C.char)(unsafe.Pointer(&info.instanceName))),
					InstanceID:   int(info.instanceId),
					UUID:         C.GoString((*C.char)(unsafe.Pointer(&info.uuid))),
				},
				ProfileInfo: MimProfileInfo{
					GDMACount:    int(info.gdmaCount),
					JPUCount:     int(info.jpuCount),
					MemorySize:   int(info.memorySize),
					MLUCoreCount: int(info.mluCoreCount),
					ProfileID:    int(info.profileId),
					ProfileName:  C.GoString((*C.char)(unsafe.Pointer(&info.profileName))),
					VPUCount:     int(info.vpuCount),
				},
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
	if infs == nil {
		return smluInfos, fmt.Errorf("malloc failed for GetAllSMluInfo")
	}
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
	// while new process may be produced between two cndevGetAllSMluInstanceInfo calls,
	// so need to loop until this corner case vanishes
	for {
		err := errorString(r)
		if err == nil {
			break
		}
		if r == C.CNDEV_ERROR_INSUFFICIENT_SPACE {
			log.Debugf("cndevGetAllSMluInstanceInfo with insufficient space with real counts %d, with slot: %d, will try with the real counts", smluCount, idx)
			newInfs := (*C.cndevSMluInfo_t)(C.realloc(unsafe.Pointer(infs), C.size_t(smluCount)*C.size_t(unsafe.Sizeof(smluInfo))))
			if newInfs == nil {
				return smluInfos, fmt.Errorf("realloc failed for cndevGetAllSMluInstanceInfo")
			}
			infs = newInfs
			r = C.cndevGetAllSMluInstanceInfo(&smluCount, infs, c.cndevHandleMap[idx])
			continue
		}
		return smluInfos, err
	}

	infos := (*[1 << 16]C.cndevSMluInfo_t)(unsafe.Pointer(infs))[:smluCount]
	for i := 0; i < int(smluCount); i++ {
		info := infos[i]
		smluInfos = append(smluInfos,
			SmluInfo{
				InstanceInfo: SmluInstanceInfo{
					InstanceID:   int(info.instanceId),
					InstanceName: C.GoString((*C.char)(unsafe.Pointer(&info.instanceName))),
					UUID:         C.GoString((*C.char)(unsafe.Pointer(&info.uuid))),
				},
				ProfileInfo: SmluProfileInfo{
					IpuTotal:    uint32(info.mluQuota[C.CNDEV_SMLU_MAX]),
					MemTotal:    uint64(info.memorySize[C.CNDEV_SMLU_MAX]),
					ProfileID:   int(info.profileId),
					ProfileName: C.GoString((*C.char)(unsafe.Pointer(&info.profileName))),
				},
				IpuUtil: uint32(info.mluQuota[C.CNDEV_SMLU_USAGE]),
				MemUsed: uint64(info.memorySize[C.CNDEV_SMLU_USAGE]),
			})
	}
	log.Debugf("smlu infos for device %d are %+v", idx, smluInfos)

	return smluInfos, nil
}

func (c *cndev) GetDeviceAddressSwaps(idx uint) (uint32, uint32, uint32, uint32, uint32, error) {
	var addressSwap C.cndevAddressSwap_t
	addressSwap.version = C.CNDEV_VERSION_6
	r := C.cndevGetAddressSwaps(&addressSwap, c.cndevHandleMap[idx])
	return uint32(addressSwap.correctCounts), uint32(addressSwap.uncorrectCounts), uint32(addressSwap.histogram[C.CNDEV_AVAILABILITY_XLABLE_NONE]), uint32(addressSwap.histogram[C.CNDEV_AVAILABILITY_XLABLE_PARTIAL]), uint32(addressSwap.histogram[C.CNDEV_AVAILABILITY_XLABLE_MAX]), errorString(r)
}

func (c *cndev) GetDeviceOsMemory(idx uint) (int64, int64, error) {
	var deviceOsMemInfo C.cndevDeviceOsMemoryInfo_t
	r := C.cndevGetDeviceOsMemoryUsageV2(&deviceOsMemInfo, c.cndevHandleMap[idx])
	return int64(deviceOsMemInfo.deviceSystemMemoryUsed), int64(deviceOsMemInfo.deviceSystemMemoryTotal), errorString(r)
}

func (c *cndev) GetDeviceBAR4MemoryInfo(idx uint) (uint64, uint64, uint64, error) {
	var memInfo C.cndevBAR4Memory_t
	r := C.cndevGetBAR4MemoryInfo(&memInfo, c.cndevHandleMap[idx])
	return uint64(memInfo.total_size), uint64(memInfo.used_size), uint64(memInfo.free_size), errorString(r)
}

func (c *cndev) GetDeviceClusterCount(idx uint) (int, error) {
	var cardClusterNum C.cndevCardClusterCount_t
	cardClusterNum.version = C.CNDEV_VERSION_5
	r := C.cndevGetClusterCount(&cardClusterNum, c.cndevHandleMap[idx])
	return int(cardClusterNum.count), errorString(r)
}

func (c *cndev) GetDeviceCndevVersion() (uint, uint, uint, error) {
	var cndevVersion C.cndevLibVersionInfo_t
	cndevVersion.version = C.CNDEV_VERSION_5
	r := C.cndevGetLibVersion(&cndevVersion)
	return uint(cndevVersion.libMajorVersion), uint(cndevVersion.libMinorVersion), uint(cndevVersion.libBuildVersion), errorString(r)
}

func (c *cndev) GetDeviceComputeCapability(idx uint) (uint, uint, error) {
	var major, minor C.uint
	r := C.cndevGetComputeCapability(&major, &minor, c.cndevHandleMap[idx])
	return uint(major), uint(minor), errorString(r)
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

func (c *cndev) GetDeviceCPUUtil(idx uint) (uint16, []uint8, []uint8, []uint8, []uint8, []uint8, []uint8, error) {
	var cardCPUUtil C.cndevDeviceCPUUtilizationV2_t
	cardCPUUtil.version = C.CNDEV_VERSION_6
	r := C.cndevGetDeviceCPUUtilizationV2(&cardCPUUtil, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return 0, nil, nil, nil, nil, nil, nil, err
	}
	coreCPUUtil := make([]uint8, uint8(cardCPUUtil.coreNumber))
	coreCPUUser := make([]uint8, uint8(cardCPUUtil.coreNumber))
	coreCPUNice := make([]uint8, uint8(cardCPUUtil.coreNumber))
	coreCPUSys := make([]uint8, uint8(cardCPUUtil.coreNumber))
	coreCPUSi := make([]uint8, uint8(cardCPUUtil.coreNumber))
	coreCPUHi := make([]uint8, uint8(cardCPUUtil.coreNumber))

	for i := uint8(0); i < uint8(cardCPUUtil.coreNumber); i++ {
		coreCPUUtil[i] = uint8(cardCPUUtil.coreUtilization[i])
		coreCPUUser[i] = uint8(cardCPUUtil.coreStateUtilization[i][0])
		coreCPUNice[i] = uint8(cardCPUUtil.coreStateUtilization[i][1])
		coreCPUSys[i] = uint8(cardCPUUtil.coreStateUtilization[i][2])
		coreCPUSi[i] = uint8(cardCPUUtil.coreStateUtilization[i][3])
		coreCPUHi[i] = uint8(cardCPUUtil.coreStateUtilization[i][4])
	}
	return uint16(cardCPUUtil.chipUtilization), coreCPUUtil, coreCPUUser, coreCPUNice, coreCPUSys, coreCPUSi, coreCPUHi, nil
}

func (c *cndev) GetDeviceCRCInfo(idx uint) (uint64, uint64, error) {
	var cardCRCInfo C.cndevCRCInfo_t
	cardCRCInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetCRCInfo(&cardCRCInfo, c.cndevHandleMap[idx])
	return uint64(cardCRCInfo.die2dieCRCError), uint64(cardCRCInfo.die2dieCRCErrorOverflow), errorString(r)
}

func (c *cndev) GetDeviceCurrentInfo(idx uint) (int, int, int, error) {
	var currentInfo C.cndevCurrentInfo_t
	r := C.cndevGetCurrentInfo(&currentInfo, c.cndevHandleMap[idx])
	return int(currentInfo.ipuCoreCurrent), int(currentInfo.socCurrent), int(currentInfo.hbmCurrent), errorString(r)
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

func (c *cndev) GetDeviceEccMode(idx uint) (int, int, error) {
	var currentMode, pendingMode C.cndevEccMode_t
	currentMode.version = C.CNDEV_VERSION_6
	pendingMode.version = C.CNDEV_VERSION_6
	r := C.cndevGetEccMode(&currentMode, &pendingMode, c.cndevHandleMap[idx])
	return int(currentMode.mode), int(pendingMode.mode), errorString(r)
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

func (c *cndev) GetDeviceMaxPCIeInfo(idx uint) (int, int, error) {
	var pcieInfo C.cndevPCIInfo_t
	pcieInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetMaxPCIInfo(&pcieInfo, c.cndevHandleMap[idx])
	return int(pcieInfo.maxSpeed), int(pcieInfo.maxWidth), errorString(r)
}

func (c *cndev) GetDeviceMemEccCounter(idx uint) ([]uint32, []uint64, error) {
	var memEccCounter C.cndevMemEccCounter_t
	memEccCounter.version = C.CNDEV_VERSION_6
	index := C.int(idx)
	if id, ok := c.cndevHandleMap[idx]; ok {
		index = id
	}
	r := C.cndevGetMemEccCounter(&memEccCounter, index)
	if err := errorString(r); err != nil {
		return nil, nil, err
	}
	volatile := make([]uint32, 5)
	aggregate := make([]uint64, 6)
	volatile[0] = uint32(memEccCounter.volatileSramEccSbeCounter)
	volatile[1] = uint32(memEccCounter.volatileSramEccDbeCounter)
	volatile[2] = uint32(memEccCounter.volatileSramEccParityCounter)
	volatile[3] = uint32(memEccCounter.volatileDramEccSbeCounter)
	volatile[4] = uint32(memEccCounter.volatileDramEccDbeCounter)
	aggregate[0] = uint64(memEccCounter.aggregateDramEccSbeCounter)
	aggregate[1] = uint64(memEccCounter.aggregateDramEccDbeCounter)
	aggregate[2] = uint64(memEccCounter.aggregateSramEccSbeCounter)
	aggregate[3] = uint64(memEccCounter.aggregateSramEccDbeCounter)
	aggregate[4] = uint64(memEccCounter.aggregateSramEccParityCounter)
	if bool(memEccCounter.aggregateSramEccThresholdExceeded) {
		aggregate[5] = 1
	} else {
		aggregate[5] = 0
	}
	return volatile, aggregate, nil
}

func (c *cndev) GetDeviceMemory(idx uint) (int64, int64, int64, int64, error) {
	var cardMemInfo C.cndevMemoryInfoV2_t
	index := C.int(idx)
	if id, ok := c.cndevHandleMap[idx]; ok {
		index = id
	}
	r := C.cndevGetMemoryUsageV2(&cardMemInfo, index)
	return int64(cardMemInfo.physicalMemoryUsed), int64(cardMemInfo.physicalMemoryTotal), int64(cardMemInfo.virtualMemoryUsed), int64(cardMemInfo.virtualMemoryTotal), errorString(r)
}

func (c *cndev) GetDeviceMemoryDieCount(idx uint) (int, error) {
	//this api is deprecated, keep to comply with old version
	return 0, nil
}

func (c *cndev) GetDeviceMimProfileInfo(idx uint) ([]MimProfileInfo, error) {
	var mimProfileInfos []MimProfileInfo

	for profileID := 0; profileID < C.CNDEV_MLUINSTANCE_PROFILE_COUNT; profileID++ {
		var profileInfo C.cndevMluInstanceProfileInfo_t
		profileInfo.version = C.CNDEV_VERSION_6
		r := C.cndevGetMluInstanceProfileInfo(&profileInfo, C.int(profileID), c.cndevHandleMap[idx])
		if err := errorString(r); r == C.CNDEV_ERROR_NOT_SUPPORTED || r == C.CNDEV_ERROR_NOT_FOUND {
			continue
		} else if err != nil {
			return mimProfileInfos, err
		}
		mimProfileInfos = append(mimProfileInfos,
			MimProfileInfo{
				GDMACount:    int(profileInfo.gdmaCount),
				JPUCount:     int(profileInfo.jpuCount),
				MemorySize:   int(profileInfo.memorySize),
				MLUCoreCount: int(profileInfo.mluCoreCount),
				ProfileID:    int(profileInfo.profileId),
				ProfileName:  C.GoString((*C.char)(unsafe.Pointer(&profileInfo.name))),
				VPUCount:     int(profileInfo.vpuCount),
			})
	}
	return mimProfileInfos, nil
}

func (c *cndev) GetDeviceMimProfileMaxInstanceCount(idx, profile uint) (int, error) {
	var count C.int
	r := C.cndevGetMaxMluInstanceCount(&count, C.int(profile), c.cndevHandleMap[idx])
	return int(count), errorString(r)
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

func (c *cndev) GetDeviceMLULinkErrorCounter(idx, link uint) (uint64, error) {
	var cardMLULinkErrorCounter C.cndevMLULinkErrorCounter_t
	r := C.cndevGetMLULinkErrorCounter(&cardMLULinkErrorCounter, c.cndevHandleMap[idx], C.int(link))
	return uint64(cardMLULinkErrorCounter.illegalAccessCnt), errorString(r)
}

func (c *cndev) GetDeviceMLULinkEventCounter(idx, link uint) (uint64, error) {
	var cardMLULinkEventCounter C.cndevMLULinkEventCounter_t
	r := C.cndevGetMLULinkEventCounter(&cardMLULinkEventCounter, c.cndevHandleMap[idx], C.int(link))
	return uint64(cardMLULinkEventCounter.linkDown), errorString(r)
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

func (c *cndev) GetDeviceMLULinkPPI(idx, link uint) (string, error) {
	var cardMLULinkPPI C.cndevMLULinkPPI_t
	r := C.cndevGetMLULinkPPI(&cardMLULinkPPI, c.cndevHandleMap[idx], C.int(link))
	return C.GoString((*C.char)(unsafe.Pointer(&cardMLULinkPPI.ppi))), errorString(r)
}

func (c *cndev) GetDeviceMLULinkSpeedInfo(idx, link uint) (float32, int, error) {
	var cardMLULinkSpeedInfo C.cndevMLULinkSpeed_t
	cardMLULinkSpeedInfo.version = C.CNDEV_VERSION_5
	r := C.cndevGetMLULinkSpeedInfo(&cardMLULinkSpeedInfo, c.cndevHandleMap[idx], C.int(link))
	return float32(cardMLULinkSpeedInfo.speedValue), int(cardMLULinkSpeedInfo.speedFormat), errorString(r)
}

func (c *cndev) GetDeviceMLULinkRemoteInfo(idx, link uint) (uint64, uint64, uint32, uint32, int32, uint64, string, string, string, string, error) {
	var cardMLULinkRemoteInfo C.cndevMLULinkRemoteInfo_t
	cardMLULinkRemoteInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetMLULinkRemoteInfo(&cardMLULinkRemoteInfo, c.cndevHandleMap[idx], C.int(link))
	if err := errorString(r); err != nil {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", err
	}
	var ip, mac, uuid, portName string
	ip = C.GoString((*C.char)(unsafe.Pointer(&cardMLULinkRemoteInfo.devIp)))
	mac = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", cardMLULinkRemoteInfo.mac_addr[0], cardMLULinkRemoteInfo.mac_addr[1], cardMLULinkRemoteInfo.mac_addr[2], cardMLULinkRemoteInfo.mac_addr[3], cardMLULinkRemoteInfo.mac_addr[4], cardMLULinkRemoteInfo.mac_addr[5])
	uuid = C.GoString((*C.char)(unsafe.Pointer(&cardMLULinkRemoteInfo.uuid)))
	portName = C.GoString((*C.char)(unsafe.Pointer(&cardMLULinkRemoteInfo.port_name)))
	return uint64(cardMLULinkRemoteInfo.mcSn), uint64(cardMLULinkRemoteInfo.baSn), uint32(cardMLULinkRemoteInfo.slotId), uint32(cardMLULinkRemoteInfo.portId), int32(cardMLULinkRemoteInfo.connectType), uint64(cardMLULinkRemoteInfo.ncsUUID64), ip, mac, uuid, portName, errorString(r)
}

func (c *cndev) GetDeviceMLULinkStatus(idx, link uint) (int, int, int, error) {
	var cardMLULinkStatus C.cndevMLULinkStatusV2_t
	r := C.cndevGetMLULinkStatusV2(&cardMLULinkStatus, c.cndevHandleMap[idx], C.int(link))
	return int(cardMLULinkStatus.macState), int(cardMLULinkStatus.serdesState), int(cardMLULinkStatus.presenceState), errorString(r)
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

func (c *cndev) GetDeviceOverTemperatureInfo(idx uint) (uint32, uint32, error) {
	var overTemperatureInfo C.cndevOverTemperatureInfo_t
	overTemperatureInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetOverTemperatureInfo(&overTemperatureInfo, c.cndevHandleMap[idx])
	return uint32(overTemperatureInfo.powerOffCounter), uint32(overTemperatureInfo.underClockCounter), errorString(r)
}

func (c *cndev) GetDeviceOpticalInfo(idx, link uint) (uint8, float32, float32, []float32, []float32, error) {
	var opticalInfo C.cndevOpticalInfo_t
	r := C.cndevGetOpticalInfo(&opticalInfo, c.cndevHandleMap[idx], C.int(link))
	if err := errorString(r); err != nil {
		return 0, 0, 0, nil, nil, err
	}
	txpwr := make([]float32, uint(opticalInfo.valid_lane_num))
	for i := uint(0); i < uint(opticalInfo.valid_lane_num); i++ {
		txpwr[i] = float32(opticalInfo.txpwr[i])
	}
	rxpwr := make([]float32, uint(opticalInfo.valid_lane_num))
	for i := uint(0); i < uint(opticalInfo.valid_lane_num); i++ {
		rxpwr[i] = float32(opticalInfo.rxpwr[i])
	}
	return uint8(opticalInfo.present), float32(opticalInfo.temp), float32(opticalInfo.volt), txpwr, rxpwr, errorString(r)
}

func (c *cndev) GetDeviceOverTemperatureShutdownThreshold(idx uint) (int, error) {
	type fieldVaule struct {
		fieldID     int32
		ret         int32
		value       int
		valueType   int
		latencyUsec int64
		timestamp   int64
		scopeID     uint
	}
	var fieldVauleInfo fieldVaule
	value := (*C.cndevFieldVaule_t)(C.malloc(C.size_t(unsafe.Sizeof(fieldVauleInfo))))
	if value == nil {
		return 0, fmt.Errorf("malloc failed for cndevDeviceGetFieldValues")
	}
	defer C.free(unsafe.Pointer(value))

	value.fieldId = C.cndevFieldTemperatureShutDown
	r := C.cndevDeviceGetFieldValues(c.cndevHandleMap[idx], C.int(1), value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	if err := errorString(value.ret); err != nil {
		return 0, err
	}
	shutdown := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(value)) + unsafe.Sizeof(fieldVauleInfo.fieldID) + unsafe.Sizeof(fieldVauleInfo.ret)))
	return int(shutdown), nil
}

func (c *cndev) GetDeviceOverTemperatureSlowdownThreshold(idx uint) (int, error) {
	type fieldVaule struct {
		fieldID     int32
		ret         int32
		value       int
		valueType   int
		latencyUsec int64
		timestamp   int64
		scopeID     uint
	}
	var fieldVauleInfo fieldVaule
	value := (*C.cndevFieldVaule_t)(C.malloc(C.size_t(unsafe.Sizeof(fieldVauleInfo))))
	if value == nil {
		return 0, fmt.Errorf("malloc failed for cndevDeviceGetFieldValues")
	}
	defer C.free(unsafe.Pointer(value))

	value.fieldId = C.cndevFieldTemperatureSlowDown
	r := C.cndevDeviceGetFieldValues(c.cndevHandleMap[idx], C.int(1), value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	if err := errorString(value.ret); err != nil {
		return 0, err
	}
	slowDown := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(value)) + unsafe.Sizeof(fieldVauleInfo.fieldID) + unsafe.Sizeof(fieldVauleInfo.ret)))
	return int(slowDown), nil
}

func (c *cndev) GetDeviceParityError(idx uint) (int, error) {
	var parityError C.cndevParityError_t
	parityError.version = C.CNDEV_VERSION_5
	r := C.cndevGetParityError(&parityError, c.cndevHandleMap[idx])
	return int(parityError.counter), errorString(r)
}

func (c *cndev) GetDevicePCIeInfo(idx uint) (int, uint, uint, uint16, uint16, uint, uint, uint, uint, error) {
	var pcieInfo C.cndevPCIeInfoV2_t
	r := C.cndevGetPCIeInfoV2(&pcieInfo, c.cndevHandleMap[idx])
	return int(pcieInfo.slotId), uint(pcieInfo.subsystemId), uint(pcieInfo.deviceId), uint16(pcieInfo.vendor),
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
	var pcieThroughput C.cndevPCIethroughputV2_t
	r := C.cndevGetPCIethroughputV2(&pcieThroughput, c.cndevHandleMap[idx])
	return int64(pcieThroughput.pcieRead), int64(pcieThroughput.pcieWrite), errorString(r)
}

func (c *cndev) GetDevicePower(idx uint) (int, error) {
	var devicePower C.cndevDevicePowerInfo_t
	index := C.int(idx)
	if id, ok := c.cndevHandleMap[idx]; ok {
		index = id
	}
	r := C.cndevGetDevicePowerInfo(&devicePower, index)
	return int(devicePower.usage), errorString(r)
}

func (c *cndev) GetDevicePowerManagementDefaultLimitation(idx uint) (uint16, error) {
	var powerManagementLimitation C.cndevPowerManagementLimitation_t
	powerManagementLimitation.version = C.CNDEV_VERSION_6
	r := C.cndevGetPowerManagementDefaultLimitation(&powerManagementLimitation, c.cndevHandleMap[idx])
	return uint16(powerManagementLimitation.powerLimit), errorString(r)
}

func (c *cndev) GetDevicePowerManagementLimitation(idx uint) (uint16, error) {
	var powerManagementLimitation C.cndevPowerManagementLimitation_t
	powerManagementLimitation.version = C.CNDEV_VERSION_6
	r := C.cndevGetPowerManagementLimitation(&powerManagementLimitation, c.cndevHandleMap[idx])
	return uint16(powerManagementLimitation.powerLimit), errorString(r)
}

func (c *cndev) GetDevicePowerManagementLimitRange(idx uint) (uint16, uint16, error) {
	var powerManagementLimitRange C.cndevPowerManagementLimitationRange_t
	powerManagementLimitRange.version = C.CNDEV_VERSION_6
	r := C.cndevGetPowerManagementLimitationRange(&powerManagementLimitRange, c.cndevHandleMap[idx])
	return uint16(powerManagementLimitRange.minPowerLimit), uint16(powerManagementLimitRange.maxPowerLimit), errorString(r)
}

func (c *cndev) GetDevicePerformanceThrottleReason(idx uint) (bool, error) {
	var performanceThrottleReason C.cndevPerformanceThrottleReason_t
	performanceThrottleReason.version = C.CNDEV_VERSION_6
	r := C.cndevGetPerformanceThrottleReason(&performanceThrottleReason, c.cndevHandleMap[idx])
	return performanceThrottleReason.thermalSlowdown == C.CNDEV_FEATURE_ENABLED, errorString(r)
}

func (c *cndev) GetDeviceProcessUtil(idx uint) ([]uint32, []uint32, []uint32, []uint32, []uint32, []uint32, error) {
	processCount := 1 << 4
	var util C.cndevProcessUtilization_t
	utils := (*C.cndevProcessUtilization_t)(C.malloc(C.size_t(processCount) * C.size_t(unsafe.Sizeof(util))))
	if utils == nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("malloc failed for cndevGetProcessUtilization")
	}
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
	// after cndevGetProcessUtilization call the processCount will be set to real count,
	// while new process may be produced between two cndevGetProcessUtilization calls,
	// so need to loop until this corner case vanishes
	for {
		err := errorString(r)
		if err == nil {
			break
		}
		if r == C.CNDEV_ERROR_INSUFFICIENT_SPACE {
			log.Debugf("cndevGetProcessUtilization with insufficient space with real counts %d, with slot: %d, will try with the real counts", processCount, idx)
			newUtils := (*C.cndevProcessUtilization_t)(C.realloc(unsafe.Pointer(utils), C.size_t(processCount)*C.size_t(unsafe.Sizeof(util))))
			if newUtils == nil {
				return nil, nil, nil, nil, nil, nil, fmt.Errorf("realloc failed for cndevGetProcessUtilization")
			}
			utils = newUtils
			r = C.cndevGetProcessUtilization((*C.uint)(unsafe.Pointer(&processCount)), utils, c.cndevHandleMap[idx])
			continue
		}
		return nil, nil, nil, nil, nil, nil, err
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

func (c *cndev) GetDeviceRemappedRows(idx uint) (uint32, uint32, uint32, uint32, uint32, error) {
	major, minor, build, err := c.GetDeviceDriverVersion(idx)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	v1, err := semver.NewVersion(fmt.Sprintf("%d.%d.%d", major, minor, build))
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	// cndevGetRemappedRowsV2 is only available in driver version later than 6.0.2
	v2 := semver.MustParse("6.0.2")

	if v1.GreaterThan(v2) {
		var remappedRowsV2 C.cndevRemappedRowV2_t
		remappedRowsV2.version = C.CNDEV_VERSION_6
		r := C.cndevGetRemappedRowsV2(&remappedRowsV2, c.cndevHandleMap[idx])
		if err := errorString(r); err != nil {
			return 0, 0, 0, 0, 0, err
		}
		var repairStatus C.cndevRepairStatus_t
		repairStatus.version = C.CNDEV_VERSION_6
		r = C.cndevGetRepairStatus(&repairStatus, c.cndevHandleMap[idx])
		if err := errorString(r); err != nil {
			return 0, 0, 0, 0, 0, err
		}
		var pendingRows, failedRows, retirePending uint32
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
		if bool(repairStatus.isRetirePending) {
			retirePending = 1
		} else {
			retirePending = 0
		}
		return uint32(remappedRowsV2.correctCounts), uint32(remappedRowsV2.uncorrectCounts),
			uint32(pendingRows), uint32(failedRows), uint32(retirePending), nil
	}

	var remappedRows C.cndevRemappedRow_t
	remappedRows.version = C.CNDEV_VERSION_6
	r := C.cndevGetRemappedRows(&remappedRows, c.cndevHandleMap[idx])
	return uint32(remappedRows.correctRows), uint32(remappedRows.uncorrectRows),
		uint32(remappedRows.pendingRows), uint32(remappedRows.failedRows), 0, errorString(r)
}

func (c *cndev) GetDeviceRetiredPageInfo(idx uint) (uint32, uint32, error) {
	var retiredPage C.cndevRetiredPageInfo_t
	retiredPage.version = C.CNDEV_VERSION_5
	retiredPage.cause = 0
	r := C.cndevGetRetiredPages(&retiredPage, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return 0, 0, err
	}
	single := uint32(retiredPage.pageCount)
	retiredPage.cause = 1
	r = C.cndevGetRetiredPages(&retiredPage, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return single, 0, err
	}
	double := uint32(retiredPage.pageCount)
	return single, double, errorString(r)
}

func (c *cndev) GetDeviceRetiredPagesOperation(idx uint) (int, error) {
	var retirement C.cndevRetiredPageOperation_t
	retirement.version = C.CNDEV_VERSION_6
	r := C.cndevGetRetiredPagesOperation(&retirement, c.cndevHandleMap[idx])
	return int(retirement.retirePageOption), errorString(r)
}

func (c *cndev) GetDeviceSMluProfileIDInfo(idx uint) ([]int, error) {
	var profileIDInfo C.cndevSMluProfileIdInfo_t
	profileIDInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetSMluProfileIdInfo(&profileIDInfo, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return nil, err
	}
	info := make([]int, int(profileIDInfo.count))
	for i := 0; i < int(profileIDInfo.count); i++ {
		info[i] = int(profileIDInfo.profileId[i])
	}
	return info, nil
}

func (c *cndev) GetDeviceSMluProfileInfo(idx, profile uint) (SmluProfileInfo, uint32, uint32, error) {
	var profileInfo C.cndevSMluProfileInfo_t
	var smluProfileInfo SmluProfileInfo
	profileInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetSMluProfileInfo(&profileInfo, C.int(profile), c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return smluProfileInfo, 0, 0, err
	}
	total := uint32(profileInfo.totalCapacity)
	remain := uint32(profileInfo.remainCapacity)
	smluProfileInfo = SmluProfileInfo{
		IpuTotal:    uint32(profileInfo.mluQuota[C.CNDEV_SMLU_MAX]),
		MemTotal:    uint64(profileInfo.memorySize[C.CNDEV_SMLU_MAX]),
		ProfileID:   int(profileInfo.profileId),
		ProfileName: C.GoString((*C.char)(unsafe.Pointer(&profileInfo.name))),
	}
	return smluProfileInfo, total, remain, nil
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

func (c *cndev) GetDeviceTensorUtil(idx uint) (int, error) {
	type fieldVaule struct {
		fieldID     int32
		ret         int32
		value       int
		valueType   int
		latencyUsec int64
		timestamp   int64
		scopeID     uint
	}
	var fieldVauleInfo fieldVaule
	value := (*C.cndevFieldVaule_t)(C.malloc(C.size_t(unsafe.Sizeof(fieldVauleInfo))))
	if value == nil {
		return 0, fmt.Errorf("malloc failed for cndevDeviceGetFieldValues")
	}
	defer C.free(unsafe.Pointer(value))

	value.fieldId = C.cndevFieldTensorAverageUtilization

	r := C.cndevDeviceGetFieldValues(c.cndevHandleMap[idx], C.int(1), value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	if err := errorString(value.ret); err != nil {
		return 0, err
	}
	tensorUtil := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(value)) + unsafe.Sizeof(fieldVauleInfo.fieldID) + unsafe.Sizeof(fieldVauleInfo.ret)))
	return int(tensorUtil), nil
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

func (c *cndev) GetDeviceVideoCodecUtil(idx uint) ([]int, []int, error) {
	var cardVideoCodecUtil C.cndevVideoCodecUtilization_t
	cardVideoCodecUtil.version = C.CNDEV_VERSION_5
	r := C.cndevGetVideoCodecUtilization(&cardVideoCodecUtil, c.cndevHandleMap[idx])
	if err := errorString(r); err != nil {
		return nil, nil, err
	}
	videoCodecUtil := make([]int, uint(cardVideoCodecUtil.vpuCount))
	for i := uint(0); i < uint(cardVideoCodecUtil.vpuCount); i++ {
		videoCodecUtil[i] = int(cardVideoCodecUtil.vpuCodecUtilization[i])
	}
	videoDecoderUtil := make([]int, uint(cardVideoCodecUtil.decoderCount))
	for i := uint(0); i < uint(cardVideoCodecUtil.decoderCount); i++ {
		videoDecoderUtil[i] = int(cardVideoCodecUtil.vpuCodecUtilization[i])
	}

	return videoCodecUtil, videoDecoderUtil, nil
}

func (c *cndev) GetDeviceVoltageInfo(idx uint) (int, int, int, error) {
	var voltageInfo C.cndevVoltageInfo_t
	r := C.cndevGetVoltageInfo(&voltageInfo, c.cndevHandleMap[idx])
	return int(voltageInfo.ipuCoreVoltage), int(voltageInfo.socVoltage), int(voltageInfo.hbmVddVoltage), errorString(r)
}

func (c *cndev) GetTopologyRelationship(domain, bus, device, function, rdmaDomain, rdmaBus, rdmaDevice, rdmaFunction uint) (int, error) {
	var treeNode, rdmaTreeNode *C.cndevTopologyNode_t
	version := C.CNDEV_VERSION_6
	r := C.cndevGetNodeByBDF(C.int(version), &treeNode, C.uint(domain), C.uint(bus), C.uint(device), C.uint(function))
	if err := errorString(r); err != nil {
		return 0, err
	}
	r = C.cndevGetNodeByBDF(C.int(version), &rdmaTreeNode, C.uint(rdmaDomain), C.uint(rdmaBus), C.uint(rdmaDevice), C.uint(rdmaFunction))
	if err := errorString(r); err != nil {
		return 0, err
	}

	var topoRelationship C.cndevTopologyRelationship_t
	topoRelationship.version = C.CNDEV_VERSION_6
	r = C.cndevTopologyGetRelationshipByNode(&topoRelationship, treeNode, rdmaTreeNode)
	if err := errorString(r); err != nil {
		return 0, err
	}
	return int(topoRelationship.relation), nil
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
				XID:       fmt.Sprintf("%x", eventData.eventData),
				XIDBase10: int64(eventData.eventData),
				DeviceInfo: DeviceInfo{
					Device:            uint(eventData.device),
					ComputeInstanceID: int32(eventData.computeInstanceId),
					MLUInstanceID:     int32(eventData.mluInstanceId),
				},
			},
			Timestamp: int64(eventData.timestamp),
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
