package controller

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/operator-framework/combo/api/v1alpha1"
)

type combinationController struct {
	client.Client
	log logr.Logger
}

func (c *combinationController) manageWith(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Combination{}).
		Complete(c)
}

func (c *combinationController) Reconcile(ctx context.Context, req ctrl.Request) (reconcile.Result, error) {
	// Set up a convenient log object so we don't have to type request over and over again
	log := c.log.WithValues("request", req)
	log.Info("reconciling combination")

	in := &v1alpha1.Combination{}
	if err := c.Get(ctx, req.NamespacedName, in); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	return reconcile.Result{}, nil
}
