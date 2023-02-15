package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stardogv1beta1 "github.com/vshn/stardog-userrole-operator/api/v1beta1"
)

// InstanceReconciler reconciles an Instance object
type InstanceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=instances,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=instances/status,verbs=get;update;patch

// Reconcile verifies if Stardog instances are available and can be used by other resources
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance := &stardogv1beta1.Instance{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		r.Log.Error(err, "error getting instance")
		return ctrl.Result{}, err
	}

	apiClient, err := NewStardogAPIClientFromInstance(ctx, r.Client, instance)
	if err != nil {
		r.Log.Error(err, "error creating new Stardog API client")
		instance.Status.Conditions = []v1.Condition{
			createInstanceStatusUnavailableCondition(fmt.Sprintf("error creating new Stardog API client: %s", err)),
		}
		err = r.Status().Update(ctx, instance)
		if err != nil {
			r.Log.Error(err, "error updating status")
		}
		return ctrl.Result{}, err
	}

	// Check if Instance is available by listing its databases
	_, err = apiClient.ListDatabases(ctx)
	if err != nil {
		r.Log.Error(err, "error getting databases")
		instance.Status.Conditions = []v1.Condition{
			createInstanceStatusUnavailableCondition(fmt.Sprintf("error getting databases, user might be lacking permissions in Stardog: %s", err)),
		}
		err = r.Status().Update(ctx, instance)
		if err != nil {
			r.Log.Error(err, "error updating status")
		}
		return ctrl.Result{}, err
	}

	// Mark Instance as available
	if len(instance.Status.Conditions) == 0 || instance.Status.Conditions[0].Status != v1.ConditionTrue {
		instance.Status.Conditions = []v1.Condition{
			createInstanceStatusAvailableCondition("listing databases succeeded"),
		}
		err = r.Status().Update(ctx, instance)
		if err != nil {
			r.Log.Error(err, "error updating instance status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stardogv1beta1.Instance{}).
		Complete(r)
}
