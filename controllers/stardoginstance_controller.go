package controllers

import (
	"context"
	"fmt"
	"github.com/vshn/stardog-userrole-operator/api/v1beta1"
	stardog "github.com/vshn/stardog-userrole-operator/stardogrest/client"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users"
	"net/url"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
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
const instanceDatabasesFinalizer = "finalizer.stardog.instance.databases"

// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardoginstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stardog.vshn.ch,resources=stardoginstances/status,verbs=get;update;patch

func (r *StardogInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	namespace := req.NamespacedName
	stardogInstance := &StardogInstance{}

	err := r.Client.Get(ctx, namespace, stardogInstance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("StardogInstance not found, ignoring reconcile.", "StardogInstance", namespace)
			return ctrl.Result{Requeue: false}, nil
		}
		r.Log.Error(err, "could not retrieve StardogInstance.", "StardogInstance", namespace)
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, err
	}

	if environmentDisabled(stardogInstance) {
		return ctrl.Result{Requeue: false}, nil
	}

	sir := &StardogInstanceReconciliation{
		reconciliationContext: &ReconciliationContext{
			context:       ctx,
			conditions:    make(map[StardogConditionType]StardogCondition),
			namespace:     namespace.Namespace,
			stardogClient: stardog.NewHTTPClient(nil),
		},
		resource: stardogInstance,
	}

	return r.ReconcileStardogInstance(sir)
}

func (r *StardogInstanceReconciler) ReconcileStardogInstance(sir *StardogInstanceReconciliation) (ctrl.Result, error) {

	rc := sir.reconciliationContext
	stardogInstance := sir.resource
	r.Log.Info("reconciling", getLoggingKeysAndValuesForStardogInstance(stardogInstance)...)

	isStardogInstanceMarkedToBeDeleted := stardogInstance.GetDeletionTimestamp() != nil
	if isStardogInstanceMarkedToBeDeleted {
		r.Log.Info(fmt.Sprintf("checking if StardogInstance %s is deletable", sir.resource.Name))
		if err := r.deleteStardogInstance(sir); err != nil {
			rc.SetStatusCondition(createStatusConditionTerminating(err))
			rc.SetStatusCondition(createStatusConditionReady(false, "StardogInstance not ready"))
			return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(sir)
		}
		return ctrl.Result{Requeue: false}, nil
	}

	if err := r.validateSpecification(sir.resource.Spec); err != nil {
		rc.SetStatusCondition(createStatusConditionInvalid(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "StardogInstance not ready"))
		return ctrl.Result{Requeue: false}, r.updateStatus(sir)
	}
	rc.SetStatusIfExisting(StardogInvalid, v1.ConditionFalse)

	if err := r.validateConnection(sir); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "StardogInstance not ready"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(sir)
	}
	rc.SetStatusIfExisting(StardogErrored, v1.ConditionFalse)

	controllerutil.AddFinalizer(sir.resource, instanceUserFinalizer)
	controllerutil.AddFinalizer(sir.resource, instanceRoleFinalizer)
	controllerutil.AddFinalizer(sir.resource, instanceDatabasesFinalizer)

	if err := r.Update(sir.reconciliationContext.context, sir.resource); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "StardogInstance not ready"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(sir)
	}
	rc.SetStatusCondition(createStatusConditionReady(true, "Synchronized"))
	return ctrl.Result{Requeue: true, RequeueAfter: ReconFreq}, r.updateStatus(sir)
}

func (r *StardogInstanceReconciler) deleteStardogInstance(sir *StardogInstanceReconciliation) error {
	stardogInstance := sir.resource
	rc := sir.reconciliationContext

	if contains(stardogInstance.GetFinalizers(), instanceUserFinalizer) {
		if err := r.userFinalizer(sir); err != nil {
			return err
		}
		controllerutil.RemoveFinalizer(stardogInstance, instanceUserFinalizer)
	}
	if contains(stardogInstance.GetFinalizers(), instanceRoleFinalizer) {
		if err := r.roleFinalizer(sir); err != nil {
			return err
		}
		controllerutil.RemoveFinalizer(stardogInstance, instanceRoleFinalizer)
	}
	if contains(stardogInstance.GetFinalizers(), instanceDatabasesFinalizer) {
		if err := r.databaseFinalizer(sir); err != nil {
			return err
		}
		controllerutil.RemoveFinalizer(stardogInstance, instanceDatabasesFinalizer)
	}
	if err := r.Update(rc.context, stardogInstance); err != nil {
		return err
	}
	return nil
}

func (r *StardogInstanceReconciler) userFinalizer(sir *StardogInstanceReconciliation) error {
	resource := sir.resource
	stardogUserList := &StardogUserList{}
	err := r.Client.List(sir.reconciliationContext.context, stardogUserList, client.InNamespace(resource.Namespace))
	if err != nil {
		return fmt.Errorf("cannot list Stardog User CRDs: %v", err)
	}
	items := stardogUserList.Items
	activeUsers := make([]string, len(items))
	for _, stardogUser := range items {
		if stardogUser.Spec.StardogInstanceRef == resource.Name {
			activeUsers = append(activeUsers, stardogUser.Name)
		}
	}
	if len(items) > 0 {
		return fmt.Errorf("cannot delete StardogInstance, found %s user CRDs", activeUsers)
	}
	return nil
}

func (r *StardogInstanceReconciler) roleFinalizer(sir *StardogInstanceReconciliation) error {
	resource := sir.resource
	stardogRoleList := &StardogRoleList{}
	err := r.Client.List(sir.reconciliationContext.context, stardogRoleList, client.InNamespace(resource.Namespace))
	if err != nil {
		return fmt.Errorf("cannot list Stardog Role CRDs: %v", err)
	}
	items := stardogRoleList.Items
	activeRoles := make([]string, len(items))
	for _, stardogRole := range items {
		if stardogRole.Spec.StardogInstanceRef == resource.Name {
			activeRoles = append(activeRoles, stardogRole.Name)
		}
	}
	if len(items) > 0 {
		return fmt.Errorf("cannot delete StardogInstance, found %s role CRDs", activeRoles)
	}
	return nil
}

func (r *StardogInstanceReconciler) databaseFinalizer(sir *StardogInstanceReconciliation) error {
	resource := sir.resource
	databasesList := &v1beta1.DatabaseList{}
	err := r.Client.List(sir.reconciliationContext.context, databasesList)
	if err != nil {
		return fmt.Errorf("cannot list Stardog Databases CRDs: %v", err)
	}
	items := databasesList.Items
	activeDatabases := make([]string, len(items))
	for _, stardogDatabase := range items {
		if containsStardogInstanceRef(stardogDatabase.Spec.StardogInstanceRefs, v1beta1.NewStardogInstanceRef(resource.Name, resource.Namespace)) {
			activeDatabases = append(activeDatabases, stardogDatabase.Name)
		}
	}
	if len(activeDatabases) > 0 {
		return fmt.Errorf("cannot delete StardogInstance, found %s databases CRDs", activeDatabases)
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
	credentials := spec.AdminCredentials

	r.Log.V(1).Info("retrieving admin credentials from Secret", "secret", credentials.Namespace+"/"+credentials.SecretRef)
	auth, err := rc.initStardogClient(r.Client, *sir.resource)
	if err != nil {
		return err
	}

	_, err = rc.stardogClient.Users.IsEnabled(users.NewIsEnabledParams().WithUser("admin"), auth)
	if err != nil {
		return err
	}
	return nil
}

func (r *StardogInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&StardogInstance{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
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
