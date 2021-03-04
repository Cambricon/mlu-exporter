// Copyright 2021 Cambricon, Inc.
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

package cnpapi

// #cgo LDFLAGS: -ldl
// #include "include/cnpapi.h"
// #include "include/cnpapi_pmu_api.h"
// #include "include/cnpapi_types.h"
// #include "cnpapi_dl.h"
import "C"
import (
	"errors"
	"fmt"
	"strings"
)

const (
	PCIeWriteBytes    = "pcie__write_bytes"
	PCIeReadBytes     = "pcie__read_bytes"
	DramWriteBytes    = "dram__write_bytes"
	DramReadBytes     = "dram__read_bytes"
	MLULinkWriteBytes = "mlulink__write_bytes"
	MLULinkReadBytes  = "mlulink__read_bytes"
)

var counterNames = []string{
	PCIeWriteBytes,
	PCIeReadBytes,
	DramWriteBytes,
	DramReadBytes,
	MLULinkWriteBytes,
	MLULinkReadBytes,
}

type Cnpapi interface {
	Init() error
	PmuFlushData(idx uint) error
	GetPCIeReadBytes(idx uint) (uint64, error)
	GetPCIeWriteBytes(idx uint) (uint64, error)
	GetDramReadBytes(idx uint) (uint64, error)
	GetDramWriteBytes(idx uint) (uint64, error)
	GetMLULinkReadBytes(idx uint) (uint64, error)
	GetMLULinkWriteBytes(idx uint) (uint64, error)
}

type cnpapi struct {
	counters map[string]C.uint64_t
}

func NewCnpapiClient() Cnpapi {
	return &cnpapi{}
}

func errorString(ret C.cnpapiResult) error {
	if ret == C.CNPAPI_SUCCESS {
		return nil
	}
	var msg *C.char
	result := C.cnpapiGetResultString(ret, &msg)
	if result != C.CNPAPI_SUCCESS {
		return fmt.Errorf("cnpapiGetResultString not succeeded, result: %v, input: %v", result, ret)
	}
	return fmt.Errorf("cnpapi: %v", C.GoString(msg))
}

func (c *cnpapi) Init() error {
	r := C.cnpapiInit_dl()
	if r == C.CNPAPI_ERROR_NOT_INITIALIZED {
		return errors.New("could not load cnpapi library")
	}
	if err := errorString(r); err != nil {
		return fmt.Errorf("init err: %w", err)
	}
	if err := c.getCounters(); err != nil {
		return fmt.Errorf("getCounters err: %w", err)
	}
	r = C.cnpapiPmuSetFlushMode(C.PMU_EXPLICIT_FLUSH)
	if err := errorString(r); err != nil {
		return fmt.Errorf("cnpapiPmuSetFlushMode err: %w", err)
	}
	devs, err := c.getDeviceCount()
	if err != nil {
		return fmt.Errorf("getDeviceCount err: %w", err)
	}
	for n := uint(0); n < devs; n++ {
		if err := c.enableCounters(n); err != nil {
			return fmt.Errorf("enableCounters err: %w", err)
		}
	}
	return nil
}

func (c *cnpapi) getCounters() error {
	counters := make(map[string]C.uint64_t)
	device, err := c.getDeviceType(uint(0))
	if err != nil {
		return err
	}
	for _, name := range counterNames {
		if strings.Contains(name, "mlulink") && !supportsMLULink(device) {
			continue
		}
		var id C.uint64_t
		r := C.cnpapiPmuGetCounterIdByName(C.CString(name), &id)
		if err := errorString(r); err != nil {
			return err
		}
		counters[name] = id
	}
	c.counters = counters
	return nil
}

func supportsMLULink(device C.cnpapiDeviceType_t) bool {
	if device == C.CNPAPI_MLU290 {
		return true
	}
	return false
}

func (c *cnpapi) getDeviceType(idx uint) (C.cnpapiDeviceType_t, error) {
	var device C.cnpapiDeviceType_t
	r := C.cnpapiGetDeviceType(C.int(idx), &device)
	return device, errorString(r)
}

func (c *cnpapi) getDeviceCount() (uint, error) {
	var count C.int
	r := C.cnpapiGetDeviceCount(&count)
	if err := errorString(r); err != nil {
		return 0, err
	}
	return uint(count), nil
}

func (c *cnpapi) enableCounters(idx uint) error {
	for _, name := range counterNames {
		id := c.counters[name]
		r := C.cnpapiPmuEnableCounter(C.int(idx), id, true)
		if err := errorString(r); err != nil {
			return err
		}
	}
	return nil
}

func (c *cnpapi) PmuFlushData(idx uint) error {
	r := C.cnpapiPmuFlushData(C.int(idx))
	return errorString(r)
}

func (c *cnpapi) GetPCIeReadBytes(idx uint) (uint64, error) {
	var value C.uint64_t
	id, ok := c.counters[PCIeReadBytes]
	if !ok {
		return 0, nil
	}
	r := C.cnpapiPmuGetCounterValue(C.int(idx), id, &value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	return uint64(value), nil
}

func (c *cnpapi) GetPCIeWriteBytes(idx uint) (uint64, error) {
	var value C.uint64_t
	id, ok := c.counters[PCIeWriteBytes]
	if !ok {
		return 0, nil
	}
	r := C.cnpapiPmuGetCounterValue(C.int(idx), id, &value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	return uint64(value), nil
}

func (c *cnpapi) GetDramReadBytes(idx uint) (uint64, error) {
	var value C.uint64_t
	id, ok := c.counters[DramReadBytes]
	if !ok {
		return 0, nil
	}
	r := C.cnpapiPmuGetCounterValue(C.int(idx), id, &value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	return uint64(value), nil
}

func (c *cnpapi) GetDramWriteBytes(idx uint) (uint64, error) {
	var value C.uint64_t
	id, ok := c.counters[DramWriteBytes]
	if !ok {
		return 0, nil
	}
	r := C.cnpapiPmuGetCounterValue(C.int(idx), id, &value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	return uint64(value), nil
}

func (c *cnpapi) GetMLULinkReadBytes(idx uint) (uint64, error) {
	var value C.uint64_t
	id, ok := c.counters[MLULinkReadBytes]
	if !ok {
		return 0, nil
	}
	r := C.cnpapiPmuGetCounterValue(C.int(idx), id, &value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	return uint64(value), nil
}

func (c *cnpapi) GetMLULinkWriteBytes(idx uint) (uint64, error) {
	var value C.uint64_t
	id, ok := c.counters[MLULinkWriteBytes]
	if !ok {
		return 0, nil
	}
	r := C.cnpapiPmuGetCounterValue(C.int(idx), id, &value)
	if err := errorString(r); err != nil {
		return 0, err
	}
	return uint64(value), nil
}
