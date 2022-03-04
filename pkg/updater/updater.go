package updater

import (
	"context"
	"reflect"

	"github.com/operator-framework/combo/api/v1alpha1"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func New(client client.Client) Updater {
	return Updater{
		client: client,
	}
}

type Updater struct {
	client            client.Client
	updateStatusFuncs []UpdateStatusFunc
}

type UpdateStatusFunc func(combo *v1alpha1.CombinationStatus) bool

func (u *Updater) UpdateStatus(fs ...UpdateStatusFunc) {
	u.updateStatusFuncs = append(u.updateStatusFuncs, fs...)
}

func (u *Updater) Apply(ctx context.Context, c *v1alpha1.Combination) error {
	backoff := retry.DefaultRetry

	return retry.RetryOnConflict(backoff, func() error {
		if err := u.client.Get(ctx, client.ObjectKeyFromObject(c), c); err != nil {
			return err
		}
		needsStatusUpdate := false
		for _, f := range u.updateStatusFuncs {
			needsStatusUpdate = f(&c.Status) || needsStatusUpdate
		}
		if needsStatusUpdate {
			log.FromContext(ctx).Info("applying status changes")
			return u.client.Status().Update(ctx, c)
		}
		return nil
	})
}

func EnsureCondition(condition metav1.Condition) UpdateStatusFunc {
	return func(status *v1alpha1.CombinationStatus) bool {
		existing := meta.FindStatusCondition(status.Conditions, condition.Type)
		if existing == nil || !conditionsSemanticallyEqual(*existing, condition) {
			meta.SetStatusCondition(&status.Conditions, condition)
			return true
		}
		return false
	}
}

func EnsureEvaluations(evaluation []string) UpdateStatusFunc {
	return func(status *v1alpha1.CombinationStatus) bool {
		if reflect.DeepEqual(status.Evaluations, evaluation) {
			return false
		}
		status.Evaluations = evaluation
		return true
	}
}

func conditionsSemanticallyEqual(a, b metav1.Condition) bool {
	return a.Type == b.Type && a.Status == b.Status && a.Reason == b.Reason && a.Message == b.Message && a.ObservedGeneration == b.ObservedGeneration
}
