package v1beta1

import (
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OrganizationSpec defines the desired state of the Organization
type OrganizationSpec struct {
	//+kubebuilder:validation:required
	// Name the short name of an organization
	Name string `json:"name,omitempty"`

	//+kubebuilder:validation:required
	// DisplayName the long name of an organization
	DisplayName string `json:"displayName,omitempty"`

	//+kubebuilder:validation:required
	// DatabaseRef is the name of the Database this Organization is assigned to
	DatabaseRef string `json:"databaseRef,omitempty"`

	//+kubebuilder:validation:required
	// NamedGraphs are the suffix graph names for this organization. The prefix can be found in the Database resource.
	NamedGraphs []string `json:"namedGraphs,omitempty"`
}

// OrganizationStatus defines the observed state of the Organization
type OrganizationStatus struct {
	Name        string                      `json:"name,omitempty"`
	DisplayName string                      `json:"displayName,omitempty"`
	DatabaseRef string                      `json:"databaseRef,omitempty"`
	NamedGraphs []string                    `json:"namedGraphs,omitempty"`
	Conditions  []v1alpha1.StardogCondition `json:"conditions,omitempty"`
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
