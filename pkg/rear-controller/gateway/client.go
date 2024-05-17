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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/klog/v2"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	reservationv1alpha1 "github.com/fluidos-project/node/apis/reservation/v1alpha1"
	"github.com/fluidos-project/node/pkg/utils/getters"
	"github.com/fluidos-project/node/pkg/utils/models"
	"github.com/fluidos-project/node/pkg/utils/parseutil"
)

// TODO: move this function into the REAR Gateway package

// ReserveFlavor reserves a flavor with the given flavorID.
func (g *Gateway) ReserveFlavor(ctx context.Context, reservation *reservationv1alpha1.Reservation, flavorID string) (*models.Transaction, error) {
	err := checkLiqoReadiness(g.LiqoReady)
	if err != nil {
		return nil, err
	}

	liqoCredentials, err := getters.GetLiqoCredentials(ctx, g.client)
	if err != nil {
		klog.Errorf("Error when getting Liqo credentials: %s", err)
		return nil, err
	}

	var transaction models.Transaction

	body := models.ReserveRequest{
		FlavorID: flavorID,
		Buyer: models.NodeIdentity{
			NodeID: g.ID.NodeID,
			IP:     g.ID.IP,
			Domain: g.ID.Domain,
		},
		ClusterID: liqoCredentials.ClusterID,
		Partition: func() *models.Partition {
			if reservation.Spec.Partition != nil {
				return parseutil.ParsePartition(reservation.Spec.Partition)
			}
			return nil
		}(),
	}

	klog.Infof("Reservation %s for flavor %s", reservation.Name, flavorID)

	if reservation.Spec.Partition != nil {
		body.Partition = parseutil.ParsePartition(reservation.Spec.Partition)
	}

	selectorBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// TODO: this url should be taken from the nodeIdentity of the flavor
	bodyBytes := bytes.NewBuffer(selectorBytes)
	url := fmt.Sprintf("http://%s%s%s", reservation.Spec.Seller.IP, ReserveFlavorPath, flavorID)

	klog.Infof("Sending request to %s", url)

	resp, err := makeRequest(ctx, "POST", url, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		klog.Errorf("Received non-OK response status code: %v", resp)
		return nil, fmt.Errorf("received non-OK response status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&transaction); err != nil {
		return nil, err
	}

	// check if transaction is not correctly set
	// This behaviuor is a possible bug of the rear-controller
	// TODO: check it
	if transaction.TransactionID == "" {
		klog.Errorf("TransactionID received is empty")
		return &transaction, fmt.Errorf("transactionID is empty")
	}

	klog.Infof("Flavor %s reserved: transaction ID %s", flavorID, transaction.TransactionID)

	g.addNewTransacion(&transaction)

	return &transaction, nil
}

// PurchaseFlavor purchases a flavor with the given flavorID.
func (g *Gateway) PurchaseFlavor(ctx context.Context, transactionID string, seller nodecorev1alpha1.NodeIdentity) (*models.ResponsePurchase, error) {
	err := checkLiqoReadiness(g.LiqoReady)
	if err != nil {
		return nil, err
	}

	var purchase models.ResponsePurchase

	// Check if the transaction exists
	transaction, err := g.GetTransaction(transactionID)
	if err != nil {
		return nil, err
	}

	body := models.PurchaseRequest{
		TransactionID: transaction.TransactionID,
	}

	selectorBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyBytes := bytes.NewBuffer(selectorBytes)
	// TODO: this url should be taken from the nodeIdentity of the flavor
	url := fmt.Sprintf("http://%s%s%s", seller.IP, PurchaseFlavorPath, transactionID)

	resp, err := makeRequest(ctx, "POST", url, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&purchase); err != nil {
		return nil, err
	}

	return &purchase, nil
}

// DiscoverFlavors is a function that returns an array of Flavor that fit the Selector by performing a get request to an http server.
func (g *Gateway) DiscoverFlavors(ctx context.Context, selector *nodecorev1alpha1.FlavorSelector) ([]*nodecorev1alpha1.Flavor, error) {
	err := checkLiqoReadiness(g.LiqoReady)
	if err != nil {
		return nil, err
	}

	var s *models.Selector
	var flavorsCR []*nodecorev1alpha1.Flavor

	if selector != nil {
		s = parseutil.ParseFlavorSelector(selector)
	}

	providers := getters.GetLocalProviders(context.Background(), g.client)

	// Send the POST request to all the servers in the list
	for _, provider := range providers {
		flavor, err := discover(ctx, s, provider)
		if err != nil {
			klog.Errorf("Error when searching Flavor: %s", err)
			return nil, err
		}
		// Check if the flavor is nil
		if flavor == nil {
			klog.Infof("No Flavors found for provider %s", provider)
		} else {
			klog.Infof("Flavor found for provider %s", provider)
			flavorsCR = append(flavorsCR, flavor)
		}
	}

	klog.Infof("Found %d flavors", len(flavorsCR))
	return flavorsCR, nil
}

func discover(ctx context.Context, s *models.Selector, provider string) (*nodecorev1alpha1.Flavor, error) {
	if s != nil {
		klog.Infof("Searching Flavor with selector %v", s)
		return searchFlavorWithSelector(ctx, s, provider)
	}
	klog.Infof("Searching Flavor with no selector")
	return searchFlavor(ctx, provider)
}

func checkLiqoReadiness(b bool) error {
	if !b {
		klog.Errorf("Liqo is not ready, please check or wait for the Liqo installation")
		return fmt.Errorf("liqo is not ready, please check or wait for the Liqo installation")
	}
	return nil
}
