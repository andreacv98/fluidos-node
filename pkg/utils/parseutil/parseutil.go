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
func ParseFlavorSelector(selector *nodecorev1alpha1.FlavorSelector) *models.Selector {
	s := &models.Selector{
		FlavorType: selector.FlavorType,
	}

	if selector.MatchSelector != nil {
		s.MatchSelector = &models.MatchSelector{
			CPU:     selector.MatchSelector.CPU,
			Memory:  selector.MatchSelector.Memory,
			Pods:    selector.MatchSelector.Pods,
			Storage: selector.MatchSelector.Storage,
			Gpu: &models.GpuCharacteristics{
				Model:  selector.MatchSelector.Gpu.Model,
				Cores:  selector.MatchSelector.Gpu.Cores,
				Memory: selector.MatchSelector.Gpu.Memory,
			},
		}
	}

	if selector.RangeSelector != nil {
		s.RangeSelector = &models.RangeSelector{
			MinCPU:     selector.RangeSelector.MinCpu,
			MinMemory:  selector.RangeSelector.MinMemory,
			MinPods:    selector.RangeSelector.MinPods,
			MinStorage: selector.RangeSelector.MinStorage,
			MinGpu: &models.GpuCharacteristics{
				Model:  selector.RangeSelector.MinGpu.Model,
				Cores:  selector.RangeSelector.MinGpu.Cores,
				Memory: selector.RangeSelector.MinGpu.Memory,
			},
			MaxCPU:     selector.RangeSelector.MaxCpu,
			MaxMemory:  selector.RangeSelector.MaxMemory,
			MaxPods:    selector.RangeSelector.MaxPods,
			MaxStorage: selector.RangeSelector.MaxStorage,
			MaxGpu: &models.GpuCharacteristics{
				Model:  selector.RangeSelector.MaxGpu.Model,
				Cores:  selector.RangeSelector.MaxGpu.Cores,
				Memory: selector.RangeSelector.MaxGpu.Memory,
			},
		}
	}

	return s
}

// ParsePartition creates a Partition Object from a Partition CR.
func ParsePartition(partition *nodecorev1alpha1.Partition) *models.Partition {
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
func ParsePartitionFromObj(partition *models.Partition) *nodecorev1alpha1.Partition {
	return &nodecorev1alpha1.Partition{
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

	errParse, flavorType, flavorTypeStruct := ParseFlavorType(flavor)
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
		Timestamp:           flavor.Spec.Timestamp.Time,
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

// ParseFlavorType parses a Flavor into a the type and the unmarsheled raw value.
func ParseFlavorType(flavor *nodecorev1alpha1.Flavor) (error, nodecorev1alpha1.FlavorTypeIdentifier, interface{}) {

	var validationErr error

	switch flavor.Spec.Type.TypeIdentifier {

	case nodecorev1alpha1.Type_K8Slice:

		var k8slice nodecorev1alpha1.K8Slice
		validationErr = json.Unmarshal(flavor.Spec.Type.TypeData.Raw, &k8slice)
		return validationErr, nodecorev1alpha1.Type_K8Slice, k8slice

	case nodecorev1alpha1.Type_VM:
		// TODO: Implement VM flavor parsing
		return fmt.Errorf("flavor type %s not supported", flavor.Spec.Type.TypeIdentifier), "", nil

	case nodecorev1alpha1.Type_Service:
		// TODO: Implement Service flavor parsing
		return fmt.Errorf("flavor type %s not supported", flavor.Spec.Type.TypeIdentifier), "", nil

	default:
		return fmt.Errorf("flavor type %s not supported", flavor.Spec.Type.TypeIdentifier), "", nil
	}
}
