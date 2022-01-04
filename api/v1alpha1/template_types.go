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
package v1alpha1
*/
package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// TemplateSpec defines the desired state of a Template
type TemplateSpec struct {
	// Body is the parameterized template string.
	Body string `json:"body"`

	// Parameters is the set of strings within Body to treat as parameters.
	// +kubebuilder:validation:MinItems:=1
	Parameters []string `json:"parameters,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=combo,scope=Cluster

// Template is a custom resource that represents a parameterized set of Kubernetes manifests.
type Template struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TemplateSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// TemplateList contains a list of Template
type TemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Template `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Template{}, &TemplateList{})
}
