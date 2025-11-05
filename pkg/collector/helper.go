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

package collector

import (
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/podresources"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type chassisInfo struct {
	chassisSN          uint64
	chassisProductDate string
	chassisProductName string
	chassisVendorName  string
	chassisPartNumber  string
	chassisBmcIP       string
	nvmeNum            int
	ibNum              int
	psuNum             int
}

type pcieInfo struct {
	pcieSlot            string
	pcieSubsystem       string
	pcieDeviceID        string
	pcieVendor          string
	pcieSubsystemVendor string
	pcieID              string
	pcieSpeed           string
	pcieWidth           string
}

type mluLinkRemoteInfo struct {
	remoteMcSn      string
	remoteBaSn      string
	remoteSlotID    string
	remotePortID    string
	remoteNcsUUID64 string
	remoteIP        string
	remoteMac       string
	remoteUUID      string
	remotePortName  string
}

type labelInfo struct {
	chassisDev        int
	chassisDevInfo    cndev.ChassisDevInfo
	chassisInfo       chassisInfo
	cluster           string
	cndevVersion      string
	computeCapability string
	core              int
	cpuCore           string
	host              string
	hostIP            string
	rdmaDevice        rdmaDevice
	lane              int
	link              int
	linkVersion       string
	memoryDie         string
	mimInstanceInfo   cndev.MimInstanceInfo
	mimProfileInfo    cndev.MimProfileInfo
	mluLinkRemoteInfo mluLinkRemoteInfo
	moduleID          uint16
	pcieInfo          pcieInfo
	pid               uint32
	ppi               string
	podInfo           podresources.PodInfo
	smluInstanceInfo  cndev.SmluInstanceInfo
	smluProfileInfo   cndev.SmluProfileInfo
	stat              MLUStat
	typ               string
	unhealthReason    string
	vf                string
	xidInfo           cndev.XIDInfoWithTimestamp
}

func getLabelValues(labels []string, info labelInfo) []string {
	values := []string{}
	for _, l := range labels {
		switch l {
		case MLU:
			values = append(values, fmt.Sprintf("%d", info.stat.slot))
		case Model:
			values = append(values, info.stat.model)
		case Type:
			var typ string
			// suffix needs to be stripped except for "mlu270-x5k"
			if strings.EqualFold(info.stat.model, "mlu270-x5k") {
				typ = strings.ToLower(info.stat.model)
			} else {
				typ = strings.ToLower(strings.Split(info.stat.model, "-")[0])
			}
			if info.stat.mimEnabled && info.vf != "" {
				index, err := strconv.Atoi(info.vf)
				if err != nil {
					log.Warnf("convert vf %s with err %v", info.vf, err)
					break
				}
				for _, inf := range info.stat.mimInfos {
					if inf.InstanceInfo.InstanceID == index {
						typ = typ + ".mim-" + inf.ProfileInfo.ProfileName
						break
					}
				}
			}
			if info.stat.smluEnabled && info.vf != "" {
				index, err := strconv.Atoi(info.vf)
				if err != nil {
					log.Warnf("Convert vf %s with err %v", info.vf, err)
					break
				}
				smluInfos := info.stat.smluInfos
				if len(info.stat.smluInfos) == 0 {
					cli := cndev.NewCndevClient()
					if err := cli.Init(true); err != nil {
						log.Error(errors.Wrap(err, "Init"))
						break
					}
					infs, err := cli.GetAllSMluInfo(info.stat.slot)
					if err != nil {
						log.Warnf("Failed to get smlu info %s with err %v", info.vf, err)
						break
					}
					smluInfos = infs
				}
				if info.typ != "" {
					typ = typ + ".smlu." + info.typ
				} else {
					for _, inf := range smluInfos {
						if inf.InstanceInfo.InstanceID == index {
							typ = typ + ".smlu-" + inf.ProfileInfo.ProfileName
							break
						}
					}
				}
			}
			values = append(values, typ)
		case SN:
			values = append(values, info.stat.sn)
		case UUID:
			if info.stat.mimEnabled && info.vf != "" {
				index, err := strconv.Atoi(info.vf)
				if err != nil {
					log.Warnf("convert vf %s with err %v", info.vf, err)
					break
				}
				for _, inf := range info.stat.mimInfos {
					if inf.InstanceInfo.InstanceID == index {
						values = append(values, inf.InstanceInfo.UUID)
						break
					}
				}
				break
			}
			if info.stat.smluEnabled && info.vf != "" {
				index, err := strconv.Atoi(info.vf)
				if err != nil {
					log.Warnf("convert vf %s with err %v", info.vf, err)
					break
				}
				smluInfos := info.stat.smluInfos
				if len(info.stat.smluInfos) == 0 {
					cli := cndev.NewCndevClient()
					if err := cli.Init(true); err != nil {
						log.Error(errors.Wrap(err, "Init"))
						break
					}
					infs, err := cli.GetAllSMluInfo(info.stat.slot)
					if err != nil {
						log.Warnf("Failed to get smlu info %s with err %v", info.vf, err)
						break
					}
					smluInfos = infs
				}
				for _, inf := range smluInfos {
					if inf.InstanceInfo.InstanceID == index {
						values = append(values, inf.InstanceInfo.UUID)
						break
					}
				}
				break
			}
			values = append(values, info.stat.uuid)
		case Node:
			values = append(values, info.host)
		case NodeIP:
			values = append(values, info.hostIP)
		case Cluster:
			values = append(values, info.cluster)
		case CndevVersion:
			values = append(values, info.cndevVersion)
		case ComputeCapability:
			values = append(values, info.computeCapability)
		case Core:
			values = append(values, fmt.Sprintf("%d", info.core))
		case CPUCore:
			values = append(values, info.cpuCore)
		case MimGDMACount:
			values = append(values, fmt.Sprintf("%d", info.mimProfileInfo.GDMACount))
		case MimInstanceName:
			values = append(values, info.mimInstanceInfo.InstanceName)
		case MimInstanceID:
			values = append(values, fmt.Sprintf("%d", info.mimInstanceInfo.InstanceID))
		case MimJPUCount:
			values = append(values, fmt.Sprintf("%d", info.mimProfileInfo.JPUCount))
		case MimMemorySize:
			values = append(values, fmt.Sprintf("%d", info.mimProfileInfo.MemorySize))
		case MimMLUCoreCount:
			values = append(values, fmt.Sprintf("%d", info.mimProfileInfo.MLUCoreCount))
		case MimProfileID:
			values = append(values, fmt.Sprintf("%d", info.mimProfileInfo.ProfileID))
		case MimProfileName:
			values = append(values, info.mimProfileInfo.ProfileName)
		case MimUUID:
			values = append(values, info.mimInstanceInfo.UUID)
		case MimVPUCount:
			values = append(values, fmt.Sprintf("%d", info.mimProfileInfo.VPUCount))
		case MCU:
			values = append(values, info.stat.mcu)
		case Driver:
			values = append(values, info.stat.driver)
		case Namespace:
			values = append(values, info.podInfo.Namespace)
		case Pod:
			values = append(values, info.podInfo.Pod)
		case PPI:
			values = append(values, info.ppi)
		case Container:
			values = append(values, info.podInfo.Container)
		case VF:
			values = append(values, info.vf)
		case RDMADeviceName:
			values = append(values, info.rdmaDevice.name)
		case RDMADevicePCIeAddress:
			values = append(values, info.rdmaDevice.pcieAddress)
		case RDMADevicePCIeDomain:
			values = append(values, fmt.Sprintf("%d", info.rdmaDevice.domain))
		case RDMADevicePCIeBus:
			values = append(values, fmt.Sprintf("%d", info.rdmaDevice.bus))
		case RDMADevicePCIeDevice:
			values = append(values, fmt.Sprintf("%d", info.rdmaDevice.device))
		case RDMADevicePCIeFunction:
			values = append(values, fmt.Sprintf("%d", info.rdmaDevice.function))
		case RDMADeviceNicName:
			values = append(values, info.rdmaDevice.nicName)
		case RDMADeviceIPAddress:
			values = append(values, info.rdmaDevice.ipAddress)
		case Lane:
			values = append(values, fmt.Sprintf("%d", info.lane))
		case Link:
			values = append(values, fmt.Sprintf("%d", info.link))
		case LinkVersion:
			values = append(values, info.linkVersion)
		case ModuleID:
			values = append(values, fmt.Sprintf("%d", info.moduleID))
		case PCIeSlot:
			values = append(values, info.pcieInfo.pcieSlot)
		case PCIeSubsystem:
			values = append(values, info.pcieInfo.pcieSubsystem)
		case PCIeDeviceID:
			values = append(values, info.pcieInfo.pcieDeviceID)
		case PCIeVendor:
			values = append(values, info.pcieInfo.pcieVendor)
		case PCIeSubsystemVendor:
			values = append(values, info.pcieInfo.pcieSubsystemVendor)
		case PCIeID:
			values = append(values, info.pcieInfo.pcieID)
		case PCIeSpeed:
			values = append(values, info.pcieInfo.pcieSpeed)
		case PCIeWidth:
			values = append(values, info.pcieInfo.pcieWidth)
		case Pid:
			values = append(values, fmt.Sprintf("%d", info.pid))
		case MemoryDie:
			values = append(values, info.memoryDie)
		case SmluInstanceID:
			values = append(values, fmt.Sprintf("%d", info.smluInstanceInfo.InstanceID))
		case SmluInstanceName:
			values = append(values, info.smluInstanceInfo.InstanceName)
		case SmluIpuTotal:
			values = append(values, fmt.Sprintf("%d", info.smluProfileInfo.IpuTotal))
		case SmluMemTotal:
			values = append(values, fmt.Sprintf("%d", info.smluProfileInfo.MemTotal))
		case SmluProfileID:
			values = append(values, fmt.Sprintf("%d", info.smluProfileInfo.ProfileID))
		case SmluProfileName:
			values = append(values, info.smluProfileInfo.ProfileName)
		case SmluUUID:
			values = append(values, info.smluInstanceInfo.UUID)
		case XID:
			values = append(values, info.xidInfo.XID)
		case XIDBase10:
			values = append(values, fmt.Sprintf("%d", info.xidInfo.XIDBase10))
		case XIDTimestamp:
			values = append(values, fmt.Sprintf("%d", info.xidInfo.Timestamp))
		case XIDComputeInstanceID:
			values = append(values, fmt.Sprintf("%d", info.xidInfo.ComputeInstanceID))
		case XIDMLUInstanceID:
			values = append(values, fmt.Sprintf("%d", info.xidInfo.MLUInstanceID))
		case RemoteMcSn:
			values = append(values, info.mluLinkRemoteInfo.remoteMcSn)
		case RemoteBaSn:
			values = append(values, info.mluLinkRemoteInfo.remoteBaSn)
		case RemoteSlotID:
			values = append(values, info.mluLinkRemoteInfo.remoteSlotID)
		case RemotePortID:
			values = append(values, info.mluLinkRemoteInfo.remotePortID)
		case RemoteNcsUUID64:
			values = append(values, info.mluLinkRemoteInfo.remoteNcsUUID64)
		case RemoteIP:
			values = append(values, info.mluLinkRemoteInfo.remoteIP)
		case RemoteMac:
			values = append(values, info.mluLinkRemoteInfo.remoteMac)
		case RemoteUUID:
			values = append(values, info.mluLinkRemoteInfo.remoteUUID)
		case RemotePortName:
			values = append(values, info.mluLinkRemoteInfo.remotePortName)
		case UnhealthReason:
			values = append(values, info.unhealthReason)
		case ChassisDev:
			values = append(values, fmt.Sprintf("%d", info.chassisDev))
		case ChassisSN:
			values = append(values, fmt.Sprintf("%d", info.chassisInfo.chassisSN))
		case ChassisProductDate:
			values = append(values, info.chassisInfo.chassisProductDate)
		case ChassisProductName:
			values = append(values, info.chassisInfo.chassisProductName)
		case ChassisVendorName:
			values = append(values, info.chassisInfo.chassisVendorName)
		case ChassisPartNumber:
			values = append(values, info.chassisInfo.chassisPartNumber)
		case ChassisBmcIP:
			values = append(values, info.chassisInfo.chassisBmcIP)
		case ChassisNvmeNum:
			values = append(values, fmt.Sprintf("%d", info.chassisInfo.nvmeNum))
		case ChassisIbNum:
			values = append(values, fmt.Sprintf("%d", info.chassisInfo.ibNum))
		case ChassisPsuNum:
			values = append(values, fmt.Sprintf("%d", info.chassisInfo.psuNum))
		case ChassisDevSN:
			values = append(values, info.chassisDevInfo.SN)
		case ChassisDevModel:
			values = append(values, info.chassisDevInfo.Model)
		case ChassisDevFw:
			values = append(values, info.chassisDevInfo.Fw)
		case ChassisDevMfc:
			values = append(values, info.chassisDevInfo.Mfc)
		default:
			values = append(values, "") // configured label not applicable, add this to prevent panic
		}
	}
	log.Debugf("GetLabelValues: %+v", values)
	return values
}

type MLUStat struct {
	cndevInterfaceDisabled   map[string]bool
	mlulinkInterfaceDisabled map[int]map[string]bool
	driver                   string
	link                     int
	linkPPI                  map[int]string
	mcu                      string
	mimEnabled               bool
	mimInfos                 []cndev.MimInfo
	model                    string
	opticalPresent           map[int]uint8
	slot                     uint
	smluEnabled              bool
	smluInfos                []cndev.SmluInfo
	sn                       string
	uuid                     string
}

type MLUStatMap struct {
	StatMap   sync.Map
	InProblem atomic.Bool
}

func (m *MLUStatMap) Range(f func(key string, value MLUStat) bool) {
	m.StatMap.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(MLUStat))
	})
}

func (m *MLUStatMap) Load(key string) MLUStat {
	value, ok := m.StatMap.Load(key)
	if !ok {
		return MLUStat{}
	}
	return value.(MLUStat)
}

func collectMLUInfo(mluInfo *MLUStatMap, cli cndev.Cndev, count uint) {
	mluInfo.InProblem.Store(false)

	for i := uint(0); i < count; i++ {
		dis := make(map[string]bool)

		log.Debugf("Start slot %d GetDeviceModel", i)
		model := cli.GetDeviceModel(i)
		if model == "" {
			log.Warnf("GetDeviceModel for slot %d model is empty", i)
			mluInfo.InProblem.Store(true)
		}

		log.Debugf("Start slot %d GetDeviceUUID", i)
		uuid, err := cli.GetDeviceUUID(i)
		if err != nil || uuid == "" {
			log.Warn(errors.Wrapf(err, "GetDeviceUUID for slot %d with err or uuid is empty", i))
			mluInfo.InProblem.Store(true)
		}

		log.Debugf("Start slot %d GetDeviceSN", i)
		sn, err := cli.GetDeviceSN(i)
		if err != nil || sn == "" {
			log.Warn(errors.Wrapf(err, "GetDeviceSN for slot %d with err or sn is empty", i))
			mluInfo.InProblem.Store(true)
		}

		log.Debugf("Start slot %d GetDeviceVersion", i)
		mcuMajor, mcuMinor, mcuBuild, driverMajor, driverMinor, driverBuild, err := cli.GetDeviceVersion(i)
		if err != nil {
			log.Warn(errors.Wrapf(err, "GetDeviceVersion for slot %d", i))
			mluInfo.InProblem.Store(true)
		}

		log.Debugf("Start slot %d DeviceMimModeEnabled", i)
		var mimInfos []cndev.MimInfo
		mimEnabled, err := cli.DeviceMimModeEnabled(i)
		if err != nil {
			log.Warn(errors.Wrapf(err, "DeviceMimModeEnabled for slot %d", i))
		} else if mimEnabled {
			mimInfos, err = cli.GetAllMLUInstanceInfo(i)
			if err != nil {
				log.Warn(errors.Wrap(err, "GetAllMLUInstanceInfo"))
			}
		}

		log.Debugf("Start slot %d DeviceSmluModeEnabled", i)
		smluEnabled, err := cli.DeviceSmluModeEnabled(i)
		if err != nil {
			log.Warn(errors.Wrapf(err, "DeviceSmluModeEnabled for slot %d", i))
		}

		log.Debugf("Start slot %d GetDeviceMLULinkPortNumber", i)
		link := cli.GetDeviceMLULinkPortNumber(i)
		if link == -1 {
			log.Warnf("GetDeviceMLULinkPortNumber for slot %d", i)
			dis["mluLinkPortNumberDisabled"] = true
		}
		log.Debugf("Slot %d mlulink num %d", i, link)
		linkPPI := map[int]string{}
		opticalPresent := map[int]uint8{}
		mlulinkDis := make(map[int]map[string]bool)
		for j := 0; j < link; j++ {
			mlulinkDis[j] = make(map[string]bool)
			log.Debugf("Start slot %d link %d GetDeviceMLULinkStatus", i, j)
			if _, _, _, err = cli.GetDeviceMLULinkStatus(i, uint(j)); err != nil {
				log.Debug(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkStatus", i, j))
				mlulinkDis[j]["mluLinkStatusDisabled"] = true
			}

			log.Debugf("Start slot %d link %d GetDeviceMLULinkState", i, j)
			if _, _, err = cli.GetDeviceMLULinkState(i, uint(j)); err != nil {
				log.Debug(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkState", i, j))
				mlulinkDis[j]["mluLinkStateDisabled"] = true
			}

			log.Debugf("Start slot %d link %d GetDeviceMLULinkCapability", i, j)
			if _, _, err = cli.GetDeviceMLULinkCapability(i, uint(j)); err != nil {
				log.Debug(errors.Wrapf(err, "Slot %d GetDeviceMLULinkCapability", i))
				mlulinkDis[j]["mluLinkCapabilityDisabled"] = true
			}

			log.Debugf("Start slot %d link %d GetDeviceMLULinkPortMode", i, j)
			if _, err = cli.GetDeviceMLULinkPortMode(i, uint(j)); err != nil {
				log.Debug(errors.Wrapf(err, "Slot %d GetDeviceMLULinkPortMode", i))
				mlulinkDis[j]["mluLinkPortModeDisabled"] = true
			}

			log.Debugf("Start slot %d link %d GetDeviceMLULinkVersion", i, j)
			if _, _, _, err = cli.GetDeviceMLULinkVersion(i, uint(j)); err != nil {
				log.Debug(errors.Wrapf(err, "Slot %d GetDeviceMLULinkVersion", i))
				mlulinkDis[j]["mluLinkVersionDisabled"] = true
			}

			log.Debugf("Start slot %d link %d GetDeviceMLULinkEventCounter", i, j)
			if _, err = cli.GetDeviceMLULinkEventCounter(i, uint(j)); err != nil {
				log.Debug(errors.Wrapf(err, "Slot %d GetDeviceMLULinkEventCounter", i))
				mlulinkDis[j]["mluLinkEventCounterDisabled"] = true
			}

			log.Debugf("Start slot %d link %d GetDeviceMLULinkErrorCounter", i, j)
			if _, _, _, err = cli.GetDeviceMLULinkErrorCounter(i, uint(j)); err != nil {
				log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMLULinkErrorCounter", i))
				mlulinkDis[j]["mluLinkErrorCounterDisabled"] = true
			}

			log.Debugf("Start slot %d link %d GetDeviceMLULinkCounter", i, j)
			if _, _, _, _, _, _, _, _, _, _, _, _, _, err = cli.GetDeviceMLULinkCounter(i, uint(j)); err != nil {
				log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMLULinkCounter", i))
				mlulinkDis[j]["mluLinkCounterDisabled"] = true
			}

			log.Debugf("Start slot %d link %d GetDeviceMLULinkRemoteInfo", i, j)
			if _, _, _, _, _, _, _, _, _, _, err = cli.GetDeviceMLULinkRemoteInfo(i, uint(j)); err != nil {
				log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMLULinkErrorCounter", i))
				mlulinkDis[j]["mluLinkRemoteInfoDisabled"] = true
			}

			log.Debugf("Start slot %d link %d GetDeviceMLULinkSpeedInfo", i, j)
			if _, _, err = cli.GetDeviceMLULinkSpeedInfo(i, uint(j)); err != nil {
				log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMLULinkSpeedInfo", i))
				mlulinkDis[j]["mluLinkSpeedInfoDisabled"] = true
			}

			var ppi string
			log.Debugf("Start slot %d link %d GetDeviceMLULinkPPI", i, j)
			if ppi, err = cli.GetDeviceMLULinkPPI(i, uint(j)); err != nil {
				log.Warn(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkPPI", i, j))
				continue
			}
			if ppi != "N/A" {
				linkPPI[j] = ppi
				log.Debugf("Start slot %d link %d GetDeviceOpticalInfo", i, j)
				present, _, _, _, _, err := cli.GetDeviceOpticalInfo(i, uint(j))
				if err != nil {
					log.Warn(errors.Wrapf(err, "Slot %d link %d GetDeviceOpticalInfo", i, j))
					continue
				}
				opticalPresent[j] = present
			}
		}

		log.Debugf("Start slot %d GetDeviceCRCInfo", i)
		if _, _, err = cli.GetDeviceCRCInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceCRCInfo", i))
			dis["crcDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceECCInfo", i)
		if _, _, _, _, _, _, _, _, err = cli.GetDeviceECCInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", i))
			dis["eccDisabled"] = true
		}

		log.Debugf("Start slot %d cndevGetProcessInfo", i)
		if _, _, _, err = cli.GetDeviceProcessInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceProcessInfo", i))
			dis["processInfoDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceProcessUtil", i)
		if _, _, _, _, _, _, err = cli.GetDeviceProcessUtil(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", i))
			dis["processUtilDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceRetiredPageInfo", i)
		if _, _, err = cli.GetDeviceRetiredPageInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", i))
			dis["retiredPageDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceRemappedRows", i)
		if _, _, _, _, _, err = cli.GetDeviceRemappedRows(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceRemappedRows", i))
			dis["remappedRowsDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceHeartbeatCount", i)
		if _, err = cli.GetDeviceHeartbeatCount(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceHeartbeatCount", i))
			dis["heartBeatDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceMemEccCounter", i)
		if _, _, err = cli.GetDeviceMemEccCounter(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", i))
			dis["memEccDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceAddressSwaps", i)
		if _, _, _, _, _, err = cli.GetDeviceAddressSwaps(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceAddressSwaps", i))
			dis["addressSwapsDisabled"] = true
		}

		log.Debugf("Start slot %d GetSupportedEventTypes", i)
		if err = cli.GetSupportedEventTypes(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetSupportedEventTypes", i))
			dis["xidCallbackDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceCurrentInfo", i)
		if _, _, _, err = cli.GetDeviceCurrentInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceCurrentInfo", i))
			dis["currentInfoDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceVoltageInfo", i)
		if _, _, _, err = cli.GetDeviceVoltageInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceVoltageInfo", i))
			dis["voltageInfoDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceFrequencyStatus", i)
		if _, err = cli.GetDeviceFrequencyStatus(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceFrequencyStatus", i))
			dis["frequencyStatusDisabled"] = true
		}

		log.Debugf("Start slot %d GetDevicePerformanceThrottleReason", i)
		if _, err = cli.GetDevicePerformanceThrottleReason(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDevicePerformanceThrottleReason", i))
			dis["performanceThrottleDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceOverTemperatureInfo", i)
		if _, _, err = cli.GetDeviceOverTemperatureInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceOverTemperatureInfo", i))
			dis["overTemperatureInfoDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceOverTemperatureShutdownThreshold", i)
		if _, err = cli.GetDeviceOverTemperatureShutdownThreshold(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceOverTemperatureThreshold", i))
			dis["overTemperatureThresholdDisabled"] = true
		}

		log.Debugf("Start slot %d GetDevicePowerManagementLimitation", i)
		if _, err = cli.GetDevicePowerManagementLimitation(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDevicePowerManagementLimitation", i))
			dis["powerCurrentLimitDisabled"] = true
		}

		log.Debugf("Start slot %d GetDevicePowerManagementDefaultLimitation", i)
		if _, err = cli.GetDevicePowerManagementDefaultLimitation(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDevicePowerManagementDefaultLimitation", i))
			dis["powerDefaultLimitDisabled"] = true
		}

		log.Debugf("Start slot %d GetDevicePowerManagementLimitRange", i)
		if _, _, err = cli.GetDevicePowerManagementLimitRange(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDevicePowerManagementLimitRange", i))
			dis["powerLimitRangeDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceBAR4MemoryInfo", i)
		if _, _, _, err = cli.GetDeviceBAR4MemoryInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceBAR4MemoryInfo", i))
			dis["bar4MemoryInfoDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceCPUUtil", i)
		if _, _, _, _, _, _, _, err = cli.GetDeviceCPUUtil(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceCPUUtil", i))
			dis["cpuUtilDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceMaxPCIeInfo", i)
		if _, _, err = cli.GetDeviceMaxPCIeInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMaxPCIeInfo", i))
			dis["maxPCIeInfoDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceEccMode", i)
		if _, _, err = cli.GetDeviceEccMode(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceEccMode", i))
			dis["eccModeDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceRetiredPagesOperation", i)
		if _, err = cli.GetDeviceRetiredPagesOperation(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceRetiredPagesOperation", i))
			dis["retiredPagesOperationDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceComputeCapability", i)
		if _, _, err = cli.GetDeviceComputeCapability(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceComputeCapability", i))
			dis["computeCapabilityDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceTensorUtil", i)
		if _, err = cli.GetDeviceTensorUtil(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceTensorUtil", i))
			dis["tensorUtilDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceComputeMode", i)
		if _, err = cli.GetDeviceComputeMode(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceComputeMode", i))
			dis["computeModeDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceTemperature", i)
		if _, _, _, _, _, err = cli.GetDeviceTemperature(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceTemperature", i))
			dis["temperatureDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceUtil", i)
		if _, _, err = cli.GetDeviceUtil(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceUtil", i))
			dis["deviceUtilDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceDDRInfo", i)
		if _, _, err = cli.GetDeviceDDRInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceDDRInfo", i))
			dis["ddrInfoDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceFrequency", i)
		if _, _, _, err = cli.GetDeviceFrequency(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceFrequency", i))
			dis["frequencyDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceOsMemory", i)
		if _, _, err = cli.GetDeviceOsMemory(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceOsMemory", i))
			dis["osMemoryDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceHealth", i)
		if _, _, _, _, err = cli.GetDeviceHealth(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceHealth", i))
			dis["healthDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceFanSpeed", i)
		if _, err = cli.GetDeviceFanSpeed(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceFanSpeed", i))
			dis["fanSpeedDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceImageCodecUtil", i)
		if _, err = cli.GetDeviceImageCodecUtil(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceImageCodecUtil", i))
			dis["imageCodecUtilDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceMemory", i)
		if _, _, _, _, _, err = cli.GetDeviceMemory(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMemory", i))
			dis["deviceMemoryDisabled"] = true
		}

		log.Debugf("Start slot %d GetDevicePCIeInfo", i)
		if _, _, _, _, _, _, _, _, _, _, err = cli.GetDevicePCIeInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDevicePCIeInfo", i))
			dis["pcieInfoDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceNUMANodeID", i)
		if _, err = cli.GetDeviceNUMANodeID(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceNUMANodeID", i))
			dis["numaDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceParityError", i)
		if _, err = cli.GetDeviceParityError(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceParityError", i))
			dis["parityErrorDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceCurrentPCIeInfo", i)
		if _, _, err = cli.GetDeviceCurrentPCIeInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceCurrentPCIeInfo", i))
			dis["currentPCIeInfoDisabled"] = true
		}

		log.Debugf("Start slot %d GetDevicePCIeThroughput", i)
		if _, _, err = cli.GetDevicePCIeThroughput(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDevicePCIeThroughput", i))
			dis["pcieThroughputDisabled"] = true
		}

		log.Debugf("Start slot %d GetDevicePCIeReplayCount", i)
		if _, err = cli.GetDevicePCIeReplayCount(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDevicePCIeReplayCount", i))
			dis["pcieRelayCountDisabled"] = true
		}

		log.Debugf("Start slot %d GetDevicePower", i)
		if _, err = cli.GetDevicePower(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDevicePower", i))
			dis["devicePowerDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceTinyCoreUtil", i)
		if _, err = cli.GetDeviceTinyCoreUtil(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceTinyCoreUtil", i))
			dis["tinyCoreDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceCndevVersion", i)
		if _, _, _, err = cli.GetDeviceCndevVersion(); err != nil {
			log.Warn(errors.Wrap(err, "GetDeviceCndevVersion"))
			dis["cndevVersionDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceVideoCodecUtil", i)
		if _, _, err = cli.GetDeviceVideoCodecUtil(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceVideoCodecUtil", i))
			dis["videoCodecUtilDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceVfState", i)
		if _, err = cli.GetDeviceVfState(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceVfState", i))
			dis["vfStateDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceRepairStatus", i)
		if _, _, _, err = cli.GetDeviceRepairStatus(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceRepairStatus", i))
			dis["repairStatusDisabled"] = true
		}

		log.Debugf("Start slot %d GetDeviceChassisInfo", i)
		if _, _, _, _, _, _, _, _, _, err = cli.GetDeviceChassisInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceChassisInfo", i))
			dis["chassisInfo"] = true
		}

		log.Debugf("Start slot %d GetDeviceActivity", i)
		if _, err = cli.GetDeviceActivity(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceActivity", i))
			dis["activityDisabled"] = true
		}

		mluInfo.StatMap.Store(uuid, MLUStat{
			cndevInterfaceDisabled:   dis,
			mlulinkInterfaceDisabled: mlulinkDis,
			driver:                   calcVersion(driverMajor, driverMinor, driverBuild),
			link:                     link,
			linkPPI:                  linkPPI,
			mcu:                      calcVersion(mcuMajor, mcuMinor, mcuBuild),
			mimEnabled:               mimEnabled,
			mimInfos:                 mimInfos,
			model:                    model,
			opticalPresent:           opticalPresent,
			slot:                     i,
			smluEnabled:              smluEnabled,
			sn:                       sn,
			uuid:                     uuid,
		})

		if (len(dis) != 0 || len(mlulinkDis) != 0) && i == 0 {
			log.Warnf("Some cndev interfaces are not supported: %v, %v", dis, mlulinkDis)
		}
	}
}

func calcVersion(major uint, minor uint, build uint) string {
	return fmt.Sprintf("v%d.%d.%d", major, minor, build)
}

func getRDMAPCIeAddress(deviceName string) string {
	ueventPath := fmt.Sprintf("/sys/class/infiniband/%s/device/uevent", deviceName)
	data, err := os.ReadFile(ueventPath)
	if err != nil {
		log.Errorf("Error reading uevent file for %s: %v", deviceName, err)
		return ""
	}
	// eg: PCI_SLOT_NAME=0000:21:00.0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "PCI_SLOT_NAME=") {
			return strings.TrimPrefix(line, "PCI_SLOT_NAME=")
		}
	}
	return ""
}

func parseRDMAPCIeAddress(pcieAddress string) (uint, uint, uint, uint, error) {
	// PCIe address format: 0000:XX:YY.Z
	parts := strings.Split(pcieAddress, ":")
	if len(parts) != 3 {
		return 0, 0, 0, 0, fmt.Errorf("invalid PCIe address format")
	}
	domain, err := strconv.ParseInt(parts[0], 16, 0)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid domain part")
	}
	bus, err := strconv.ParseInt(parts[1], 16, 0)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid bus part")
	}
	subParts := strings.Split(parts[2], ".")
	if len(subParts) != 2 {
		return 0, 0, 0, 0, fmt.Errorf("invalid PCIe address format")
	}
	device, err := strconv.ParseInt(subParts[0], 16, 0)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	function, err := strconv.ParseInt(subParts[1], 16, 0)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid function part")
	}
	return uint(domain), uint(bus), uint(device), uint(function), nil
}

func getRDMAPCIeInfo() []rdmaDevice {
	files, err := os.ReadDir("/sys/class/infiniband/")
	if err != nil {
		log.Warnf("Error reading /sys/class/infiniband/: %v", err)
		return nil
	}
	var devices []rdmaDevice
	for _, file := range files {
		if file.Type()&fs.ModeSymlink != 0 {
			path := filepath.Join("/sys/class/infiniband", file.Name(), "device/net")
			f, err := os.ReadDir(path)
			var nicName string
			if err == nil && len(f) > 0 {
				nicName = f[0].Name()
			} else {
				log.Warnf("Error reading /sys/class/infiniband/%s/device/net: %v", file.Name(), err)
				continue
			}

			pcieAddress := getRDMAPCIeAddress(file.Name())
			if pcieAddress == "" {
				log.Warnf("No RDMA PCIe address found for device %s", file.Name())
				continue
			}

			domain, bus, device, function, err := parseRDMAPCIeAddress(pcieAddress)
			if err != nil {
				log.Warnf("Error parsing RDMA PCIe address for %s: %v", file.Name(), err)
				continue
			}

			ipAddress := getIPAddressByNICName(nicName)

			devices = append(devices, rdmaDevice{
				name:        file.Name(),
				pcieAddress: pcieAddress,
				domain:      domain,
				bus:         bus,
				device:      device,
				function:    function,
				nicName:     nicName,
				ipAddress:   ipAddress,
			})
		}
	}

	if len(devices) == 0 {
		return nil
	}

	sort.Slice(devices, func(i, j int) bool {
		if devices[i].domain != devices[j].domain {
			return devices[i].domain < devices[j].domain
		}
		if devices[i].bus != devices[j].bus {
			return devices[i].bus < devices[j].bus
		}
		if devices[i].device != devices[j].device {
			return devices[i].device < devices[j].device
		}
		return devices[i].function < devices[j].function
	})

	return devices
}

func getIPAddressByNICName(nicName string) string {
	var ipAddress string
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Warnf("Error getting network interfaces: %v", err)
		return ipAddress
	}

	for _, iface := range interfaces {
		if iface.Name == nicName {
			addrs, err := iface.Addrs()
			if err != nil {
				log.Warnf("Error getting IP addresses for interface %s: %v", nicName, err)
				return ipAddress
			}
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if v.IP.To4() != nil {
						ipAddress = v.IP.String()
						break
					}
				}
			}
			break
		}
	}
	return ipAddress
}

func fetchMLUCounts() (uint, error) {
	targetVendorID := uint16(0xcabc) // cambricon mlu vendor ID is 0xcabc
	targetClassBase := uint8(0x12)   // cambricon mlu class code base is 0x12
	pciDevicesPath := "/sys/bus/pci/devices"

	readHexFile := func(path string) (uint64, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return 0, err
		}
		s := strings.TrimSpace(string(data))
		s = strings.TrimPrefix(s, "0x")
		val, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return 0, err
		}
		return val, nil
	}

	entries, err := os.ReadDir(pciDevicesPath)
	if err != nil {
		log.Errorf("Can't read pci dir: %v", err)
		return 0, err
	}

	var count uint
	for _, entry := range entries {
		devicePath := filepath.Join(pciDevicesPath, entry.Name())
		if _, err := os.Stat(filepath.Join(devicePath, "physfn")); err == nil {
			log.Debugf("Skip SR-IOV VF device: %s", devicePath)
			continue
		}
		vendorID, err := readHexFile(filepath.Join(devicePath, "vendor"))
		if err != nil {
			log.Warnf("Can't read vendor file: %v", err)
			continue
		}
		if uint16(vendorID) != targetVendorID {
			log.Debugf("VendorID not match 0x%x", vendorID)
			continue
		}
		classCode, err := readHexFile(filepath.Join(devicePath, "class"))
		if err != nil {
			log.Warnf("Can't read class file: %v", err)
			continue
		}
		classBase := uint8((classCode >> 16) & 0xFF)
		if classBase == targetClassBase {
			log.Debugf("Find mlu device: %s with vendorID 0x%x,classCode: 0x%x, classBase: 0x%x", devicePath, vendorID, classCode, classBase)
			count++
		}
	}

	log.Debugf("Find %d mlu devices", count)
	return count, nil
}

func isDriverRunning(counts uint, cli cndev.Cndev) bool {
	for i := range counts {
		_, good, running, _, err := cli.GetDeviceHealth(i)
		if err != nil {
			log.Warn(errors.Wrapf(err, "GetDeviceHealth %d", i))
			return false
		}
		if !good {
			log.Warnf("MLU device %d health maybe in problem, ignoring at init", i)
		}
		if !running {
			log.Warnf("MLU device %d driver is not running", i)
			return false
		}
	}
	return true
}

func EnsureMLUAllOK(cli cndev.Cndev, mluInfo *MLUStatMap, ignoreMissingLabels bool) {
	log.Infof("Start to ensure mlu driver status is ok")
	i := 1
	for {
		if i < 60000 {
			i = i << 1
		}
		time.Sleep(time.Duration(min(i, 60000)) * time.Millisecond)

		if i == 2 && ignoreMissingLabels {
			if err := cli.Init(false); err != nil {
				log.Errorf("Init cndev client failed with err: %v", err)
				continue
			}
		} else {
			if err := cli.Init(true); err != nil {
				log.Errorf("Init cndev client failed with err: %v", err)
				continue
			}
		}

		log.Debug("Start GetDeviceCount")
		counts, err := cli.GetDeviceCount()
		if err != nil {
			log.Error(errors.Wrap(err, "GetDeviceCount"))
			continue
		}
		if counts == 0 {
			log.Warn("No MLU device found with GetDeviceCount")
			continue
		}
		log.Debugf("Devie counts: %d", counts)

		realCounts, err := fetchMLUCounts()
		if err != nil {
			log.Error(errors.Wrap(err, "fetchMLUCounts"))
			continue
		}
		log.Debugf("RealCounts is :%d ", realCounts)

		if counts != realCounts {
			log.Warnf("MLU device count not match, counts: %d, realCounts: %d", counts, realCounts)
			continue
		}
		log.Debugf("MLU device count match, count is %d", counts)

		if err := cli.GenerateDeviceHandleMap(counts); err != nil {
			log.Panicf("Generate Device Handle Map failed, this should never happen, counts: %d, err: %v", counts, err)
		}

		if !isDriverRunning(counts, cli) {
			log.Warn("MLU driver is in problem, please check the device")
			if !ignoreMissingLabels {
				log.Debug("MLU driver is not running, now as ignoreMissingLabels is false, should try to get driver status again")
				continue
			}
		}

		collectMLUInfo(mluInfo, cli, counts)

		if mluInfo.InProblem.Load() {
			log.Warn("MLU labels are missing, will try to get them again")
			if !ignoreMissingLabels {
				log.Debug("MLU labels are missing, now as ignoreMissingLabels is false, should try to get labels again")
				continue
			}
		} else {
			log.Info("Driver of MLU devices are all running")
		}

		return
	}
}

func EnsureCndevLib() error {
	var src string
	x86Src := "/host/usr/lib/x86_64-linux-gnu/libcndev.so"
	armSrc := "/host/usr/lib/aarch64-linux-gnu/libcndev.so"
	lib64Src := "/host/usr/lib64/libcndev.so"
	if _, err := os.Stat(x86Src); err == nil {
		src = x86Src
		log.Infof("Found libcndev.so on host: %s", x86Src)
	} else if _, err := os.Stat(armSrc); err == nil {
		src = armSrc
		log.Infof("Found libcndev.so on host: %s", armSrc)
	} else if _, err := os.Stat(lib64Src); err == nil {
		src = lib64Src
		log.Infof("Found libcndev.so on host: %s", lib64Src)
	} else {
		log.Info("Found no libcndev.so on host, use default")
		return nil
	}

	dst := "/usr/lib/libcndev.so"
	log.Info("Found libcndev.so in host, try to copy to /usr/lib")
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Errorf("Can't open %s", src)
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		log.Errorf("Can't create %s", dst)
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
