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

package metrics

const (
	Cndev        = "cndev"
	PodResources = "podresources"

	Temperature   = "temperature"
	BoardHealth   = "board_health"
	MemTotal      = "physical_memory_total"
	MemUsed       = "physical_memory_used"
	MemUtil       = "memory_utilization"
	BoardUtil     = "board_utilization"
	CoreUtil      = "core_utilization"
	FanSpeed      = "fan_speed"
	BoardPower    = "board_power"
	BoardCapacity = "board_capacity"
	BoardUsage    = "board_usage"
	BoardVersion  = "board_version"

	BoardAllocated         = "board_allocated"
	MLUContainer           = "mlu_container"
	ContainerMLUUtil       = "container_mlu_utilization"
	ContainerMLUVFUtil     = "container_mlu_vf_utilization"
	ContainerMLUMemUtil    = "container_mlu_memory_utilization"
	ContainerMLUBoardPower = "container_mlu_board_power"

	Slot     = "slot"
	Model    = "model"
	MluType  = "mluType"
	SN       = "sn"
	Hostname = "nodeName"
	Cluster  = "cluster"
	Core     = "core"
	MCU      = "mcu"
	Driver   = "driver"

	Namespace = "namespace"
	Pod       = "pod"
	Container = "container"
	VFID      = "vfID"
)
