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
	"fmt"
	"strconv"
	"sync"
	"unsafe"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
)

type MimInstanceInfo struct {
	InstanceName string
	UUID         string
	InstanceID   int
}

type MimProfileInfo struct {
	ProfileName  string
	GDMACount    int
	JPUCount     int
	MemorySize   int
	MLUCoreCount int
	ProfileID    int
	VPUCount     int
}

type MimInfo struct {
	ProfileInfo    MimProfileInfo
	InstanceInfo   MimInstanceInfo
	PlacementStart int
	PlacementSize  int
}

type SmluInstanceInfo struct {
	InstanceName string
	UUID         string
	InstanceID   int
}

type SmluProfileInfo struct {
	ProfileName string
	MemTotal    uint64
	ProfileID   int
	IpuTotal    uint32
}

type SmluInfo struct {
	InstanceInfo SmluInstanceInfo
	ProfileInfo  SmluProfileInfo
	MemUsed      uint64
	IpuUtil      uint32
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

type ChassisDevInfo struct {
	SN    string
	Model string
	Fw    string
	Mfc   string
}

type Cndev interface {
	Init(healthCheck bool) error
	Release() error

	DeviceMimModeEnabled(idx uint) (bool, error)
	DeviceSmluModeEnabled(idx uint) (bool, error)
	GenerateDeviceHandleMap(count uint) error
	GetAllMLUInstanceInfo(idx uint) ([]MimInfo, error)
	GetAllSMluInfo(idx uint) ([]SmluInfo, error)
	GetDeviceActivity(idx uint) (int, error)
	GetDeviceAddressSwaps(idx uint) (uint32, uint32, uint32, uint32, uint32, error)
	GetDeviceBAR4MemoryInfo(idx uint) (uint64, uint64, uint64, error)
	GetDeviceChassisInfo(idx uint) (uint64, string, string, string, string, string, []ChassisDevInfo, []ChassisDevInfo, []ChassisDevInfo, error)
	GetDeviceClusterCount(idx uint) (int, error)
	GetDeviceCndevVersion() (uint, uint, uint, error)
	GetDeviceComputeCapability(idx uint) (uint, uint, error)
	GetDeviceComputeMode(idx uint) (uint, error)
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
	GetDeviceHealth(idx uint) (int, bool, bool, []string, error)
	GetDeviceHeartbeatCount(idx uint) (uint32, error)
	GetDeviceImageCodecUtil(idx uint) ([]int, error)
	GetDeviceMaxPCIeInfo(idx uint) (int, int, error)
	GetDeviceMemEccCounter(idx uint) ([]uint32, []uint64, error)
	GetDeviceMemory(idx uint) (int64, int64, int64, int64, int64, error)
	GetDeviceMemoryDieCount(idx uint) (int, error)
	GetDeviceMimProfileInfo(idx uint) ([]MimProfileInfo, error)
	GetDeviceMimProfileMaxInstanceCount(idx, profile uint) (int, error)
	GetDeviceMLULinkCapability(idx, link uint) (uint, uint, error)
	GetDeviceMLULinkCounter(idx, link uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error)
	GetDeviceMLULinkErrorCounter(idx, link uint) (uint64, uint64, uint64, error)
	GetDeviceMLULinkEventCounter(idx, link uint) (uint64, error)
	GetDeviceMLULinkPortMode(idx, link uint) (int, error)
	GetDeviceMLULinkPortNumber(idx uint) int
	GetDeviceMLULinkPPI(idx, link uint) (string, error)
	GetDeviceMLULinkRemoteInfo(idx, link uint) (uint64, uint64, uint32, uint32, int32, uint64, string, string, string, string, error)
	GetDeviceMLULinkSpeedInfo(idx, link uint) (float32, int, error)
	GetDeviceMLULinkState(idx, link uint) (int, int, error)
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
	GetDevicePCIeInfo(idx uint) (int, uint, uint, uint16, uint16, uint, uint, uint, uint, uint16, error)
	GetDevicePCIeReplayCount(idx uint) (uint32, error)
	GetDevicePCIeThroughput(idx uint) (int64, int64, error)
	GetDevicePerformanceThrottleReason(idx uint) (bool, error)
	GetDevicePower(idx uint) (int, error)
	GetDevicePowerManagementDefaultLimitation(idx uint) (uint16, error)
	GetDevicePowerManagementLimitation(idx uint) (uint16, error)
	GetDevicePowerManagementLimitRange(idx uint) (uint16, uint16, error)
	GetDeviceProcessInfo(idx uint) ([]uint32, []uint64, []uint64, error)
	GetDeviceProcessUtil(idx uint) ([]uint32, []uint32, []uint32, []uint32, []uint32, []uint32, error)
	GetDeviceRemappedRows(idx uint) (uint32, uint32, uint32, uint32, uint32, error)
	GetDeviceRepairStatus(idx uint) (bool, bool, bool, error)
	GetDeviceRetiredPageInfo(idx uint) (uint32, uint32, error)
	GetDeviceRetiredPagesOperation(idx uint) (int, error)
	GetDeviceSMluProfileIDInfo(idx uint) ([]int, error)
	GetDeviceSMluProfileInfo(idx, profile uint) (SmluProfileInfo, uint32, uint32, error)
	GetDeviceSN(idx uint) (string, error)
	GetDeviceTemperature(idx uint) (int, int, int, []int, []int, error)
	GetDeviceTensorUtil(idx uint) (int, error)
	GetDeviceTinyCoreUtil(idx uint) ([]int, error)
	GetDeviceUtil(idx uint) (int, []int, error)
	GetDeviceFrequency(idx uint) (int, int, int, error)
	GetDeviceUUID(idx uint) (string, error)
	GetDeviceVersion(idx uint) (uint, uint, uint, uint, uint, uint, error)
	GetDeviceVfState(idx uint) (int, error)
	GetDeviceVideoCodecUtil(idx uint) ([]int, []int, error)
	GetDeviceVoltageInfo(idx uint) (int, int, int, error)
	GetDeviceFrequencyStatus(idx uint) (int, error)
	RegisterEventsHandleAndWait(slots []int, ch chan XIDInfoWithTimestamp) error
	GetTopologyRelationship(domain1, bus1, device1, function1, domain2, bus2, device2, function2 uint) (int, error)
	GetSupportedEventTypes(idx uint) error
}

var (
	cndevGlobalHandleMap = &sync.Map{}
)

type cndev struct {
	cndevHandleMap *sync.Map
}

func (c *cndev) Load(key uint) C.cndevDevice_t {
	value, ok := c.cndevHandleMap.Load(key)
	if !ok {
		log.Warnf("Cndev handle not found for key: %d, this should never happen, directly use key as value", key)
		return C.cndevDevice_t(int32(key))
	}
	return value.(C.cndevDevice_t)
}

func NewCndevClient() Cndev {
	return &cndev{
		cndevHandleMap: cndevGlobalHandleMap,
	}
}

func (c *cndev) Init(healthCheck bool) error {
	if healthCheck {
		return (errorString(C.cndevInit(C.int(0))))
	}
	return errorString(dl.cndevInit())
}

func (c *cndev) Release() error {
	r := dl.cndevRelease()
	return errorString(r)
}

func (c *cndev) DeviceMimModeEnabled(idx uint) (bool, error) {
	if ret := dl.checkExist("cndevGetMimMode"); ret != C.CNDEV_SUCCESS {
		return false, errorString(ret)
	}

	var mode C.cndevMimMode_t
	mode.version = C.CNDEV_VERSION_6
	r := C.cndevGetMimMode(&mode, c.Load(idx))
	return mode.mimMode == C.CNDEV_FEATURE_ENABLED, errorString(r)
}

func (c *cndev) DeviceSmluModeEnabled(idx uint) (bool, error) {
	if ret := dl.checkExist("cndevGetSMLUMode"); ret != C.CNDEV_SUCCESS {
		return false, errorString(ret)
	}

	var mode C.cndevSMLUMode_t
	mode.version = C.CNDEV_VERSION_6
	r := C.cndevGetSMLUMode(&mode, c.Load(idx))
	return mode.smluMode == C.CNDEV_FEATURE_ENABLED, errorString(r)
}

func (c *cndev) GetAllMLUInstanceInfo(idx uint) ([]MimInfo, error) {
	if ret := dl.checkExist("cndevGetAllMluInstanceInfo"); ret != C.CNDEV_SUCCESS {
		return nil, errorString(ret)
	}

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
	r := C.cndevGetAllMluInstanceInfo(&miCount, infs, c.Load(idx))
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
			r = C.cndevGetAllMluInstanceInfo(&miCount, infs, c.Load(idx))
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
	if ret := dl.checkExist("cndevGetAllSMluInstanceInfo"); ret != C.CNDEV_SUCCESS {
		return nil, errorString(ret)
	}

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
	r := C.cndevGetAllSMluInstanceInfo(&smluCount, infs, c.Load(idx))
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
			r = C.cndevGetAllSMluInstanceInfo(&smluCount, infs, c.Load(idx))
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

func (c *cndev) GetDeviceActivity(idx uint) (int, error) {
	if ret := dl.checkExist("cndevDeviceGetFieldValues"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

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

	value.fieldId = C.cndevFieldMLUUtilization

	r := C.cndevDeviceGetFieldValues(c.Load(idx), C.int(1), value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	if err := errorString(value.ret); err != nil {
		return 0, err
	}
	activity := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(value)) + unsafe.Sizeof(fieldVauleInfo.fieldID) + unsafe.Sizeof(fieldVauleInfo.ret)))
	return int(activity), nil
}

func (c *cndev) GetDeviceAddressSwaps(idx uint) (uint32, uint32, uint32, uint32, uint32, error) {
	if ret := dl.checkExist("cndevGetAddressSwaps"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, 0, 0, errorString(ret)
	}

	var addressSwap C.cndevAddressSwap_t
	addressSwap.version = C.CNDEV_VERSION_6
	r := C.cndevGetAddressSwaps(&addressSwap, c.Load(idx))
	return uint32(addressSwap.correctCounts), uint32(addressSwap.uncorrectCounts), uint32(addressSwap.histogram[C.CNDEV_AVAILABILITY_XLABLE_NONE]), uint32(addressSwap.histogram[C.CNDEV_AVAILABILITY_XLABLE_PARTIAL]), uint32(addressSwap.histogram[C.CNDEV_AVAILABILITY_XLABLE_MAX]), errorString(r)
}

func (c *cndev) GetDeviceOsMemory(idx uint) (int64, int64, error) {
	if ret := dl.checkExist("cndevGetDeviceOsMemoryUsageV2"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var deviceOsMemInfo C.cndevDeviceOsMemoryInfo_t
	r := C.cndevGetDeviceOsMemoryUsageV2(&deviceOsMemInfo, c.Load(idx))
	return int64(deviceOsMemInfo.deviceSystemMemoryUsed), int64(deviceOsMemInfo.deviceSystemMemoryTotal), errorString(r)
}

func (c *cndev) GetDeviceBAR4MemoryInfo(idx uint) (uint64, uint64, uint64, error) {
	if ret := dl.checkExist("cndevGetBAR4MemoryInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, errorString(ret)
	}

	var memInfo C.cndevBAR4Memory_t
	r := C.cndevGetBAR4MemoryInfo(&memInfo, c.Load(idx))
	return uint64(memInfo.total_size), uint64(memInfo.used_size), uint64(memInfo.free_size), errorString(r)
}

func (c *cndev) GetDeviceChassisInfo(idx uint) (uint64, string, string, string, string, string, []ChassisDevInfo, []ChassisDevInfo, []ChassisDevInfo, error) {
	if ret := dl.checkExist("cndevGetChassisInfoV2"); ret != C.CNDEV_SUCCESS {
		return 0, "", "", "", "", "", nil, nil, nil, errorString(ret)
	}

	var chassisInfo C.cndevChassisInfoV2_t
	chassisInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetChassisInfoV2(&chassisInfo, c.Load(idx))
	if err := errorString(r); err != nil {
		return 0, "", "", "", "", "", nil, nil, nil, err
	}
	nvme := make([]ChassisDevInfo, uint8(chassisInfo.nvmeSsdNum))
	ib := make([]ChassisDevInfo, uint8(chassisInfo.ibBoardNum))
	psu := make([]ChassisDevInfo, uint8(chassisInfo.psuNum))
	for i := uint8(0); i < uint8(chassisInfo.nvmeSsdNum); i++ {
		nvme = append(nvme,
			ChassisDevInfo{
				SN:    C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.nvmeInfo[i].nvmeSn[0]))),
				Model: C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.nvmeInfo[i].nvmeModel[0]))),
				Fw:    C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.nvmeInfo[i].nvmeFw[0]))),
				Mfc:   C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.nvmeInfo[i].nvmeMfc[0]))),
			})
	}
	log.Debugf("chassis nvme infos for device %d num: %d info: %+v", idx, chassisInfo.nvmeSsdNum, nvme)
	for i := uint8(0); i < uint8(chassisInfo.ibBoardNum); i++ {
		ib = append(ib,
			ChassisDevInfo{
				SN:    C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.ibInfo[i].ibSn[0]))),
				Model: C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.ibInfo[i].ibModel[0]))),
				Fw:    C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.ibInfo[i].ibFw[0]))),
				Mfc:   C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.ibInfo[i].ibMfc[0]))),
			})
	}
	log.Debugf("chassis ib infos for device %d num: %d info: %+v", idx, chassisInfo.ibBoardNum, ib)
	for i := uint8(0); i < uint8(chassisInfo.psuNum); i++ {
		psu = append(psu,
			ChassisDevInfo{
				SN:    C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.psuInfo[i].psuSn[0]))),
				Model: C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.psuInfo[i].psuModel[0]))),
				Fw:    C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.psuInfo[i].psuFw[0]))),
				Mfc:   C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.psuInfo[i].psuMfc[0]))),
			})
	}
	log.Debugf("chassis psu for device %d num: %d info: %+v", idx, chassisInfo.psuNum, psu)

	return uint64(chassisInfo.chassisSn), C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.chassisProductDate))), C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.chassisProductName))),
		C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.chassisVendorName))), C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.chassisPartNumber))),
		C.GoString((*C.char)(unsafe.Pointer(&chassisInfo.bmcIP[0]))), nvme, ib, psu, errorString(r)
}

func (c *cndev) GetDeviceClusterCount(idx uint) (int, error) {
	if ret := dl.checkExist("cndevGetClusterCount"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var cardClusterNum C.cndevCardClusterCount_t
	cardClusterNum.version = C.CNDEV_VERSION_6
	r := C.cndevGetClusterCount(&cardClusterNum, c.Load(idx))
	return int(cardClusterNum.count), errorString(r)
}

func (c *cndev) GetDeviceCndevVersion() (uint, uint, uint, error) {
	if ret := dl.checkExist("cndevGetLibVersion"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, errorString(ret)
	}

	var cndevVersion C.cndevLibVersionInfo_t
	cndevVersion.version = C.CNDEV_VERSION_6
	r := C.cndevGetLibVersion(&cndevVersion)
	return uint(cndevVersion.libMajorVersion), uint(cndevVersion.libMinorVersion), uint(cndevVersion.libBuildVersion), errorString(r)
}

func (c *cndev) GetDeviceComputeCapability(idx uint) (uint, uint, error) {
	if ret := dl.checkExist("cndevGetComputeCapability"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var major, minor C.uint
	r := C.cndevGetComputeCapability(&major, &minor, c.Load(idx))
	return uint(major), uint(minor), errorString(r)
}

func (c *cndev) GetDeviceComputeMode(idx uint) (uint, error) {
	if ret := dl.checkExist("cndevGetComputeMode"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var computeMode C.cndevComputeMode_t
	computeMode.version = C.CNDEV_VERSION_6
	r := C.cndevGetComputeMode(&computeMode, c.Load(idx))
	return uint(computeMode.mode), errorString(r)
}

func (c *cndev) GetDeviceCoreNum(idx uint) (uint, error) {
	if ret := dl.checkExist("cndevGetCoreCount"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}
	var cardCoreNum C.cndevCardCoreCount_t
	cardCoreNum.version = C.CNDEV_VERSION_6
	r := C.cndevGetCoreCount(&cardCoreNum, c.Load(idx))
	return uint(cardCoreNum.count), errorString(r)
}

func (c *cndev) GetDeviceCount() (uint, error) {
	if ret := dl.checkExist("cndevGetDeviceCount"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var cardInfo C.cndevCardInfo_t
	cardInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetDeviceCount(&cardInfo)
	return uint(cardInfo.number), errorString(r)
}

func (c *cndev) GetDeviceCPUUtil(idx uint) (uint16, []uint8, []uint8, []uint8, []uint8, []uint8, []uint8, error) {
	if ret := dl.checkExist("cndevGetDeviceCPUUtilizationV2"); ret != C.CNDEV_SUCCESS {
		return 0, nil, nil, nil, nil, nil, nil, errorString(ret)
	}

	var cardCPUUtil C.cndevDeviceCPUUtilizationV2_t
	cardCPUUtil.version = C.CNDEV_VERSION_6
	r := C.cndevGetDeviceCPUUtilizationV2(&cardCPUUtil, c.Load(idx))
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
	if ret := dl.checkExist("cndevGetCRCInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var cardCRCInfo C.cndevCRCInfo_t
	cardCRCInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetCRCInfo(&cardCRCInfo, c.Load(idx))
	return uint64(cardCRCInfo.die2dieCRCError), uint64(cardCRCInfo.die2dieCRCErrorOverflow), errorString(r)
}

func (c *cndev) GetDeviceCurrentInfo(idx uint) (int, int, int, error) {
	if ret := dl.checkExist("cndevGetCurrentInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, errorString(ret)
	}

	var currentInfo C.cndevCurrentInfo_t
	r := C.cndevGetCurrentInfo(&currentInfo, c.Load(idx))
	return int(currentInfo.ipuCoreCurrent), int(currentInfo.socCurrent), int(currentInfo.hbmCurrent), errorString(r)
}

func (c *cndev) GetDeviceCurrentPCIeInfo(idx uint) (int, int, error) {
	if ret := dl.checkExist("cndevGetCurrentPCIInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var currentPCIeInfo C.cndevCurrentPCIInfo_t
	currentPCIeInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetCurrentPCIInfo(&currentPCIeInfo, c.Load(idx))
	return int(currentPCIeInfo.currentSpeed), int(currentPCIeInfo.currentWidth), errorString(r)
}

func (c *cndev) GetDeviceDDRInfo(idx uint) (int, float64, error) {
	if ret := dl.checkExist("cndevGetDDRInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var cardDDRInfo C.cndevDDRInfo_t
	cardDDRInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetDDRInfo(&cardDDRInfo, c.Load(idx))
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
	if ret := dl.checkExist("cndevGetVersionInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, errorString(ret)
	}

	var cardVersionInfo C.cndevVersionInfo_t
	cardVersionInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetVersionInfo(&cardVersionInfo, c.Load(idx))
	return uint(cardVersionInfo.driverMajorVersion), uint(cardVersionInfo.driverMinorVersion), uint(cardVersionInfo.driverBuildVersion), errorString(r)
}

func (c *cndev) GetDeviceECCInfo(idx uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error) {
	if ret := dl.checkExist("cndevGetECCInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, 0, 0, 0, 0, 0, errorString(ret)
	}

	var cardECCInfo C.cndevECCInfo_t
	cardECCInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetECCInfo(&cardECCInfo, c.Load(idx))
	return uint64(cardECCInfo.addressForbiddenError), uint64(cardECCInfo.correctedError), uint64(cardECCInfo.multipleError), uint64(cardECCInfo.multipleMultipleError), uint64(cardECCInfo.multipleOneError), uint64(cardECCInfo.oneBitError), uint64(cardECCInfo.totalError), uint64(cardECCInfo.uncorrectedError), errorString(r)
}

func (c *cndev) GetDeviceEccMode(idx uint) (int, int, error) {
	if ret := dl.checkExist("cndevGetEccMode"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var currentMode, pendingMode C.cndevEccMode_t
	currentMode.version = C.CNDEV_VERSION_6
	pendingMode.version = C.CNDEV_VERSION_6
	r := C.cndevGetEccMode(&currentMode, &pendingMode, c.Load(idx))
	return int(currentMode.mode), int(pendingMode.mode), errorString(r)
}

func (c *cndev) GetDeviceFanSpeed(idx uint) (int, error) {
	if ret := dl.checkExist("cndevGetFanSpeedInfo"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var cardFanSpeed C.cndevFanSpeedInfo_t
	cardFanSpeed.version = C.CNDEV_VERSION_6
	r := C.cndevGetFanSpeedInfo(&cardFanSpeed, c.Load(idx))
	return int(cardFanSpeed.fanSpeed), errorString(r)
}

func (c *cndev) GetDeviceFrequency(idx uint) (int, int, int, error) {
	if ret := dl.checkExist("cndevGetFrequencyInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, errorString(ret)
	}

	var cardFrequency C.cndevFrequencyInfo_t
	cardFrequency.version = C.CNDEV_VERSION_6
	r := C.cndevGetFrequencyInfo(&cardFrequency, c.Load(idx))
	if err := errorString(r); err != nil {
		return 0, 0, 0, err
	}
	return int(cardFrequency.boardFreq), int(cardFrequency.ddrFreq), int(cardFrequency.boardDefaultFreq), nil
}

func (c *cndev) GetDeviceHealth(idx uint) (int, bool, bool, []string, error) {
	if ret := dl.checkExist("cndevGetCardHealthStateV2"); ret != C.CNDEV_SUCCESS {
		return 0, false, false, nil, errorString(ret)
	}

	var cardHealthState C.cndevCardHealthStateV2_t
	r := C.cndevGetCardHealthStateV2(&cardHealthState, c.Load(idx))
	unhealth := make([]string, 0, cardHealthState.incident_count)
	for i := 0; i < int(cardHealthState.incident_count); i++ {
		errorCode := cardHealthState.incidents[i].error.code
		unhealth = append(unhealth, healthErrorCodeToString(errorCode))
	}
	return int(cardHealthState.health), cardHealthState.deviceState == C.CNDEV_HEALTH_STATE_DEVICE_GOOD, cardHealthState.driverState == C.CNDEV_HEALTH_STATE_DRIVER_RUNNING, unhealth, errorString(r)
}

func (c *cndev) GetDeviceHeartbeatCount(idx uint) (uint32, error) {
	if ret := dl.checkExist("cndevGetCardHeartbeatCount"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var heartbeatCount C.cndevCardHeartbeatCount_t
	heartbeatCount.version = C.CNDEV_VERSION_6
	r := C.cndevGetCardHeartbeatCount(&heartbeatCount, c.Load(idx))
	if err := errorString(r); err != nil {
		return 0, err
	}
	return uint32(heartbeatCount.heartbeatCount), nil
}

func (c *cndev) GetDeviceImageCodecUtil(idx uint) ([]int, error) {
	if ret := dl.checkExist("cndevGetImageCodecUtilization"); ret != C.CNDEV_SUCCESS {
		return nil, errorString(ret)
	}

	var cardImageCodecUtil C.cndevImageCodecUtilization_t
	cardImageCodecUtil.version = C.CNDEV_VERSION_6
	r := C.cndevGetImageCodecUtilization(&cardImageCodecUtil, c.Load(idx))
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
	if ret := dl.checkExist("cndevGetMaxPCIInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var pcieInfo C.cndevPCIInfo_t
	pcieInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetMaxPCIInfo(&pcieInfo, c.Load(idx))
	return int(pcieInfo.maxSpeed), int(pcieInfo.maxWidth), errorString(r)
}

func (c *cndev) GetDeviceMemEccCounter(idx uint) ([]uint32, []uint64, error) {
	if ret := dl.checkExist("cndevGetMemEccCounter"); ret != C.CNDEV_SUCCESS {
		return nil, nil, errorString(ret)
	}

	var memEccCounter C.cndevMemEccCounter_t
	memEccCounter.version = C.CNDEV_VERSION_6
	r := C.cndevGetMemEccCounter(&memEccCounter, c.Load(idx))
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

func (c *cndev) GetDeviceMemory(idx uint) (int64, int64, int64, int64, int64, error) {
	if ret := dl.checkExist("cndevGetMemoryUsageV2"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, 0, 0, errorString(ret)
	}

	var cardMemInfo C.cndevMemoryInfoV2_t
	r := C.cndevGetMemoryUsageV2(&cardMemInfo, c.Load(idx))
	return int64(cardMemInfo.physicalMemoryUsed), int64(cardMemInfo.globalMemory), int64(cardMemInfo.virtualMemoryUsed), int64(cardMemInfo.virtualMemoryTotal), int64(cardMemInfo.reservedMemory), errorString(r)
}

func (c *cndev) GetDeviceMemoryDieCount(idx uint) (int, error) {
	//this api is deprecated, keep to comply with old version
	log.Debugf("GetDeviceMemoryDieCount for slot:%d, this is deprecated, ignore", idx)
	return 0, nil
}

func (c *cndev) GetDeviceMimProfileInfo(idx uint) ([]MimProfileInfo, error) {
	if ret := dl.checkExist("cndevGetMluInstanceProfileInfo"); ret != C.CNDEV_SUCCESS {
		return nil, errorString(ret)
	}

	var mimProfileInfos []MimProfileInfo

	for profileID := 0; profileID < C.CNDEV_MLUINSTANCE_PROFILE_COUNT; profileID++ {
		var profileInfo C.cndevMluInstanceProfileInfo_t
		profileInfo.version = C.CNDEV_VERSION_6
		r := C.cndevGetMluInstanceProfileInfo(&profileInfo, C.int(profileID), c.Load(idx))
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
	if ret := dl.checkExist("cndevGetMaxMluInstanceCount"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var count C.int
	r := C.cndevGetMaxMluInstanceCount(&count, C.int(profile), c.Load(idx))
	return int(count), errorString(r)
}

func (c *cndev) GetDeviceMLULinkCapability(idx, link uint) (uint, uint, error) {
	if ret := dl.checkExist("cndevGetMLULinkCapability"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var cardMLULinkCapability C.cndevMLULinkCapability_t
	cardMLULinkCapability.version = C.CNDEV_VERSION_6
	r := C.cndevGetMLULinkCapability(&cardMLULinkCapability, c.Load(idx), C.int(link))
	return uint(cardMLULinkCapability.p2pTransfer), uint(cardMLULinkCapability.interlakenSerdes), errorString(r)
}

func (c *cndev) GetDeviceMLULinkCounter(idx, link uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error) {
	if ret := dl.checkExist("cndevGetMLULinkCounter"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, errorString(ret)
	}

	var cardMLULinkCount C.cndevMLULinkCounter_t
	cardMLULinkCount.version = C.CNDEV_VERSION_6
	r := C.cndevGetMLULinkCounter(&cardMLULinkCount, c.Load(idx), C.int(link))
	return uint64(cardMLULinkCount.cntrReadByte), uint64(cardMLULinkCount.cntrReadPackage), uint64(cardMLULinkCount.cntrWriteByte), uint64(cardMLULinkCount.cntrWritePackage),
		uint64(cardMLULinkCount.errCorrected), uint64(cardMLULinkCount.errCRC24), uint64(cardMLULinkCount.errCRC32), uint64(cardMLULinkCount.errEccDouble),
		uint64(cardMLULinkCount.errFatal), uint64(cardMLULinkCount.errReplay), uint64(cardMLULinkCount.errUncorrected), uint64(cardMLULinkCount.cntrCnpPackage),
		uint64(cardMLULinkCount.cntrPfcPackage), errorString(r)
}

func (c *cndev) GetDeviceMLULinkErrorCounter(idx, link uint) (uint64, uint64, uint64, error) {
	if ret := dl.checkExist("cndevGetMLULinkErrorCounter"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, errorString(ret)
	}

	var cardMLULinkErrorCounter C.cndevMLULinkErrorCounter_t
	r := C.cndevGetMLULinkErrorCounter(&cardMLULinkErrorCounter, c.Load(idx), C.int(link))
	return uint64(cardMLULinkErrorCounter.illegalAccessCnt), uint64(cardMLULinkErrorCounter.correctFecCnt), uint64(cardMLULinkErrorCounter.uncorrectFecCnt), errorString(r)
}

func (c *cndev) GetDeviceMLULinkEventCounter(idx, link uint) (uint64, error) {
	if ret := dl.checkExist("cndevGetMLULinkEventCounter"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var cardMLULinkEventCounter C.cndevMLULinkEventCounter_t
	r := C.cndevGetMLULinkEventCounter(&cardMLULinkEventCounter, c.Load(idx), C.int(link))
	return uint64(cardMLULinkEventCounter.linkDown), errorString(r)
}

func (c *cndev) GetDeviceMLULinkPortMode(idx, link uint) (int, error) {
	if ret := dl.checkExist("cndevGetMLULinkPortMode"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var cardMLULinkPortMode C.cndevMLULinkPortMode_t
	cardMLULinkPortMode.version = C.CNDEV_VERSION_6
	r := C.cndevGetMLULinkPortMode(&cardMLULinkPortMode, c.Load(idx), C.int(link))
	return int(cardMLULinkPortMode.mode), errorString(r)
}

func (c *cndev) GetDeviceMLULinkPortNumber(idx uint) int {
	if ret := dl.checkExist("cndevGetMLULinkPortNumber"); ret != C.CNDEV_SUCCESS {
		return -1
	}

	return int(C.cndevGetMLULinkPortNumber(c.Load(idx)))
}

func (c *cndev) GetDeviceMLULinkPPI(idx, link uint) (string, error) {
	if ret := dl.checkExist("cndevGetMLULinkPPI"); ret != C.CNDEV_SUCCESS {
		return "", errorString(ret)
	}

	var cardMLULinkPPI C.cndevMLULinkPPI_t
	r := C.cndevGetMLULinkPPI(&cardMLULinkPPI, c.Load(idx), C.int(link))
	return C.GoString((*C.char)(unsafe.Pointer(&cardMLULinkPPI.ppi))), errorString(r)
}

func (c *cndev) GetDeviceMLULinkSpeedInfo(idx, link uint) (float32, int, error) {
	if ret := dl.checkExist("cndevGetMLULinkSpeedInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var cardMLULinkSpeedInfo C.cndevMLULinkSpeed_t
	cardMLULinkSpeedInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetMLULinkSpeedInfo(&cardMLULinkSpeedInfo, c.Load(idx), C.int(link))
	return float32(cardMLULinkSpeedInfo.speedValue), int(cardMLULinkSpeedInfo.speedFormat), errorString(r)
}

func (c *cndev) GetDeviceMLULinkRemoteInfo(idx, link uint) (uint64, uint64, uint32, uint32, int32, uint64, string, string, string, string, error) {
	if ret := dl.checkExist("cndevGetMLULinkRemoteInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", errorString(ret)
	}

	var cardMLULinkRemoteInfo C.cndevMLULinkRemoteInfo_t
	cardMLULinkRemoteInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetMLULinkRemoteInfo(&cardMLULinkRemoteInfo, c.Load(idx), C.int(link))
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

func (c *cndev) GetDeviceMLULinkState(idx, link uint) (int, int, error) {
	if ret := dl.checkExist("cndevSetMLULinkState"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var ibCmdInfo, obCmdInfo C.int
	r := C.cndevGetMLULinkState(&ibCmdInfo, &obCmdInfo, c.Load(idx), C.int(link))
	return int(ibCmdInfo), int(obCmdInfo), errorString(r)
}

func (c *cndev) GetDeviceMLULinkStatus(idx, link uint) (int, int, int, error) {
	if ret := dl.checkExist("cndevGetMLULinkStatusV2"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, errorString(ret)
	}

	var cardMLULinkStatus C.cndevMLULinkStatusV2_t
	r := C.cndevGetMLULinkStatusV2(&cardMLULinkStatus, c.Load(idx), C.int(link))
	return int(cardMLULinkStatus.macState), int(cardMLULinkStatus.serdesState), int(cardMLULinkStatus.presenceState), errorString(r)
}

func (c *cndev) GetDeviceMLULinkVersion(idx, link uint) (uint, uint, uint, error) {
	if ret := dl.checkExist("cndevGetMLULinkVersion"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, errorString(ret)
	}

	var cardMLULinkVersion C.cndevMLULinkVersion_t
	cardMLULinkVersion.version = C.CNDEV_VERSION_6
	r := C.cndevGetMLULinkVersion(&cardMLULinkVersion, c.Load(idx), C.int(link))
	return uint(cardMLULinkVersion.majorVersion), uint(cardMLULinkVersion.minorVersion), uint(cardMLULinkVersion.buildVersion), errorString(r)
}

func (c *cndev) GetDeviceMLUInstanceByID(idx uint, instanceID int) (uint, error) {
	if ret := dl.checkExist("cndevGetMluInstanceById"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var handle C.cndevMluInstance_t
	r := C.cndevGetMluInstanceById(&handle, C.int(instanceID), c.Load(idx))
	return uint(handle), errorString(r)
}

func (c *cndev) GetDeviceModel(idx uint) string {
	if ret := dl.checkExist("cndevGetCardNameStringByDevId"); ret != C.CNDEV_SUCCESS {
		return ""
	}

	return C.GoString(C.cndevGetCardNameStringByDevId(c.Load(idx)))
}

func (c *cndev) GetDeviceNUMANodeID(idx uint) (int, error) {
	if ret := dl.checkExist("cndevGetNUMANodeIdByDevId"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var cardNUMANodeID C.cndevNUMANodeId_t
	cardNUMANodeID.version = C.CNDEV_VERSION_6
	r := C.cndevGetNUMANodeIdByDevId(&cardNUMANodeID, c.Load(idx))
	return int(cardNUMANodeID.nodeId), errorString(r)
}

func (c *cndev) GetDeviceOverTemperatureInfo(idx uint) (uint32, uint32, error) {
	if ret := dl.checkExist("cndevGetOverTemperatureInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var overTemperatureInfo C.cndevOverTemperatureInfo_t
	overTemperatureInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetOverTemperatureInfo(&overTemperatureInfo, c.Load(idx))
	return uint32(overTemperatureInfo.powerOffCounter), uint32(overTemperatureInfo.underClockCounter), errorString(r)
}

func (c *cndev) GetDeviceOpticalInfo(idx, link uint) (uint8, float32, float32, []float32, []float32, error) {
	if ret := dl.checkExist("cndevGetOpticalInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, nil, nil, errorString(ret)
	}

	var opticalInfo C.cndevOpticalInfo_t
	r := C.cndevGetOpticalInfo(&opticalInfo, c.Load(idx), C.int(link))
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
	if ret := dl.checkExist("cndevDeviceGetFieldValues"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

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
	r := C.cndevDeviceGetFieldValues(c.Load(idx), C.int(1), value)
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
	if ret := dl.checkExist("cndevDeviceGetFieldValues"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

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
	r := C.cndevDeviceGetFieldValues(c.Load(idx), C.int(1), value)
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
	if ret := dl.checkExist("cndevGetParityError"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var parityError C.cndevParityError_t
	parityError.version = C.CNDEV_VERSION_6
	r := C.cndevGetParityError(&parityError, c.Load(idx))
	return int(parityError.counter), errorString(r)
}

func (c *cndev) GetDevicePCIeInfo(idx uint) (int, uint, uint, uint16, uint16, uint, uint, uint, uint, uint16, error) {
	if ret := dl.checkExist("cndevGetPCIeInfoV2"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, errorString(ret)
	}

	var pcieInfo C.cndevPCIeInfoV2_t
	r := C.cndevGetPCIeInfoV2(&pcieInfo, c.Load(idx))
	return int(pcieInfo.slotId), uint(pcieInfo.subsystemId), uint(pcieInfo.deviceId), uint16(pcieInfo.vendor),
		uint16(pcieInfo.subsystemVendor), uint(pcieInfo.domain), uint(pcieInfo.bus), uint(pcieInfo.device),
		uint(pcieInfo.function), uint16(pcieInfo.moduleId), errorString(r)
}

func (c *cndev) GetDevicePCIeReplayCount(idx uint) (uint32, error) {
	if ret := dl.checkExist("cndevGetPcieReplayCounter"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var pcieReplay C.cndevPcieReplayCounter_t
	pcieReplay.version = C.CNDEV_VERSION_6
	r := C.cndevGetPcieReplayCounter(&pcieReplay, c.Load(idx))
	return uint32(pcieReplay.counter), errorString(r)
}

func (c *cndev) GetDevicePCIeThroughput(idx uint) (int64, int64, error) {
	if ret := dl.checkExist("cndevGetPCIethroughputV2"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var pcieThroughput C.cndevPCIethroughputV2_t
	r := C.cndevGetPCIethroughputV2(&pcieThroughput, c.Load(idx))
	return int64(pcieThroughput.pcieRead), int64(pcieThroughput.pcieWrite), errorString(r)
}

func (c *cndev) GetDevicePower(idx uint) (int, error) {
	if ret := dl.checkExist("cndevGetDevicePowerInfo"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var devicePower C.cndevDevicePowerInfo_t

	r := C.cndevGetDevicePowerInfo(&devicePower, c.Load(idx))
	return int(devicePower.usage), errorString(r)
}

func (c *cndev) GetDevicePowerManagementDefaultLimitation(idx uint) (uint16, error) {
	if ret := dl.checkExist("cndevGetPowerManagementDefaultLimitation"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var powerManagementLimitation C.cndevPowerManagementLimitation_t
	powerManagementLimitation.version = C.CNDEV_VERSION_6
	r := C.cndevGetPowerManagementDefaultLimitation(&powerManagementLimitation, c.Load(idx))
	return uint16(powerManagementLimitation.powerLimit), errorString(r)
}

func (c *cndev) GetDevicePowerManagementLimitation(idx uint) (uint16, error) {
	if ret := dl.checkExist("cndevGetPowerManagementLimitation"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var powerManagementLimitation C.cndevPowerManagementLimitation_t
	powerManagementLimitation.version = C.CNDEV_VERSION_6
	r := C.cndevGetPowerManagementLimitation(&powerManagementLimitation, c.Load(idx))
	return uint16(powerManagementLimitation.powerLimit), errorString(r)
}

func (c *cndev) GetDevicePowerManagementLimitRange(idx uint) (uint16, uint16, error) {
	if ret := dl.checkExist("cndevGetPowerManagementLimitationRange"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var powerManagementLimitRange C.cndevPowerManagementLimitationRange_t
	powerManagementLimitRange.version = C.CNDEV_VERSION_6
	r := C.cndevGetPowerManagementLimitationRange(&powerManagementLimitRange, c.Load(idx))
	return uint16(powerManagementLimitRange.minPowerLimit), uint16(powerManagementLimitRange.maxPowerLimit), errorString(r)
}

func (c *cndev) GetDevicePerformanceThrottleReason(idx uint) (bool, error) {
	if ret := dl.checkExist("cndevGetPerformanceThrottleReason"); ret != C.CNDEV_SUCCESS {
		return false, errorString(ret)
	}

	var performanceThrottleReason C.cndevPerformanceThrottleReason_t
	performanceThrottleReason.version = C.CNDEV_VERSION_6
	r := C.cndevGetPerformanceThrottleReason(&performanceThrottleReason, c.Load(idx))
	return performanceThrottleReason.thermalSlowdown == C.CNDEV_FEATURE_ENABLED, errorString(r)
}

func (c *cndev) GetDeviceProcessInfo(idx uint) ([]uint32, []uint64, []uint64, error) {
	if ret := dl.checkExist("cndevGetProcessInfo"); ret != C.CNDEV_SUCCESS {
		return nil, nil, nil, errorString(ret)
	}

	processCount := 1 << 4
	var util C.cndevProcessInfo_t
	utils := (*C.cndevProcessInfo_t)(C.malloc(C.size_t(processCount) * C.size_t(unsafe.Sizeof(util))))
	if utils == nil {
		return nil, nil, nil, fmt.Errorf("malloc failed for cndevGetProcessInfo")
	}
	defer func() {
		if utils != nil {
			C.free(unsafe.Pointer(utils))
		}
	}()

	utils.version = C.CNDEV_VERSION_6
	r := C.cndevGetProcessInfo((*C.uint)(unsafe.Pointer(&processCount)), utils, c.Load(idx))
	if err := errorString(r); err != nil && r != C.CNDEV_ERROR_INSUFFICIENT_SPACE {
		return nil, nil, nil, err
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
			log.Debugf("cndevGetProcessInfo with insufficient space with real counts %d, with slot: %d, will try with the real counts", processCount, idx)
			newUtils := (*C.cndevProcessInfo_t)(C.realloc(unsafe.Pointer(utils), C.size_t(processCount)*C.size_t(unsafe.Sizeof(util))))
			if newUtils == nil {
				return nil, nil, nil, fmt.Errorf("realloc failed for cndevGetProcessInfo")
			}
			utils = newUtils
			r = C.cndevGetProcessInfo((*C.uint)(unsafe.Pointer(&processCount)), utils, c.Load(idx))
			continue
		}
		return nil, nil, nil, err
	}

	pids := make([]uint32, processCount)
	physicalMemUsed := make([]uint64, processCount)
	virtualMemUsed := make([]uint64, processCount)
	results := (*[1 << 16]C.cndevProcessInfo_t)(unsafe.Pointer(utils))[:processCount]
	for i := 0; i < processCount; i++ {
		pids[i] = uint32(results[i].pid)
		physicalMemUsed[i] = uint64(results[i].physicalMemoryUsed)
		virtualMemUsed[i] = uint64(results[i].virtualMemoryUsed)
	}

	return pids, physicalMemUsed, virtualMemUsed, nil
}

func (c *cndev) GetDeviceProcessUtil(idx uint) ([]uint32, []uint32, []uint32, []uint32, []uint32, []uint32, error) {
	if ret := dl.checkExist("cndevGetProcessUtilization"); ret != C.CNDEV_SUCCESS {
		return nil, nil, nil, nil, nil, nil, errorString(ret)
	}

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

	utils.version = C.CNDEV_VERSION_6
	r := C.cndevGetProcessUtilization((*C.uint)(unsafe.Pointer(&processCount)), utils, c.Load(idx))
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
			r = C.cndevGetProcessUtilization((*C.uint)(unsafe.Pointer(&processCount)), utils, c.Load(idx))
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
		if ret := dl.checkExist("cndevGetRemappedRowsV2", "cndevGetRepairStatus"); ret != C.CNDEV_SUCCESS {
			return 0, 0, 0, 0, 0, errorString(ret)
		}

		var remappedRowsV2 C.cndevRemappedRowV2_t
		remappedRowsV2.version = C.CNDEV_VERSION_6
		r := C.cndevGetRemappedRowsV2(&remappedRowsV2, c.Load(idx))
		if err := errorString(r); err != nil {
			return 0, 0, 0, 0, 0, err
		}
		var repairStatus C.cndevRepairStatus_t
		repairStatus.version = C.CNDEV_VERSION_6
		r = C.cndevGetRepairStatus(&repairStatus, c.Load(idx))
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

	if ret := dl.checkExist("cndevGetRemappedRows"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, 0, 0, err
	}

	var remappedRows C.cndevRemappedRow_t
	remappedRows.version = C.CNDEV_VERSION_6
	r := C.cndevGetRemappedRows(&remappedRows, c.Load(idx))
	return uint32(remappedRows.correctRows), uint32(remappedRows.uncorrectRows),
		uint32(remappedRows.pendingRows), uint32(remappedRows.failedRows), 0, errorString(r)
}

func (c *cndev) GetDeviceRepairStatus(idx uint) (bool, bool, bool, error) {
	if ret := dl.checkExist("cndevGetRepairStatus"); ret != C.CNDEV_SUCCESS {
		return false, false, false, errorString(ret)
	}

	var repairStatus C.cndevRepairStatus_t
	repairStatus.version = C.CNDEV_VERSION_6
	r := C.cndevGetRepairStatus(&repairStatus, c.Load(idx))
	return bool(repairStatus.isPending), bool(repairStatus.isFailure), bool(repairStatus.isRetirePending), errorString(r)
}

func (c *cndev) GetDeviceRetiredPageInfo(idx uint) (uint32, uint32, error) {
	if ret := dl.checkExist("cndevGetRetiredPages"); ret != C.CNDEV_SUCCESS {
		return 0, 0, errorString(ret)
	}

	var retiredPage C.cndevRetiredPageInfo_t
	retiredPage.version = C.CNDEV_VERSION_6
	retiredPage.cause = 0
	r := C.cndevGetRetiredPages(&retiredPage, c.Load(idx))
	if err := errorString(r); err != nil {
		return 0, 0, err
	}
	single := uint32(retiredPage.pageCount)
	retiredPage.cause = 1
	r = C.cndevGetRetiredPages(&retiredPage, c.Load(idx))
	if err := errorString(r); err != nil {
		return single, 0, err
	}
	double := uint32(retiredPage.pageCount)
	return single, double, errorString(r)
}

func (c *cndev) GetDeviceRetiredPagesOperation(idx uint) (int, error) {
	if ret := dl.checkExist("cndevGetRetiredPagesOperation"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var retirement C.cndevRetiredPageOperation_t
	retirement.version = C.CNDEV_VERSION_6
	r := C.cndevGetRetiredPagesOperation(&retirement, c.Load(idx))
	return int(retirement.retirePageOption), errorString(r)
}

func (c *cndev) GetDeviceSMluProfileIDInfo(idx uint) ([]int, error) {
	if ret := dl.checkExist("cndevGetSMluProfileIdInfo"); ret != C.CNDEV_SUCCESS {
		return nil, errorString(ret)
	}

	var profileIDInfo C.cndevSMluProfileIdInfo_t
	profileIDInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetSMluProfileIdInfo(&profileIDInfo, c.Load(idx))
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
	if ret := dl.checkExist("cndevGetSMluProfileInfo"); ret != C.CNDEV_SUCCESS {
		return SmluProfileInfo{}, 0, 0, errorString(ret)
	}

	var profileInfo C.cndevSMluProfileInfo_t
	var smluProfileInfo SmluProfileInfo
	profileInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetSMluProfileInfo(&profileInfo, C.int(profile), c.Load(idx))
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
	if ret := dl.checkExist("cndevGetCardSN"); ret != C.CNDEV_SUCCESS {
		return "", errorString(ret)
	}

	var cardSN C.cndevCardSN_t
	cardSN.version = C.CNDEV_VERSION_6
	r := C.cndevGetCardSN(&cardSN, c.Load(idx))
	if err := errorString(r); err != nil {
		return "", err
	}
	sn := fmt.Sprintf("%x", uint64(cardSN.sn))
	return sn, nil
}

func (c *cndev) GetDeviceTemperature(idx uint) (int, int, int, []int, []int, error) {
	if ret := dl.checkExist("cndevGetTemperatureInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, nil, nil, errorString(ret)
	}

	var cardTemperature C.cndevTemperatureInfo_t
	cardTemperature.version = C.CNDEV_VERSION_6
	cluster, err := c.GetDeviceClusterCount(idx)
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}
	r := C.cndevGetTemperatureInfo(&cardTemperature, c.Load(idx))
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
	if ret := dl.checkExist("cndevDeviceGetFieldValues"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

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

	r := C.cndevDeviceGetFieldValues(c.Load(idx), C.int(1), value)
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
	if ret := dl.checkExist("cndevGetTinyCoreUtilization"); ret != C.CNDEV_SUCCESS {
		return nil, errorString(ret)
	}

	var cardTinyCoreUtil C.cndevTinyCoreUtilization_t
	cardTinyCoreUtil.version = C.CNDEV_VERSION_6
	r := C.cndevGetTinyCoreUtilization(&cardTinyCoreUtil, c.Load(idx))
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
	if ret := dl.checkExist("cndevGetDeviceUtilizationInfo"); ret != C.CNDEV_SUCCESS {
		return 0, nil, errorString(ret)
	}

	var cardUtil C.cndevUtilizationInfo_t
	cardUtil.version = C.CNDEV_VERSION_6
	core, err := c.GetDeviceCoreNum(idx)
	if err != nil {
		return 0, nil, err
	}
	r := C.cndevGetDeviceUtilizationInfo(&cardUtil, c.Load(idx))
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
	if ret := dl.checkExist("cndevGetUUID"); ret != C.CNDEV_SUCCESS {
		return "", errorString(ret)
	}

	var uuidInfo C.cndevUUID_t
	uuidInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetUUID(&uuidInfo, c.Load(idx))
	if err := errorString(r); err != nil {
		return "", err
	}
	return C.GoString((*C.char)(unsafe.Pointer(&uuidInfo.uuid))), nil
}

func (c *cndev) GetDeviceVersion(idx uint) (uint, uint, uint, uint, uint, uint, error) {
	if ret := dl.checkExist("cndevGetVersionInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, 0, 0, 0, errorString(ret)
	}

	var versionInfo C.cndevVersionInfo_t
	versionInfo.version = C.CNDEV_VERSION_6
	r := C.cndevGetVersionInfo(&versionInfo, c.Load(idx))
	return uint(versionInfo.mcuMajorVersion), uint(versionInfo.mcuMinorVersion), uint(versionInfo.mcuBuildVersion), uint(versionInfo.driverMajorVersion), uint(versionInfo.driverMinorVersion), uint(versionInfo.driverBuildVersion), errorString(r)
}

// GetDeviceVfState returns cndevCardVfState_t.vfState, which is bitmap of mlu vf state. See cndev user doc for more details.
func (c *cndev) GetDeviceVfState(idx uint) (int, error) {
	if ret := dl.checkExist("cndevGetCardVfState"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var vfstate C.cndevCardVfState_t
	vfstate.version = C.CNDEV_VERSION_6
	r := C.cndevGetCardVfState(&vfstate, c.Load(idx))
	log.Debugf("slot:%d, vfstate is:%d", idx, int(vfstate.vfState))
	return int(vfstate.vfState), errorString(r)
}

func (c *cndev) GetDeviceVideoCodecUtil(idx uint) ([]int, []int, error) {
	if ret := dl.checkExist("cndevGetVideoCodecUtilization"); ret != C.CNDEV_SUCCESS {
		return nil, nil, errorString(ret)
	}

	var cardVideoCodecUtil C.cndevVideoCodecUtilization_t
	cardVideoCodecUtil.version = C.CNDEV_VERSION_6
	r := C.cndevGetVideoCodecUtilization(&cardVideoCodecUtil, c.Load(idx))
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
	if ret := dl.checkExist("cndevGetVoltageInfo"); ret != C.CNDEV_SUCCESS {
		return 0, 0, 0, errorString(ret)
	}

	var voltageInfo C.cndevVoltageInfo_t
	r := C.cndevGetVoltageInfo(&voltageInfo, c.Load(idx))
	return int(voltageInfo.ipuCoreVoltage), int(voltageInfo.socVoltage), int(voltageInfo.hbmVddVoltage), errorString(r)
}

func (c *cndev) GetDeviceFrequencyStatus(idx uint) (int, error) {
	if ret := dl.checkExist("cndevGetMLUFrequencyStatus"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

	var freqStatus C.cndevMLUFrequencyStatus_t
	r := C.cndevGetMLUFrequencyStatus(&freqStatus, c.Load(idx))
	return int(freqStatus.mluFrequencyLockStatus), errorString(r)
}

func (c *cndev) GetTopologyRelationship(domain, bus, device, function, rdmaDomain, rdmaBus, rdmaDevice, rdmaFunction uint) (int, error) {
	if ret := dl.checkExist("cndevGetNodeByBDF", "cndevTopologyGetRelationshipByNode"); ret != C.CNDEV_SUCCESS {
		return 0, errorString(ret)
	}

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
	if ret := dl.checkExist("cndevGetSupportedEventTypes"); ret != C.CNDEV_SUCCESS {
		return errorString(ret)
	}

	types := C.ulonglong(C.cndevEventTypeAll)
	r := C.cndevGetSupportedEventTypes(&types, c.Load(idx))

	return errorString(r)
}

func (c *cndev) RegisterEventsHandleAndWait(slots []int, ch chan XIDInfoWithTimestamp) error {
	if ret := dl.checkExist("cndevEventHandleCreate", "cndevRegisterEvents"); ret != C.CNDEV_SUCCESS {
		return errorString(ret)
	}

	var handle C.cndevEventHandle
	r := C.cndevEventHandleCreate(&handle)
	if err := errorString(r); err != nil {
		return err
	}

	for _, idx := range slots {
		r := C.cndevRegisterEvents(handle, C.cndevEventTypeAll, c.Load(uint(idx)))
		if err := errorString(r); err != nil {
			log.Errorf("register event error: %v for slot %d", err, idx)
		}
	}

	go waitEvents(handle, ch)
	return nil
}

func waitEvents(handle C.cndevEventHandle, ch chan XIDInfoWithTimestamp) {
	if ret := dl.checkExist("cndevEventWait"); ret != C.CNDEV_SUCCESS {
		return
	}

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
				XID:       fmt.Sprintf("0x%016x", eventData.eventData),
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
	if ret := dl.checkExist("cndevGetErrorString"); ret != C.CNDEV_SUCCESS {
		return fmt.Errorf("cndev: not support cndevGetErrorString")
	}

	if ret == C.CNDEV_SUCCESS {
		return nil
	}
	err := C.GoString(C.cndevGetErrorString(ret))
	return fmt.Errorf("cndev: %v", err)
}

func (c *cndev) GenerateDeviceHandleMap(count uint) error {
	if ret := dl.checkExist("cndevGetDeviceHandleByIndex"); ret != C.CNDEV_SUCCESS {
		return errorString(ret)
	}

	for i := uint(0); i < count; i++ {
		var handle C.cndevDevice_t
		r := C.cndevGetDeviceHandleByIndex(C.int(i), &handle)
		if errorString(r) != nil {
			return errorString(r)
		}
		cndevGlobalHandleMap.Store(uint(i), handle)
	}

	return nil
}

func healthErrorCodeToString(code C.cndevHealthError_t) string {
	switch code {
	case C.CNDEV_FR_OK:
		return "CNDEV_FR_OK"
	case C.CNDEV_FR_UNKNOWN:
		return "CNDEV_FR_UNKNOWN"
	case C.CNDEV_FR_INTERNAL:
		return "CNDEV_FR_INTERNAL"
	case C.CNDEV_FR_XID_ERROR:
		return "CNDEV_FR_XID_ERROR"
	case C.CNDEV_FR_PCIE_REPLAY_RATE:
		return "CNDEV_FR_PCIE_REPLAY_RATE"
	case C.CNDEV_FR_PCIE_GENERATION:
		return "CNDEV_FR_PCIE_REPLAY_RATE"
	case C.CNDEV_FR_PCIE_WIDTH:
		return "CNDEV_FR_PCIE_WIDTH"
	case C.CNDEV_FR_LOW_BANDWIDTH:
		return "CNDEV_FR_LOW_BANDWIDTH"
	case C.CNDEV_FR_HIGH_LATENCY:
		return "CNDEV_FR_HIGH_LATENCY"
	case C.CNDEV_FR_MLULINK_DOWN:
		return "CNDEV_FR_MLULINK_DOWN"
	case C.CNDEV_FR_MLULINK_REPLAY:
		return "CNDEV_FR_MLULINK_REPLAY"
	case C.CNDEV_FR_MLULINK_ERROR_THRESHOLD:
		return "CNDEV_FR_MLULINK_ERROR_THRESHOLD"
	case C.CNDEV_FR_MLULINK_CRC_ERROR_THRESHOLD:
		return "CNDEV_FR_MLULINK_CRC_ERROR_THRESHOLD"
	case C.CNDEV_FR_DIE2DIE_CRC_ERROR_THRESHOLD:
		return "CNDEV_FR_DIE2DIE_CRC_ERROR_THRESHOLD"
	case C.CNDEV_FR_SBE_VIOLATION:
		return "CNDEV_FR_SBE_VIOLATION"
	case C.CNDEV_FR_DBE_VIOLATION:
		return "CNDEV_FR_DBE_VIOLATION"
	case C.CNDEV_FR_SBE_THRESHOLD_VIOLATION:
		return "CNDEV_FR_SBE_THRESHOLD_VIOLATION"
	case C.CNDEV_FR_DBE_THRESHOLD_VIOLATION:
		return "CNDEV_FR_DBE_THRESHOLD_VIOLATION"
	case C.CNDEV_FR_ROW_REMAP_FAILURE:
		return "CNDEV_FR_ROW_REMAP_FAILURE"
	case C.CNDEV_FR_PENDING_ROW_REMAP:
		return "CNDEV_FR_PENDING_ROW_REMAP"
	case C.CNDEV_FR_ADDRESSSWAP_FAILURE:
		return "CNDEV_FR_ADDRESSSWAP_FAILURE"
	case C.CNDEV_FR_PENDING_ADDRESSSWAP:
		return "CNDEV_FR_PENDING_ADDRESSSWAP"
	case C.CNDEV_FR_REPAIR_FAILURE:
		return "CNDEV_FR_REPAIR_FAILURE"
	case C.CNDEV_FR_PENDING_PAGE_REPAIR:
		return "CNDEV_FR_PENDING_PAGE_REPAIR"
	case C.CNDEV_FR_PENDING_PAGE_RETIREMENTS:
		return "CNDEV_FR_PENDING_PAGE_RETIREMENTS"
	case C.CNDEV_FR_PAGE_REPAIR_RESOURCE_THRESHOLD:
		return "CNDEV_FR_PAGE_REPAIR_RESOURCE_THRESHOLD"
	case C.CNDEV_FR_TEMP_VIOLATION:
		return "CNDEV_FR_TEMP_VIOLATION"
	case C.CNDEV_FR_THERMAL_VIOLATION:
		return "CNDEV_FR_THERMAL_VIOLATION"
	case C.CNDEV_FR_DEVICE_COUNT_MISMATCH:
		return "CNDEV_FR_DEVICE_COUNT_MISMATCH"
	case C.CNDEV_FR_DRIVER_ERROR:
		return "CNDEV_FR_DRIVER_ERROR"
	case C.CNDEV_FR_DEVSYS_ERROR:
		return "CNDEV_FR_DEVSYS_ERROR"
	default:
		return fmt.Sprintf("CNDEV_FR_UNKNOWN_CODE_%d", code)
	}
}
