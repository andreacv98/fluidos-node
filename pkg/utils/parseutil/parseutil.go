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
	"k8s.io/klog/v2"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	reservationv1alpha1 "github.com/fluidos-project/node/apis/reservation/v1alpha1"
	"github.com/fluidos-project/node/pkg/utils/models"
)

// ParseFlavorSelector parses FlavorSelector into a Selector.
func ParseFlavorSelector(selector *nodecorev1alpha1.Selector) (models.Selector, error) {
	// Parse the Selector
	klog.Infof("Parsing the selector %s", selector.SelectorTypeIdentifier)
	selectorIdentifier, selectorStruct, err := nodecorev1alpha1.ParseSolverSelector(selector)
	if err != nil {
		return nil, err
	}

	klog.Infof("Selector type: %s", selectorIdentifier)

	switch selectorIdentifier {
	case nodecorev1alpha1.Type_K8Slice:
		// Force casting of selectorStruct to K8Slice
		selectorStruct := selectorStruct.(nodecorev1alpha1.K8SliceSelector)
		klog.Info("Forced casting of selectorStruct to K8Slice")
		// Print the selectorStruct
		klog.Infof("SelectorStruct: %v", selectorStruct)
		// Generate the model for the K8Slice selector
		k8SliceSelector, err := parseK8SliceFilters(&selectorStruct)
		if err != nil {
			return nil, err
		}

		klog.Infof("K8SliceSelector: %v", k8SliceSelector)

		return *k8SliceSelector, nil

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

func parseK8SliceFilters(k8sSelector *nodecorev1alpha1.K8SliceSelector) (*models.K8SliceSelector, error) {

	var cpuFilterModel, memoryFilterModel, podsFilterModel, storageFilterModel models.ResourceQuantityFilter

	// Parse the CPU filter
	if k8sSelector.CpuFilter != nil {
		klog.Info("Parsing the CPU filter")
		switch k8sSelector.CpuFilter.Name {
		case nodecorev1alpha1.TypeMatchFilter:
			klog.Info("Parsing the CPU filter as a MatchFilter")
			// Unmarshal the data into a ResourceMatchSelector
			var cpuFilter nodecorev1alpha1.ResourceMatchSelector
			err := json.Unmarshal(k8sSelector.CpuFilter.Data.Raw, &cpuFilter)
			if err != nil {
				return nil, err
			}

			cpuFilterData := models.ResourceQuantityMatchFilter{
				Value: cpuFilter.Value.DeepCopy(),
			}
			// Marshal the CPU filter data into JSON
			cpuFilterDataJSON, err := json.Marshal(cpuFilterData)
			if err != nil {
				return nil, err
			}

			// Generate the model for the CPU filter
			cpuFilterModel = models.ResourceQuantityFilter{
				Name:	models.MatchFilter,
				Data:	cpuFilterDataJSON,
			}
			klog.Infof("CPU filter model: %v", cpuFilterModel)
		case nodecorev1alpha1.TypeRangeFilter:
			klog.Info("Parsing the CPU filter as a RangeFilter")
			// Unmarshal the data into a ResourceRangeSelector
			var cpuFilter nodecorev1alpha1.ResourceRangeSelector
			err := json.Unmarshal(k8sSelector.CpuFilter.Data.Raw, &cpuFilter)
			if err != nil {
				return nil, err
			}

			var cpuFilterMinCopy *resource.Quantity
			if cpuFilter.Min != nil {
				cpuFilterMinCopyData := cpuFilter.Min.DeepCopy()
				cpuFilterMinCopy = &cpuFilterMinCopyData
			}
			cpuFilterMinCopy = nil

			var cpuFilterMaxCopy *resource.Quantity
			if cpuFilter.Max != nil {
				cpuFilterMaxCopyData := cpuFilter.Max.DeepCopy()
				cpuFilterMaxCopy = &cpuFilterMaxCopyData
			}
			cpuFilterMaxCopy = nil

			cpuFilterData := models.ResourceQuantityRangeFilter{
				Min: cpuFilterMinCopy,
				Max: cpuFilterMaxCopy,
			}
			// Marshal the CPU filter data into JSON
			cpuFilterDataJSON, err := json.Marshal(cpuFilterData)
			if err != nil {
				return nil, err
			}

			// Generate the model for the CPU filter
			cpuFilterModel = models.ResourceQuantityFilter{
				Name:	models.RangeFilter,
				Data:	cpuFilterDataJSON,
			}
			klog.Infof("CPU filter model: %v", cpuFilterModel)
		default:
			return nil, fmt.Errorf("unknown filter type")
		}
	}

	if k8sSelector.MemoryFilter != nil {
		klog.Info("Parsing the Memory filter")
		// Parse the Memory filter
		switch k8sSelector.MemoryFilter.Name {
		case nodecorev1alpha1.TypeMatchFilter:
			klog.Info("Parsing the Memory filter as a MatchFilter")
			// Unmarshal the data into a ResourceMatchSelector
			var memoryFilter nodecorev1alpha1.ResourceMatchSelector
			err := json.Unmarshal(k8sSelector.MemoryFilter.Data.Raw, &memoryFilter)
			if err != nil {
				return nil, err
			}

			memoryFilterData := models.ResourceQuantityMatchFilter{
				Value: memoryFilter.Value.DeepCopy(),
			}
			// Marshal the Memory filter data into JSON
			memoryFilterDataJSON, err := json.Marshal(memoryFilterData)
			if err != nil {
				return nil, err
			}

			// Generate the model for the Memory filter
			memoryFilterModel = models.ResourceQuantityFilter{
				Name:	models.MatchFilter,
				Data:	memoryFilterDataJSON,
			}
			klog.Infof("Memory filter model: %v", memoryFilterModel)
		case nodecorev1alpha1.TypeRangeFilter:
			klog.Info("Parsing the Memory filter as a RangeFilter")
			// Unmarshal the data into a ResourceRangeSelector
			var memoryFilter nodecorev1alpha1.ResourceRangeSelector
			err := json.Unmarshal(k8sSelector.MemoryFilter.Data.Raw, &memoryFilter)
			if err != nil {
				return nil, err
			}

			var memoryFilterMinCopy *resource.Quantity
			if memoryFilter.Min != nil {
				memoryFilterMinCopyData := memoryFilter.Min.DeepCopy()
				memoryFilterMinCopy = &memoryFilterMinCopyData
			}
			memoryFilterMinCopy = nil

			var memoryFilterMaxCopy *resource.Quantity
			if memoryFilter.Max != nil {
				memoryFilterMaxCopyData := memoryFilter.Max.DeepCopy()
				memoryFilterMaxCopy = &memoryFilterMaxCopyData
			}
			memoryFilterMaxCopy = nil

			memoryFilterData := models.ResourceQuantityRangeFilter{
				Min: memoryFilterMinCopy,
				Max: memoryFilterMaxCopy,
			}
			// Marshal the Memory filter data into JSON
			memoryFilterDataJSON, err := json.Marshal(memoryFilterData)
			if err != nil {
				return nil, err
			}

			// Generate the model for the Memory filter
			memoryFilterModel = models.ResourceQuantityFilter{
				Name:	models.RangeFilter,
				Data:	memoryFilterDataJSON,
			}
			klog.Infof("Memory filter model: %v", memoryFilterModel)
		default:
			return nil, fmt.Errorf("unknown filter type")
		}
	}

	if k8sSelector.PodsFilter != nil {
		klog.Info("Parsing the Pods filter")
		// Parse the Pods filter
		switch k8sSelector.PodsFilter.Name {
		case nodecorev1alpha1.TypeMatchFilter:
			klog.Info("Parsing the Pods filter as a MatchFilter")
			// Unmarshal the data into a ResourceMatchSelector
			var podsFilter nodecorev1alpha1.ResourceMatchSelector
			err := json.Unmarshal(k8sSelector.PodsFilter.Data.Raw, &podsFilter)
			if err != nil {
				return nil, err
			}

			podsFilterData := models.ResourceQuantityMatchFilter{
				Value: podsFilter.Value.DeepCopy(),
			}
			// Marshal the Pods filter data into JSON
			podsFilterDataJSON, err := json.Marshal(podsFilterData)
			if err != nil {
				return nil, err
			}

			// Generate the model for the Pods filter
			podsFilterModel = models.ResourceQuantityFilter{
				Name:	models.MatchFilter,
				Data:	podsFilterDataJSON,
			}
			klog.Infof("Pods filter model: %v", podsFilterModel)
		case nodecorev1alpha1.TypeRangeFilter:
			klog.Info("Parsing the Pods filter as a RangeFilter")
			// Unmarshal the data into a ResourceRangeSelector
			var podsFilter nodecorev1alpha1.ResourceRangeSelector
			err := json.Unmarshal(k8sSelector.PodsFilter.Data.Raw, &podsFilter)
			if err != nil {
				return nil, err
			}

			var podsFilterMinCopy *resource.Quantity
			if podsFilter.Min != nil {
				podsFilterMinCopyData := podsFilter.Min.DeepCopy()
				podsFilterMinCopy = &podsFilterMinCopyData
			}
			podsFilterMinCopy = nil

			var podsFilterMaxCopy *resource.Quantity
			if podsFilter.Max != nil {
				podsFilterMaxCopyData := podsFilter.Max.DeepCopy()
				podsFilterMaxCopy = &podsFilterMaxCopyData
			}
			podsFilterMaxCopy = nil

			podsFilterData := models.ResourceQuantityRangeFilter{
				Min: podsFilterMinCopy,
				Max: podsFilterMaxCopy,
			}
			// Marshal the Pods filter data into JSON
			podsFilterDataJSON, err := json.Marshal(podsFilterData)
			if err != nil {
				return nil, err
			}

			// Generate the model for the Pods filter
			podsFilterModel = models.ResourceQuantityFilter{
				Name:	models.RangeFilter,
				Data:	podsFilterDataJSON,
			}
			klog.Infof("Pods filter model: %v", podsFilterModel)
		default:
			return nil, fmt.Errorf("unknown filter type")
		}
	}

	if k8sSelector.StorageFilter != nil {
		klog.Info("Parsing the Storage filter")
		// Parse the Storage filter
		switch k8sSelector.StorageFilter.Name {
		case nodecorev1alpha1.TypeMatchFilter:
			klog.Info("Parsing the Storage filter as a MatchFilter")
			// Unmarshal the data into a ResourceMatchSelector
			var storageFilter nodecorev1alpha1.ResourceMatchSelector
			err := json.Unmarshal(k8sSelector.StorageFilter.Data.Raw, &storageFilter)
			if err != nil {
				return nil, err
			}

			storageFilterData := models.ResourceQuantityMatchFilter{
				Value: storageFilter.Value.DeepCopy(),
			}
			// Marshal the Storage filter data into JSON
			storageFilterDataJSON, err := json.Marshal(storageFilterData)
			if err != nil {
				return nil, err
			}

			// Generate the model for the Storage filter
			storageFilterModel = models.ResourceQuantityFilter{
				Name:	models.MatchFilter,
				Data:	storageFilterDataJSON,
			}
			klog.Infof("Storage filter model: %v", storageFilterModel)
		case nodecorev1alpha1.TypeRangeFilter:
			klog.Info("Parsing the Storage filter as a RangeFilter")
			// Unmarshal the data into a ResourceRangeSelector
			var storageFilter nodecorev1alpha1.ResourceRangeSelector
			err := json.Unmarshal(k8sSelector.StorageFilter.Data.Raw, &storageFilter)
			if err != nil {
				return nil, err
			}

			var storageFilterMinCopy *resource.Quantity
			if storageFilter.Min != nil {
				storageFilterMinCopyData := storageFilter.Min.DeepCopy()
				storageFilterMinCopy = &storageFilterMinCopyData
			}
			storageFilterMinCopy = nil

			var storageFilterMaxCopy *resource.Quantity
			if storageFilter.Max != nil {
				storageFilterMaxCopyData := storageFilter.Max.DeepCopy()
				storageFilterMaxCopy = &storageFilterMaxCopyData
			}
			storageFilterMaxCopy = nil

			storageFilterData := models.ResourceQuantityRangeFilter{
				Min: storageFilterMinCopy,
				Max: storageFilterMaxCopy,
			}
			// Marshal the Storage filter data into JSON
			storageFilterDataJSON, err := json.Marshal(storageFilterData)
			if err != nil {
				return nil, err
			}

			// Generate the model for the Storage filter
			storageFilterModel = models.ResourceQuantityFilter{
				Name:	models.RangeFilter,
				Data:	storageFilterDataJSON,
			}
			klog.Infof("Storage filter model: %v", storageFilterModel)
		default:
			return nil, fmt.Errorf("unknown filter type")
		}
	}

	// Generate the model for the K8Slice selector
	k8SliceSelector := models.K8SliceSelector{
		Cpu: func() *models.ResourceQuantityFilter {
			if k8sSelector.CpuFilter != nil {
				return &cpuFilterModel
			} else {
				return nil
			}
		}(),
		Memory: func() *models.ResourceQuantityFilter {
			if k8sSelector.MemoryFilter != nil {
				return &memoryFilterModel
			} else {
				return nil
			}
		}(),
		Pods: func() *models.ResourceQuantityFilter {
			if k8sSelector.PodsFilter != nil {
				return &podsFilterModel
			} else {
				return nil
			}
		}(),
		Storage: func() *models.ResourceQuantityFilter {
			if k8sSelector.StorageFilter != nil {
				return &storageFilterModel
			} else {
				return nil
			}
		}(),
	}

	return &k8SliceSelector, nil
}

// ParsePartition creates a Partition Object from a Partition CR.
func ParsePartition(partition *nodecorev1alpha1.Partition) *models.Partition {
	// Parse the Partition
	partitionType, partitionStruct, err := nodecorev1alpha1.ParsePartition(partition)
	if err != nil {
		return nil
	}

	switch partitionType {
	case nodecorev1alpha1.Type_K8Slice:
		// Force casting of partitionStruct to K8Slice
		partitionStruct := partitionStruct.(nodecorev1alpha1.K8SlicePartition)
		k8slicePartitionJSON := models.K8SlicePartition{
			CPU:    partitionStruct.CPU,
			Memory: partitionStruct.Memory,
			Pods:   partitionStruct.Pods,
			Gpu: func() *models.GpuCharacteristics {
				if partitionStruct.Gpu != nil {
					return &models.GpuCharacteristics{
						Model:  partitionStruct.Gpu.Model,
						Cores:  partitionStruct.Gpu.Cores,
						Memory: partitionStruct.Gpu.Memory,
					}
				}
				return nil
			}(),
			Storage: partitionStruct.Storage,
		}
		// Marshal the K8Slice partition to JSON
		partitionData, err := json.Marshal(k8slicePartitionJSON)
		if err != nil {
			klog.Errorf("Error marshalling the K8Slice partition: %s", err)
			return nil
		}
		return &models.Partition{
			Name: models.K8SliceNameDefault,
			Data: partitionData,
		}
		// TODO: Implement the other partition types, if any
	default:
		return nil
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
		modelFlavorTypeData := models.K8Slice{
			Characteristics: models.K8SliceCharacteristics{
				Cpu:    flavorTypeStruct.Characteristics.Cpu,
				Memory: flavorTypeStruct.Characteristics.Memory,
				Pods:   flavorTypeStruct.Characteristics.Pods,
				Gpu: func() *models.GpuCharacteristics {
					if flavorTypeStruct.Characteristics.Gpu != nil {
						return &models.GpuCharacteristics{
							Model:  flavorTypeStruct.Characteristics.Gpu.Model,
							Cores:  flavorTypeStruct.Characteristics.Gpu.Cores,
							Memory: flavorTypeStruct.Characteristics.Gpu.Memory,
						}
					}
					return nil
				}(),
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

		// Encode the K8Slice data into JSON
		encodedFlavorTypeData, err := json.Marshal(modelFlavorTypeData)
		if err != nil {
			klog.Errorf("Error encoding the K8Slice data: %s", err)
			return nil
		}

		modelFlavorType = models.FlavorType{
			Name: models.K8SliceNameDefault,
			Data: encodedFlavorTypeData,
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
				partition := ParsePartition(contract.Spec.Partition)
				return partition
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
