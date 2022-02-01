/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/kind/pkg/errors"

	combov1alpha1 "github.com/operator-framework/combo/api/v1alpha1"
	combinationPkg "github.com/operator-framework/combo/pkg/combination"
	templatePkg "github.com/operator-framework/combo/pkg/template"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CombinationReconciler reconciles a Combination object
type CombinationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// mapTemplateToCombinations is responsible for taking the template object and finding all associated
// combinations that should be requeued. This should only happen whenever a template is changed in someway.
func (c *CombinationReconciler) mapTemplateToCombinations(template client.Object) []reconcile.Request {
	if template == nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := log.FromContext(ctx)

	requests := []reconcile.Request{}

	// Retrieve and validate the template's name
	templateName := template.GetName()

	// Find all of the combinations that rely on this template
	combinationList := combov1alpha1.CombinationList{}
	if err := c.List(ctx, &combinationList, &client.ListOptions{}); err != nil {
		return requests
	}

	//  Enqueue reliant combinations for updates
	for _, combination := range combinationList.Items {
		if combination.Spec.Template == templateName {
			log.Info(fmt.Sprintf("enqueueing %s combination in response to associated %s template being updated", combination.Name, templateName))
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{Name: combination.Name},
			})
		}
	}

	return requests
}

// formatArguments takes the arguments for the combination and formats them ito what the combination package
// is expecting
func formatArguments(arguments []combov1alpha1.Argument) map[string][]string {
	formattedArguments := map[string][]string{}
	for _, argument := range arguments {
		formattedArguments[argument.Key] = argument.Values
	}
	return formattedArguments
}

//+kubebuilder:rbac:groups=combo.io,resources=combinations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=combo.io,resources=combinations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=combo.io,resources=combinations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Combination object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (c *CombinationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Set up a convenient log object so we don’t have to type request over and over again
	log := log.FromContext(ctx)

	log.V(1).Info("new combination event inbound")

	// Attempt to retrieve the requested combination CR
	combination := &combov1alpha1.Combination{}
	err := c.Get(ctx, req.NamespacedName, combination)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// If combination is being deleted, remove from queue
	if !combination.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("combination is being deleted, ignoring event")
		return reconcile.Result{}, nil
	}

	// Remove any previous evaluation in case of failure
	combination.Status.Evaluations = []string{}

	// Attempt to retrieve the template referenced in the combination CR
	template := &combov1alpha1.Template{}
	if err := c.Get(ctx, types.NamespacedName{Name: combination.Spec.Template}, template); err != nil {
		combination.SetStatusCondition(metav1.Condition{
			Type:               combov1alpha1.TypeInvalid,
			Status:             metav1.ConditionFalse,
			Reason:             combov1alpha1.ReasonTemplateNotFound,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Message:            fmt.Sprintf("failed to retrieve %s template: %s", combination.Spec.Template, err.Error()),
		})
		return reconcile.Result{}, errors.NewAggregate([]error{err, c.Status().Update(ctx, combination)})
	}

	// Build combination stream to be utilized in template builder
	comboStream := combinationPkg.NewStream(
		combinationPkg.WithArgs(formatArguments(combination.Spec.Arguments)),
		combinationPkg.WithSolveAhead(),
	)

	// Create a new template builder
	builder, err := templatePkg.NewBuilder(strings.NewReader(template.Spec.Body), comboStream)
	if err != nil {
		combination.SetStatusCondition(metav1.Condition{
			Type:               combov1alpha1.TypeInvalid,
			Status:             metav1.ConditionFalse,
			Reason:             combov1alpha1.ReasonTemplateBodyInvalid,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Message:            fmt.Sprintf("failed to construct a builder out of %s template body:  %s", combination.Spec.Template, err.Error()),
		})
		return reconcile.Result{}, errors.NewAggregate([]error{err, c.Status().Update(ctx, combination)})
	}

	// Build the manifest combinations
	generatedManifests, err := builder.Build(ctx)
	if err != nil {
		combination.SetStatusCondition(metav1.Condition{
			Type:               combov1alpha1.TypeInvalid,
			Status:             metav1.ConditionFalse,
			Reason:             combov1alpha1.ReasonEvaluationsInvalid,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Message:            fmt.Sprintf("failed to generate manifest %s combinations: %s", combination.Spec.Template, err.Error()),
		})
		return reconcile.Result{}, errors.NewAggregate([]error{err, c.Status().Update(ctx, combination)})
	}

	// Add the combination evaluations and update Status
	combination.Status.Evaluations = generatedManifests
	combination.SetStatusCondition(metav1.Condition{
		Type:               combov1alpha1.TypeFinished,
		Status:             metav1.ConditionTrue,
		Reason:             combov1alpha1.ReasonProcessed,
		LastTransitionTime: metav1.NewTime(time.Now()),
		Message:            "evaluations successfully processed",
	})

	// Return and update the combination's status
	return reconcile.Result{}, c.Status().Update(ctx, combination)
}

// SetupWithManager sets up the controller with the Manager.
func (c *CombinationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	templateHandler := handler.EnqueueRequestsFromMapFunc(c.mapTemplateToCombinations)

	return ctrl.NewControllerManagedBy(mgr).
		For(&combov1alpha1.Combination{}).
		Watches(&source.Kind{Type: &combov1alpha1.Template{}}, templateHandler).
		Complete(c)
}
