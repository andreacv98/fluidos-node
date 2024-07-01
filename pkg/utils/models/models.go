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

package models

import (
	"encoding/json"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
)

// Flavor represents a Flavor object with its characteristics and policies.
type Flavor struct {
	FlavorID            string       `json:"flavorID"`
	ProviderID          string       `json:"providerID"`
	Type                FlavorType   `json:"type"`
	NetworkPropertyType string       `json:"networkPropertyType,omitempty"`
	Timestamp           time.Time    `json:"timestamp"`
	Location            *Location     `json:"location,omitempty"`
	Price               Price        `json:"price"`
	Owner               NodeIdentity `json:"owner"`
	Availability        bool         `json:"availability"`
}

// FlavorType represents the type of a Flavor.
type FlavorType struct {
	Name FlavorTypeName  `json:"name"`
	Data json.RawMessage `json:"data"`
}

type FlavorTypeData interface {
	GetFlavorTypeName() FlavorTypeName
}

type FlavorTypeName string

type CarbonFootprint struct {
	Embodied    int `json:"embodied"`
	Operational int `json:"operational"`
}

const (
	K8SliceNameDefault FlavorTypeName = "k8slice"
	VMNameDefault      FlavorTypeName = "vm"
	ServiceNameDefault FlavorTypeName = "service"
	SensorNameDefault  FlavorTypeName = "sensor"
)

// Location represents the location of a Flavor, with latitude, longitude, altitude, and additional notes.
type Location struct {
	Latitude        string `json:"latitude,omitempty"`
	Longitude       string `json:"longitude,omitempty"`
	Country         string `json:"altitude,omitempty"`
	City            string `json:"city,omitempty"`
	AdditionalNotes string `json:"additionalNotes,omitempty"`
}

type NodeIdentityAdditionalInfo struct {
	LiqoID string `json:"liqoID,omitempty"`
}

// NodeIdentity represents the owner of a Flavor, with associated ID, IP, and domain name.
type NodeIdentity struct {
	NodeID                string                      `json:"ID"`
	IP                    string                      `json:"IP"`
	Domain                string                      `json:"domain"`
	AdditionalInformation *NodeIdentityAdditionalInfo `json:"additionalInformation,omitempty"`
}

// Price represents the price of a Flavor, with the amount, currency, and period associated.
type Price struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Period   string `json:"period"`
}

// Selector represents the criteria for selecting Flavors.
type Selector interface {
	GetSelectorType() FlavorTypeName
}

type K8SliceSelector struct {
	Cpu     *ResourceQuantityFilter	`scheme:"cpu,omitempty"`
	Memory  *ResourceQuantityFilter	`scheme:"memory,omitempty"`
	Pods    *ResourceQuantityFilter	`scheme:"pods,omitempty"`
	Storage *ResourceQuantityFilter	`scheme:"storage,omitempty"`
}

func (ks K8SliceSelector) GetSelectorType() FlavorTypeName {
	return K8SliceNameDefault
}

type ResourceQuantityFilter struct {
	Name FilterType			`scheme:"name"`
	Data json.RawMessage	`scheme:"data"`
}

type ResourceQuantityFilterData interface {
	GetFilterType() FilterType
}

type FilterType string

const (
	MatchFilter FilterType = "Match"
	RangeFilter FilterType = "Range"
)

type ResourceQuantityMatchFilter struct {
	Value resource.Quantity `scheme:"value,omitempty"`
}

func (fq ResourceQuantityMatchFilter) GetFilterType() FilterType {
	return MatchFilter
}

type ResourceQuantityRangeFilter struct {
	Min *resource.Quantity	`scheme:"min,omitempty"`
	Max *resource.Quantity	`scheme:"max,omitempty"`
}

func (fq ResourceQuantityRangeFilter) GetFilterType() FilterType {
	return RangeFilter
}

// MapToFlavorTypeName maps a nodecorev1alpha1.FlavorTypeIdentifier to a models.FlavorTypeName.
func MapToFlavorTypeName(flavorType nodecorev1alpha1.FlavorTypeIdentifier) FlavorTypeName {
	switch flavorType {
	case nodecorev1alpha1.Type_K8Slice:
		return K8SliceNameDefault
	case nodecorev1alpha1.Type_VM:
		return VMNameDefault
	case nodecorev1alpha1.Type_Service:
		return ServiceNameDefault
	default:
		return ""
	}
}

// MapFromFlavorTypeName maps a models.FlavorTypeName to a nodecorev1alpha1.FlavorTypeIdentifier.
func MapFromFlavorTypeName(flavorType FlavorTypeName) nodecorev1alpha1.FlavorTypeIdentifier {
	switch flavorType {
	case K8SliceNameDefault:
		return nodecorev1alpha1.Type_K8Slice
	case VMNameDefault:
		return nodecorev1alpha1.Type_VM
	case ServiceNameDefault:
		return nodecorev1alpha1.Type_Service
	default:
		return ""
	}
}

// MapToFilterType maps a nodecorev1alpha1.FilterType to a models.FilterType.
func MapToFilterType(filterType nodecorev1alpha1.FilterType) FilterType {
	switch filterType {
	case nodecorev1alpha1.TypeMatchFilter:
		return MatchFilter
	case nodecorev1alpha1.TypeRangeFilter:
		return RangeFilter
	default:
		return ""
	}
}

// MapFromFilterType maps a models.FilterType to a nodecorev1alpha1.FilterType.
func MapFromFilterType(filterType FilterType) nodecorev1alpha1.FilterType {
	switch filterType {
	case MatchFilter:
		return nodecorev1alpha1.TypeMatchFilter
	case RangeFilter:
		return nodecorev1alpha1.TypeRangeFilter
	default:
		return ""
	}
}
