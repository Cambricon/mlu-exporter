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

package host

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Host interface {
	GetCPUStats() (float64, float64, error)
	GetMemoryStats() (float64, float64, error)
}

type host struct {
}

func NewHostClient() Host {
	return &host{}
}

func (h *host) GetCPUStats() (float64, float64, error) {
	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	info := strings.Split(string(contents), "\n")[0]
	var cpu string
	var user, nice, sys, idle, iowait, irq, softirq, steal, guest, guestnice float64
	_, err = fmt.Sscanf(info, "%s %v %v %v %v %v %v %v %v %v %v", &cpu, &user, &nice, &sys, &idle, &iowait, &irq, &softirq, &steal, &guest, &guestnice)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to scan %s, err: %w", info, err)
	}
	idle = idle + iowait
	noneIdle := user + nice + sys + irq + softirq + steal
	total := idle + noneIdle
	log.Debugf("CPU Total: %v ,CPU Idle: %v", total, idle)
	return total, idle, err
}

func (h *host) GetMemoryStats() (float64, float64, error) {
	contents, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(contents), "\n")
	var total, free float64
	_, err = fmt.Sscanf(lines[0], "MemTotal: %v kB", &total)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to scan %s, err: %w", lines[0], err)
	}
	_, err = fmt.Sscanf(lines[1], "MemFree: %v kB", &free)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to scan %s, err: %w", lines[1], err)
	}
	log.Debugf("Memory Total: %v ,Memory Free: %v", total, free)
	return total, free, err
}
