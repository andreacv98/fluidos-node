// Copyright 2022-2023 FLUIDOS Project
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

package grpc

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	reservationv1alpha1 "github.com/fluidos-project/node/apis/reservation/v1alpha1"
)

func getContractResourcesByClusterID(cl client.Client, clusterID string) (map[string]*resource.Quantity, error) {
	var contracts reservationv1alpha1.ContractList

	if err := cl.List(context.Background(), &contracts, client.MatchingFields{"spec.buyerClusterID": clusterID}); err != nil {
		if client.IgnoreNotFound(err) != nil {
			klog.Errorf("Error when listing Contracts: %s", err)
			return nil, err
		}
	}

	if len(contracts.Items) == 0 {
		klog.Errorf("No contracts found for cluster %s", clusterID)
		return nil, fmt.Errorf("no contracts found for cluster %s", clusterID)
	}

	if len(contracts.Items) > 1 {
		resources := multipleContractLogic(contracts.Items)
		return resources, nil
	}

	contract := contracts.Items[0]

	return addResources(make(map[string]*resource.Quantity), contract.Spec.Partition), nil
}

func multipleContractLogic(contracts []reservationv1alpha1.Contract) map[string]*resource.Quantity {
	resources := make(map[string]*resource.Quantity)
	for i := range contracts {
		if contracts[i].Spec.Partition != nil {
			resources = addResources(resources, contracts[i].Spec.Partition)
		} else {
			klog.Errorf("Contract %s has no partition", contracts[i].Name)
		}

	}
	return resources
}

// This function adds the resources of a contract to the existing resourceList.
func addResources(resources map[string]*resource.Quantity, partition *nodecorev1alpha1.Partition) map[string]*resource.Quantity {

	// Parse partition
	partitionType, paritionData, err := nodecorev1alpha1.ParsePartition(partition)
	if err != nil {
		klog.Errorf("Error when parsing partition: %s", err)
		return nil
	}
	switch partitionType {
	case nodecorev1alpha1.Type_K8Slice:

		for key, value := range mapQuantityToResourceList(paritionData.(*nodecorev1alpha1.K8SlicePartition)) {
			if prevRes, ok := resources[key]; !ok {
				resources[key] = value
			} else {
				prevRes.Add(*value)
				resources[key] = prevRes
			}
		}
	default:
		klog.Errorf("Partition type %s not supported", partitionType)
		return nil
	}
	return resources
}

func mapQuantityToResourceList(partition *nodecorev1alpha1.K8SlicePartition) map[string]*resource.Quantity {
	resources := make(map[string]*resource.Quantity)
	resources[corev1.ResourceCPU.String()] = &partition.CPU
	resources[corev1.ResourceMemory.String()] = &partition.Memory
	resources[corev1.ResourcePods.String()] = &partition.Pods
	resources[corev1.ResourceStorage.String()] = partition.Storage
	resources[corev1.ResourceEphemeralStorage.String()] = partition.Storage
	return resources
}
