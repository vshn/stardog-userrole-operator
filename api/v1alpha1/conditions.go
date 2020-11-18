package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StardogCondition describes a status condition of a StardogRole
type StardogCondition struct {
	Type               StardogConditionType   `json:"type"`
	Status             corev1.ConditionStatus `json:"status"`
	LastTransitionTime metav1.Time            `json:"lastTransitionTime,omitempty"`
	Reason             string                 `json:"reason,omitempty"`
	Message            string                 `json:"message,omitempty"`
}

type StardogConditionType string

type StardogConditionMap map[StardogConditionType]StardogCondition

const (
	// StardogReady tracks if the Stardog has been successfully reconciled.
	StardogReady StardogConditionType = "Ready"
	// StardogErrored is given when the object could not be reconciled with Stardog.
	StardogErrored StardogConditionType = "Errored"
	// StardogInvalid is given when the the object contains invalid properties. The object will not be further
	// reconciled until the issue is fixed.
	StardogInvalid StardogConditionType = "Invalid"
	// StardogTerminating is given when the the Stardog resource is to be deleted but the object's finalizers cannot
	// be cleared for a reason.
	StardogTerminating StardogConditionType = "StardogTerminating"

	ReasonFailed      = "SynchronizationFailed"
	ReasonSucceeded   = "SynchronizationSucceeded"
	ReasonSpecInvalid = "InvalidSpec"
	ReasonTerminating = "StardogTerminating"
)
