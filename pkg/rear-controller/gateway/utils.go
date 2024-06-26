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

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"

	"github.com/fluidos-project/node/pkg/utils/models"
)

// selectorToQueryParams converts a selector to a query string.
func selectorToQueryParams(selector models.Selector) (string, error) {
	var values url.Values
	var encoder = schema.NewEncoder()

	switch selector.GetSelectorType() {
	case models.K8SliceNameDefault:
		err := encoder.Encode(selector.(models.K8SliceSelector), values)
		if err != nil {
			return "", err
		}

		return values.Encode(), nil
	default:
		return "", fmt.Errorf("unsupported selector type")
	}
}

// queryParamToSelector converts a query string to a selector.
func queryParamToSelector(queryValues url.Values, selectorType models.FlavorTypeName) (models.Selector, error) {

	var decoder = schema.NewDecoder()

	switch selectorType {
	case models.K8SliceNameDefault:
		var selector models.K8SliceSelector
		err := decoder.Decode(&selector, queryValues)
		if err != nil {
			return nil, err
		}
		return selector, nil
	default:
		return nil, fmt.Errorf("unsupported selector type")
	}
}

// GetTransaction returns a transaction from the transactions map.
func (g *Gateway) GetTransaction(transactionID string) (models.Transaction, error) {
	transaction, exists := g.Transactions[transactionID]
	if !exists {
		return models.Transaction{}, fmt.Errorf("transaction not found")
	}
	return transaction, nil
}

// SearchTransaction returns a transaction from the transactions map.
func (g *Gateway) SearchTransaction(buyerID, flavorID string) (*models.Transaction, bool) {
	for _, t := range g.Transactions {
		if t.Buyer.NodeID == buyerID && t.FlavorID == flavorID {
			return &t, true
		}
	}
	return &models.Transaction{}, false
}

// addNewTransacion add a new transaction to the transactions map.
func (g *Gateway) addNewTransacion(transaction *models.Transaction) {
	g.Transactions[transaction.TransactionID] = *transaction
}

// removeTransaction removes a transaction from the transactions map.
func (g *Gateway) removeTransaction(transactionID string) {
	delete(g.Transactions, transactionID)
}

// handleError handles errors by sending an error response.
func handleError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

// encodeResponse encodes the response as JSON and writes it to the response writer.
func encodeResponse(w http.ResponseWriter, data interface{}) {
	encodeResponseStatusCode(w, data, http.StatusOK)
}

func encodeResponseStatusCode(w http.ResponseWriter, data interface{}, statusCode int) {
	resp, err := json.Marshal(data)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}
