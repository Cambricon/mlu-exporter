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

const (
	Cndev        = "cndev"
	PodResources = "podresources"
	Cnpapi       = "cnpapi"
	Host         = "host"

	Temperature         = "temperature"
	Health              = "health"
	MemTotal            = "memory_total"
	MemUsed             = "memory_used"
	MemUtil             = "memory_utilization"
	Util                = "utilization"
	CoreUtil            = "core_utilization"
	FanSpeed            = "fan_speed"
	PowerUsage          = "power_usage"
	Capacity            = "capacity"
	Usage               = "usage"
	Version             = "version"
	VirtualMemUtil      = "virtual_memory_utilization"
	VirtualFunctionUtil = "virtual_function_utilization"

	Allocated = "allocated"
	Container = "container"

	PCIeRead     = "pcie_read"
	PCIeWrite    = "pcie_write"
	DramRead     = "dram_read"
	DramWrite    = "dram_write"
	MLULinkRead  = "mlulink_read"
	MLULinkWrite = "mlulink_write"

	HostCPUTotal = "host_cpu_total"
	HostCPUIdle  = "host_cpu_idle"
	HostMemTotal = "host_memory_total"
	HostMemFree  = "host_memory_free"

	MLU     = "mlu"
	Model   = "model"
	Type    = "type"
	SN      = "sn"
	UUID    = "uuid"
	Node    = "node"
	Cluster = "cluster"
	Core    = "core"
	MCU     = "mcu"
	Driver  = "driver"

	Namespace = "namespace"
	Pod       = "pod"
	VF        = "vf"
)
