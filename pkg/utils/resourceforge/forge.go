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

package resourceforge

import (
	"encoding/json"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"

	advertisementv1alpha1 "github.com/fluidos-project/node/apis/advertisement/v1alpha1"
	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	reservationv1alpha1 "github.com/fluidos-project/node/apis/reservation/v1alpha1"
	"github.com/fluidos-project/node/pkg/utils/flags"
	"github.com/fluidos-project/node/pkg/utils/models"
	"github.com/fluidos-project/node/pkg/utils/namings"
	"github.com/fluidos-project/node/pkg/utils/parseutil"
	"github.com/fluidos-project/node/pkg/utils/tools"
)

// ForgeDiscovery creates a Discovery CR from a FlavorSelector and a solverID.
func ForgeDiscovery(selector *nodecorev1alpha1.FlavorSelector, solverID string) *advertisementv1alpha1.Discovery {
	return &advertisementv1alpha1.Discovery{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namings.ForgeDiscoveryName(solverID),
			Namespace: flags.FluidoNamespace,
		},
		Spec: advertisementv1alpha1.DiscoverySpec{
			Selector: func() *nodecorev1alpha1.FlavorSelector {
				if selector != nil {
					return selector
				}
				return nil
			}(),
			SolverID:  solverID,
			Subscribe: false,
		},
	}
}

// ForgePeeringCandidate creates a PeeringCandidate CR from a Flavor and a Discovery.
func ForgePeeringCandidate(flavorPeeringCandidate *nodecorev1alpha1.Flavor,
	solverID string, available bool) (pc *advertisementv1alpha1.PeeringCandidate) {
	pc = &advertisementv1alpha1.PeeringCandidate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namings.ForgePeeringCandidateName(flavorPeeringCandidate.Name),
			Namespace: flags.FluidoNamespace,
		},
		Spec: advertisementv1alpha1.PeeringCandidateSpec{
			Flavor: nodecorev1alpha1.Flavor{
				ObjectMeta: metav1.ObjectMeta{
					Name:      flavorPeeringCandidate.Name,
					Namespace: flavorPeeringCandidate.Namespace,
				},
				Spec: flavorPeeringCandidate.Spec,
			},
			Available: available,
		},
	}
	pc.Spec.SolverID = solverID
	return
}

// ForgeReservation creates a Reservation CR from a PeeringCandidate.
func ForgeReservation(pc *advertisementv1alpha1.PeeringCandidate,
	partition *nodecorev1alpha1.Partition,
	ni nodecorev1alpha1.NodeIdentity) *reservationv1alpha1.Reservation {
	solverID := pc.Spec.SolverID
	reservation := &reservationv1alpha1.Reservation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namings.ForgeReservationName(solverID),
			Namespace: flags.FluidoNamespace,
		},
		Spec: reservationv1alpha1.ReservationSpec{
			SolverID: solverID,
			Buyer:    ni,
			Seller: nodecorev1alpha1.NodeIdentity{
				Domain: pc.Spec.Flavor.Spec.Owner.Domain,
				NodeID: pc.Spec.Flavor.Spec.Owner.NodeID,
				IP:     pc.Spec.Flavor.Spec.Owner.IP,
			},
			PeeringCandidate: nodecorev1alpha1.GenericRef{
				Name:      pc.Name,
				Namespace: pc.Namespace,
			},
			Reserve:  true,
			Purchase: true,
			Partition: func() *nodecorev1alpha1.Partition {
				if partition != nil {
					return partition
				}
				return nil
			}(),
		},
	}
	if partition != nil {
		reservation.Spec.Partition = partition
	}
	return reservation
}

// ForgeContract creates a Contract CR.
func ForgeContract(flavor *nodecorev1alpha1.Flavor, transaction *models.Transaction,
	lc *nodecorev1alpha1.LiqoCredentials) *reservationv1alpha1.Contract {
	return &reservationv1alpha1.Contract{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namings.ForgeContractName(flavor.Name),
			Namespace: flags.FluidoNamespace,
		},
		Spec: reservationv1alpha1.ContractSpec{
			Flavor: *flavor,
			Buyer: nodecorev1alpha1.NodeIdentity{
				Domain: transaction.Buyer.Domain,
				IP:     transaction.Buyer.IP,
				NodeID: transaction.Buyer.NodeID,
			},
			BuyerClusterID:    transaction.ClusterID,
			Seller:            flavor.Spec.Owner,
			SellerCredentials: *lc,
			TransactionID:     transaction.TransactionID,
			Partition: func() *nodecorev1alpha1.Partition {
				if transaction.Partition != nil {
					return parseutil.ParsePartitionFromObj(transaction.Partition)
				}
				return nil
			}(),
			ExpirationTime:   time.Now().Add(flags.ExpirationContract).Format(time.RFC3339),
			ExtraInformation: nil,
		},
		Status: reservationv1alpha1.ContractStatus{
			Phase: nodecorev1alpha1.PhaseStatus{
				Phase:     nodecorev1alpha1.PhaseActive,
				StartTime: tools.GetTimeNow(),
			},
		},
	}
}

// ForgeFlavorFromMetrics creates a new flavor custom resource from the metrics of the node.
func ForgeFlavorFromMetrics(node *models.NodeInfo, ni nodecorev1alpha1.NodeIdentity, ownerReferences []metav1.OwnerReference) (flavor *nodecorev1alpha1.Flavor) {

	k8SliceType := nodecorev1alpha1.K8Slice{
		Characteristics: nodecorev1alpha1.K8SliceCharacteristics{
			Cpu:     node.ResourceMetrics.CPUAvailable,
			Memory:  node.ResourceMetrics.MemoryAvailable,
			Pods:    node.ResourceMetrics.PodsAvailable,
			Storage: node.ResourceMetrics.EphemeralStorage,
			Gpu: &nodecorev1alpha1.GPU{
				Model:  node.ResourceMetrics.GPU.Model,
				Cores:  node.ResourceMetrics.GPU.CoresAvailable,
				Memory: node.ResourceMetrics.GPU.MemoryAvailable,
			},
		},
		Properties: nodecorev1alpha1.Properties{},
		Policies: nodecorev1alpha1.Policies{
			Partitionability: nodecorev1alpha1.Partitionability{
				CpuMin:     parseutil.ParseQuantityFromString(flags.CPUMin),
				MemoryMin:  parseutil.ParseQuantityFromString(flags.MemoryMin),
				PodsMin:    parseutil.ParseQuantityFromString(flags.PodsMin),
				CpuStep:    parseutil.ParseQuantityFromString(flags.CPUStep),
				MemoryStep: parseutil.ParseQuantityFromString(flags.MemoryStep),
				PodsStep:   parseutil.ParseQuantityFromString(flags.PodsStep),
			},
		},
	}

	// Serialize K8SliceType to JSON
	k8SliceTypeJSON, err := json.Marshal(k8SliceType)
	if err != nil {
		klog.Errorf("Error when marshalling K8SliceType: %s", err)
		return nil
	}

	return &nodecorev1alpha1.Flavor{
		ObjectMeta: metav1.ObjectMeta{
			Name:            namings.ForgeFlavorName("", ni.Domain),
			Namespace:       flags.FluidoNamespace,
			OwnerReferences: ownerReferences,
		},
		Spec: nodecorev1alpha1.FlavorSpec{
			ProviderID: ni.NodeID,
			FlavorType: nodecorev1alpha1.FlavorType{
				TypeIdentifier: nodecorev1alpha1.Type_K8Slice,
				TypeData:       runtime.RawExtension{Raw: k8SliceTypeJSON},
			},
			Owner: ni,
			Price: nodecorev1alpha1.Price{
				Amount:   flags.AMOUNT,
				Currency: flags.CURRENCY,
				Period:   flags.PERIOD,
			},
			Availability: true,
			// TODO: test without network property and location
			NetworkPropertyType: "networkProperty",
			Location: &nodecorev1alpha1.Location{
				Latitude:        "10",
				Longitude:       "58",
				Country:         "Italy",
				City:            "Turin",
				AdditionalNotes: "None",
			},
		},
	}
}

// ForgeFlavorFromRef creates a new flavor starting from a Reference Flavor and the new Characteristics.
func ForgeFlavorFromRef(f *nodecorev1alpha1.Flavor, flavorType *nodecorev1alpha1.FlavorType) (flavor *nodecorev1alpha1.Flavor) {
	return &nodecorev1alpha1.Flavor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namings.ForgeFlavorName(string(f.Spec.FlavorType.TypeIdentifier), f.Spec.Owner.Domain),
			Namespace: flags.FluidoNamespace,
		},
		Spec: nodecorev1alpha1.FlavorSpec{
			ProviderID:          f.Spec.ProviderID,
			FlavorType:          *flavorType,
			Owner:               f.Spec.Owner,
			Price:               f.Spec.Price,
			Availability:        true,
			NetworkPropertyType: f.Spec.NetworkPropertyType,
			Location:            f.Spec.Location,
		},
	}
}

// FORGER FUNCTIONS FROM OBJECTS

// ForgeTransactionObj creates a new Transaction object.
func ForgeTransactionObj(id string, req *models.ReserveRequest) *models.Transaction {
	return &models.Transaction{
		TransactionID: id,
		Buyer:         req.Buyer,
		ClusterID:     req.ClusterID,
		FlavorID:      req.FlavorID,
		Partition: func() *models.Partition {
			if req.Partition != nil {
				return req.Partition
			}
			return nil
		}(),
		StartTime: tools.GetTimeNow(),
	}
}

// ForgeContractObj creates a new Contract object.
func ForgeContractObj(contract *reservationv1alpha1.Contract) models.Contract {
	return models.Contract{
		ContractID:     contract.Name,
		Flavor:         *parseutil.ParseFlavor(&contract.Spec.Flavor),
		Buyer:          parseutil.ParseNodeIdentity(contract.Spec.Buyer),
		BuyerClusterID: contract.Spec.BuyerClusterID,
		Seller:         parseutil.ParseNodeIdentity(contract.Spec.Seller),
		SellerCredentials: models.LiqoCredentials{
			ClusterID:   contract.Spec.SellerCredentials.ClusterID,
			ClusterName: contract.Spec.SellerCredentials.ClusterName,
			Token:       contract.Spec.SellerCredentials.Token,
			Endpoint:    contract.Spec.SellerCredentials.Endpoint,
		},
		Partition: func() *models.Partition {
			if contract.Spec.Partition != nil {
				return parseutil.ParsePartition(contract.Spec.Partition)
			}
			return nil
		}(),
		TransactionID:  contract.Spec.TransactionID,
		ExpirationTime: contract.Spec.ExpirationTime,
		ExtraInformation: func() map[string]string {
			if contract.Spec.ExtraInformation != nil {
				return contract.Spec.ExtraInformation
			}
			return nil
		}(),
	}
}

// ForgeResponsePurchaseObj creates a new response purchase.
func ForgeResponsePurchaseObj(contract *models.Contract) *models.ResponsePurchase {
	return &models.ResponsePurchase{
		Contract: *contract,
		Status:   "Completed",
	}
}

// ForgeContractFromObj creates a Contract from a reservation.
func ForgeContractFromObj(contract *models.Contract) *reservationv1alpha1.Contract {
	return &reservationv1alpha1.Contract{
		ObjectMeta: metav1.ObjectMeta{
			Name:      contract.ContractID,
			Namespace: flags.FluidoNamespace,
		},
		Spec: reservationv1alpha1.ContractSpec{
			Flavor: *ForgeFlavorFromObj(&contract.Flavor),
			Buyer: nodecorev1alpha1.NodeIdentity{
				Domain: contract.Buyer.Domain,
				IP:     contract.Buyer.IP,
				NodeID: contract.Buyer.NodeID,
			},
			BuyerClusterID: contract.BuyerClusterID,
			Seller: nodecorev1alpha1.NodeIdentity{
				NodeID: contract.Seller.NodeID,
				IP:     contract.Seller.IP,
				Domain: contract.Seller.Domain,
			},
			SellerCredentials: nodecorev1alpha1.LiqoCredentials{
				ClusterID:   contract.SellerCredentials.ClusterID,
				ClusterName: contract.SellerCredentials.ClusterName,
				Token:       contract.SellerCredentials.Token,
				Endpoint:    contract.SellerCredentials.Endpoint,
			},
			TransactionID: contract.TransactionID,
			Partition: func() *nodecorev1alpha1.Partition {
				if contract.Partition != nil {
					return parseutil.ParsePartitionFromObj(contract.Partition)
				}
				return nil
			}(),
			ExpirationTime: contract.ExpirationTime,
			ExtraInformation: func() map[string]string {
				if contract.ExtraInformation != nil {
					return contract.ExtraInformation
				}
				return nil
			}(),
		},
		Status: reservationv1alpha1.ContractStatus{
			Phase: nodecorev1alpha1.PhaseStatus{
				Phase:     nodecorev1alpha1.PhaseActive,
				StartTime: tools.GetTimeNow(),
			},
		},
	}
}

// ForgeTransactionFromObj creates a transaction from a Transaction object.
func ForgeTransactionFromObj(transaction *models.Transaction) *reservationv1alpha1.Transaction {
	return &reservationv1alpha1.Transaction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      transaction.TransactionID,
			Namespace: flags.FluidoNamespace,
		},
		Spec: reservationv1alpha1.TransactionSpec{
			FlavorID:  transaction.FlavorID,
			StartTime: transaction.StartTime,
			Buyer: nodecorev1alpha1.NodeIdentity{
				Domain: transaction.Buyer.Domain,
				IP:     transaction.Buyer.IP,
				NodeID: transaction.Buyer.NodeID,
			},
			ClusterID: transaction.ClusterID,
			Partition: func() *nodecorev1alpha1.Partition {
				if transaction.Partition != nil {
					return parseutil.ParsePartitionFromObj(transaction.Partition)
				}
				return nil
			}(),
		},
	}
}

// ForgeFlavorFromObj creates a Flavor CR from a Flavor Object (REAR).
func ForgeFlavorFromObj(flavor *models.Flavor) *nodecorev1alpha1.Flavor {

	var flavorType nodecorev1alpha1.FlavorType

	switch flavor.Type.GetFlavorTypeName() {
	case models.K8SliceNameDefault:
		flavorTypeData := nodecorev1alpha1.K8Slice{
			Characteristics: nodecorev1alpha1.K8SliceCharacteristics{
				Cpu:     flavor.Type.(models.K8Slice).Characteristics.Cpu,
				Memory:  flavor.Type.(models.K8Slice).Characteristics.Memory,
				Pods:    flavor.Type.(models.K8Slice).Characteristics.Pods,
				Storage: flavor.Type.(models.K8Slice).Characteristics.Storage,
				Gpu: &nodecorev1alpha1.GPU{
					Model:  flavor.Type.(models.K8Slice).Characteristics.Gpu.Model,
					Cores:  flavor.Type.(models.K8Slice).Characteristics.Gpu.Cores,
					Memory: flavor.Type.(models.K8Slice).Characteristics.Gpu.Memory,
				},
			},
			Properties: nodecorev1alpha1.Properties{
				Latency:           flavor.Type.(models.K8Slice).Properties.Latency,
				SecurityStandards: flavor.Type.(models.K8Slice).Properties.SecurityStandards,
				CarbonFootprint: &nodecorev1alpha1.CarbonFootprint{
					Embodied:    flavor.Type.(models.K8Slice).Properties.CarbonFootprint.Embodied,
					Operational: flavor.Type.(models.K8Slice).Properties.CarbonFootprint.Operational,
				},
			},
			Policies: nodecorev1alpha1.Policies{
				Partitionability: nodecorev1alpha1.Partitionability{
					CpuMin:     flavor.Type.(models.K8Slice).Policies.Partitionability.CpuMin,
					MemoryMin:  flavor.Type.(models.K8Slice).Policies.Partitionability.MemoryMin,
					PodsMin:    flavor.Type.(models.K8Slice).Policies.Partitionability.PodsMin,
					CpuStep:    flavor.Type.(models.K8Slice).Policies.Partitionability.CpuStep,
					MemoryStep: flavor.Type.(models.K8Slice).Policies.Partitionability.MemoryStep,
					PodsStep:   flavor.Type.(models.K8Slice).Policies.Partitionability.PodsStep,
				},
			},
		}
		flavorTypeDataJSON, err := json.Marshal(flavorTypeData)
		if err != nil {
			klog.Errorf("Error when marshalling K8SliceType: %s", err)
			return nil
		}
		flavorType = nodecorev1alpha1.FlavorType{
			TypeIdentifier: nodecorev1alpha1.Type_K8Slice,
			TypeData:       runtime.RawExtension{Raw: flavorTypeDataJSON},
		}

	default:
		klog.Errorf("Flavor type not recognized")
		return nil
	}
	f := &nodecorev1alpha1.Flavor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flavor.FlavorID,
			Namespace: flags.FluidoNamespace,
		},
		Spec: nodecorev1alpha1.FlavorSpec{
			ProviderID: flavor.Owner.NodeID,
			FlavorType: flavorType,
			Owner: nodecorev1alpha1.NodeIdentity{
				Domain: flavor.Owner.Domain,
				IP:     flavor.Owner.IP,
				NodeID: flavor.Owner.NodeID,
			},
			Price: nodecorev1alpha1.Price{
				Amount:   flavor.Price.Amount,
				Currency: flavor.Price.Currency,
				Period:   flavor.Price.Period,
			},
			Availability:        flavor.Availability,
			NetworkPropertyType: flavor.NetworkPropertyType,
			Location: &nodecorev1alpha1.Location{
				Latitude:        flavor.Location.Latitude,
				Longitude:       flavor.Location.Longitude,
				Country:         flavor.Location.Country,
				City:            flavor.Location.City,
				AdditionalNotes: flavor.Location.AdditionalNotes,
			},
		},
	}
	return f
}

// ForgePartition creates a Partition from a FlavorSelector.
func ForgePartition(selector *nodecorev1alpha1.FlavorSelector) *nodecorev1alpha1.Partition {
	return &nodecorev1alpha1.Partition{
		CPU:     selector.RangeSelector.MinCpu,
		Memory:  selector.RangeSelector.MinMemory,
		Pods:    selector.RangeSelector.MinPods,
		Storage: selector.RangeSelector.MinStorage,
		Gpu: &nodecorev1alpha1.GPU{
			Model:  selector.RangeSelector.MinGpu.Model,
			Cores:  selector.RangeSelector.MinGpu.Cores,
			Memory: selector.RangeSelector.MinGpu.Memory,
		},
	}
}

// ForgeAllocation creates an Allocation from a Contract.
func ForgeAllocation(contract *reservationv1alpha1.Contract, intentID, nodeName string,
	destination nodecorev1alpha1.Destination, nodeType nodecorev1alpha1.NodeType) *nodecorev1alpha1.Allocation {
	return &nodecorev1alpha1.Allocation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namings.ForgeAllocationName(contract.Spec.Flavor.Name),
			Namespace: flags.FluidoNamespace,
		},
		Spec: nodecorev1alpha1.AllocationSpec{
			RemoteClusterID: func() string {
				if nodeType == nodecorev1alpha1.Node {
					return contract.Spec.BuyerClusterID
				}
				return contract.Spec.SellerCredentials.ClusterID
			}(),
			IntentID:    intentID,
			NodeName:    nodeName,
			Type:        nodeType,
			Destination: destination,
			Forwarding:  false,
			Contract: nodecorev1alpha1.GenericRef{
				Name:      contract.Name,
				Namespace: contract.Namespace,
			},
		},
	}
}
