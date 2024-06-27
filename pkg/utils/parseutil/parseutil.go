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

package parseutil

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	reservationv1alpha1 "github.com/fluidos-project/node/apis/reservation/v1alpha1"
	"github.com/fluidos-project/node/pkg/utils/models"
)

// ParseFlavorSelector parses FlavorSelector into a Selector.
func ParseFlavorSelector(selector *nodecorev1alpha1.Selector) (models.Selector, error) {
	// Parse the Selector
	selectorIdentifier, selectorStruct, err := nodecorev1alpha1.ParseSolverSelector(selector)
	if err != nil {
		return nil, err
	}

	switch selectorIdentifier {
	case nodecorev1alpha1.Type_K8Slice:
		// Force casting of selectorStruct to K8Slice
		selectorStruct := selectorStruct.(nodecorev1alpha1.K8SliceSelector)

		// Generate the model for the K8Slice selector
		k8SliceSelector, err := ParseK8SliceFilters(&selectorStruct)
		if err != nil {
			return nil, err
		}

		return k8SliceSelector, nil

	case nodecorev1alpha1.Type_VM:
		// Force casting of selectorStruct to VM
		// TODO: Implement the parsing of the VM selector
		return nil, fmt.Errorf("VM selector not implemented")

	case nodecorev1alpha1.Type_Service:
		// Force casting of selectorStruct to Service
		// TODO: Implement the parsing of the Service selector
		return nil, fmt.Errorf("service selector not implemented")

	default:
		return nil, fmt.Errorf("unknown selector type")

	}
}
func ParseK8SliceFilters(k8sSelector *nodecorev1alpha1.K8SliceSelector) (*models.K8SliceSelector, error) {

	var cpuFilterModel, memoryFilterModel, podsFilterModel, storageFilterModel models.ResourceQuantityFilter

	// Parse the CPU filter
	switch k8sSelector.CpuFilter.FilterType {
	case nodecorev1alpha1.TypeMatchFilter:
		// Unmarshal the data into a ResourceMatchSelector
		var cpuFilter nodecorev1alpha1.ResourceMatchSelector
		err := json.Unmarshal(k8sSelector.CpuFilter.Data.Raw, &cpuFilter)
		if err != nil {
			return nil, err
		}

		// Generate the model for the CPU filter
		cpuFilterModel = models.ResourceQuantityMatchFilter{
			Value: cpuFilter.Value.DeepCopy(),
		}
	case nodecorev1alpha1.TypeRangeFilter:
		// Unmarshal the data into a ResourceRangeSelector
		var cpuFilter nodecorev1alpha1.ResourceRangeSelector
		err := json.Unmarshal(k8sSelector.CpuFilter.Data.Raw, &cpuFilter)
		if err != nil {
			return nil, err
		}

		// Generate the model for the CPU filter
		cpuFilterModel = models.ResourceQuantityRangeFilter{
			Min: cpuFilter.Min.DeepCopy(),
			Max: cpuFilter.Max.DeepCopy(),
		}
	default:
		return nil, fmt.Errorf("unknown filter type")
	}

	// Parse the Memory filter
	switch k8sSelector.MemoryFilter.FilterType {
	case nodecorev1alpha1.TypeMatchFilter:
		// Unmarshal the data into a ResourceMatchSelector
		var memoryFilter nodecorev1alpha1.ResourceMatchSelector
		err := json.Unmarshal(k8sSelector.MemoryFilter.Data.Raw, &memoryFilter)
		if err != nil {
			return nil, err
		}

		// Generate the model for the Memory filter
		memoryFilterModel = models.ResourceQuantityMatchFilter{
			Value: memoryFilter.Value.DeepCopy(),
		}
	case nodecorev1alpha1.TypeRangeFilter:
		// Unmarshal the data into a ResourceRangeSelector
		var memoryFilter nodecorev1alpha1.ResourceRangeSelector
		err := json.Unmarshal(k8sSelector.MemoryFilter.Data.Raw, &memoryFilter)
		if err != nil {
			return nil, err
		}

		// Generate the model for the Memory filter
		memoryFilterModel = models.ResourceQuantityRangeFilter{
			Min: memoryFilter.Min.DeepCopy(),
			Max: memoryFilter.Max.DeepCopy(),
		}
	default:
		return nil, fmt.Errorf("unknown filter type")
	}

	// Parse the Pods filter
	switch k8sSelector.PodsFilter.FilterType {
	case nodecorev1alpha1.TypeMatchFilter:
		// Unmarshal the data into a ResourceMatchSelector
		var podsFilter nodecorev1alpha1.ResourceMatchSelector
		err := json.Unmarshal(k8sSelector.PodsFilter.Data.Raw, &podsFilter)
		if err != nil {
			return nil, err
		}

		// Generate the model for the Pods filter
		podsFilterModel = models.ResourceQuantityMatchFilter{
			Value: podsFilter.Value.DeepCopy(),
		}
	case nodecorev1alpha1.TypeRangeFilter:
		// Unmarshal the data into a ResourceRangeSelector
		var podsFilter nodecorev1alpha1.ResourceRangeSelector
		err := json.Unmarshal(k8sSelector.PodsFilter.Data.Raw, &podsFilter)
		if err != nil {
			return nil, err
		}

		// Generate the model for the Pods filter
		podsFilterModel = models.ResourceQuantityRangeFilter{
			Min: podsFilter.Min.DeepCopy(),
			Max: podsFilter.Max.DeepCopy(),
		}
	default:
		return nil, fmt.Errorf("unknown filter type")
	}

	// Parse the Storage filter
	switch k8sSelector.StorageFilter.FilterType {
	case nodecorev1alpha1.TypeMatchFilter:
		// Unmarshal the data into a ResourceMatchSelector
		var storageFilter nodecorev1alpha1.ResourceMatchSelector
		err := json.Unmarshal(k8sSelector.StorageFilter.Data.Raw, &storageFilter)
		if err != nil {
			return nil, err
		}

		// Generate the model for the Storage filter
		storageFilterModel = models.ResourceQuantityMatchFilter{
			Value: storageFilter.Value.DeepCopy(),
		}
	case nodecorev1alpha1.TypeRangeFilter:
		// Unmarshal the data into a ResourceRangeSelector
		var storageFilter nodecorev1alpha1.ResourceRangeSelector
		err := json.Unmarshal(k8sSelector.StorageFilter.Data.Raw, &storageFilter)
		if err != nil {
			return nil, err
		}

		// Generate the model for the Storage filter
		storageFilterModel = models.ResourceQuantityRangeFilter{
			Min: storageFilter.Min.DeepCopy(),
			Max: storageFilter.Max.DeepCopy(),
		}
	default:
		return nil, fmt.Errorf("unknown filter type")
	}

	// Generate the model for the K8Slice selector
	k8SliceSelector := models.K8SliceSelector{
		Cpu:     cpuFilterModel,
		Memory:  memoryFilterModel,
		Pods:    podsFilterModel,
		Storage: storageFilterModel,
	}

	return &k8SliceSelector, nil

}

// ParsePartition creates a Partition Object from a Partition CR.
func ParsePartition(partition *nodecorev1alpha1.K8SlicePartition) *models.Partition {
	return &models.Partition{
		CPU:     partition.CPU,
		Memory:  partition.Memory,
		Pods:    partition.Pods,
		Storage: partition.Storage,
		Gpu: models.GpuCharacteristics{
			Model:  partition.Gpu.Model,
			Cores:  partition.Gpu.Cores,
			Memory: partition.Gpu.Memory,
		},
	}
}

// ParsePartitionFromObj creates a Partition CR from a Partition Object.
func ParsePartitionFromObj(partition *models.Partition) *nodecorev1alpha1.K8SlicePartition {
	return &nodecorev1alpha1.K8SlicePartition{
		CPU:    partition.CPU,
		Memory: partition.Memory,
		Pods:   partition.Pods,
		Gpu: &nodecorev1alpha1.GPU{
			Model:  partition.Gpu.Model,
			Cores:  partition.Gpu.Cores,
			Memory: partition.Gpu.Memory,
		},
		Storage: partition.Storage,
	}
}

// ParseNodeIdentity creates a NodeIdentity Object from a NodeIdentity CR.
func ParseNodeIdentity(node nodecorev1alpha1.NodeIdentity) models.NodeIdentity {
	return models.NodeIdentity{
		NodeID: node.NodeID,
		IP:     node.IP,
		Domain: node.Domain,
	}
}

// ParseFlavor creates a Flavor Object from a Flavor CR.
func ParseFlavor(flavor *nodecorev1alpha1.Flavor) *models.Flavor {

	var modelFlavor models.Flavor

	flavorType, flavorTypeStruct, errParse := nodecorev1alpha1.ParseFlavorType(flavor)
	if errParse != nil {
		return nil
	}

	var modelFlavorType models.FlavorType

	switch flavorType {
	case nodecorev1alpha1.Type_K8Slice:
		// Force casting of flavorTypeStruct to K8Slice
		flavorTypeStruct := flavorTypeStruct.(nodecorev1alpha1.K8Slice)
		modelFlavorType = models.K8Slice{
			Name: models.K8SliceNameDefault,
			Characteristics: models.K8SliceCharacteristics{
				Cpu:    flavorTypeStruct.Characteristics.Cpu,
				Memory: flavorTypeStruct.Characteristics.Memory,
				Pods:   flavorTypeStruct.Characteristics.Pods,
				Gpu: models.GpuCharacteristics{
					Model:  flavorTypeStruct.Characteristics.Gpu.Model,
					Cores:  flavorTypeStruct.Characteristics.Gpu.Cores,
					Memory: flavorTypeStruct.Characteristics.Gpu.Memory,
				},
				Storage: flavorTypeStruct.Characteristics.Storage,
			},
			Properties: models.K8SliceProperties{
				Latency:           flavorTypeStruct.Properties.Latency,
				SecurityStandards: flavorTypeStruct.Properties.SecurityStandards,
				CarbonFootprint: models.CarbonFootprint{
					Embodied:    flavorTypeStruct.Properties.CarbonFootprint.Embodied,
					Operational: flavorTypeStruct.Properties.CarbonFootprint.Operational,
				},
			},
			Policies: models.K8SlicePolicies{
				Partitionability: models.K8SlicePartitionability{
					CpuMin:     flavorTypeStruct.Policies.Partitionability.CpuMin,
					MemoryMin:  flavorTypeStruct.Policies.Partitionability.MemoryMin,
					PodsMin:    flavorTypeStruct.Policies.Partitionability.PodsMin,
					CpuStep:    flavorTypeStruct.Policies.Partitionability.CpuStep,
					MemoryStep: flavorTypeStruct.Policies.Partitionability.MemoryStep,
					PodsStep:   flavorTypeStruct.Policies.Partitionability.PodsStep,
				},
			},
		}
	}

	modelFlavor = models.Flavor{
		FlavorID:            flavor.Name,
		Type:                modelFlavorType,
		ProviderID:          flavor.Spec.ProviderID,
		NetworkPropertyType: flavor.Spec.NetworkPropertyType,
		Timestamp:           flavor.CreationTimestamp.Time,
		Location: models.Location{
			Latitude:        flavor.Spec.Location.Latitude,
			Longitude:       flavor.Spec.Location.Longitude,
			Country:         flavor.Spec.Location.Country,
			City:            flavor.Spec.Location.City,
			AdditionalNotes: flavor.Spec.Location.AdditionalNotes,
		},
		Price: models.Price{
			Amount:   flavor.Spec.Price.Amount,
			Currency: flavor.Spec.Price.Currency,
			Period:   flavor.Spec.Price.Period,
		},
		Owner:        ParseNodeIdentity(flavor.Spec.Owner),
		Availability: flavor.Spec.Availability,
	}

	return &modelFlavor
}

// ParseContract creates a Contract Object.
func ParseContract(contract *reservationv1alpha1.Contract) *models.Contract {
	return &models.Contract{
		ContractID:     contract.Name,
		Flavor:         *ParseFlavor(&contract.Spec.Flavor),
		Buyer:          ParseNodeIdentity(contract.Spec.Buyer),
		BuyerClusterID: contract.Spec.BuyerClusterID,
		TransactionID:  contract.Spec.TransactionID,
		Partition: func() *models.Partition {
			if contract.Spec.Partition != nil {
				return ParsePartition(contract.Spec.Partition)
			}
			return nil
		}(),
		Seller: ParseNodeIdentity(contract.Spec.Seller),
		SellerCredentials: models.LiqoCredentials{
			ClusterID:   contract.Spec.SellerCredentials.ClusterID,
			ClusterName: contract.Spec.SellerCredentials.ClusterName,
			Token:       contract.Spec.SellerCredentials.Token,
			Endpoint:    contract.Spec.SellerCredentials.Endpoint,
		},
		ExpirationTime:   contract.Spec.ExpirationTime,
		ExtraInformation: contract.Spec.ExtraInformation,
	}
}

// ParseQuantityFromString parses a string into a resource.Quantity.
func ParseQuantityFromString(s string) resource.Quantity {
	i, err := resource.ParseQuantity(s)
	if err != nil {
		return *resource.NewQuantity(0, resource.DecimalSI)
	}
	return i
}
