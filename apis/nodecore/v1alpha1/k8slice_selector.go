// Copyright 2022-2024 FLUIDOS Project
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
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	TypeMatchSelector TypeSelector = "Match"
	TypeRangeSelector TypeSelector = "Range"
)

type TypeSelector string

const (
	ResourceTypeCpu     ResourceType = "Cpu"
	ResourceTypeMemory  ResourceType = "Memory"
	ResourceTypePods    ResourceType = "Pods"
	ResourceTypeStorage ResourceType = "Storage"
)

type ResourceType string

type ResourceSelector struct {
	// Type of the resource filter
	Type TypeSelector `json:"type"`
	// Filter data
	Data runtime.RawExtension `json:"data"`
}

type CpuMatchSelector struct {
	// Exact value of the CPU
	Value resource.Quantity `json:"value"`
}

func (CpuMatchSelector) GetTypeCpuSelector() TypeSelector {
	return TypeMatchSelector
}

type CpuRangeSelector struct {
	// Minimum value of the CPU
	Min resource.Quantity `json:"min"`
	// Maximum value of the CPU
	Max resource.Quantity `json:"max"`
}

func (CpuRangeSelector) GetTypeCpuSelector() TypeSelector {
	return TypeRangeSelector
}

type MemoryMatchSelector struct {
	// Exact value of the memory
	Value resource.Quantity `json:"value"`
}

func (MemoryMatchSelector) GetTypeMemorySelector() TypeSelector {
	return TypeMatchSelector
}

type MemoryRangeSelector struct {
	// Minimum value of the memory
	Min resource.Quantity `json:"min"`
	// Maximum value of the memory
	Max resource.Quantity `json:"max"`
}

func (MemoryRangeSelector) GetTypeMemorySelector() TypeSelector {
	return TypeRangeSelector
}

type PodsMatchSelector struct {
	// Exact value of the pods
	Value resource.Quantity `json:"value"`
}

func (PodsMatchSelector) GetTypePodsSelector() TypeSelector {
	return TypeMatchSelector
}

type PodsRangeSelector struct {
	// Minimum value of the pods
	Min resource.Quantity `json:"min"`
	// Maximum value of the pods
	Max resource.Quantity `json:"max"`
}

func (PodsRangeSelector) GetTypePodsSelector() TypeSelector {
	return TypeRangeSelector
}

type StorageMatchSelector struct {
	// Exact value of the storage
	Value resource.Quantity `json:"value"`
}

func (StorageMatchSelector) GetTypeStorageSelector() TypeSelector {
	return TypeMatchSelector
}

type StorageRangeSelector struct {
	// Minimum value of the storage
	Min resource.Quantity `json:"min"`
	// Maximum value of the storage
	Max resource.Quantity `json:"max"`
}

func (StorageRangeSelector) GetTypeStorageSelector() TypeSelector {
	return TypeRangeSelector
}

type K8SliceFilter struct {
	CpuFilter     ResourceSelector `json:"cpuFilter"`
	MemoryFilter  ResourceSelector `json:"memoryFilter"`
	PodsFilter    ResourceSelector `json:"podsFilter"`
	StorageFilter ResourceSelector `json:"storageFilter"`
}

func (*K8SliceFilter) GetFlavorTypeFilter() FlavorTypeIdentifier {
	return Type_K8Slice
}

func ParseK8SliceResourceSelector(resourceType ResourceType, rf *ResourceSelector) (TypeSelector, interface{}, error) {
	var validationErr error

	switch resourceType {
	case ResourceTypeCpu:
		switch rf.Type {
		case TypeMatchSelector:
			// TODO
			var CpuMatchSelector CpuMatchSelector
			validationErr = json.Unmarshal(rf.Data.Raw, &CpuMatchSelector)
			return CpuMatchSelector.GetTypeCpuSelector(), CpuMatchSelector, validationErr
		case TypeRangeSelector:
			// TODO
		default:
			return "", nil, fmt.Errorf("resource selector type %s not supported", rf.Type)
		}
	case ResourceTypeMemory:
		switch rf.Type {
		case TypeMatchSelector:
			// TODO
		case TypeRangeSelector:
			// TODO
		default:
			return "", nil, fmt.Errorf("resource selector type %s not supported", rf.Type)
		}
	case ResourceTypePods:
		switch rf.Type {
		case TypeMatchSelector:
			// TODO
		case TypeRangeSelector:
			// TODO
		default:
			return "", nil, fmt.Errorf("resource selector type %s not supported", rf.Type)
		}
	case ResourceTypeStorage:
		switch rf.Type {
		case TypeMatchSelector:
			// TODO
		case TypeRangeSelector:
			// TODO
		default:
			return "", nil, fmt.Errorf("resource selector type %s not supported", rf.Type)
		}
	default:
		return "", nil, fmt.Errorf("resource type %s not supported", resourceType)
	}

	return "", nil, fmt.Errorf("resource type %s not supported", resourceType)
}
