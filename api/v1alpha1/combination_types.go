package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	// Phase represents a human-readable description of where the combinations is
	Phase string `json:"phase,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=combo,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
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

func init() {
	SchemeBuilder.Register(&Combination{}, &CombinationList{})
}
