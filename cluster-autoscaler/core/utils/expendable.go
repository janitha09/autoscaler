/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

// FilterOutExpendableAndSplit filters out expendable pods and splits into:
//   - waiting for lower priority pods preemption
//   - other pods.
func FilterOutExpendableAndSplit(unschedulableCandidates []*apiv1.Pod, expendablePodsPriorityCutoff int) ([]*apiv1.Pod, []*apiv1.Pod) {
	var unschedulableNonExpendable []*apiv1.Pod
	var waitingForLowerPriorityPreemption []*apiv1.Pod
	for _, pod := range unschedulableCandidates {
		if pod.Spec.Priority != nil && int(*pod.Spec.Priority) < expendablePodsPriorityCutoff {
			klog.V(4).Infof("Pod %s has priority below %d (%d) and will scheduled when enough resources is free. Ignoring in scale up.", pod.Name, expendablePodsPriorityCutoff, *pod.Spec.Priority)
		} else if nominatedNodeName := pod.Status.NominatedNodeName; nominatedNodeName != "" {
			waitingForLowerPriorityPreemption = append(waitingForLowerPriorityPreemption, pod)
			klog.V(4).Infof("Pod %s will be scheduled after low priority pods are preempted on %s. Ignoring in scale up.", pod.Name, nominatedNodeName)
		} else {
			unschedulableNonExpendable = append(unschedulableNonExpendable, pod)
		}
	}
	return unschedulableNonExpendable, waitingForLowerPriorityPreemption
}

// FilterOutExpendablePods filters out expendable pods.
func FilterOutExpendablePods(pods []*apiv1.Pod, expendablePodsPriorityCutoff int) []*apiv1.Pod {
	var result []*apiv1.Pod
	for _, pod := range pods {
		if pod.Spec.Priority == nil || int(*pod.Spec.Priority) >= expendablePodsPriorityCutoff {
			result = append(result, pod)
		}
	}
	return result
}
