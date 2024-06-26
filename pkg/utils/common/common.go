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

package common

import (
	"fmt"

	"k8s.io/klog/v2"

	advertisementv1alpha1 "github.com/fluidos-project/node/apis/advertisement/v1alpha1"
	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	reservationv1alpha1 "github.com/fluidos-project/node/apis/reservation/v1alpha1"
	"github.com/fluidos-project/node/pkg/utils/models"
	"github.com/fluidos-project/node/pkg/utils/namings"
	"github.com/fluidos-project/node/pkg/utils/parseutil"
)

// FilterFlavorsBySelector returns the Flavor CRs in the cluster that match the selector.
func FilterFlavorsBySelector(flavors []nodecorev1alpha1.Flavor, selector models.Selector) ([]nodecorev1alpha1.Flavor, error) {
	var flavorsSelected []nodecorev1alpha1.Flavor

	// Get the Flavors that match the selector
	for i := range flavors {
		f := flavors[i]
		// TODO: Not very strong and nice comparison, needs to be reviewed and improved
		if string(f.Spec.FlavorType.TypeIdentifier) == string(selector.GetSelectorType()) {
			if FilterFlavor(selector, &f) {
				flavorsSelected = append(flavorsSelected, f)
			}
		}
	}

	return flavorsSelected, nil
}

func FilterFlavor(selector models.Selector, f *nodecorev1alpha1.Flavor) bool {

	flavorTypeIdentifier, flavorTypeData, err := nodecorev1alpha1.ParseFlavorType(f)

	if err != nil {
		klog.Errorf("error parsing flavor type: %v", err)
		return false
	}

	switch flavorTypeIdentifier {
	case nodecorev1alpha1.Type_K8Slice:
		// Check if selector type matches flavor type
		if selector.GetSelectorType() != models.K8SliceNameDefault {
			klog.Errorf("selector type %s does not match flavor type %s", selector.GetSelectorType(), models.K8SliceNameDefault)
			return false
		}
		// Cast the selector to a K8Slice selector
		k8sliceSelector := selector.(*models.K8SliceSelector)
		return FilterFlavorK8Slice(k8sliceSelector, flavorTypeData.(*nodecorev1alpha1.K8Slice))
	default:
		// Flavor type not supported
		klog.Errorf("flavor type %s not supported", f.Spec.FlavorType.TypeIdentifier)
		return false
	}
}

// FilterFlavorK8Slice filters the K8Slice Flavor CRs in the cluster that match the selector.
func FilterFlavorK8Slice(selector *models.K8SliceSelector, flavorTypeK8Slice *nodecorev1alpha1.K8Slice) bool {

	// CPU Filter
	if selector.Cpu != nil {
		// Check if the flavor matches the CPU filter
		switch selector.Cpu.GetFilterType() {
		case models.MatchFilter:
			// Cast the CPU filter to a match filter
			cpuFilter := selector.Cpu.(*models.ResourceQuantityMatchFilter)
			// Check if the flavor CPU matches the filter
			if flavorTypeK8Slice.Characteristics.Cpu.Cmp(cpuFilter.Value) != 0 {
				klog.Infof("CPU Filter: %d - Flavor CPU: %d", cpuFilter.Value, flavorTypeK8Slice.Characteristics.Cpu)
				return false
			}
		case models.RangeFilter:
			// Cast the CPU filter to a range filter
			cpuFilter := selector.Cpu.(*models.ResourceQuantityRangeFilter)
			// Check if the flavor CPU is within the range
			if flavorTypeK8Slice.Characteristics.Cpu.Cmp(cpuFilter.Min) < 0 || flavorTypeK8Slice.Characteristics.Cpu.Cmp(cpuFilter.Max) > 0 {
				klog.Infof("CPU Filter: %d-%d - Flavor CPU: %d", cpuFilter.Min, cpuFilter.Max, flavorTypeK8Slice.Characteristics.Cpu)
				return false
			}
		}
	}

	// Memory Filter
	if selector.Memory != nil {
		// Check if the flavor matches the Memory filter
		switch selector.Memory.GetFilterType() {
		case models.MatchFilter:
			memoryFilter := selector.Memory.(*models.ResourceQuantityMatchFilter)
			if flavorTypeK8Slice.Characteristics.Memory.Cmp(memoryFilter.Value) != 0 {
				klog.Infof("Memory Filter: %d - Flavor Memory: %d", memoryFilter.Value, flavorTypeK8Slice.Characteristics.Memory)
				return false
			}
		case models.RangeFilter:
			memoryFilter := selector.Memory.(*models.ResourceQuantityRangeFilter)
			if flavorTypeK8Slice.Characteristics.Memory.Cmp(memoryFilter.Min) < 0 || flavorTypeK8Slice.Characteristics.Memory.Cmp(memoryFilter.Max) > 0 {
				klog.Infof("Memory Filter: %d-%d - Flavor Memory: %d", memoryFilter.Min, memoryFilter.Max, flavorTypeK8Slice.Characteristics.Memory)
				return false
			}
		}
	}

	// Pods Filter
	if selector.Pods != nil {
		// Check if the flavor matches the Pods filter
		switch selector.Pods.GetFilterType() {
		case models.MatchFilter:
			podsFilter := selector.Pods.(*models.ResourceQuantityMatchFilter)
			if flavorTypeK8Slice.Characteristics.Pods.Cmp(podsFilter.Value) != 0 {
				klog.Infof("Pods Filter: %d - Flavor Pods: %d", podsFilter.Value, flavorTypeK8Slice.Characteristics.Pods)
				return false
			}
		case models.RangeFilter:
			podsFilter := selector.Pods.(*models.ResourceQuantityRangeFilter)
			if flavorTypeK8Slice.Characteristics.Pods.Cmp(podsFilter.Min) < 0 || flavorTypeK8Slice.Characteristics.Pods.Cmp(podsFilter.Max) > 0 {
				klog.Infof("Pods Filter: %d-%d - Flavor Pods: %d", podsFilter.Min, podsFilter.Max, flavorTypeK8Slice.Characteristics.Pods)
				return false
			}
		}
	}

	// Storage Filter
	if selector.Storage != nil {
		// Check if the flavor matches the Storage filter
		switch selector.Storage.GetFilterType() {
		case models.MatchFilter:
			storageFilter := selector.Storage.(*models.ResourceQuantityMatchFilter)
			if flavorTypeK8Slice.Characteristics.Storage.Cmp(storageFilter.Value) != 0 {
				klog.Infof("Storage Filter: %d - Flavor Storage: %d", storageFilter.Value, flavorTypeK8Slice.Characteristics.Storage)
				return false
			}
		case models.RangeFilter:
			storageFilter := selector.Storage.(*models.ResourceQuantityRangeFilter)
			if flavorTypeK8Slice.Characteristics.Storage.Cmp(storageFilter.Min) < 0 || flavorTypeK8Slice.Characteristics.Storage.Cmp(storageFilter.Max) > 0 {
				klog.Infof("Storage Filter: %d-%d - Flavor Storage: %d", storageFilter.Min, storageFilter.Max, flavorTypeK8Slice.Characteristics.Storage)
				return false
			}
		}
	}

	return true

}

// FilterPeeringCandidate filters the peering candidate based on the solver's flavor selector.
func FilterPeeringCandidate(selector *nodecorev1alpha1.Selector, pc *advertisementv1alpha1.PeeringCandidate) bool {
	s := parseutil.ParseFlavorSelector(selector)
	return FilterFlavor(s, &pc.Spec.Flavor)
}

// CheckSelector ia a func to check if the syntax of the Selector is right.
// Strict and range syntax cannot be used together.
func CheckSelector(selector models.Selector) error {
	// Parse the selector to check the syntax
	switch selector.GetSelectorType() {
	case models.K8SliceNameDefault:
		k8sliceSelector := selector.(*models.K8SliceSelector)
		klog.Infof("Checking K8Slice selector: %v", k8sliceSelector)
		// Nothing is compulsory in the K8Slice selector
		return nil
	default:
		return fmt.Errorf("selector type %s not supported", selector.GetSelectorType())
	}
}

// SOLVER PHASE SETTERS

// DiscoveryStatusCheck checks the status of the discovery.
func DiscoveryStatusCheck(solver *nodecorev1alpha1.Solver, discovery *advertisementv1alpha1.Discovery) {
	if discovery.Status.Phase.Phase == nodecorev1alpha1.PhaseSolved {
		klog.Infof("Discovery %s has found candidates: %s", discovery.Name, discovery.Status.PeeringCandidateList)
		solver.Status.FindCandidate = nodecorev1alpha1.PhaseSolved
		solver.Status.DiscoveryPhase = nodecorev1alpha1.PhaseSolved
		solver.SetPhase(nodecorev1alpha1.PhaseRunning, "Solver has completed the Discovery phase")
	}
	if discovery.Status.Phase.Phase == nodecorev1alpha1.PhaseFailed {
		klog.Infof("Discovery %s has failed. Reason: %s", discovery.Name, discovery.Status.Phase.Message)
		klog.Infof("Peering candidate not found, Solver %s failed", solver.Name)
		solver.Status.FindCandidate = nodecorev1alpha1.PhaseFailed
		solver.Status.DiscoveryPhase = nodecorev1alpha1.PhaseFailed
	}
	if discovery.Status.Phase.Phase == nodecorev1alpha1.PhaseTimeout {
		klog.Infof("Discovery %s has timed out", discovery.Name)
		solver.Status.FindCandidate = nodecorev1alpha1.PhaseTimeout
		solver.Status.DiscoveryPhase = nodecorev1alpha1.PhaseTimeout
		solver.SetPhase(nodecorev1alpha1.PhaseTimeout, "Discovery has expired before finding a candidate")
	}
	if discovery.Status.Phase.Phase == nodecorev1alpha1.PhaseRunning {
		klog.Infof("Discovery %s is running", discovery.Name)
		solver.SetDiscoveryStatus(nodecorev1alpha1.PhaseRunning)
	}
	if discovery.Status.Phase.Phase == nodecorev1alpha1.PhaseIdle {
		klog.Infof("Discovery %s is idle", discovery.Name)
		solver.SetDiscoveryStatus(nodecorev1alpha1.PhaseIdle)
	}
}

// ReservationStatusCheck checks the status of the reservation.
func ReservationStatusCheck(solver *nodecorev1alpha1.Solver, reservation *reservationv1alpha1.Reservation) {
	klog.Infof("Reservation %s is in phase %s", reservation.Name, reservation.Status.Phase.Phase)
	flavorName := namings.RetrieveFlavorNameFromPC(reservation.Spec.PeeringCandidate.Name)
	if reservation.Status.Phase.Phase == nodecorev1alpha1.PhaseSolved {
		klog.Infof("Reservation %s has reserved and purchase the flavor %s", reservation.Name, flavorName)
		solver.Status.ReservationPhase = nodecorev1alpha1.PhaseSolved
		solver.Status.ReserveAndBuy = nodecorev1alpha1.PhaseSolved
		solver.Status.Contract = reservation.Status.Contract
		solver.SetPhase(nodecorev1alpha1.PhaseRunning, "Reservation: Flavor reserved and purchased")
	}
	if reservation.Status.Phase.Phase == nodecorev1alpha1.PhaseFailed {
		klog.Infof("Reservation %s has failed. Reason: %s", reservation.Name, reservation.Status.Phase.Message)
		solver.Status.ReservationPhase = nodecorev1alpha1.PhaseFailed
		solver.Status.ReserveAndBuy = nodecorev1alpha1.PhaseFailed
		solver.SetPhase(nodecorev1alpha1.PhaseFailed, "Reservation: Flavor reservation and purchase failed")
	}
	if reservation.Status.Phase.Phase == nodecorev1alpha1.PhaseRunning {
		if reservation.Status.ReservePhase == nodecorev1alpha1.PhaseRunning {
			klog.Infof("Reservation %s is running", reservation.Name)
			solver.SetPhase(nodecorev1alpha1.PhaseRunning, "Reservation: Reserve is running")
		}
		if reservation.Status.PurchasePhase == nodecorev1alpha1.PhaseRunning {
			klog.Infof("Purchasing %s is running", reservation.Name)
			solver.SetPhase(nodecorev1alpha1.PhaseRunning, "Reservation: Purchase is running")
		}
	}
	if reservation.Status.Phase.Phase == nodecorev1alpha1.PhaseIdle {
		klog.Infof("Reservation %s is idle", reservation.Name)
		solver.SetReservationStatus(nodecorev1alpha1.PhaseIdle)
	}
}
