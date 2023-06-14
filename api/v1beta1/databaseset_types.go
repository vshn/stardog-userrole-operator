package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatabaseSpec defines the desired state of the Database
type DatabaseSetSpec struct {
	//+kubebuilder:validation:optional
	NamedGraphPrefix string `json:"namedGraphPrefix,omitempty"`

	//+kubebuilder:validation:required
	// Instances contains the references to the Stardog instances the database should exist in
	Instances []string `json:"stardogInstanceRefs,omitempty"`
}

// DatabaseSetStatus defines the observed state of the Database
type DatabaseSetStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DatabaseSet is the Schema for the databasesets API
type DatabaseSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSetSpec   `json:"spec,omitempty"`
	Status DatabaseSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DatabaseSetList contains a list of DatabaseSet
type DatabaseSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DatabaseSet{}, &DatabaseSetList{})
}
