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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var allocationlog = logf.Log.WithName("allocation-resource")

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (r *Allocation) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//nolint:lll // kubebuilder directives are too long, but they must be on the same line
//+kubebuilder:webhook:path=/mutate-nodecore-fluidos-eu-v1alpha1-allocation,mutating=true,failurePolicy=fail,sideEffects=None,groups=nodecore.fluidos.eu,resources=allocations,verbs=create;update,versions=v1alpha1,name=mallocation.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Allocation{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *Allocation) Default() {
	allocationlog.Info("DEFAULT WEBHOOK")
	allocationlog.Info("default", "name", r.Name)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//nolint:lll // kubebuilder directives are too long, but they must be on the same line
//+kubebuilder:webhook:path=/validate-nodecore-fluidos-eu-v1alpha1-allocation,mutating=false,failurePolicy=fail,sideEffects=None,groups=nodecore.fluidos.eu,resources=allocations,verbs=create;update,versions=v1alpha1,name=vallocation.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Allocation{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *Allocation) ValidateCreate() (admission.Warnings, error) {
	allocationlog.Info("VALIDATE CREATE WEBHOOK")
	allocationlog.Info("validate create", "name", r.Name)

	if err := validateAllocation(r); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *Allocation) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	allocationlog.Info("VALIDATE UPDATE WEBHOOK")
	allocationlog.Info("validate update", "name", r.Name)

	allocationlog.Info("old", "old", old)

	if err := validateAllocation(r); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *Allocation) ValidateDelete() (admission.Warnings, error) {
	allocationlog.Info("VALIDATE DELETE WEBHOOK")
	allocationlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func validateAllocation(allocation *Allocation) error {
	if allocation == nil {
		return nil
	}

	// Validate creation of Allocation checking AllocationType->typeIdentifier matches the struct inside the AllocationType->TypeData
	typeIdentifier, _, err := ParseResourceReference(&allocation.Spec.ResourceReference)
	if err != nil {
		return err
	}
	switch typeIdentifier {
	case TypeK8Slice:
		allocationlog.Info("Allocation Flavor Type is K8Slice")
	case TypeVM:
		allocationlog.Info("Allocation Flavor Type is VM")
	case TypeService:
		allocationlog.Info("Allocation Flavor Type is Service")
	default:
		allocationlog.Info("Allocation Flavor Type is not valid")
	}

	return nil
}
