package conditions

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ProcessingCombinationsCondition = metav1.Condition{
		Type:    "InProgress",
		Status:  "Unknown",
		Message: "Combination is not yet ready.",
		Reason:  "Processing",
	}
	TemplateNotFoundCondition = metav1.Condition{
		Type:    "Invalid",
		Status:  "False",
		Message: "Combination cannot be processed as template was not retrievable.",
		Reason:  "TemplateNotFound",
	}
	TemplateBodyInvalid = metav1.Condition{
		Type:    "Invalid",
		Status:  "False",
		Message: "Combination cannot be processed as template body was not valid.",
		Reason:  "TemplateBodyInvalid",
	}
	ManifestGenerationFailed = metav1.Condition{
		Type:    "Finished",
		Status:  "False",
		Message: "Generation of the manifest combinations failed.",
		Reason:  "EvaluationsInvalid",
	}
	ProccessedCombinationsCondition = metav1.Condition{
		Type:    "Finished",
		Status:  "True",
		Message: "Combination has processed successfully.",
		Reason:  "Processed",
	}
)

func NewConditions(transitionTime time.Time, err error, conditions ...metav1.Condition) []metav1.Condition {
	newConditions := []metav1.Condition{}
	for _, condition := range conditions {
		if err != nil {
			condition.Message += fmt.Sprintf("Error: %s", err.Error())
		}

		condition.LastTransitionTime = metav1.NewTime(transitionTime)

		newConditions = append(newConditions, condition)
	}
	return newConditions
}
