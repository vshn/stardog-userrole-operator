package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vshn/stardog-userrole-operator/stardogrest"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
)

// StardogRoleReconciler reconciles a StardogRole object
type StardogRoleReconciler struct {
	client.Client
	Log               logr.Logger
	Scheme            *runtime.Scheme
	ReconcileInterval time.Duration
}

const roleFinalizer = "finalizer.stardog.roles"

// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardogroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardogroles/status,verbs=get;update;patch

func (r *StardogRoleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	namespace := req.NamespacedName
	stardogRole := &StardogRole{}

	err := r.Client.Get(ctx, namespace, stardogRole)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("StardogRole not found, ignoring reconcile.", "StardogRole", namespace)
			return ctrl.Result{Requeue: false}, nil
		}
		r.Log.Error(err, "Could not retrieve StardogRole.", "StardogRole", namespace)
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, err
	}

	srr := &StardogRoleReconciliation{
		reconciliationContext: &ReconciliationContext{
			context:       ctx,
			conditions:    make(map[StardogConditionType]StardogCondition),
			namespace:     namespace.Namespace,
			stardogClient: stardogrest.NewExtendedBaseClient(),
		},
		resource: stardogRole,
	}

	return r.ReconcileStardogRole(srr)
}

func (r *StardogRoleReconciler) ReconcileStardogRole(srr *StardogRoleReconciliation) (ctrl.Result, error) {

	rc := srr.reconciliationContext
	stardogRole := srr.resource
	r.Log.Info("reconciling", getLoggingKeysAndValuesForStardogRole(stardogRole)...)

	isStardogRoleMarkedToBeDeleted := stardogRole.GetDeletionTimestamp() != nil
	if isStardogRoleMarkedToBeDeleted {
		if err := r.deleteStardogRole(srr); err != nil {
			rc.SetStatusCondition(createStatusConditionTerminating(err))
			rc.SetStatusCondition(createStatusConditionReady(false, "StardogRole cannot be deleted"))
			return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(srr)
		}
		return ctrl.Result{Requeue: false}, nil
	}

	if err := r.validateSpecification(srr.resource); err != nil {
		rc.SetStatusCondition(createStatusConditionInvalid(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Specification cannot be validated"))
		return ctrl.Result{Requeue: false}, r.updateStatus(srr)
	}
	rc.SetStatusIfExisting(StardogInvalid, v1.ConditionFalse)

	if err := r.syncRole(srr); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Synchronization failed"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(srr)
	}
	rc.SetStatusIfExisting(StardogErrored, v1.ConditionFalse)

	r.Log.V(1).Info("adding Finalizer for the StardogRole")
	controllerutil.AddFinalizer(srr.resource, roleFinalizer)

	if err := r.Update(srr.reconciliationContext.context, srr.resource); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Cannot update role"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(srr)
	}
	rc.SetStatusCondition(createStatusConditionReady(true, "Synchronized"))
	return ctrl.Result{Requeue: true, RequeueAfter: ReconFreq}, r.updateStatus(srr)
}

func (r *StardogRoleReconciler) syncRole(srr *StardogRoleReconciliation) error {
	spec := srr.resource.Spec
	ctx := srr.reconciliationContext.context
	namespace := srr.reconciliationContext.namespace
	roleName := spec.RoleName
	if roleName == "" {
		roleName = srr.resource.Name
	}

	r.Log.V(1).Info("init Stardog Client from ", "ref", spec.StardogInstanceRef)
	err := srr.reconciliationContext.initStardogClientFromRef(r.Client, spec.StardogInstanceRef)
	if err != nil {
		return err
	}

	stardogClient := srr.reconciliationContext.stardogClient

	r.Log.Info("synchronizing role", "role", roleName)
	roles, err := stardogClient.ListRoles(ctx)
	if err != nil {
		return fmt.Errorf("cannot list current roles in %s: %v", namespace, err)
	}
	if !contains(*roles.Roles, roleName) {
		_, err = stardogClient.CreateRole(ctx, &stardogrest.Rolename{Rolename: &roleName})
		if err != nil {
			return fmt.Errorf("cannot create role in %s/%s: %v", namespace, roleName, err)
		}
	}

	var existingPermissions []stardogrest.Permission
	if contains(*roles.Roles, roleName) {
		r.Log.V(1).Info("adding permissions to role", "role", roleName)
		permissionsObject, err := stardogClient.ListRolePermissions(ctx, roleName)
		if err != nil {
			return fmt.Errorf("cannot list permissions for role %s in %s: %v", roleName, namespace, err)
		}
		existingPermissions = *permissionsObject.Permissions
	}

	var permissionErrors []error
	permissions := spec.Permissions
	for _, existingPermission := range existingPermissions {
		if !containsStardogPermission(permissions, existingPermission) {
			_, err := stardogClient.RemoveRolePermission(ctx, roleName, existingPermission)
			if err != nil {
				permissionErrors = append(permissionErrors, err)
			}
		}
	}

	for _, permission := range permissions {
		if !containsOperatorPermission(existingPermissions, permission) {
			_, err := stardogClient.AddRolePermission(ctx, roleName, stardogrest.Permission{
				Action:       &permission.Action,
				ResourceType: &permission.ResourceType,
				Resource:     &permission.Resources,
			})
			if err != nil {
				permissionErrors = append(permissionErrors, err)
			}
		}
	}

	if len(permissionErrors) > 0 {
		aggregateErrors := errors.NewAggregate(permissionErrors)
		return fmt.Errorf("cannot add all permissions to role %s in %s: %s", roleName, namespace, aggregateErrors)
	}

	return nil
}

func (r *StardogRoleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&StardogRole{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

func (r *StardogRoleReconciler) validateSpecification(stardogRole *StardogRole) error {
	r.Log.V(1).Info("validating StardogRoleSpec")
	spec := stardogRole.Spec
	if spec.StardogInstanceRef == "" {
		return fmt.Errorf(".spec.StardogInstanceRef is required")
	}
	if len(spec.Permissions) > 0 {
		for i, permission := range spec.Permissions {
			if permission.Action == "" {
				return fmt.Errorf(".spec.Permissions[%d].Action is required", i)
			}
			if permission.ResourceType == "" {
				return fmt.Errorf(".spec.Permissions[%d].ResourceType is required", i)
			}
			if len(permission.Resources) == 0 {
				return fmt.Errorf(".spec.Permissions[%d].Resources at least one resource is required", i)
			}
		}
	}
	return nil
}

func (r *StardogRoleReconciler) updateStatus(srr *StardogRoleReconciliation) error {
	cfg := srr.resource
	status := cfg.Status
	// Once we are on Kubernetes 0.19, we can use metav1.Conditions, but for now, we have to implement our helpers on
	// our own.
	status.Conditions = mergeWithExistingConditions(status.Conditions, srr.reconciliationContext.conditions)
	cfg.Status = status
	err := r.Client.Status().Update(srr.reconciliationContext.context, cfg)
	if err != nil {
		r.Log.Error(err, "could not update StardogRole", getLoggingKeysAndValuesForStardogRole(cfg)...)
		return err
	}
	r.Log.Info("updated StardogRole status", getLoggingKeysAndValuesForStardogRole(cfg)...)
	return nil
}

func (r *StardogRoleReconciler) deleteStardogRole(srr *StardogRoleReconciliation) error {
	r.Log.Info(fmt.Sprintf("checking if StardogRole %s is deletable", srr.resource.Name))
	stardogRole := srr.resource
	if err := r.finalize(srr, roleFinalizer); err != nil {
		return err
	}
	controllerutil.RemoveFinalizer(stardogRole, roleFinalizer)
	err := r.Update(srr.reconciliationContext.context, stardogRole)
	if err != nil {
		return fmt.Errorf("cannot update StardogRole CRD: %v", err)
	}
	return nil
}

func (r *StardogRoleReconciler) finalize(srr *StardogRoleReconciliation, finalizer string) error {
	ctx := srr.reconciliationContext.context
	spec := srr.resource.Spec
	namespace := srr.reconciliationContext.namespace

	r.Log.V(1).Info("setup Stardog Client from ", "ref", spec.StardogInstanceRef)
	err := srr.reconciliationContext.initStardogClientFromRef(r.Client, spec.StardogInstanceRef)
	if err != nil {
		return err
	}

	stardogClient := srr.reconciliationContext.stardogClient

	role := spec.RoleName
	if role == "" {
		role = srr.resource.Name
	}
	rolesObject, err := stardogClient.ListRoleUsers(ctx, role)
	if err != nil {
		return fmt.Errorf("cannot get current list of roles in %s: %v", namespace, err)
	}

	if len(*rolesObject.Users) > 0 {
		return fmt.Errorf("cannot delete role %s as it is used by %s users in %s", role, strings.Join(*rolesObject.Users, ","), namespace)
	}

	_, err = srr.reconciliationContext.stardogClient.RemoveRole1(ctx, role, &[]bool{false}[0])
	if err != nil {
		return fmt.Errorf("cannot remove Stardog Role %s/%s: %v", namespace, role, err)
	}
	return nil
}

func getLoggingKeysAndValuesForStardogRole(stardogRole *StardogRole) []interface{} {
	return []interface{}{
		"StardogRole", stardogRole.Namespace + "/" + stardogRole.Name,
	}
}
