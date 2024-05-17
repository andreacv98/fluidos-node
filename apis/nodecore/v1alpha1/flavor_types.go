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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	K8S FlavorType = "k8s-fluidos"
)

type FlavorType string

type Characteristics struct {

	// Architecture is the architecture of the Flavor.
	Architecture string `json:"architecture"`

	// CPU is the number of CPU cores of the Flavor.
	Cpu resource.Quantity `json:"cpu"`

	// Memory is the amount of RAM of the Flavor.
	Memory resource.Quantity `json:"memory"`

	// Pods is the maximum number of pods of the Flavor.
	Pods resource.Quantity `json:"pods"`

	// GPU is the number of GPU cores of the Flavor.
	Gpu resource.Quantity `json:"gpu,omitempty"`

	// EphemeralStorage is the amount of ephemeral storage of the Flavor.
	EphemeralStorage resource.Quantity `json:"ephemeral-storage,omitempty"`

	// PersistentStorage is the amount of persistent storage of the Flavor.
	PersistentStorage resource.Quantity `json:"persistent-storage,omitempty"`
}

type Policy struct {

	// Partitionable contains the partitioning properties of the Flavor.
	Partitionable *Partitionable `json:"partitionable,omitempty"`

	// Aggregatable contains the aggregation properties of the Flavor.
	Aggregatable *Aggregatable `json:"aggregatable,omitempty"`
}

// Partitionable represents the partitioning properties of a Flavor, such as the minimum and incremental values of CPU and RAM.
type Partitionable struct {
	// CpuMin is the minimum requirable number of CPU cores of the Flavor.
	CpuMin resource.Quantity `json:"cpuMin"`

	// MemoryMin is the minimum requirable amount of RAM of the Flavor.
	MemoryMin resource.Quantity `json:"memoryMin"`

	// PodsMin is the minimum requirable number of pods of the Flavor.
	PodsMin resource.Quantity `json:"podsMin"`

	// CpuStep is the incremental value of CPU cores of the Flavor.
	CpuStep resource.Quantity `json:"cpuStep"`

	// MemoryStep is the incremental value of RAM of the Flavor.
	MemoryStep resource.Quantity `json:"memoryStep"`

	// PodsStep is the incremental value of pods of the Flavor.
	PodsStep resource.Quantity `json:"podsStep"`
}

// Aggregatable represents the aggregation properties of a Flavor, such as the minimum instance count.
type Aggregatable struct {
	// MinCount is the minimum requirable number of instances of the Flavor.
	MinCount int `json:"minCount"`

	// MaxCount is the maximum requirable number of instances of the Flavor.
	MaxCount int `json:"maxCount"`
}

type Price struct {

	// Amount is the amount of the price.
	Amount string `json:"amount"`

	// Currency is the currency of the price.
	Currency string `json:"currency"`

	// Period is the period of the price.
	Period string `json:"period"`
}

type OptionalFields struct {

	// Availability is the availability flag of the Flavor.
	// It is a field inherited from the REAR Protocol specifications.
	Availability bool `json:"availability,omitempty"`

	// WorkerID is the ID of the worker that provides the Flavor.
	WorkerID string `json:"workerID,omitempty"`
}

// FlavorSpec defines the desired state of Flavor
type FlavorSpec struct {
	// This specs are based on the REAR Protocol specifications.

	// ProviderID is the ID of the FLUIDOS Node ID that provides this Flavor.
	// It can correspond to ID of the owner FLUIDOS Node or to the ID of a FLUIDOS SuperNode that represents the entry point to a FLUIDOS Domain
	ProviderID string `json:"providerID"`

	// Type is the type of the Flavor. Currently, only K8S is supported.
	Type FlavorType `json:"type"`

	// Characteristics contains the characteristics of the Flavor.
	// They are based on the type of the Flavor and can change depending on it. In this case, the type is K8S so the characteristics are CPU, Memory, GPU and EphemeralStorage.
	Characteristics Characteristics `json:"characteristics"`

	// Policy contains the policy of the Flavor. The policy describes the partitioning and aggregation properties of the Flavor.
	Policy Policy `json:"policy"`

	// Owner contains the identity info of the owner of the Flavor. It can be unknown if the Flavor is provided by a reseller or a third party.
	Owner NodeIdentity `json:"owner"`

	// Price contains the price model of the Flavor.
	Price Price `json:"price"`

	// This field is used to specify the optional fields that can be retrieved from the Flavor.
	// In the future it will be expanded to include more optional fields defined in the REAR Protocol or custom ones.
	OptionalFields OptionalFields `json:"optionalFields"`
}

// FlavorStatus defines the observed state of Flavor.
type FlavorStatus struct {

	// This field represents the expiration time of the Flavor. It is used to determine when the Flavor is no longer valid.
	ExpirationTime string `json:"expirationTime"`

	// This field represents the creation time of the Flavor.
	CreationTime string `json:"creationTime"`

	// This field represents the last update time of the Flavor.
	LastUpdateTime string `json:"lastUpdateTime"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// +kubebuilder:printcolumn:name="Provider ID",type=string,JSONPath=`.spec.providerID`
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.type`
// +kubebuilder:printcolumn:name="CPU",type=string,priority=1,JSONPath=`.spec.characteristics.cpu`
// +kubebuilder:printcolumn:name="Memory",type=string,priority=1,JSONPath=`.spec.characteristics.memory`
// +kubebuilder:printcolumn:name="Owner Name",type=string,priority=1,JSONPath=`.spec.owner.nodeID`
// +kubebuilder:printcolumn:name="Owner Domain",type=string,priority=1,JSONPath=`.spec.owner.domain`
// +kubebuilder:printcolumn:name="Available",type=boolean,JSONPath=`.spec.optionalFields.availability`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// Flavor is the Schema for the flavors API.
type Flavor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FlavorSpec   `json:"spec,omitempty"`
	Status FlavorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FlavorList contains a list of Flavor
type FlavorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Flavor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Flavor{}, &FlavorList{})
}
