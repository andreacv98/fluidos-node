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
func FilterFlavorsBySelector(flavors []nodecorev1alpha1.Flavor, selector *models.Selector) ([]nodecorev1alpha1.Flavor, error) {
	var flavorsSelected []nodecorev1alpha1.Flavor

	// Get the Flavors that match the selector
	for i := range flavors {
		f := flavors[i]
		if string(f.Spec.FlavorType.TypeIdentifier) == selector.FlavorType {
			if FilterFlavor(selector, &f) {
				flavorsSelected = append(flavorsSelected, f)
			}
		}
	}

	return flavorsSelected, nil
}

func FilterFlavor(selector *models.Selector, f *nodecorev1alpha1.Flavor) bool {

	err, FlavorTypeIdentifier, flavorTypeData := parseutil.ParseFlavorType(f)

	if err != nil {
		klog.Errorf("error parsing flavor type: %v", err)
		return false
	}

	switch FlavorTypeIdentifier {
	case nodecorev1alpha1.Type_K8Slice:
		return FilterFlavorK8Slice(selector, flavorTypeData.(*nodecorev1alpha1.K8Slice))
	default:
		// Flavor type not supported
		klog.Errorf("flavor type %s not supported", f.Spec.FlavorType.TypeIdentifier)
		return false
	}
}

// FilterFlavorK8Slice filters the K8Slice Flavor CRs in the cluster that match the selector.
func FilterFlavorK8Slice(selector *models.Selector, flavorTypeK8Slice *nodecorev1alpha1.K8Slice) bool {

	if selector.MatchSelector != nil {
		return FilterFlavorK8Slice(selector, flavorTypeK8Slice)
	}

	if selector.RangeSelector != nil && selector.MatchSelector == nil {
		return filterK8SliceByRangeSelector(selector, flavorTypeK8Slice)
	}

	return true

}

func filterK8SliceByMatchSelector(selector *models.Selector, k8sliceFlavor *nodecorev1alpha1.K8Slice) bool {
	if selector.MatchSelector.CPU.CmpInt64(0) == 0 && k8sliceFlavor.Characteristics.Cpu.Cmp(selector.MatchSelector.CPU) != 0 {
		klog.Infof("MatchSelector Cpu: %d - Flavor Cpu: %d", selector.MatchSelector.CPU, k8sliceFlavor.Characteristics.Cpu)
		return false
	}

	if selector.MatchSelector.Memory.CmpInt64(0) == 0 && k8sliceFlavor.Characteristics.Memory.Cmp(selector.MatchSelector.Memory) != 0 {
		klog.Infof("MatchSelector Memory: %d - Flavor Memory: %d", selector.MatchSelector.Memory, k8sliceFlavor.Characteristics.Memory)
		return false
	}

	if selector.MatchSelector.Pods.CmpInt64(0) == 0 && k8sliceFlavor.Characteristics.Pods.Cmp(selector.MatchSelector.Pods) != 0 {
		klog.Infof("MatchSelector Pods: %d - Flavor Pods: %d", selector.MatchSelector.Pods, k8sliceFlavor.Characteristics.Pods)
		return false
	}

	if selector.MatchSelector.Storage.CmpInt64(0) == 0 &&
		k8sliceFlavor.Characteristics.Storage.Cmp(selector.MatchSelector.Storage) != 0 {
		klog.Infof("MatchSelector EphemeralStorage: %d - Flavor EphemeralStorage: %d",
			selector.MatchSelector.Storage, k8sliceFlavor.Characteristics.Storage)
		return false
	}

	if selector.MatchSelector.Gpu != nil && k8sliceFlavor.Characteristics.Gpu == nil {
		if selector.MatchSelector.Gpu.Cmp(*k8sliceFlavor.Characteristics.Gpu) != 0 {
			klog.Infof("MatchSelector GPU: %d - Flavor GPU: %d", selector.MatchSelector.Gpu, k8sliceFlavor.Characteristics.Gpu)
			return false
		}
	}

	return true
}

func filterK8SliceByRangeSelector(selector *models.Selector, k8sliceFlavor *nodecorev1alpha1.K8Slice) bool {
	if selector.RangeSelector.MinCPU.CmpInt64(0) != 0 && k8sliceFlavor.Characteristics.Cpu.Cmp(selector.RangeSelector.MinCPU) < 0 {
		klog.Infof("RangeSelector MinCpu: %d - Flavor Cpu: %d", selector.RangeSelector.MinCPU, k8sliceFlavor.Characteristics.Cpu)
		return false
	}

	if selector.RangeSelector.MinMemory.CmpInt64(0) != 0 && k8sliceFlavor.Characteristics.Memory.Cmp(selector.RangeSelector.MinMemory) < 0 {
		klog.Infof("RangeSelector MinMemory: %d - Flavor Memory: %d", selector.RangeSelector.MinMemory, k8sliceFlavor.Characteristics.Memory)
		return false
	}

	if selector.RangeSelector.MinPods.CmpInt64(0) != 0 && k8sliceFlavor.Characteristics.Pods.Cmp(selector.RangeSelector.MinPods) < 0 {
		klog.Infof("RangeSelector MinPods: %d - Flavor Pods: %d", selector.RangeSelector.MinPods, k8sliceFlavor.Characteristics.Pods)
		return false
	}

	if selector.RangeSelector.MinStorage.CmpInt64(0) != 0 && k8sliceFlavor.Characteristics.Storage.Cmp(selector.RangeSelector.MinStorage) < 0 {
		klog.Infof("RangeSelector MinEph: %d - Flavor EphemeralStorage: %d", selector.RangeSelector.MinStorage, k8sliceFlavor.Characteristics.Storage)
		return false
	}

	if selector.RangeSelector.MinGpu != nil && k8sliceFlavor.Characteristics.Gpu != nil {
		if selector.RangeSelector.MinGpu.Cmp(*k8sliceFlavor.Characteristics.Gpu) > 0 {
			klog.Infof("RangeSelector MinGpu: %d - Flavor Gpu: %d", selector.RangeSelector.MinGpu, k8sliceFlavor.Characteristics.Gpu)
			return false
		}
	}

	if selector.RangeSelector.MaxCPU.CmpInt64(0) != 0 && k8sliceFlavor.Characteristics.Cpu.Cmp(selector.RangeSelector.MaxCPU) > 0 {
		return false
	}

	if selector.RangeSelector.MaxMemory.CmpInt64(0) != 0 && k8sliceFlavor.Characteristics.Memory.Cmp(selector.RangeSelector.MaxMemory) > 0 {
		return false
	}

	if selector.RangeSelector.MaxPods.CmpInt64(0) != 0 && k8sliceFlavor.Characteristics.Pods.Cmp(selector.RangeSelector.MaxPods) > 0 {
		return false
	}

	if selector.RangeSelector.MaxStorage.CmpInt64(0) != 0 && k8sliceFlavor.Characteristics.Storage.Cmp(selector.RangeSelector.MaxStorage) > 0 {
		return false
	}

	if selector.RangeSelector.MaxGpu != nil && k8sliceFlavor.Characteristics.Gpu != nil {
		if selector.RangeSelector.MaxGpu.Cmp(*k8sliceFlavor.Characteristics.Gpu) < 0 {
			klog.Infof("RangeSelector MaxGpu: %d - Flavor Gpu: %d", selector.RangeSelector.MaxGpu, k8sliceFlavor.Characteristics.Gpu)
			return false
		}
	}

	return true
}

// FilterPeeringCandidate filters the peering candidate based on the solver's flavor selector.
func FilterPeeringCandidate(selector *nodecorev1alpha1.FlavorSelector, pc *advertisementv1alpha1.PeeringCandidate) bool {
	s := parseutil.ParseFlavorSelector(selector)
	return FilterFlavor(s, &pc.Spec.Flavor)
}

// CheckSelector ia a func to check if the syntax of the Selector is right.
// Strict and range syntax cannot be used together.
func CheckSelector(selector *models.Selector) error {
	if selector.MatchSelector != nil && selector.RangeSelector != nil {
		return fmt.Errorf("selector syntax error: strict and range syntax cannot be used together")
	}
	return nil
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
