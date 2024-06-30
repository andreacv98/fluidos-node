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

package gateway

var Routes = struct {
	Flavors        string
	K8SliceFlavors string
	VMFlavors      string
	ServiceFlavors string
	Reserve        string
	Purchase       string
}{
	Flavors:        "/api/v1/flavors",
	K8SliceFlavors: "/api/v1/flavors/k8slice",
	VMFlavors:      "/api/v1/flavors/vm",
	ServiceFlavors: "/api/v1/flavors/service",
	Reserve:        "/api/v1/reservations",
	Purchase:       "/api/v1/transactions/{transactionID}/purchase",
}
