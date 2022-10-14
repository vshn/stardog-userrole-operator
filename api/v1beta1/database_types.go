package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatabaseSpec defines the desired state of the Database
type DatabaseSpec struct {
	//+kubebuilder:validation:required
	// Instances contains the references to the Stardog instances the database should exist in
	Instances []StardogInstanceRef `json:"stardogInstanceRef,omitempty"`

	//+kubebuilder:validation:optional
	NamedGraphPrefix string `json:"url,omitempty"`
}

// StardogInstanceRef defines a reference to a Stardog instance
type StardogInstanceRef struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// DatabaseStatus defines the observed state of the Database
type DatabaseStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Database is the Schema for the databases API
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSpec   `json:"spec,omitempty"`
	Status DatabaseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DatabaseList contains a list of Database
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}
