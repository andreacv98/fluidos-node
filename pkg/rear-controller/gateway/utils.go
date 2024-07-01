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
	"reflect"

	"github.com/gorilla/schema"
	"k8s.io/klog/v2"

	"github.com/fluidos-project/node/pkg/utils/models"
)

// selectorToQueryParams converts a selector to a query string.
func selectorToQueryParams(selector models.Selector) (string, error) {
	var values = url.Values{}
	var encoder = schema.NewEncoder()

	switch selector.GetSelectorType() {
	case models.K8SliceNameDefault:
		k8sliceSelector := selector.(models.K8SliceSelector)
		klog.Infof("k8sliceSelector: %v", k8sliceSelector)
		encoder.RegisterEncoder(&models.ResourceQuantityFilter{}, EncoderResourceQuantityFilter)
		klog.Infof("encoder: %v", encoder)
		err := encoder.Encode(k8sliceSelector, values)
		if err != nil {
			return "", err
		}

		klog.Infof("Encoded selector: %v", values.Encode())

		return values.Encode(), nil
	default:
		return "", fmt.Errorf("unsupported selector type")
	}
}

// queryParamToSelector converts a query string to a selector.
func queryParamToSelector(queryValues url.Values, selectorType models.FlavorTypeName) (models.Selector, error) {

	var decoder = schema.NewDecoder()

	klog.Infof("queryValues: %v", queryValues)

	switch selectorType {
	case models.K8SliceNameDefault:
		var selector models.K8SliceSelector
		decoder.RegisterConverter(models.ResourceQuantityFilter{}, ConverterResourceQuantityFilter)
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

func EncoderResourceQuantityFilter(v reflect.Value) string {
	filter, ok := v.Interface().(*models.ResourceQuantityFilter)
	if filter == nil {
		return ""
	}
	if !ok {
		klog.Info("Not a ResourceQuantityFilter")
		return ""
	}

	values := url.Values{}

	var dataMap map[string]string

	switch filter.Name {
	case models.MatchFilter:
		var matchFilter models.ResourceQuantityMatchFilter
		if err := json.Unmarshal(filter.Data, &matchFilter); err != nil {
			return ""
		}
		dataMap = map[string]string{
			"[match]": matchFilter.Value.String(),
		}
	case models.RangeFilter:
		var rangeFilter models.ResourceQuantityRangeFilter
		if err := json.Unmarshal(filter.Data, &rangeFilter); err != nil {
			return ""
		}
		subDataMap := map[string]string{}
		if rangeFilter.Min != nil {
			subDataMap["min"] = rangeFilter.Min.String()
		}
		if rangeFilter.Max != nil {
			subDataMap["max"] = rangeFilter.Max.String()
		}
		dataMap = map[string]string{
			"[range]": "[min]="+subDataMap["min"] + "," + "[max]="+subDataMap["max"],
		}
	default:
		return ""
	}

	for k, v := range dataMap {
		values.Set(k, v)
	}
	klog.Infof("Encoded Resource Selector selector: %v", values.Encode())

	klog.Infof("Encoded selector: %v", values.Encode())

	return values.Encode()
}

func decodeResourceQuantityFilter(values url.Values, key string) (*models.ResourceQuantityFilter, error) {
	if values.Get(key) == "" {
		return nil, nil
	}

	filter := &models.ResourceQuantityFilter{
		Name: models.FilterType(values.Get(key + "[name]")),
	}

	switch filter.Name {
	case models.MatchFilter:
		var matchFilter models.ResourceQuantityMatchFilter
		if err := json.Unmarshal([]byte(values.Get(key + "[data]")), &matchFilter); err != nil {
			return nil, fmt.Errorf("error unmarshaling MatchFilter: %v", err)
		}
		filter.Data = json.RawMessage(values.Get(key + "[data]"))

	case models.RangeFilter:
		var rangeFilter models.ResourceQuantityRangeFilter
		if err := json.Unmarshal([]byte(values.Get(key + "[data]")), &rangeFilter); err != nil {
			return nil, fmt.Errorf("error unmarshaling RangeFilter: %v", err)
		}
		data, err := json.Marshal(rangeFilter)
		if err != nil {
			return nil, fmt.Errorf("error marshaling RangeFilter: %v", err)
		}
		filter.Data = json.RawMessage(data)

	default:
		return nil, fmt.Errorf("unsupported filter type: %s", filter.Name)
	}

	return filter, nil
}

func ConverterResourceQuantityFilter(v string) reflect.Value {
	selector := &models.K8SliceSelector{}

	values, err := url.ParseQuery(v)

	if err != nil {
		return reflect.ValueOf(nil)
	}

	selector.Cpu, err = decodeResourceQuantityFilter(values, "cpu")
	if err != nil {
		return reflect.ValueOf(nil)
	}
	selector.Memory, err = decodeResourceQuantityFilter(values, "memory")
	if err != nil {
		return reflect.ValueOf(nil)
	}
	selector.Pods, err = decodeResourceQuantityFilter(values, "pods")
	if err != nil {
		return reflect.ValueOf(nil)
	}
	selector.Storage, err = decodeResourceQuantityFilter(values, "storage")
	if err != nil {
		return reflect.ValueOf(nil)
	}

	return reflect.ValueOf(selector)
}
