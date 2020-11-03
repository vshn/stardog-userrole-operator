/*
Licensed under the Apache License, Version 2.0 (the "License");
http://www.apache.org/licenses/LICENSE-2.0
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StardogInstanceSpec defines the desired state of StardogInstance
type StardogInstanceSpec struct {
	// ServerUrl describes the url of the Stardog Instace
	ServerUrl string `json:"serverUrl,omitempty"`
	// AdminCredentials references the credentials that gives administrative access to the Stardog instance.
	// This is used by the Operator to make changes in the roles, permissions and users.
	AdminCredentials StardogUserCredentialsSpec `json:"adminCredentials,omitempty"`
}

// StardogInstanceStatus defines the observed state of StardogInstance
type StardogInstanceStatus struct {
	// Conditions contain the states of the StardogInstance. A StardogInstance is considered Ready when the Admin user can make authorized REST API calls.
	Conditions []StardogCondition `json:"conditions,omitempty" patchStrategy:"merge"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// StardogInstance contains information about a Stardog server or cluster.
type StardogInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StardogInstanceSpec   `json:"spec,omitempty"`
	Status StardogInstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StardogInstanceList contains a list of StardogInstance
type StardogInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StardogInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StardogInstance{}, &StardogInstanceList{})
}
