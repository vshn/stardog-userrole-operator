package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StardogInstanceSpec defines the desired state of StardogInstance
type StardogInstanceSpec struct {
	// ServerUrl describes the url of the Stardog Instance
	// +kubebuilder:validation:Required
	ServerUrl string `json:"serverUrl,omitempty"`
	// AdminCredentials references the credentials that gives administrative access to the Stardog instance.
	// This is used by the Operator to make changes in the roles, permissions and users.
	// +kubebuilder:validation:Required
	AdminCredentials StardogUserCredentialsSpec `json:"adminCredentials,omitempty"`
	// Disabled whether this instance is disabled or enabled for operator to recycle resources
	Disabled bool `json:"disabled,omitempty"`
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
