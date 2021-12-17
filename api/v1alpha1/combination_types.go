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

import (
	metautils "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	TypeInvalid    = "Invalid"
	TypeFinished   = "Finished"
	TypeInProgress = "InProgress"

	ReasonProcessing          = "Processing"
	ReasonTemplateNotFound    = "TemplateNotFound"
	ReasonTemplateBodyInvalid = "TemplateBodyInvalid"
	ReasonEvaluationsInvalid  = "EvaluationsInvalid"
	ReasonProcessed           = "Processed"
)

// CombinationSpec defines arguments that replace parameters within the given template
type CombinationSpec struct {
	// Template is the name of the template to evaluate.
	// +kubebuilder:validation:Pattern=[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*
	Template string `json:"template"`

	// Arguments contains the list of values to use for each parameter in the combination.
	// +kubebuilder:validation:MinItems:=1
	Arguments []Argument `json:"arguments,omitempty"`
}

// Argument defines a key and values for it that will be replaced in a template
type Argument struct {
	// Key defines what is going to be replaced in the template
	Key string `json:"key"`

	// Values defines the options to replace the defined key
	// +kubebuilder:validation:MinItems:=1
	Values []string `json:"values"`
}

// CombinationStatus defines the observed state of Combination
type CombinationStatus struct {
	// Conditions represents the current condition of the Combination.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Represents the evaluation to this combination once processed
	Evaluations []string `json:"evaluations,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=combo,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Combination is the Schema for a combination
type Combination struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CombinationSpec   `json:"spec"`
	Status CombinationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CombinationList contains a list of Combination
type CombinationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Combination `json:"items"`
}

// SetCondition sets the condition if it has not already been set
func (c *Combination) SetStatusCondition(condition metav1.Condition) {
	metautils.SetStatusCondition(&c.Status.Conditions, condition)
}

func init() {
	SchemeBuilder.Register(&Combination{}, &CombinationList{})
}
