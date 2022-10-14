package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stardogv1beta1 "github.com/vshn/stardog-userrole-operator/api/v1beta1"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=instances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=instances/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.6.4/pkg/reconcile
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	r.Log = r.Log.WithValues("instance", req.NamespacedName)

	instance := &stardogv1beta1.Instance{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		r.Log.Error(err, "error getting instance")
		instance.Status.Available = v1.Condition{
			Status: v1.ConditionFalse, Reason: "Unavailable",
			Message:            fmt.Sprintf("instance could not be retrieved: %s", err),
			LastTransitionTime: v1.NewTime(time.Now()),
		}
		err = r.Status().Update(ctx, instance)
		if err != nil {
			r.Log.Error(err, "error updating status")
		}
		return ctrl.Result{}, err
	}

	client, err := getStardogApiClient(ctx, r.Client, instance)
	if err != nil {
		r.Log.Error(err, "error creating API client")
		instance.Status.Available = v1.Condition{
			Status: v1.ConditionUnknown, Reason: "NoApiClient",
			Message:            fmt.Sprintf("error creating API client: %s", err),
			LastTransitionTime: v1.NewTime(time.Now()),
		}
		err = r.Status().Update(ctx, instance)
		if err != nil {
			r.Log.Error(err, "error updating status")
		}
		return ctrl.Result{}, err
	}

	_, err = client.ListDatabases(ctx)
	if err != nil {
		r.Log.Error(err, "error getting databases")
		instance.Status.Available = v1.Condition{
			Status: v1.ConditionFalse, Reason: "Unavailable",
			Message:            fmt.Sprintf("error getting databases, user might be lacking permissions in Stardog: %s", err),
			LastTransitionTime: v1.NewTime(time.Now()),
		}
		err = r.Status().Update(ctx, instance)
		if err != nil {
			r.Log.Error(err, "error updating status")
		}
		return ctrl.Result{}, err
	}

	if instance.Status.Available.Status != v1.ConditionTrue {
		instance.Status.Available = v1.Condition{
			Status: v1.ConditionTrue, Reason: "Available",
			Message:            "Instance became available",
			LastTransitionTime: v1.NewTime(time.Now()),
		}
		err = r.Status().Update(ctx, instance)
		if err != nil {
			r.Log.Error(err, "error updating status")
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
