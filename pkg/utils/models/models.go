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
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
)

// Flavor represents a Flavor object with its characteristics and policies.
type Flavor struct {
	FlavorID            string       `json:"flavorID"`
	ProviderID          string       `json:"providerID"`
	Type                FlavorType   `json:"type"`
	NetworkPropertyType string       `json:"networkPropertyType,omitempty"`
	Timestamp           time.Time    `json:"timestamp"`
	Location            Location     `json:"location,omitempty"`
	Price               Price        `json:"price"`
	Owner               NodeIdentity `json:"owner"`
	Availability        bool         `json:"availability"`
}

// FlavorType represents the type of a Flavor.
type FlavorType interface {
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
type Selector struct {
	FlavorType    string         `json:"type,omitempty"`
	RangeSelector *RangeSelector `json:"rangeSelector,omitempty"`
	MatchSelector *MatchSelector `json:"matchSelector,omitempty"`
}

// MatchSelector represents the criteria for selecting Flavors through a strict match.
type MatchSelector struct {
	CPU     resource.Quantity   `json:"cpu,omitempty"`
	Memory  resource.Quantity   `json:"memory,omitempty"`
	Pods    resource.Quantity   `json:"pods,omitempty"`
	Storage resource.Quantity   `json:"storage,omitempty"`
	Gpu     *GpuCharacteristics `json:"gpu,omitempty"`
}

// RangeSelector represents the criteria for selecting Flavors through a range.
type RangeSelector struct {
	MinCPU     resource.Quantity   `json:"minCpu,omitempty"`
	MinMemory  resource.Quantity   `json:"minMemory,omitempty"`
	MinPods    resource.Quantity   `json:"minPods,omitempty"`
	MinStorage resource.Quantity   `json:"minStorage,omitempty"`
	MinGpu     *GpuCharacteristics `json:"minGpu,omitempty"`
	MaxCPU     resource.Quantity   `json:"maxCpu,omitempty"`
	MaxMemory  resource.Quantity   `json:"maxMemory,omitempty"`
	MaxPods    resource.Quantity   `json:"maxPods,omitempty"`
	MaxStorage resource.Quantity   `json:"maxStorage,omitempty"`
	MaxGpu     *GpuCharacteristics `json:"maxGpu,omitempty"`
}
