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

package gateway

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	reservationv1alpha1 "github.com/fluidos-project/node/apis/reservation/v1alpha1"
	"github.com/fluidos-project/node/pkg/utils/common"
	"github.com/fluidos-project/node/pkg/utils/flags"
	"github.com/fluidos-project/node/pkg/utils/getters"
	"github.com/fluidos-project/node/pkg/utils/models"
	"github.com/fluidos-project/node/pkg/utils/namings"
	"github.com/fluidos-project/node/pkg/utils/parseutil"
	"github.com/fluidos-project/node/pkg/utils/resourceforge"
	"github.com/fluidos-project/node/pkg/utils/services"
	"github.com/fluidos-project/node/pkg/utils/tools"
)

// getFlavors gets all the flavors CRs from the cluster.
func (g *Gateway) getFlavors(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	klog.Infof("Processing request for getting all Flavors...")

	flavors, err := services.GetAllFlavors(g.client)
	if err != nil {
		klog.Errorf("Error getting all the Flavor CRs: %s", err)
		http.Error(w, "Error getting all the Flavor CRs", http.StatusInternalServerError)
		return
	}

	klog.Infof("Found %d Flavors in the cluster", len(flavors))

	availableFlavors := make([]nodecorev1alpha1.Flavor, 0)

	// Filtering only the available flavors
	for i := range flavors {
		if !flavors[i].Spec.Availability {
			availableFlavors = append(availableFlavors, flavors[i])
		}
	}

	klog.Infof("Available Flavors: %d", len(availableFlavors))
	if len(availableFlavors) == 0 {
		klog.Infof("No available Flavors found")
		// Return content for empty list
		emptyList := make([]*nodecorev1alpha1.Flavor, 0)
		encodeResponseStatusCode(w, emptyList, http.StatusNoContent)
		return
	}

	// Select the flavor with the max CPU
	max := resource.MustParse("0")
	index := 0
	for i := range availableFlavors {
		flavorTypeIdentifier, flavorTypeData, err := nodecorev1alpha1.ParseFlavorType(&availableFlavors[i])
		if err != nil {
			klog.Errorf("Error parsing the Flavor type: %s", err)
			http.Error(w, "Error parsing the Flavor type", http.StatusInternalServerError)
			return
		}
		if flavorTypeIdentifier == nodecorev1alpha1.Type_K8Slice {

			k8Slice := flavorTypeData.(nodecorev1alpha1.K8Slice)

			if k8Slice.Characteristics.Cpu.Cmp(max) == 1 {
				max = k8Slice.Characteristics.Cpu
				index = i
			}

		}

	}

	selected := *flavors[index].DeepCopy()

	klog.Infof("Flavor %s selected - Parsing...", selected.Name)
	parsed := parseutil.ParseFlavor(&selected)

	klog.Infof("Flavor parsed: %v", parsed)

	// Encode the Flavor as JSON and write it to the response writer
	encodeResponse(w, parsed)
}

// getFlavorBySelectorHandler gets the flavor CRs from the cluster that match the selector.
func (g *Gateway) getK8SliceFlavorsBySelector(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	klog.Infof("Processing request for getting K8Slice Flavors by selector...")

	// build the selector from the url query parameters
	selector, err := queryParamToSelector(r.URL.Query(), models.SensorNameDefault)
	if err != nil {
		klog.Errorf("Error building the selector from the URL query parameters: %s", err)
		http.Error(w, "Error building the selector from the URL query parameters", http.StatusBadRequest)
		return
	}

	flavors, err := services.GetAllFlavors(g.client)
	if err != nil {
		klog.Errorf("Error getting all the Flavor CRs: %s", err)
		http.Error(w, "Error getting all the Flavor CRs", http.StatusInternalServerError)
		return
	}

	klog.Infof("Found %d Flavors in the cluster", len(flavors))

	availableFlavors := make([]nodecorev1alpha1.Flavor, 0)

	// Filtering only the available flavors
	for i := range flavors {
		if flavors[i].Spec.Availability {
			availableFlavors = append(availableFlavors, flavors[i])
		}
	}

	klog.Infof("Available Flavors: %d", len(availableFlavors))
	if len(availableFlavors) == 0 {
		klog.Infof("No available Flavors found")
		// Return content for empty list
		emptyList := make([]*nodecorev1alpha1.Flavor, 0)
		encodeResponseStatusCode(w, emptyList, http.StatusNoContent)
		return
	}

	klog.Infof("Checking selector syntax...")
	if err := common.CheckSelector(selector); err != nil {
		klog.Errorf("Error checking the selector syntax: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	klog.Infof("Filtering Flavors by selector...")
	flavorsSelected, err := common.FilterFlavorsBySelector(availableFlavors, selector)
	if err != nil {
		http.Error(w, "Error getting the Flavors by selector", http.StatusInternalServerError)
		return
	}

	klog.Infof("Flavors found that match the selector are: %d", len(flavorsSelected))

	if len(flavorsSelected) == 0 {
		klog.Infof("No matching Flavors found")
		// Return content for empty list
		emptyList := make([]*nodecorev1alpha1.Flavor, 0)
		encodeResponse(w, emptyList)
		return
	}

	// Select the flavor with the max CPU
	max := resource.MustParse("0")
	index := 0

	for i := range availableFlavors {
		flavorTypeIdentifier, flavorTypeData, err := nodecorev1alpha1.ParseFlavorType(&availableFlavors[i])
		if err != nil {
			klog.Errorf("Error parsing the Flavor type: %s", err)
			http.Error(w, "Error parsing the Flavor type", http.StatusInternalServerError)
			return
		}
		if flavorTypeIdentifier == nodecorev1alpha1.Type_K8Slice {

			k8Slice := flavorTypeData.(nodecorev1alpha1.K8Slice)

			if k8Slice.Characteristics.Cpu.Cmp(max) == 1 {
				max = k8Slice.Characteristics.Cpu
				index = i
			}

		}

	}

	selected := *flavorsSelected[index].DeepCopy()

	klog.Infof("Flavor %s selected - Parsing...", selected.Name)
	parsed := parseutil.ParseFlavor(&selected)

	klog.Infof("Flavor parsed: %v", parsed)

	// Encode the Flavor as JSON and write it to the response writer
	encodeResponse(w, parsed)
}

// reserveFlavor reserves a Flavor by its flavorID.
func (g *Gateway) reserveFlavor(w http.ResponseWriter, r *http.Request) {
	// Get the flavorID value from the URL parameters
	params := mux.Vars(r)
	flavorID := params["flavorID"]
	var transaction *models.Transaction
	var request models.ReserveRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		klog.Errorf("Error decoding the ReserveRequest: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	klog.Infof("Partition: %v", *request.Partition)

	if flavorID != request.FlavorID {
		klog.Infof("Mismatch body & param: %s != %s", flavorID, request.FlavorID)
		http.Error(w, "Mismatch body & param", http.StatusConflict)
		return
	}

	// Check if the Transaction already exists
	t, found := g.SearchTransaction(request.Buyer.NodeID, flavorID)
	if found {
		t.StartTime = tools.GetTimeNow()
		transaction = t
		g.addNewTransacion(t)
	}

	if !found {
		klog.Infof("Reserving flavor %s started", flavorID)

		flavor, _ := services.GetFlavorByID(flavorID, g.client)
		if flavor == nil {
			http.Error(w, "Flavor not found", http.StatusNotFound)
			return
		}

		// Create a new transaction ID
		transactionID, err := namings.ForgeTransactionID()
		if err != nil {
			http.Error(w, "Error generating transaction ID", http.StatusInternalServerError)
			return
		}

		// Create a new transaction
		transaction = resourceforge.ForgeTransactionObj(transactionID, &request)

		// Add the transaction to the transactions map
		g.addNewTransacion(transaction)
	}

	klog.Infof("Transaction %s reserved", transaction.TransactionID)

	encodeResponse(w, transaction)
}

// purchaseFlavor is an handler for purchasing a Flavor.
func (g *Gateway) purchaseFlavor(w http.ResponseWriter, r *http.Request) {
	// Get the flavorID value from the URL parameters
	params := mux.Vars(r)
	transactionID := params["transactionID"]
	var purchase models.PurchaseRequest

	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if transactionID != purchase.TransactionID {
		klog.Infof("Mismatch body & param")
		http.Error(w, "Mismatch body & param", http.StatusConflict)
		return
	}

	klog.Infof("Purchasing request for transaction %s", purchase.TransactionID)

	// Retrieve the transaction from the transactions map
	transaction, err := g.GetTransaction(purchase.TransactionID)
	if err != nil {
		klog.Errorf("Error getting the Transaction: %s", err)
		http.Error(w, "Error getting the Transaction", http.StatusInternalServerError)
		return
	}

	klog.Infof("Flavor requested: %s", transaction.FlavorID)

	if tools.CheckExpiration(transaction.StartTime, flags.ExpirationTransaction) {
		klog.Infof("Transaction %s expired", transaction.TransactionID)
		http.Error(w, "Error: transaction Timeout", http.StatusRequestTimeout)
		g.removeTransaction(transaction.TransactionID)
		return
	}

	var contractList reservationv1alpha1.ContractList
	var contract reservationv1alpha1.Contract

	// Check if the Contract with the same TransactionID already exists
	if err := g.client.List(context.Background(), &contractList, client.MatchingFields{"spec.transactionID": purchase.TransactionID}); err != nil {
		if client.IgnoreNotFound(err) != nil {
			klog.Errorf("Error when listing Contracts: %s", err)
			http.Error(w, "Error when listing Contracts", http.StatusInternalServerError)
			return
		}
	}

	if len(contractList.Items) > 0 {
		klog.Infof("Contract already exists for transaction %s", purchase.TransactionID)
		contract = contractList.Items[0]
		// Create a contract object to be returned with the response
		contractObject := parseutil.ParseContract(&contract)
		// create a response purchase
		responsePurchase := resourceforge.ForgeResponsePurchaseObj(contractObject)
		// Respond with the response purchase as JSON
		encodeResponse(w, responsePurchase)
		return
	}

	klog.Infof("Performing purchase of flavor %s...", transaction.FlavorID)

	// Remove the transaction from the transactions map
	g.removeTransaction(transaction.TransactionID)

	klog.Infof("Flavor %s successfully purchased!", transaction.FlavorID)

	// Get the flavor sold for creating the contract
	flavorSold, err := services.GetFlavorByID(transaction.FlavorID, g.client)
	if err != nil {
		klog.Errorf("Error getting the Flavor by ID: %s", err)
		http.Error(w, "Error getting the Flavor by ID", http.StatusInternalServerError)
		return
	}

	liqoCredentials, err := getters.GetLiqoCredentials(context.Background(), g.client)
	if err != nil {
		klog.Errorf("Error getting Liqo Credentials: %s", err)
		http.Error(w, "Error getting Liqo Credentials", http.StatusInternalServerError)
		return
	}

	// Create a new contract
	klog.Infof("Creating a new contract...")
	contract = *resourceforge.ForgeContract(flavorSold, &transaction, liqoCredentials)
	err = g.client.Create(context.Background(), &contract)
	if err != nil {
		klog.Errorf("Error creating the Contract: %s", err)
		http.Error(w, "Error creating the Contract: "+err.Error(), http.StatusInternalServerError)
		return
	}

	klog.Infof("Contract created!")

	// Create a contract object to be returned with the response
	contractObject := parseutil.ParseContract(&contract)
	// create a response purchase
	responsePurchase := resourceforge.ForgeResponsePurchaseObj(contractObject)

	klog.Infof("Contract %v", *contractObject.Partition)

	// Create allocation
	klog.Infof("Creating allocation...")
	allocation := *resourceforge.ForgeAllocation(&contract, "", "", nodecorev1alpha1.Remote, nodecorev1alpha1.Node)
	err = g.client.Create(context.Background(), &allocation)
	if err != nil {
		klog.Errorf("Error creating the Allocation: %s", err)
		http.Error(w, "Contract created but we ran into an error while allocating the resources", http.StatusInternalServerError)
		return
	}

	klog.Infof("Response purchase %v", *responsePurchase.Contract.Partition)

	// Respond with the response purchase as JSON
	encodeResponse(w, responsePurchase)
}
