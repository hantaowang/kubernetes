/*
Copyright 2015 The Kubernetes Authors.

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

package cache

import (
	"strconv"
	"k8s.io/api/core/v1"
)

// CreateNodeNameToInfoMap obtains a list of pods and pivots that list into a map where the keys are node names
// and the values are the aggregated information for that node.
func CreateNodeNameToInfoMap(pods []*v1.Pod, nodes []*v1.Node) map[string]*NodeInfo {
	nodeNameToInfo := make(map[string]*NodeInfo)
	for _, pod := range pods {
		nodeName := pod.Spec.NodeName
		if _, ok := nodeNameToInfo[nodeName]; !ok {
			nodeNameToInfo[nodeName] = NewNodeInfo()
		}
		nodeNameToInfo[nodeName].AddPod(pod)
	}
	for _, node := range nodes {
		if _, ok := nodeNameToInfo[node.Name]; !ok {
			nodeNameToInfo[node.Name] = NewNodeInfo()
		}
		nodeNameToInfo[node.Name].SetNode(node)
	}
	return nodeNameToInfo
}

// gets the network request for a pod based on its annotation, or 0 if no such limit exists
var charToMultiplier = map[string]int64{
	"K": 1000,
	"M": 1000000,
	"G": 1000000000,
	"T": 1000000000000,
}

func GetNetworkRequest(pod *v1.Pod) int64 {
	requestStr, hasRequirement := pod.Annotations["netsys.io/network-bandwidth"]
	if !hasRequirement {
		return 0
	}
	if len(requestStr) == 0 {
		return 0
	}
	lastChar := requestStr[len(requestStr) - 1]
	var multiplier int64 = 1
	if x, ok := charToMultiplier[string(lastChar)]; ok {
		multiplier = x
	}
	val, err := strconv.ParseInt(requestStr[:len(requestStr)-1], 10, 32)
	if err != nil {
		return 0
	}
	return val * multiplier
}
