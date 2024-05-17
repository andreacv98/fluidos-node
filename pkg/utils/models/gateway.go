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

// PurchaseRequest is the request model for purchasing a Flavor.
type PurchaseRequest struct {
	TransactionID string `json:"transactionID"`
}

// ResponsePurchase contain information after purchase a Flavor.
type ResponsePurchase struct {
	Contract Contract `json:"contract"`
	Status   string   `json:"status"`
}

// ReserveRequest is the request model for reserving a Flavor.
type ReserveRequest struct {
	FlavorID  string       `json:"flavorID"`
	Buyer     NodeIdentity `json:"buyerID"`
	ClusterID string       `json:"clusterID"`
	Partition *Partition   `json:"partition,omitempty"`
}
