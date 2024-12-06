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
	"strconv"
	"strings"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/podresources"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

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
	cluster           string
	cndevVersion      string
	computeCapability string
	core              int
	cpuCore           string
	host              string
	lane              int
	link              int
	linkVersion       string
	memoryDie         string
	mimInstanceInfo   cndev.MimInstanceInfo
	mimProfileInfo    cndev.MimProfileInfo
	mluLinkRemoteInfo mluLinkRemoteInfo
	pcieInfo          pcieInfo
	pid               uint32
	podInfo           podresources.PodInfo
	smluInstanceInfo  cndev.SmluInstanceInfo
	smluProfileInfo   cndev.SmluProfileInfo
	stat              MLUStat
	typ               string
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
					if err := cli.Init(false); err != nil {
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
					if err := cli.Init(false); err != nil {
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
		case Container:
			values = append(values, info.podInfo.Container)
		case VF:
			values = append(values, info.vf)
		case Lane:
			values = append(values, fmt.Sprintf("%d", info.lane))
		case Link:
			values = append(values, fmt.Sprintf("%d", info.link))
		case LinkVersion:
			values = append(values, info.linkVersion)
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
			values = append(values, fmt.Sprintf("%d", info.xidInfo.XID))
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
		default:
			values = append(values, "") // configured label not applicable, add this to prevent panic
		}
	}
	log.Debugf("getLabelValues: %+v", values)
	return values
}

type MLUStat struct {
	cndevInterfaceDisabled map[string]bool
	driver                 string
	link                   int
	linkActive             map[int]bool
	mcu                    string
	mimEnabled             bool
	mimInfos               []cndev.MimInfo
	model                  string
	opticalPresent         map[int]uint8
	slot                   uint
	smluEnabled            bool
	smluInfos              []cndev.SmluInfo
	sn                     string
	uuid                   string
}

func CollectMLUInfo(cli cndev.Cndev) map[string]MLUStat {
	num, err := cli.GetDeviceCount()
	if err != nil {
		log.Error(errors.Wrap(err, "GetDeviceCount"))
	}
	info := make(map[string]MLUStat)
	for i := uint(0); i < num; i++ {
		dis := make(map[string]bool)
		model := cli.GetDeviceModel(i)
		uuid, err := cli.GetDeviceUUID(i)
		if err != nil {
			log.Error(errors.Wrap(err, "GetDeviceUUID"))
		}
		sn, err := cli.GetDeviceSN(i)
		if err != nil {
			log.Error(errors.Wrap(err, "GetDeviceSN"))
		}
		mcuMajor, mcuMinor, mcuBuild, driverMajor, driverMinor, driverBuild, err := cli.GetDeviceVersion(i)
		if err != nil {
			log.Error(errors.Wrap(err, "GetDeviceVersion"))
		}
		mimEnabled, err := cli.DeviceMimModeEnabled(i)
		if err != nil {
			log.Warn(errors.Wrap(err, "DeviceMimModeEnabled"))
		}
		var mimInfos []cndev.MimInfo
		if mimEnabled {
			mimInfos, err = cli.GetAllMLUInstanceInfo(i)
			if err != nil {
				log.Warn(errors.Wrap(err, "GetAllMLUInstanceInfo"))
			}
		}
		smluEnabled, err := cli.DeviceSmluModeEnabled(i)
		if err != nil {
			log.Warn(errors.Wrap(err, "DeviceMimModeEnabled"))
		}
		link := cli.GetDeviceMLULinkPortNumber(i)
		linkActive := map[int]bool{}
		opticalPresent := map[int]uint8{}
		for j := 0; j < link; j++ {
			active, _, _, err := cli.GetDeviceMLULinkStatus(i, uint(j))
			if err != nil {
				log.Debug(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkStatus", i, j))
				continue
			}
			linkActive[j] = active != 0
			if _, err = cli.GetDeviceMLULinkEventCounter(i, uint(j)); err != nil {
				log.Debug(errors.Wrapf(err, "Slot %d GetDeviceMLULinkEventCounter", i))
				dis["mluLinkEventCounterDisabled"] = true
			}
			if _, err = cli.GetDeviceMLULinkErrorCounter(i, uint(j)); err != nil {
				log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMLULinkErrorCounter", i))
				dis["mluLinkErrorCounterDisabled"] = true
			}
			if _, _, _, _, _, _, _, _, _, _, err = cli.GetDeviceMLULinkRemoteInfo(i, uint(j)); err != nil {
				log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMLULinkErrorCounter", i))
				dis["mluLinkRemoteInfoDisabled"] = true
			}
			present, _, _, _, _, err := cli.GetDeviceOpticalInfo(i, uint(j))
			if err != nil {
				log.Warn(errors.Wrapf(err, "Slot %d link %d GetDeviceOpticalInfo", i, j))
				continue
			}
			opticalPresent[j] = present
		}
		if _, _, err = cli.GetDeviceCRCInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceCRCInfo", i))
			dis["crcDisabled"] = true
		}
		if _, _, _, _, _, _, _, _, err = cli.GetDeviceECCInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", i))
			dis["eccDisabled"] = true
		}
		if _, _, _, _, _, _, err = cli.GetDeviceProcessUtil(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", i))
			dis["processUtilDisabled"] = true
		}
		if _, _, err = cli.GetDeviceRetiredPageInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", i))
			dis["retiredPageDisabled"] = true
		}
		if _, _, _, _, _, err = cli.GetDeviceRemappedRows(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceRemappedRows", i))
			dis["remappedRowsDisabled"] = true
		}
		if _, err = cli.GetDeviceHeartbeatCount(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceHeartbeatCount", i))
			dis["heartBeatDisabled"] = true
		}
		if _, _, err = cli.GetDeviceMemEccCounter(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", i))
			dis["memEccDisabled"] = true
		}
		if _, _, _, _, _, err = cli.GetDeviceAddressSwaps(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceAddressSwaps", i))
			dis["addressSwapsDisabled"] = true
		}
		if err = cli.GetSupportedEventTypes(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetSupportedEventTypes", i))
			dis["xidCallbackDisabled"] = true
		}
		if _, _, _, err = cli.GetDeviceCurrentInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceCurrentInfo", i))
			dis["currentInfoDisabled"] = true
		}
		if _, _, _, err = cli.GetDeviceVoltageInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceVoltageInfo", i))
			dis["voltageInfoDisabled"] = true
		}
		if _, err = cli.GetDevicePerformanceThrottleReason(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDevicePerformanceThrottleReason", i))
			dis["performanceThrottleDisabled"] = true
		}
		if _, _, err = cli.GetDeviceOverTemperatureInfo(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceOverTemperatureInfo", i))
			dis["overTemperatureInfoDisabled"] = true
		}
		if _, err = cli.GetDeviceOverTemperatureShutdownThreshold(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDeviceOverTemperatureThreshold", i))
			dis["overTemperatureThresholdDisabled"] = true
		}
		if _, err = cli.GetDevicePowerManagementLimitation(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDevicePowerManagementLimitation", i))
			dis["powerCurrentLimitDisabled"] = true
		}
		if _, err = cli.GetDevicePowerManagementDefaultLimitation(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDevicePowerManagementDefaultLimitation", i))
			dis["powerDefaultLimitDisabled"] = true
		}
		if _, _, err = cli.GetDevicePowerManagementLimitRange(i); err != nil {
			log.Debug(errors.Wrapf(err, "Slot %d GetDevicePowerManagementLimitRange", i))
			dis["powerLimitRangeDisabled"] = true
		}
		if _, _, _, err = cli.GetDeviceBAR4MemoryInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceBAR4MemoryInfo", i))
			dis["bar4MemoryInfoDisabled"] = true
		}
		if _, _, _, _, _, _, _, err = cli.GetDeviceCPUUtil(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceCPUUtil", i))
			dis["cpuUtilDisabled"] = true
		}
		if _, _, err = cli.GetDeviceMaxPCIeInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMaxPCIeInfo", i))
			dis["maxPCIeInfoDisabled"] = true
		}
		if _, _, err = cli.GetDeviceEccMode(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceEccMode", i))
			dis["eccModeDisabled"] = true
		}
		if _, err = cli.GetDeviceRetiredPagesOperation(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceRetiredPagesOperation", i))
			dis["retiredPagesOperationDisabled"] = true
		}
		if _, _, err = cli.GetDeviceComputeCapability(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceComputeCapability", i))
			dis["computeCapabilityDisabled"] = true
		}

		info[uuid] = MLUStat{
			cndevInterfaceDisabled: dis,
			driver:                 calcVersion(driverMajor, driverMinor, driverBuild),
			link:                   link,
			linkActive:             linkActive,
			mcu:                    calcVersion(mcuMajor, mcuMinor, mcuBuild),
			mimEnabled:             mimEnabled,
			mimInfos:               mimInfos,
			model:                  model,
			opticalPresent:         opticalPresent,
			slot:                   i,
			smluEnabled:            smluEnabled,
			sn:                     sn,
			uuid:                   uuid,
		}
		if len(dis) != 0 && i == 0 {
			log.Warnf("Some cndev interfaces are not supported: %v", dis)
		}
	}
	log.Debugf("collectSharedInfo: %+v", info)
	return info
}

func calcVersion(major uint, minor uint, build uint) string {
	return fmt.Sprintf("v%d.%d.%d", major, minor, build)
}
