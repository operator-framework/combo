package controller

import (
	"context"
	"fmt"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/go-logr/logr"
	"github.com/operator-framework/combo/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	combinationPkg "github.com/operator-framework/combo/pkg/combination"
	templatePkg "github.com/operator-framework/combo/pkg/template"
	"github.com/operator-framework/combo/pkg/updater"
)

const (
	ReferencedTemplateLabel = "combo.ReferencedTemplate"
)

type combinationController struct {
	client.Client
	log logr.Logger
}

// manageWith creates a new instance of this controller
func (c *combinationController) manageWith(mgr ctrl.Manager, verbosity int) error {
	c.log = c.log.V(verbosity)
	templateHandler := handler.EnqueueRequestsFromMapFunc(c.mapTemplateToCombinations)

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Combination{}).
		Watches(&source.Kind{Type: &v1alpha1.Template{}}, templateHandler).
		Complete(c)
}

// mapTemplateToCombinations is responsible for taking the template object and finding all associated
// combinations that should be requeued. This should only happen whenever a template is changed in someway.
func (c *combinationController) mapTemplateToCombinations(template client.Object) []reconcile.Request {
	if template == nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requests := []reconcile.Request{}

	// Retrieve and validate the template's name
	templateName := template.GetName()

	// Find all of the combinations that rely on this template
	combinationList := v1alpha1.CombinationList{}
	if err := c.List(ctx, &combinationList, &client.ListOptions{}); err != nil {
		return requests
	}

	//  Enqueue reliant combinations for updates
	for _, combination := range combinationList.Items {
		if combination.Spec.Template == templateName {
			c.log.Info(fmt.Sprintf("enqueueing %s combination in response to associated %s template being updated", combination.Name, templateName))
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{Name: combination.Name},
			})
		}
	}

	return requests
}

// Reconcile manages incoming combination CR's and processes them accordingly
func (c *combinationController) Reconcile(ctx context.Context, req ctrl.Request) (reconcile.Result, error) {
	// Set up a convenient log object so we don’t have to type request over and over again
	log := c.log.WithValues("request", req)

	log.V(1).Info("new combination event inbound")

	// Attempt to retrieve the requested combination CR
	combination := &v1alpha1.Combination{}
	err := c.Get(ctx, req.NamespacedName, combination)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// If combination is being deleted, remove from queue
	if !combination.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("combination is being deleted, ignoring event")
		return reconcile.Result{}, nil
	}

	// create update client and defer updates until exiting the control loop
	u := updater.New(c.Client)
	defer func() {
		if err := u.Apply(ctx, combination); err != nil {
			log.Error(err, "failed to update status")
		}
	}()

	// Remove any previous evaluation in case of failure
	combination.Status.Evaluations = []string{}

	// Attempt to retrieve the template referenced in the combination CR
	template := &v1alpha1.Template{}
	if err := c.Get(ctx, types.NamespacedName{Name: combination.Spec.Template}, template); err != nil {
		u.UpdateStatus(updater.EnsureCondition(metav1.Condition{
			Type:    v1alpha1.TypeInvalid,
			Status:  metav1.ConditionTrue,
			Reason:  v1alpha1.ReasonTemplateNotFound,
			Message: fmt.Sprintf("failed to retrieve %s template: %s", combination.Spec.Template, err.Error()),
		}))
		return reconcile.Result{}, err
	}

	// Build combination stream to be utilized in template builder
	comboStream := combinationPkg.NewStream(
		combinationPkg.WithArgs(formatArguments(combination.Spec.Arguments)),
		combinationPkg.WithSolveAhead(),
	)

	// Create a new template builder
	builder, err := templatePkg.NewBuilder(strings.NewReader(template.Spec.Body), comboStream)
	if err != nil {
		u.UpdateStatus(updater.EnsureCondition(metav1.Condition{
			Type:    v1alpha1.TypeInvalid,
			Status:  metav1.ConditionTrue,
			Reason:  v1alpha1.ReasonTemplateBodyInvalid,
			Message: fmt.Sprintf("failed to construct a builder out of %s template body:  %s", combination.Spec.Template, err.Error()),
		}))
		return reconcile.Result{}, err
	}

	// Build the manifest combinations
	generatedManifests, err := builder.Build(ctx)
	if err != nil {
		u.UpdateStatus(updater.EnsureCondition(metav1.Condition{
			Type:    v1alpha1.TypeInvalid,
			Status:  metav1.ConditionTrue,
			Reason:  v1alpha1.ReasonEvaluationsInvalid,
			Message: fmt.Sprintf("failed to generate manifest %s combinations: %s", combination.Spec.Template, err.Error()),
		}))
		return reconcile.Result{}, err
	}

	// Add the combination evaluations and update Status
	u.UpdateStatus(updater.EnsureCondition(metav1.Condition{
		Type:    v1alpha1.TypeFinished,
		Status:  metav1.ConditionTrue,
		Reason:  v1alpha1.ReasonProcessed,
		Message: "evaluations successfully processed",
	}), updater.EnsureEvaluations(generatedManifests))

	// If manifest application is indicated, create generated manifests on-cluster.
	if combination.Spec.Apply {
		// update status to indicate processing
		combination.SetStatusCondition(metav1.Condition{
			Type:               v1alpha1.TypeInProgress,
			Status:             metav1.ConditionTrue,
			Reason:             v1alpha1.ReasonProcessing,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Message:            "application of evaluations to the cluster in progress",
		})

		for _, manifest := range generatedManifests {
			dec := yaml.NewYAMLOrJSONDecoder(strings.NewReader(manifest), 10)
			obj := &unstructured.Unstructured{}
			err := dec.Decode(obj)
			if err != nil {
				combination.SetStatusCondition(metav1.Condition{
					Type:               v1alpha1.TypeInvalid,
					Status:             metav1.ConditionFalse,
					Reason:             v1alpha1.ReasonApplyFailed,
					LastTransitionTime: metav1.NewTime(time.Now()),
					Message:            fmt.Sprintf("failed to decode manifest %s to apply to cluster: %s", obj.GetName(), err.Error()),
				})
				return reconcile.Result{}, errors.NewAggregate([]error{err, c.Status().Update(ctx, combination)})
			}
			err = c.Create(ctx, obj)
			if err != nil && !k8serrors.IsAlreadyExists(err) {
				combination.SetStatusCondition(metav1.Condition{
					Type:               v1alpha1.TypeInvalid,
					Status:             metav1.ConditionFalse,
					Reason:             v1alpha1.ReasonApplyFailed,
					LastTransitionTime: metav1.NewTime(time.Now()),
					Message:            fmt.Sprintf("failed to apply manifest %s to cluster: %s", obj.GetName(), err.Error()),
				})
				return reconcile.Result{}, errors.NewAggregate([]error{err, c.Status().Update(ctx, combination)})
			}
		}
		// update status once all manifests have been applied
		combination.SetStatusCondition(metav1.Condition{
			Type:               v1alpha1.TypeFinished,
			Status:             metav1.ConditionTrue,
			Reason:             v1alpha1.ReasonApplySucceeded,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Message:            "applied all evaluations to the cluster",
		})
	}

	// Return and update the combination's status
	return reconcile.Result{}, err
}

// formatArguments takes the arguments for the combination and formats them ito what the combination package
// is expecting
func formatArguments(arguments []v1alpha1.Argument) map[string][]string {
	formattedArguments := map[string][]string{}
	for _, argument := range arguments {
		formattedArguments[argument.Key] = argument.Values
	}
	return formattedArguments
}
