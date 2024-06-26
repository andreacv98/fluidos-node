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
	TypeMatchFilter FilterType = "Match"
	TypeRangeFilter FilterType = "Range"
)

type FilterType string

const (
	TypeResourceQuantity ValueType = "ResourceQuantity"
	TypeInt              ValueType = "Int"
	TypeString           ValueType = "String"
)

type ValueType string

type ResourceQuantityFilter struct {
	// FilterType indicates the type of the filter
	FilterType FilterType `json:"filterType"`
	// Filter data
	Data runtime.RawExtension `json:"data"`
}

type ResourceMatchSelector struct {
	Value resource.Quantity `json:"value"`
}

type ResourceRangeSelector struct {
	Min resource.Quantity `json:"min"`
	Max resource.Quantity `json:"max"`
}

// ParseResourceQuantityFilter parses a ResourceQuantityFilter into a FilterType and the corresponding filter data.
func ParseResourceQuantityFilter(rqf *ResourceQuantityFilter) (FilterType, interface{}, error) {
	var validationErr error

	switch rqf.FilterType {
	case TypeMatchFilter:
		// Unmarshal the data into a ResourceMatchSelector
		var rms ResourceMatchSelector
		validationErr = json.Unmarshal(rqf.Data.Raw, &rms)
		return TypeMatchFilter, rms, validationErr
	case TypeRangeFilter:
		// Unmarshal the data into a ResourceRangeSelector
		var rrs ResourceRangeSelector
		validationErr = json.Unmarshal(rqf.Data.Raw, &rrs)
		return TypeRangeFilter, rrs, validationErr
	default:
		return "", nil, fmt.Errorf("unknown filter type %s", rqf.FilterType)
	}
}
