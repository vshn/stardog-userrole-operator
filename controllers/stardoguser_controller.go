/*


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
	"github.com/vshn/stardog-userrole-operator/stardogrest"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
)

// StardogUserReconciler reconciles a StardogUser object
type StardogUserReconciler struct {
	client.Client
	Log               logr.Logger
	Scheme            *runtime.Scheme
	ReconcileInterval time.Duration
}

const userFinalizer = "finalizer.stardog.users"

// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardogusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardogusers/status,verbs=get;update;patch

func (r *StardogUserReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	namespace := req.NamespacedName
	ctx := context.Background()
	stardogUser := &StardogUser{}

	err := r.Client.Get(ctx, namespace, stardogUser)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("StardogUser not found, ignoring reconcile.", "StardogUser", namespace)
			return ctrl.Result{Requeue: false}, nil
		}
		r.Log.Error(err, "Could not retrieve StardogUser.", "StardogUser", namespace)
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, err
	}

	sur := &StardogUserReconciliation{
		reconciliationContext: &ReconciliationContext{
			context:       ctx,
			conditions:    make(map[StardogConditionType]StardogCondition),
			stardogClient: stardogrest.NewExtendedBaseClient(),
			namespace:     namespace.Namespace,
		},
		resource: stardogUser,
	}

	return r.ReconcileStardogUser(sur)
}

func (r *StardogUserReconciler) ReconcileStardogUser(sur *StardogUserReconciliation) (ctrl.Result, error) {

	rc := sur.reconciliationContext
	stardogUser := sur.resource
	r.Log.Info("reconciling", getLoggingKeysAndValuesForStardogUser(stardogUser)...)

	isStardogUserMarkedToBeDeleted := stardogUser.GetDeletionTimestamp() != nil
	if isStardogUserMarkedToBeDeleted {
		if err := r.deleteStardogUser(sur); err != nil {
			rc.SetStatusCondition(createStatusConditionTerminating(err))
			rc.SetStatusCondition(createStatusConditionReady(false, "StardogInstance cannot be deleted"))
			return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(sur)
		}
		return ctrl.Result{Requeue: false}, nil
	}

	if err := r.validateSpecification(&sur.resource.Spec); err != nil {
		rc.SetStatusCondition(createStatusConditionInvalid(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Specification cannot be validated"))
		return ctrl.Result{Requeue: false}, r.updateStatus(sur)
	}
	rc.SetStatusIfExisting(StardogInvalid, v1.ConditionFalse)

	if err := r.syncUser(sur); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Synchronization failed"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(sur)
	}

	if missingAtLeastOne(sur.resource.GetFinalizers(), userFinalizer) {
		r.Log.V(1).Info("adding Finalizers for the StardogUser")
		controllerutil.AddFinalizer(sur.resource, userFinalizer)
	}

	if err := r.Update(sur.reconciliationContext.context, sur.resource); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Cannot update User"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(sur)
	}
	rc.SetStatusIfExisting(StardogErrored, v1.ConditionFalse)
	rc.SetStatusCondition(createStatusConditionReady(true, "Synchronized"))
	return ctrl.Result{Requeue: true, RequeueAfter: ReconFreq}, r.updateStatus(sur)
}

func (r *StardogUserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&StardogUser{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

func (r *StardogUserReconciler) deleteStardogUser(sur *StardogUserReconciliation) error {
	r.Log.Info(fmt.Sprintf("deleting StardogUser %s", sur.resource.Name))
	stardogUser := sur.resource
	if err := r.finalize(sur); err != nil {
		return err
	}
	controllerutil.RemoveFinalizer(stardogUser, userFinalizer)
	err := r.Update(sur.reconciliationContext.context, stardogUser)
	if err != nil {
		return err
	}
	return nil
}

func (r *StardogUserReconciler) finalize(sur *StardogUserReconciliation) error {
	rc := sur.reconciliationContext
	ctx := rc.context
	spec := sur.resource.Spec
	namespace := rc.namespace

	r.Log.V(1).Info("setup Stardog Client from ", "ref", spec.StardogInstanceRef)
	err := rc.initStardogClientFromRef(r.Client, spec.StardogInstanceRef)
	if err != nil {
		return err
	}

	r.Log.V(1).Info("retrieving user credentials from Secret", "secret", spec.Credentials.Namespace+"/"+spec.Credentials.SecretRef)
	username, _, err := rc.getCredentials(r.Client, spec.Credentials, namespace)
	if err != nil {
		return err
	}

	_, err = rc.stardogClient.RemoveUser(ctx, username)
	if err != nil {
		return fmt.Errorf("cannot remove Stardog user %s/%s: %v", namespace, username, err)
	}
	return nil
}

func (r *StardogUserReconciler) validateSpecification(spec *StardogUserSpec) error {
	r.Log.V(1).Info("validating StardogUserSpec")
	if spec.Credentials.SecretRef == "" {
		return fmt.Errorf(".spec.Credentials.SecretRef is required")
	}
	if spec.StardogInstanceRef == "" {
		return fmt.Errorf(".spec.StardogInstanceRef is required")
	}
	return nil
}

func (r *StardogUserReconciler) syncUser(sur *StardogUserReconciliation) error {
	rc := sur.reconciliationContext
	spec := sur.resource.Spec
	userCredentials := spec.Credentials
	ctx := rc.context
	namespace := rc.namespace

	r.Log.V(1).Info("init Stardog Client from ", "ref", spec.StardogInstanceRef)
	err := rc.initStardogClientFromRef(r.Client, spec.StardogInstanceRef)
	if err != nil {
		return err
	}

	r.Log.V(1).Info("retrieving user credentials from Secret", "secret", userCredentials.Namespace+"/"+userCredentials.SecretRef)
	username, password, err := rc.getCredentials(r.Client, spec.Credentials, namespace)
	if err != nil {
		return err
	}

	superuser := false
	user := stardogrest.User{
		Username:  &username,
		Password:  &[]string{password},
		Superuser: &superuser,
	}
	stardogClient := rc.stardogClient
	usersObject, err := stardogClient.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("cannot get current list of users in %s: %v", namespace, err)
	}

	users := *usersObject.Users
	if len(users) > 0 && contains(users, username) {
		r.Log.V(1).Info("user already exists", "username", username)
		_, err := stardogClient.ChangePassword(ctx, username, stardogrest.Password{Password: &password})
		if err != nil {
			return fmt.Errorf("cannot change password for %s/%s: %v", namespace, *user.Username, err)
		}
	} else {
		r.Log.V(1).Info("creating user", "username", username)
		_, err = stardogClient.CreateUser(ctx, &user)
		if err != nil {
			return fmt.Errorf("cannot create user in %s/%s: %v", namespace, *user.Username, err)
		}
	}

	rolesObject, err := stardogClient.ListUserRoles(ctx, username)
	if err != nil {
		return fmt.Errorf("cannot get list of roles from %s/%s: %v", namespace, *user.Username, err)
	}

	var roleErrors []error
	roles := spec.Roles
	existingRoles := *rolesObject.Roles
	for _, role := range roles {
		if !contains(existingRoles, role) {
			_, err := stardogClient.AddRole(ctx, *user.Username, stardogrest.Rolename{Rolename: &role})
			if err != nil {
				roleErrors = append(roleErrors, err)
			}
		}
	}

	for _, existingRole := range existingRoles {
		if !contains(roles, existingRole) {
			_, err := stardogClient.RemoveRole(ctx, *user.Username, existingRole)
			if err != nil {
				roleErrors = append(roleErrors, err)
			}
		}
	}

	if len(roleErrors) > 0 {
		aggregateErrors := errors.NewAggregate(roleErrors)
		return fmt.Errorf("some roles have not been added to user %s in %s: %s", username, namespace, aggregateErrors)
	}

	return nil
}

func (r *StardogUserReconciler) updateStatus(sur *StardogUserReconciliation) error {
	cfg := sur.resource
	status := cfg.Status
	// Once we are on Kubernetes 0.19, we can use metav1.Conditions, but for now, we have to implement our helpers on
	// our own.
	status.Conditions = mergeWithExistingConditions(status.Conditions, sur.reconciliationContext.conditions)
	cfg.Status = status
	err := r.Client.Status().Update(sur.reconciliationContext.context, cfg)
	if err != nil {
		r.Log.Error(err, "could not update StardogUser", getLoggingKeysAndValuesForStardogUser(cfg)...)
		return err
	}
	r.Log.Info("updated StardogUser status", getLoggingKeysAndValuesForStardogUser(cfg)...)
	return nil
}

func getLoggingKeysAndValuesForStardogUser(stardogUser *StardogUser) []interface{} {
	return []interface{}{
		"StardogUser", stardogUser.Namespace + "/" + stardogUser.Name,
	}
}
