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
	"k8s.io/apimachinery/pkg/api/resource"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	reservationv1alpha1 "github.com/fluidos-project/node/apis/reservation/v1alpha1"
	"github.com/fluidos-project/node/pkg/utils/models"
)

// ParseFlavorSelector parses FlavorSelector into a Selector.
func ParseFlavorSelector(selector *nodecorev1alpha1.FlavorSelector) *models.Selector {
	s := &models.Selector{
		Architecture: selector.Architecture,
		FlavorType:   selector.FlavorType,
	}

	if selector.MatchSelector != nil {
		s.MatchSelector = &models.MatchSelector{
			CPU:              selector.MatchSelector.CPU,
			Memory:           selector.MatchSelector.Memory,
			Pods:             selector.MatchSelector.Pods,
			EphemeralStorage: selector.MatchSelector.EphemeralStorage,
			Storage:          selector.MatchSelector.Storage,
			Gpu:              selector.MatchSelector.Gpu,
		}
	}

	if selector.RangeSelector != nil {
		s.RangeSelector = &models.RangeSelector{
			MinCPU:     selector.RangeSelector.MinCpu,
			MinMemory:  selector.RangeSelector.MinMemory,
			MinPods:    selector.RangeSelector.MinPods,
			MinEph:     selector.RangeSelector.MinEph,
			MinStorage: selector.RangeSelector.MinStorage,
			MinGpu:     selector.RangeSelector.MinGpu,
			MaxCPU:     selector.RangeSelector.MaxCpu,
			MaxMemory:  selector.RangeSelector.MaxMemory,
			MaxPods:    selector.RangeSelector.MaxPods,
			MaxEph:     selector.RangeSelector.MaxEph,
			MaxStorage: selector.RangeSelector.MaxStorage,
			MaxGpu:     selector.RangeSelector.MaxGpu,
		}
	}

	return s
}

// ParsePartition creates a Partition Object from a Partition CR.
func ParsePartition(partition *nodecorev1alpha1.Partition) *models.Partition {
	return &models.Partition{
		CPU:              partition.CPU,
		Memory:           partition.Memory,
		Pods:             partition.Pods,
		EphemeralStorage: partition.EphemeralStorage,
		Storage:          partition.Storage,
		Gpu:              partition.Gpu,
	}
}

// ParsePartitionFromObj creates a Partition CR from a Partition Object.
func ParsePartitionFromObj(partition *models.Partition) *nodecorev1alpha1.Partition {
	return &nodecorev1alpha1.Partition{
		Architecture:     partition.Architecture,
		CPU:              partition.CPU,
		Memory:           partition.Memory,
		Pods:             partition.Pods,
		Gpu:              partition.Gpu,
		Storage:          partition.Storage,
		EphemeralStorage: partition.EphemeralStorage,
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
	return &models.Flavor{
		FlavorID:   flavor.Name,
		Type:       string(flavor.Spec.Type),
		ProviderID: flavor.Spec.ProviderID,
		Characteristics: models.Characteristics{
			Architecture:      flavor.Spec.Characteristics.Architecture,
			CPU:               flavor.Spec.Characteristics.Cpu,
			Memory:            flavor.Spec.Characteristics.Memory,
			Pods:              flavor.Spec.Characteristics.Pods,
			PersistentStorage: flavor.Spec.Characteristics.PersistentStorage,
			EphemeralStorage:  flavor.Spec.Characteristics.EphemeralStorage,
			Gpu:               flavor.Spec.Characteristics.Gpu,
		},
		Owner: ParseNodeIdentity(flavor.Spec.Owner),
		Policy: models.Policy{
			Partitionable: func() *models.Partitionable {
				if flavor.Spec.Policy.Partitionable != nil {
					return &models.Partitionable{
						CPUMinimum:    flavor.Spec.Policy.Partitionable.CpuMin,
						MemoryMinimum: flavor.Spec.Policy.Partitionable.MemoryMin,
						PodsMinimum:   flavor.Spec.Policy.Partitionable.PodsMin,
						CPUStep:       flavor.Spec.Policy.Partitionable.CpuStep,
						MemoryStep:    flavor.Spec.Policy.Partitionable.MemoryStep,
						PodsStep:      flavor.Spec.Policy.Partitionable.PodsStep,
					}
				}
				return nil
			}(),
			Aggregatable: func() *models.Aggregatable {
				if flavor.Spec.Policy.Aggregatable != nil {
					return &models.Aggregatable{
						MinCount: flavor.Spec.Policy.Aggregatable.MinCount,
						MaxCount: flavor.Spec.Policy.Aggregatable.MaxCount,
					}
				}
				return nil
			}(),
		},
		Price: models.Price{
			Amount:   flavor.Spec.Price.Amount,
			Currency: flavor.Spec.Price.Currency,
			Period:   flavor.Spec.Price.Period,
		},
		OptionalFields: models.OptionalFields{
			Availability: flavor.Spec.OptionalFields.Availability,
			WorkerID:     flavor.Spec.OptionalFields.WorkerID,
		},
	}
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
