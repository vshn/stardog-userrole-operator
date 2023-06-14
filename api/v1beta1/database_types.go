package v1beta1

import (
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatabaseSpec defines the desired state of the Database
type DatabaseSpec struct {
	//+kubebuilder:validation:required
	DatabaseName string `json:"databaseName,omitempty"`

	//+kubebuilder:validation:optional
	Options string `json:"options,omitempty"`

	//+kubebuilder:validation:required
	// Instance contains the reference to the Stardog instance the database should exist in
	StardogInstanceRef string `json:"stardogInstanceRef,omitempty"`

	//+kubebuilder:validation:required
	NamedGraphPrefix string `json:"namedGraphPrefix,omitempty"`
}

// DatabaseStatus defines the observed state of the Database
type DatabaseStatus struct {
	DatabaseName       string                      `json:"databaseName,omitempty"`
	Options            string                      `json:"options,omitempty"`
	StardogInstanceRef string                      `json:"stardogInstanceRef,omitempty"`
	Conditions         []v1alpha1.StardogCondition `json:"conditions,omitempty"`
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
