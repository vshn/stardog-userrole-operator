package v1beta1

import (
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StardogInstanceRef contains name and namespace for a stardog instance
type StardogInstanceRef struct {
	//+kubebuilder:validation:required
	Name string `json:"name,omitempty"`

	//+kubebuilder:validation:required
	Namespace string `json:"namespace,omitempty"`
}

// DatabaseSpec defines the desired state of the Database
type DatabaseSpec struct {
	//+kubebuilder:validation:required
	// DatabaseName the database name that has to be created in the Stardog server
	DatabaseName string `json:"databaseName,omitempty"`

	//+kubebuilder:validation:optional
	// AddUserForNonHiddenGraphs a dynamically managed user of this db with custom permissions
	// Mainly used to not have access to hidden graphs
	AddUserForNonHiddenGraphs string `json:"addUserForNonHiddenGraphs,omitempty"`

	//+kubebuilder:validation:optional
	// Options is the Stardog configuration options for this database. Only json input is valid.
	Options string `json:"options,omitempty"`

	//+kubebuilder:validation:required
	// StardogInstanceRefs contains the reference to the Stardog instance the database should exist in
	StardogInstanceRefs []StardogInstanceRef `json:"stardogInstanceRefs,omitempty"`

	//+kubebuilder:validation:required
	// NamedGraphPrefix a prefix for a Stardog Named Graph.
	NamedGraphPrefix string `json:"namedGraphPrefix,omitempty"`
}

// DatabaseStatus defines the observed state of the Database
type DatabaseStatus struct {
	DatabaseName              string                      `json:"databaseName,omitempty"`
	AddUserForNonHiddenGraphs string                      `json:"addUserForNonHiddenGraphs,omitempty"`
	NamedGraphPrefix          string                      `json:"namedGraphPrefix,omitempty"`
	Options                   string                      `json:"options,omitempty"`
	StardogInstanceRefs       []StardogInstanceRef        `json:"stardogInstanceRef,omitempty"`
	Conditions                []v1alpha1.StardogCondition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

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

// NewStardogInstanceRef creates a new StardogInstanceRef from name and namespace
func NewStardogInstanceRef(name, namespace string) StardogInstanceRef {
	return StardogInstanceRef{
		Name:      name,
		Namespace: namespace,
	}
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}
