/*
Licensed under the Apache License, Version 2.0 (the "License");
http://www.apache.org/licenses/LICENSE-2.0
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stardogv1alpha1 "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
)

// StardogInstanceReconciler reconciles a StardogInstance object
type StardogInstanceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardoginstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardoginstances/status,verbs=get;update;patch

func (r *StardogInstanceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("stardoginstance", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *StardogInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stardogv1alpha1.StardogInstance{}).
		Complete(r)
}
