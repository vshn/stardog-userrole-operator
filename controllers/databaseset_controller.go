package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	stardogv1beta1 "github.com/vshn/stardog-userrole-operator/api/v1beta1"
)

// DatabaseSetReconciler reconciles a DatabaseSet object
type DatabaseSetReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databasesets,verbs=get;list;watch
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databasesets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases/status,verbs=get;update;patch

// Reconcile manages the Database objects for a DatabaseSet
func (r *DatabaseSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	databaseSet := &stardogv1beta1.DatabaseSet{}
	err := r.Get(ctx, req.NamespacedName, databaseSet)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Create Database for every instance
	for _, instance := range databaseSet.Spec.Instances {
		err = r.createDatabase(ctx, databaseSet, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// Creates a database from a DatabaseSet for a specific Instance ref
func (r *DatabaseSetReconciler) createDatabase(ctx context.Context, databaseSet *stardogv1beta1.DatabaseSet, instanceRef stardogv1beta1.StardogInstanceRef) error {
	combinedName := fmt.Sprintf("%s-%s", databaseSet.Name, instanceRef.Name)
	found := &stardogv1beta1.Database{}
	err := r.Get(ctx, types.NamespacedName{Name: combinedName, Namespace: databaseSet.Namespace}, found)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		} else {
			database := &stardogv1beta1.Database{
				ObjectMeta: v1.ObjectMeta{
					Name:      combinedName,
					Namespace: databaseSet.Namespace,
				},
				Spec: stardogv1beta1.DatabaseSpec{
					DatabaseName:     databaseSet.Name,
					InstanceRef:      instanceRef,
					NamedGraphPrefix: databaseSet.Spec.NamedGraphPrefix,
				},
			}

			// Ensure cleanup on deletion of DatabaseSet
			err = controllerutil.SetControllerReference(databaseSet, database, r.Scheme)
			if err != nil {
				return err
			}

			if err = r.Create(ctx, database); err != nil {
				return err
			}
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabaseSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stardogv1beta1.DatabaseSet{}).
		Owns(&stardogv1beta1.Database{}).
		Complete(r)
}
