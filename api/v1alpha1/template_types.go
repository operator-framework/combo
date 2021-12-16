package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// TemplateSpec defines the desired state of Template
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
