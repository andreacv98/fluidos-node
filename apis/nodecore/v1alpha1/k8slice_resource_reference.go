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

// K8SliceResourceReference is the reference to the resource that is being allocated by the K8Slice
type K8SliceResourceReference struct {
	// VirtualNodeName is the name of the virtual node that linked to the allocation of the K8Slice
	VirtualNodeName string `json:"virtualNodeName"`
}
