package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stardogv1beta1 "github.com/vshn/stardog-userrole-operator/api/v1beta1"
	"github.com/vshn/stardog-userrole-operator/pkg/stardogapi"
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

	instance := &stardogv1beta1.Instance{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		r.Log.Error(err, "error getting instance")
		return ctrl.Result{}, err
	}

	credentialSecret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: instance.Spec.AdminCredentialRef.Name}, credentialSecret)
	if err != nil {
		r.Log.Error(err, "error getting instance credentials")
		return ctrl.Result{}, err
	}

	apiClient := stardogapi.NewClient(instance.Spec.AdminCredentialRef.Key, string(credentialSecret.Data[instance.Spec.AdminCredentialRef.Key]), instance.Spec.URL)

	_, err = apiClient.ListDatabases(ctx)
	if err != nil {
		r.Log.Error(err, "error getting databases")
		instance.Status.Conditions = []v1.Condition{{
			Type:   "Available",
			Status: v1.ConditionFalse, Reason: "Unavailable",
			Message:            fmt.Sprintf("error getting databases, user might be lacking permissions in Stardog: %s", err),
			LastTransitionTime: v1.NewTime(time.Now()),
		}}
		err = r.Status().Update(ctx, instance)
		if err != nil {
			r.Log.Error(err, "error updating status")
		}
		return ctrl.Result{}, err
	}

	if len(instance.Status.Conditions) == 0 || instance.Status.Conditions[0].Status != v1.ConditionTrue {
		instance.Status.Conditions = []v1.Condition{{
			Type:   "Available",
			Status: v1.ConditionTrue, Reason: "Available",
			Message:            "Instance became available",
			LastTransitionTime: v1.NewTime(time.Now()),
		}}
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
