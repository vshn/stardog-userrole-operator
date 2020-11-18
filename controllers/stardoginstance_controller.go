/*
Licensed under the Apache License, Version 2.0 (the "License");
http://www.apache.org/licenses/LICENSE-2.0
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/vshn/stardog-userrole-operator/stardogrest/stardogrestapi"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	"github.com/vshn/stardog-userrole-operator/stardogrest"
)

// StardogInstanceReconciler reconciles a StardogInstance object
type StardogInstanceReconciler struct {
	client.Client
	Log               logr.Logger
	Scheme            *runtime.Scheme
	ReconcileInterval time.Duration
}

const instanceUserFinalizer = "finalizer.stardog.instance.users"
const instanceRoleFinalizer = "finalizer.stardog.instance.roles"

// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardoginstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardoginstances/status,verbs=get;update;patch

func (r *StardogInstanceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	namespace := req.NamespacedName
	ctx := context.Background()
	stardogInstance := &StardogInstance{}

	err := r.Client.Get(ctx, namespace, stardogInstance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("StardogInstance not found, ignoring reconcile.", "StardogInstance", namespace)
			return ctrl.Result{Requeue: false}, nil
		}
		r.Log.Error(err, "could not retrieve StardogInstance.", "StardogInstance", namespace)
		return ctrl.Result{Requeue: true, RequeueAfter: r.ReconcileInterval}, err
	}

	rc := &ReconciliationContext{
		context:       ctx,
		conditions:    make(map[StardogConditionType]StardogCondition),
		namespace:     namespace.Namespace,
		stardogClient: stardogrest.NewExtendedBaseClient(),
	}
	sir := &StardogInstanceReconciliation{
		reconciliationContext: rc,
		resource:              stardogInstance,
	}

	r.Log.Info("reconciling", getLoggingKeysAndValuesForStardogInstance(stardogInstance)...)

	isStardogInstanceMarkedToBeDeleted := stardogInstance.GetDeletionTimestamp() != nil
	if isStardogInstanceMarkedToBeDeleted {
		r.Log.Info(fmt.Sprintf("checking if StardogInstance %s is deletable", sir.resource.Name))
		if r.deleteStardogInstance(sir) != nil {
			rc.SetStatusCondition(createStatusConditionTerminating(err))
			rc.SetStatusCondition(createStatusConditionReady(false, "StardogInstance not ready"))
			return ctrl.Result{Requeue: true}, r.updateStatus(sir)
		}
		return ctrl.Result{Requeue: false}, r.updateStatus(sir)
	}

	if r.validateSpecification(sir.resource.Spec) != nil {
		rc.SetStatusCondition(createStatusConditionInvalid(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "StardogInstance not ready"))
		return ctrl.Result{}, r.updateStatus(sir)
	}
	rc.SetStatusIfExisting(StardogInvalid, v1.ConditionFalse)

	if r.validateConnection(sir) != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "StardogInstance not ready"))
		return ctrl.Result{}, r.updateStatus(sir)
	}
	rc.SetStatusIfExisting(StardogErrored, v1.ConditionFalse)

	controllerutil.AddFinalizer(sir.resource, instanceUserFinalizer)
	controllerutil.AddFinalizer(sir.resource, instanceRoleFinalizer)

	if r.Update(sir.reconciliationContext.context, sir.resource) != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "StardogInstance not ready"))
		return ctrl.Result{Requeue: true}, r.updateStatus(sir)
	}
	rc.SetStatusCondition(createStatusConditionReady(true, "Synchronized"))
	return ctrl.Result{Requeue: false}, r.updateStatus(sir)
}

func (r *StardogInstanceReconciler) deleteStardogInstance(sir *StardogInstanceReconciliation) error {
	stardogInstance := sir.resource
	rc := sir.reconciliationContext
	credentials := sir.resource.Spec.AdminCredentials

	r.Log.V(1).Info("retrieving admin credentials from Secret", "secret", credentials.Namespace+"/"+credentials.SecretRef)
	username, password, err := rc.getCredentials(r.Client, credentials, rc.namespace)
	if err != nil {
		return err
	}
	stardogClient := rc.stardogClient
	stardogClient.SetConnection(sir.resource.Spec.ServerUrl, username, password)

	if contains(stardogInstance.GetFinalizers(), instanceUserFinalizer) {
		if err := r.userFinalizer(stardogClient, rc.context); err != nil {
			return err
		}
		controllerutil.RemoveFinalizer(stardogInstance, instanceUserFinalizer)
	}
	if contains(stardogInstance.GetFinalizers(), instanceRoleFinalizer) {
		if err := r.roleFinalizer(stardogClient, rc.context); err != nil {
			return err
		}
		controllerutil.RemoveFinalizer(stardogInstance, instanceRoleFinalizer)
	}
	if err := r.Update(rc.context, stardogInstance); err != nil {
		return err
	}
	return nil
}

func (r *StardogInstanceReconciler) userFinalizer(stardogClient stardogrestapi.ExtendedBaseClientAPI, ctx context.Context) error {
	users, err := stardogClient.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("cannot list Stardog users: %v", err)
	}
	if len(*users.Users) > 0 {
		return fmt.Errorf("cannot delete StardogInstance, found %d users", len(*users.Users))
	}
	return nil
}

func (r *StardogInstanceReconciler) roleFinalizer(stardogClient stardogrestapi.ExtendedBaseClientAPI, ctx context.Context) error {
	roles, err := stardogClient.ListRoles(ctx)
	if err != nil {
		return fmt.Errorf("cannot list Stardog roles: %v", err)
	}
	if len(*roles.Roles) > 0 {
		return fmt.Errorf("cannot delete StardogInstance, found %d remaining role(s)", len(*roles.Roles))
	}
	return nil
}

func (r *StardogInstanceReconciler) validateSpecification(spec StardogInstanceSpec) error {
	r.Log.V(1).Info("validating StardogInstanceSpec")
	if spec.ServerUrl == "" {
		return fmt.Errorf(".spec.ServerUrl is required")
	}
	if spec.AdminCredentials.SecretRef == "" {
		return fmt.Errorf(".spec.AdminCredentials.SecretRef is required")
	}
	_, err := url.ParseRequestURI(spec.ServerUrl)
	if err != nil {
		return fmt.Errorf(".spec.ServerUrl is not a valid URL: %v", err)
	}
	return nil
}

func (r *StardogInstanceReconciler) validateConnection(sir *StardogInstanceReconciliation) error {
	r.Log.Info(fmt.Sprintf("verifying connection to Stardog API %s", sir.resource.Spec.ServerUrl))
	rc := sir.reconciliationContext
	spec := sir.resource.Spec
	ctx := rc.context
	credentials := spec.AdminCredentials
	stardogClient := rc.stardogClient

	r.Log.V(1).Info("retrieving admin credentials from Secret", "secret", credentials.Namespace+"/"+credentials.SecretRef)
	username, password, err := rc.getCredentials(r.Client, credentials, rc.namespace)
	if err != nil {
		return err
	}
	stardogClient.SetConnection(spec.ServerUrl, username, password)
	_, err = stardogClient.IsEnabled(ctx, username)
	if err != nil {
		return err
	}
	return nil
}

func (r *StardogInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&StardogInstance{}).
		Complete(r)
}

func (r *StardogInstanceReconciler) updateStatus(sir *StardogInstanceReconciliation) error {
	cfg := sir.resource
	status := cfg.Status
	// Once we are on Kubernetes 0.19, we can use metav1.Conditions, but for now, we have to implement our helpers on
	// our own.
	status.Conditions = mergeWithExistingConditions(status.Conditions, sir.reconciliationContext.conditions)
	cfg.Status = status
	err := r.Client.Status().Update(sir.reconciliationContext.context, cfg)
	if err != nil {
		r.Log.Error(err, "could not update StardogInstance", getLoggingKeysAndValuesForStardogInstance(cfg)...)
		return err
	}
	r.Log.Info("updated StardogInstance status", getLoggingKeysAndValuesForStardogInstance(cfg)...)
	return nil
}

func getLoggingKeysAndValuesForStardogInstance(stardogInstance *StardogInstance) []interface{} {
	return []interface{}{
		"StardogInstance", stardogInstance.Namespace + "/" + stardogInstance.Name,
	}
}
