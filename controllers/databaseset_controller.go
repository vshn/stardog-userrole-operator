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

//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databasesets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databasesets/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DatabaseSet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.6.4/pkg/reconcile
func (r *DatabaseSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()

	databaseSet := &stardogv1beta1.DatabaseSet{}
	err := r.Client.Get(ctx, req.NamespacedName, databaseSet)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	for _, instance := range databaseSet.Spec.Instances {
		combinedName := fmt.Sprintf("%s-%s", req.Name, instance.Name)
		found := &stardogv1beta1.Database{}
		err = r.Get(ctx, types.NamespacedName{Name: combinedName, Namespace: databaseSet.Namespace}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				database := &stardogv1beta1.Database{
					ObjectMeta: v1.ObjectMeta{
						Name:      combinedName,
						Namespace: databaseSet.Namespace,
					},
					Spec: stardogv1beta1.DatabaseSpec{
						DatabaseName:     databaseSet.Name,
						InstanceRef:      instance,
						NamedGraphPrefix: databaseSet.Spec.NamedGraphPrefix,
					},
				}

				err = controllerutil.SetControllerReference(databaseSet, database, r.Scheme)
				if err != nil {
					return ctrl.Result{}, err
				}

				if err = r.Create(ctx, database); err != nil {
					return ctrl.Result{}, err
				}
			} else {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabaseSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stardogv1beta1.DatabaseSet{}).
		Owns(&stardogv1beta1.Database{}).
		Complete(r)
}
