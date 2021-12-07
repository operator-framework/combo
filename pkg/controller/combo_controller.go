package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	combinationPkg "github.com/operator-framework/combo/pkg/combination"
	comboConditions "github.com/operator-framework/combo/pkg/conditions"
	templatePkg "github.com/operator-framework/combo/pkg/template"

	"github.com/operator-framework/combo/api/v1alpha1"
)

type combinationController struct {
	client.Client
	log logr.Logger
}

func (c *combinationController) manageWith(mgr ctrl.Manager, version int) error {
	c.log = c.log.V(version)
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Combination{}).
		Complete(c)
}

func (c *combinationController) Reconcile(ctx context.Context, req ctrl.Request) (result reconcile.Result, deferredErr error) {
	// Set up a convenient log object so we don't have to type request over and over again
	log := c.log.WithValues("request", req)

	log.Info("new combination inbound")

	combination := &v1alpha1.Combination{}
	err := c.Get(ctx, req.NamespacedName, combination)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Update the status whenever the reconciliation completes
	defer func() {
		if updateErr := c.Status().Update(ctx, combination); updateErr != nil {
			deferredErr = fmt.Errorf("error updating status of combination: %w", updateErr)
		}
	}()

	log.Info(fmt.Sprintf("combination %s successfully loaded in reconciler", combination.Name))

	templateQuery := types.NamespacedName{Name: combination.Spec.Template}

	template := &v1alpha1.Template{}
	if err := c.Get(ctx, templateQuery, template); err != nil {
		combination.Status.Conditions = comboConditions.NewConditions(
			time.Now(),
			err,
			comboConditions.TemplateNotFoundCondition)
		return reconcile.Result{}, fmt.Errorf("failed to retrieve %v template: %w", combination.Spec.Template, err)
	}

	comboStream := combinationPkg.NewStream(
		combinationPkg.WithArgs(formatArguments(combination.Spec.Arguments)),
		combinationPkg.WithSolveAhead(),
	)

	builder, err := templatePkg.NewBuilder(strings.NewReader(template.Spec.Body), comboStream)
	if err != nil {
		combination.Status.Conditions = comboConditions.NewConditions(
			time.Now(),
			err,
			comboConditions.TemplateBodyInvalid)
		return reconcile.Result{}, fmt.Errorf("failed to construct a builder out of %s template body: %w", template.Name, err)
	}

	log.Info(fmt.Sprintf("template %s for combination %s successfully loaded in reconciler", template.Name, combination.Name))

	generatedManifests, err := builder.Build(ctx)
	if err != nil {
		combination.Status.Conditions = comboConditions.NewConditions(
			time.Now(),
			err,
			comboConditions.ManifestGenerationFailed)
		return reconcile.Result{}, fmt.Errorf("failed to generate manifest %s combinations: %w", combination.Name, err)
	}

	log.Info(fmt.Sprintf("manifest combinations for %s successfully generated!", combination.Name))

	combination.Status = v1alpha1.CombinationStatus{
		Evaluation: generatedManifests,
		Conditions: comboConditions.NewConditions(
			time.Now(),
			nil,
			comboConditions.ProccessedCombinationsCondition),
	}

	log.Info(fmt.Sprintf("reconciliation of %s combination complete!", combination.Name))

	return reconcile.Result{}, nil
}

func formatArguments(arguments []v1alpha1.Argument) map[string][]string {
	formattedArguments := map[string][]string{}

	for _, argument := range arguments {
		formattedArguments[argument.Key] = argument.Values
	}

	return formattedArguments
}
