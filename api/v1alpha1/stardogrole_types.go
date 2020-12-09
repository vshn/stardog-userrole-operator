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

// StardogRoleSpec defines the desired state of StardogRole
type StardogRoleSpec struct {
	// RoleName describes (overrides) the name of a role that will be maintained in a Stardog instance.
	// Defaults to .metadata.name.
	// +kubebuilder:validation:Optional
	RoleName string `json:"roleName,omitempty"`
	// StardogInstanceRef references the StardogInstance object in which the role is maintained.
	// +kubebuilder:validation:Required
	StardogInstanceRef string `json:"stardogInstanceRef,omitempty"`
	// Permissions lists the permissions assigned to a role
	// +kubebuilder:validation:Optional
	Permissions []StardogPermissionSpec `json:"permissions,omitempty"`
}

// StardogPermissionSpec defines a Stardog permission assigned to a Role
type StardogPermissionSpec struct {
	// Action describes the action a specific permission is assigned to
	// +kubebuilder:validation:Enum=ALL;CREATE;DELETE;READ;WRITE;GRANT;REVOKE;EXECUTE
	// +kubebuilder:validation:Required
	Action string `json:"action,omitempty"`
	// ResourceType describes the type of resource a specific permission is assigned to
	// +kubebuilder:validation:Enum=DB;USER;ROLE;ADMIN;METADATA;NAMEDGRAPH;VIRTUALGRAPH;ICVCONSTRAINTS
	// +kubebuilder:validation:Required
	ResourceType string `json:"resourceType,omitempty"`
	// Resources is a list of permission objects that get each targeted by the action and resource type properties
	// +kubebuilder:validation:Required
	Resources []string `json:"resources,omitempty"`
}

// StardogRoleStatus defines the observed state of StardogRole
type StardogRoleStatus struct {
	// Conditions contain the states of the StardogRole. A StardogRole is considered Ready when the role has been
	// persisted to Stardog DB.
	Conditions []StardogCondition `json:"conditions,omitempty" patchStrategy:"merge"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// StardogRole is the Schema for the stardogroles API
type StardogRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StardogRoleSpec   `json:"spec,omitempty"`
	Status StardogRoleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StardogRoleList contains a list of StardogRole
type StardogRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StardogRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StardogRole{}, &StardogRoleList{})
}
