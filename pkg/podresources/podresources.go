/*Copyright (c) 2020, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package podresources

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	podresourcesapi "k8s.io/kubernetes/pkg/kubelet/apis/podresources/v1alpha1"
)

type PodResources interface {
	GetDeviceToPodInfo() (map[string]PodInfo, error)
}

type podResources struct {
	timeout  time.Duration
	socket   string
	reources []string
	maxSize  int
}

func NewPodResourcesClient(timeout time.Duration, socket string, resources []string, maxSize int) PodResources {
	return &podResources{
		timeout:  timeout,
		socket:   socket,
		reources: resources,
		maxSize:  maxSize,
	}
}

type PodInfo struct {
	Pod       string
	Namespace string
	Container string
}

func connectToServer(socket string, timeout time.Duration, maxSize int) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, socket, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxSize)),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failure connecting to %s: %v", socket, err)
	}
	return conn, nil
}

func listPods(socket string, timeout time.Duration, maxSize int) (*podresourcesapi.ListPodResourcesResponse, error) {
	conn, err := connectToServer(socket, timeout, maxSize)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := podresourcesapi.NewPodResourcesListerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := client.List(ctx, &podresourcesapi.ListPodResourcesRequest{})
	if err != nil {
		return nil, fmt.Errorf("failure getting pod resources %v", err)
	}
	return resp, nil
}

func contains(set []string, s string) bool {
	for i := range set {
		if strings.Contains(s, set[i]) {
			return true
		}
	}
	return false
}

func getDeviceToPodInfo(pods podresourcesapi.ListPodResourcesResponse, resources []string) map[string]PodInfo {
	m := make(map[string]PodInfo)
	for _, pod := range pods.GetPodResources() {
		for _, container := range pod.GetContainers() {
			for _, device := range container.GetDevices() {
				if !contains(resources, device.GetResourceName()) {
					continue
				}
				podInfo := PodInfo{
					Pod:       pod.GetName(),
					Namespace: pod.GetNamespace(),
					Container: container.GetName(),
				}
				for _, uuid := range device.GetDeviceIds() {
					m[uuid] = podInfo
				}
			}
		}
	}
	return m
}

func (k *podResources) GetDeviceToPodInfo() (map[string]PodInfo, error) {
	pods, err := listPods(k.socket, k.timeout, k.maxSize)
	if err != nil {
		return nil, fmt.Errorf("list pods: %v", err)
	}
	info := getDeviceToPodInfo(*pods, k.reources)
	return info, nil
}
