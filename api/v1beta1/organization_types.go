package v1beta1

import (
	"fmt"
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OrganizationSpec defines the desired state of an Organization
type OrganizationSpec struct {
	// +kubebuilder:validation:required
	// Name is the short name of an organization
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:required
	// DisplayName is the long name of an organization
	DisplayName string `json:"displayName,omitempty"`

	// +kubebuilder:validation:required
	// DatabaseRef is the name of the Database this Organization is assigned to
	DatabaseRef string `json:"databaseRef,omitempty"`

	// +kubebuilder:validation:required
	// NamedGraphs are the suffix graph names for this organization. The prefix can be found in the Database resource.
	// The final graphs is defined as prefix + "/" + orgName + "/" suffix
	NamedGraphs []NamedGraph `json:"namedGraphs,omitempty"`
}

// NamedGraph defines the name and if necessary add another hidden named graph for this named graph
type NamedGraph struct {
	// +kubebuilder:validation:required
	// The name of the Named Graph
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:optional
	// +kubebuilder:default=false
	// AddHidden adds another graph with the same name but with a prefix "<named-graph-name/hidden>"
	AddHidden bool `json:"addHidden,omitempty"`
}

// OrganizationStatus defines the observed state of the Organization
type OrganizationStatus struct {
	Name                string                      `json:"name,omitempty"`
	DisplayName         string                      `json:"displayName,omitempty"`
	DatabaseRef         string                      `json:"databaseRef,omitempty"`
	NamedGraphs         []NamedGraph                `json:"namedGraphs,omitempty"`
	StardogInstanceRefs []StardogInstanceRef        `json:"stardogInstanceRefs,omitempty"`
	Conditions          []v1alpha1.StardogCondition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Organization is the Schema for the organizations API
type Organization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrganizationSpec   `json:"spec,omitempty"`
	Status OrganizationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OrganizationList contains a list of Organization
type OrganizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Organization `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Organization{}, &OrganizationList{})
}

func GetNamedGraphNames(namedGraphs []NamedGraph) []string {
	names := make([]string, len(namedGraphs))
	for _, graph := range namedGraphs {
		names = append(names, graph.Name)
	}
	return names
}

func GetHiddenNamedGraphNames(namedGraphs []NamedGraph) []string {
	names := make([]string, len(namedGraphs))
	for _, graph := range namedGraphs {
		if graph.AddHidden {
			names = append(names, graph.Name)
		}
	}
	return names
}

// FindNamedGraphByName finds the NamedGraph by name from a slice of NamedGraphs
func FindNamedGraphByName(graphs []NamedGraph, name string) (NamedGraph, error) {
	for _, g := range graphs {
		if g.Name == name {
			return g, nil
		}
	}
	return NamedGraph{}, fmt.Errorf("cannot find graph from name %s", name)
}
