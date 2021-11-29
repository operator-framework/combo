package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	combinationPkg "github.com/operator-framework/combo/pkg/combination"
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

func (c *combinationController) Reconcile(ctx context.Context, req ctrl.Request) (reconcile.Result, error) {
	// Set up a convenient log object so we don't have to type request over and over again
	log := c.log.WithValues("request", req)

	log.Info("new combination inbound")

	combination := &v1alpha1.Combination{}
	if err := c.Get(ctx, req.NamespacedName, combination); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	log.Info(fmt.Sprintf("combination %s successfully loaded in reconciler", combination.Name))

	nameSpacedTemplate := types.NamespacedName{
		Name:      combination.Spec.Template,
		Namespace: combination.Namespace,
	}

	template := &v1alpha1.Template{}
	if err := c.Get(ctx, nameSpacedTemplate, template); err != nil {
		return reconcile.Result{}, err
	}

	comboStream := combinationPkg.NewStream(
		combinationPkg.WithArgs(formatArguments(combination.Spec.Arguments)),
		combinationPkg.WithSolveAhead(),
	)

	builder, err := templatePkg.NewBuilder(strings.NewReader(template.Spec.Body), comboStream)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to construct a builder out of %s template body: %w", template.Name, err)
	}

	log.Info(fmt.Sprintf("template %s for combination %s successfully loaded in reconciler", template.Name, combination.Name))

	generatedManifests, err := builder.Build(ctx)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to generate manifest %s combinations: %w", combination.Name, err)
	}

	log.Info(fmt.Sprintf("manifest combinations for %s successfully generated!", combination.Name))

	combination.Status = v1alpha1.CombinationStatus{
		Evaluation: generatedManifests,
	}

	if err = c.Status().Update(ctx, combination); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to update %s combination's status: %w", combination.Name, err)
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
