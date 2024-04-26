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
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	podresourcesapi "k8s.io/kubelet/pkg/apis/podresources/v1alpha1"
)

type PodResources interface {
	GetDeviceToPodInfo() (map[string]PodInfo, error)
}

type podResources struct {
	client    kubernetes.Interface
	host      string
	maxSize   int
	mode      string
	resources []string
	socket    string
	timeout   time.Duration
}

func NewPodResourcesClient(timeout time.Duration, socket string, resources []string, maxSize int, mode, host string, client kubernetes.Interface) PodResources {
	podResource := &podResources{
		timeout:   timeout,
		socket:    socket,
		resources: resources,
		maxSize:   maxSize,
		mode:      mode,
	}
	if mode == "dynamic-smlu" {
		podResource.host = host
		podResource.client = client
	}
	return podResource
}

type PodInfo struct {
	Container string
	Index     uint
	Namespace string
	Pod       string
	VF        string
}

func connectToServer(socket string, timeout time.Duration, maxSize int) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	creds := insecure.NewCredentials()
	dialer := func(ctx context.Context, address string) (net.Conn, error) {
		return net.DialTimeout("unix", address, timeout)
	}
	conn, err := grpc.DialContext(ctx, socket,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxSize)),
		grpc.WithContextDialer(dialer),
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

func (p podResources) getDeviceToPodInfo(pods podresourcesapi.ListPodResourcesResponse) (map[string]PodInfo, error) {
	var podList *v1.PodList
	if p.host != "" {
		var err error
		selector := fields.SelectorFromSet(fields.Set{"spec.nodeName": p.host, "status.phase": "Running"})
		podList, err = p.client.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
			FieldSelector: selector.String(),
		})
		if err != nil {
			return nil, err
		}
	}

	m := make(map[string]PodInfo)
	if p.mode == "dynamic-smlu" && podList != nil {
		for _, pod := range podList.Items {
			if pod.ObjectMeta.Annotations == nil {
				continue
			}
			anno, ok := pod.ObjectMeta.Annotations["CAMBRICON_DSMLU_PROFILE_INSTANCE"]
			if !ok {
				continue
			}
			v := strings.Split(anno, "_")
			// 1_256_1_0 represents profileID_instanceHandle_slot_instanceID
			if len(v) != 4 {
				log.Printf("Found pod %s with invalid annotation %s, ignore it", pod.Name, anno)
				continue
			}
			index, err := strconv.Atoi(v[2])
			if err != nil {
				log.Printf("strconv value %v, %v", v[2], err)
				continue
			}
			pi := PodInfo{
				Pod:       pod.Name,
				Namespace: pod.Namespace,
				Index:     uint(index),
				VF:        v[3],
			}
			for _, c := range pod.Spec.Containers {
				for k, v := range c.Resources.Limits {
					if v.Value() > 0 && strings.HasPrefix(k.String(), "cambricon.com/") &&
						(strings.HasSuffix(k.String(), ".vcore") || strings.HasSuffix(k.String(), ".vmemory")) {
						pi.Container = c.Name
						break
					}
				}
			}
			m[v[1]] = pi
		}
		return m, nil
	}

	for _, pod := range pods.GetPodResources() {
		for _, container := range pod.GetContainers() {
			for _, device := range container.GetDevices() {
				if !contains(p.resources, device.GetResourceName()) {
					continue
				}
				podName := pod.GetName()
				podNamespace := pod.GetNamespace()
				containerName := container.GetName()
				podInfo := PodInfo{
					Pod:       podName,
					Namespace: podNamespace,
					Container: containerName,
				}
				for _, uuid := range device.GetDeviceIds() {
					m[uuid] = podInfo
				}
			}
		}
	}

	return m, nil
}

func (p podResources) GetDeviceToPodInfo() (map[string]PodInfo, error) {
	pods, err := listPods(p.socket, p.timeout, p.maxSize)
	if err != nil {
		return nil, fmt.Errorf("list pods: %v", err)
	}
	return p.getDeviceToPodInfo(*pods)
}
