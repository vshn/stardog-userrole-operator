/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StardogUserSpec defines the desired state of StardogUser
type StardogUserSpec struct {
	// StardogInstanceRef references a StardogInstance object.
	StardogInstanceRef string `json:"stardogInstanceRef,omitempty"`
	// StardogUserCredentialsSpec describes the credentials of a Stardog user
	Credentials StardogUserCredentialsSpec `json:"credentials,omitempty"`
	// Roles describe a list of StardogRoles assigned to a Stardog user. The names are referring the StardogRole metadata names, not the role name that is supposed to be in Stardog.
	Roles []string `json:"roles,omitempty"`
}

// StardogUserCredentialsSpec specifies the password of a Stardog user
type StardogUserCredentialsSpec struct {
	// Namespace specifies the namespace of the Secret referenced in SecretRef.
	// Defaults to .metadata.namespace.
	Namespace string `json:"namespace,omitempty"`
	// SecretRef references the v1/Secret name which contains the "username" and "password" keys.
	SecretRef string `json:"secretRef,omitempty"`
}

// StardogUserStatus defines the observed state of StardogUser
type StardogUserStatus struct {
	// Conditions contain the states of the StardogUser. A StardogUser is considered Ready when the user has been
	// persisted to Stardog DB.
	Conditions []StardogCondition `json:"conditions,omitempty" patchStrategy:"merge"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// StardogUser is the Schema for the stardogusers API
type StardogUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StardogUserSpec   `json:"spec,omitempty"`
	Status StardogUserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StardogUserList contains a list of StardogUser
type StardogUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StardogUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StardogUser{}, &StardogUserList{})
}
