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

type labelInfo struct {
	cluster     string
	core        int
	cpuCore     string
	host        string
	link        int
	linkVersion string
	memoryDie   string
	pcieInfo    pcieInfo
	pid         uint32
	podInfo     podresources.PodInfo
	stat        MLUStat
	typ         string
	vf          string
	xidInfo     cndev.XIDInfoWithTimestamp
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
					if inf.InstanceID == index {
						typ = typ + ".mim-" + inf.Name
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
						if inf.InstanceID == index {
							typ = typ + ".smlu-" + inf.Name
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
					if inf.InstanceID == index {
						values = append(values, inf.UUID)
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
					if inf.InstanceID == index {
						values = append(values, inf.UUID)
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
		case Core:
			values = append(values, fmt.Sprintf("%d", info.core))
		case CPUCore:
			values = append(values, info.cpuCore)
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
		case XID:
			values = append(values, fmt.Sprintf("%d", info.xidInfo.XID))
		case XIDTimestamp:
			values = append(values, info.xidInfo.Timestamp)
		case XIDComputeInstanceID:
			values = append(values, fmt.Sprintf("%d", info.xidInfo.ComputeInstanceID))
		case XIDMLUInstanceID:
			values = append(values, fmt.Sprintf("%d", info.xidInfo.MLUInstanceID))
		default:
			values = append(values, "") // configured label not applicable, add this to prevent panic
		}
	}
	log.Debugf("getLabelValues: %+v", values)
	return values
}

type MLUStat struct {
	crcDisabled          bool
	driver               string
	eccDisabled          bool
	heartBeatDisabled    bool
	link                 int
	linkActive           map[int]bool
	mcu                  string
	memEccDisabled       bool
	mimEnabled           bool
	mimInfos             []cndev.MimInfo
	model                string
	processUtilDisabled  bool
	remappedRowsDisabled bool
	retiredPageDisabled  bool
	slot                 uint
	smluEnabled          bool
	smluInfos            []cndev.SmluInfo
	sn                   string
	uuid                 string
	xidCallbackDisabled  bool
}

func CollectMLUInfo(cli cndev.Cndev) map[string]MLUStat {
	num, err := cli.GetDeviceCount()
	if err != nil {
		log.Error(errors.Wrap(err, "GetDeviceCount"))
	}
	info := make(map[string]MLUStat)
	for i := uint(0); i < num; i++ {
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
		for j := 0; j < link; j++ {
			active, _, err := cli.GetDeviceMLULinkStatus(i, uint(j))
			if err != nil {
				log.Warn(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkStatus", i, j))
				continue
			}
			linkActive[j] = active != 0
		}
		var crcDisabled, eccDisabled, processUtilDisabled, retiredPageDisabled, remappedRowsDisabled bool
		if _, _, err = cli.GetDeviceCRCInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceCRCInfo", i))
			crcDisabled = true
		}
		if _, _, _, _, _, _, _, _, err = cli.GetDeviceECCInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", i))
			eccDisabled = true
		}
		if _, _, _, _, _, _, err = cli.GetDeviceProcessUtil(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", i))
			processUtilDisabled = true
		}
		if _, _, err = cli.GetDeviceRetiredPageInfo(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", i))
			retiredPageDisabled = true
		}
		if _, _, _, _, err = cli.GetDeviceRemappedRows(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceRemappedRows", i))
			remappedRowsDisabled = true
		}
		var heartBeatDisabled, memEccDisabled bool
		if _, err = cli.GetDeviceHeartbeatCount(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceHeartbeatCount", i))
			heartBeatDisabled = true
		}
		if _, _, _, _, _, err = cli.GetDeviceMemEccCounter(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", i))
			memEccDisabled = true
		}

		var xidCallbackDisabled bool
		if err = cli.GetSupportedEventTypes(i); err != nil {
			log.Warn(errors.Wrapf(err, "Slot %d GetSupportedEventTypes", i))
			xidCallbackDisabled = true
		}

		info[uuid] = MLUStat{
			crcDisabled:          crcDisabled,
			driver:               calcVersion(driverMajor, driverMinor, driverBuild),
			eccDisabled:          eccDisabled,
			heartBeatDisabled:    heartBeatDisabled,
			link:                 link,
			linkActive:           linkActive,
			mcu:                  calcVersion(mcuMajor, mcuMinor, mcuBuild),
			memEccDisabled:       memEccDisabled,
			mimEnabled:           mimEnabled,
			mimInfos:             mimInfos,
			model:                model,
			processUtilDisabled:  processUtilDisabled,
			remappedRowsDisabled: remappedRowsDisabled,
			retiredPageDisabled:  retiredPageDisabled,
			slot:                 i,
			smluEnabled:          smluEnabled,
			sn:                   sn,
			uuid:                 uuid,
			xidCallbackDisabled:  xidCallbackDisabled,
		}
	}
	log.Debugf("collectSharedInfo: %+v", info)
	return info
}

func calcVersion(major uint, minor uint, build uint) string {
	return fmt.Sprintf("v%d.%d.%d", major, minor, build)
}
