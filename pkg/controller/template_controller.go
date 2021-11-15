package controller

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/operator-framework/combo/api/v1alpha1"
)

type templateController struct {
	client.Client
	log logr.Logger
}

func (t *templateController) manageWith(mgr ctrl.Manager, version int) error {
	t.log = t.log.V(version)
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Template{}).
		Complete(t)
}

func (t *templateController) Reconcile(ctx context.Context, req ctrl.Request) (reconcile.Result, error) {
	// Set up a convenient log object so we don't have to type request over and over again
	log := t.log.WithValues("request", req)
	log.V(1).Info("reconciling template")
	return reconcile.Result{}, nil
}
