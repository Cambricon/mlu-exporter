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
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubelet/pkg/apis/podresources/v1alpha1"
)

func TestGetDeviceToPodInfo(t *testing.T) {
	node := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testnode",
		},
		Spec: v1.NodeSpec{},
	}

	pods := []*v1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testpod1",
				Namespace: "ns",
				Annotations: map[string]string{
					"CAMBRICON_DSMLU_PROFILE_INSTANCE": "1_256_0_1",
				},
			},
			Spec: v1.PodSpec{
				NodeName: "testnode",
				Containers: []v1.Container{
					{
						Name: "testcontainer",
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceName("cambricon.com/mlu.vcore"):   *resource.NewQuantity(1, resource.DecimalSI),
								v1.ResourceName("cambricon.com/mlu.vmemory"): *resource.NewQuantity(1, resource.DecimalSI),
							},
						},
					},
				},
			},
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testpod2",
				Namespace: "ns",
				Annotations: map[string]string{
					"CAMBRICON_DSMLU_PROFILE_INSTANCE": "2_512_0_2",
				},
			},
			Spec: v1.PodSpec{
				NodeName: "testnode",
				Containers: []v1.Container{
					{
						Name: "testcontainer",
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceName("cambricon.com/mlu.vcore"):   *resource.NewQuantity(2, resource.DecimalSI),
								v1.ResourceName("cambricon.com/mlu.vmemory"): *resource.NewQuantity(2, resource.DecimalSI),
							},
						},
					},
				},
			},
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testpod1",
				Namespace: "ns",
			},
			Spec: v1.PodSpec{
				NodeName: "testnode",
				Containers: []v1.Container{
					{
						Name: "testcontainer",
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceName("cambricon.com/mlu.vcore"):   *resource.NewQuantity(1, resource.DecimalSI),
								v1.ResourceName("cambricon.com/mlu.vmemory"): *resource.NewQuantity(1, resource.DecimalSI),
							},
						},
					},
				},
			},
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
			},
		},
	}

	fakeClient := fake.NewSimpleClientset(node, pods[0], pods[1])
	p := &podResources{
		client: fakeClient,
		host:   "testnode",
		mode:   "dynamic-smlu",
	}

	expectInfo := map[string]PodInfo{
		"256": {
			Index:     0,
			Namespace: "ns",
			Pod:       "testpod1",
			Container: "testcontainer",
			VF:        "1",
		},
		"512": {
			Index:     0,
			Namespace: "ns",
			Pod:       "testpod2",
			Container: "testcontainer",
			VF:        "2",
		},
	}

	podInfo, err := p.getDeviceToPodInfo(v1alpha1.ListPodResourcesResponse{})
	assert.NoError(t, err)
	if len(podInfo) != len(expectInfo) {
		t.Error("actual not equal to expect")
	}
	for k, v := range expectInfo {
		if podInfo[k] != v {
			t.Error("actual not equal to expect")
		}
	}
}
