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
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	Type_K8Slice FlavorTypeIdentifier = "K8Slice"
	Type_VM      FlavorTypeIdentifier = "VM"
	Type_Service FlavorTypeIdentifier = "Service"
)

type FlavorTypeIdentifier string

type FlavorType struct {
	// Type of the Flavor.
	TypeIdentifier FlavorTypeIdentifier `json:"typeIdentifier"`
	// Raw is the raw value of the Flavor.
	TypeData runtime.RawExtension `json:"typeData"`
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

type Price struct {
	// Amount is the amount of the price.
	Amount string `json:"amount"`

	// Currency is the currency of the price.
	Currency string `json:"currency"`

	// Period is the period of the price.
	Period string `json:"period"`
}

// Location represents the location of a Flavor.
type Location struct {
	// Latitude is the latitude of the location.
	Latitude string `json:"latitude,omitempty"`

	// Longitude is the longitude of the location.
	Longitude string `json:"longitude,omitempty"`

	// Country is the country of the location.
	Country string `json:"country,omitempty"`

	// City is the city of the location.
	City string `json:"city,omitempty"`

	// AdditionalNotes are additional notes of the location.
	AdditionalNotes string `json:"additionalNotes,omitempty"`
}

// FlavorSpec defines the desired state of Flavor
type FlavorSpec struct {
	// This specs are based on the REAR Protocol specifications.

	// ProviderID is the ID of the FLUIDOS Node ID that provides this Flavor.
	// It can correspond to ID of the owner FLUIDOS Node or to the ID of a FLUIDOS SuperNode that represents the entry point to a FLUIDOS Domain
	ProviderID string `json:"providerID"`

	// Type is the type of the Flavor.
	Type FlavorType `json:"type"`

	// Timestamp is the timestamp of the Flavor.
	Timestamp metav1.Time `json:"timestamp"`

	// Owner contains the identity info of the owner of the Flavor. It can be unknown if the Flavor is provided by a reseller or a third party.
	Owner NodeIdentity `json:"owner"`

	// Price contains the price model of the Flavor.
	Price Price `json:"price"`

	// Availability is the availability flag of the Flavor.
	Availability bool `json:"availability"`

	// NetworkPropertyType is the network property type of the Flavor.
	NetworkPropertyType string `json:"networkPropertyType,omitempty"`

	// Location is the location of the Flavor.
	Location *Location `json:"location,omitempty"`
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
// +kubebuilder:resource:shortName=fl
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
