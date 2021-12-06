package conditions

import (
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
	ProccessedCombinationsCondition = metav1.Condition{
		Type:    "Finished",
		Status:  "True",
		Message: "Combination has processed successfully.",
		Reason:  "Processed",
	}
)

func NewConditions(transitionTime time.Time, conditions ...metav1.Condition) []metav1.Condition {
	var newConditions []metav1.Condition
	for _, condition := range conditions {
		condition.LastTransitionTime = metav1.NewTime(transitionTime)
		newConditions = append(newConditions, condition)
	}
	return newConditions
}
