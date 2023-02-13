package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/sethvargo/go-password/password"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	stardogv1beta1 "github.com/vshn/stardog-userrole-operator/api/v1beta1"
	"github.com/vshn/stardog-userrole-operator/pkg/stardogapi"
	stardogapiutil "github.com/vshn/stardog-userrole-operator/pkg/stardogapi/util"

	types "k8s.io/apimachinery/pkg/types"
)

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases/status,verbs=get;update;patch

// Reconcile manages the Stardog resources for a Database object
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	database := &stardogv1beta1.Database{}
	err := r.Get(ctx, req.NamespacedName, database)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Generate and save credentials
	secretName := fmt.Sprintf("%s-%s-credentials", database.Spec.DatabaseName, database.Spec.InstanceRef.Name)
	readName := fmt.Sprintf("%s-read", database.Spec.DatabaseName)
	writeName := fmt.Sprintf("%s-write", database.Spec.DatabaseName)

	err = r.createCredentials(ctx, database, secretName, readName, writeName)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Create Stardog resources
	instance := &stardogv1beta1.Instance{}
	err = r.Get(ctx, types.NamespacedName{Name: database.Spec.InstanceRef.Name, Namespace: database.Namespace}, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	apiClient, err := NewStardogAPIClientFromInstance(ctx, r.Client, instance)
	if err != nil {
		r.Log.Error(err, "error creating new Stardog API client")
		return ctrl.Result{}, err
	}

	// Delete Stardog resources
	finalizer := "stardog.vshn.ch/finalizer"
	if database.DeletionTimestamp.IsZero() {
		// add finalizer
		if !controllerutil.ContainsFinalizer(database, finalizer) {
			controllerutil.AddFinalizer(database, finalizer)
			if err := r.Update(ctx, database); err != nil {
				r.Log.Error(err, "error updating finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(database, finalizer) {
			if err := r.handleDeletion(ctx, apiClient, database, []string{readName, writeName}, []string{readName, writeName}); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(database, finalizer)
			if err := r.Update(ctx, database); err != nil {
				r.Log.Error(err, "error updating finalizer")
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	// Create database if it doesn't exist
	liveDatabases, err := apiClient.ListDatabases(ctx)
	if err != nil {
		r.Log.Error(err, "error listing databases")
		return ctrl.Result{}, err
	}
	if !slices.Contains(liveDatabases, database.Spec.DatabaseName) {
		err = apiClient.CreateDatabase(ctx, database.Spec.DatabaseName, nil)
		if err != nil {
			r.Log.Error(err, "error creating database")
			return ctrl.Result{}, err
		}
		r.Log.Info("created Stardog database", "name", database.Spec.DatabaseName)
	}

	// Get user credentials from Secret
	secret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: database.Namespace}, secret)
	if err != nil {
		return ctrl.Result{}, err
	}

	users := []stardogapi.UserCredentials{
		{Name: readName, Password: string(secret.Data[readName])},
		{Name: writeName, Password: string(secret.Data[writeName])},
	}

	err = stardogapiutil.CreateUsers(ctx, apiClient, users)
	if err != nil {
		r.Log.Error(err, "error creating users", "users", users)
		return ctrl.Result{}, err
	}

	roles := []string{readName, writeName}
	err = stardogapiutil.CreateRoles(ctx, apiClient, roles)
	if err != nil {
		r.Log.Error(err, "error creating roles", "roles", roles)
		return ctrl.Result{}, err
	}

	readPermissions := []stardogapi.Permission{
		{ResourceType: "db", Action: "READ", Resources: []string{database.Spec.DatabaseName}},
		{ResourceType: "metadata", Action: "READ", Resources: []string{database.Spec.DatabaseName}},
	}

	writePermissions := []stardogapi.Permission{
		{ResourceType: "db", Action: "WRITE", Resources: []string{database.Spec.DatabaseName}},
		{ResourceType: "metadata", Action: "WRITE", Resources: []string{database.Spec.DatabaseName}},
	}

	err = stardogapiutil.AddPermissions(ctx, apiClient, readName, readPermissions)
	if err != nil {
		r.Log.Error(err, "adding permission to role failed", "role", readName, "permission", readPermissions)
		return ctrl.Result{}, err
	}
	err = stardogapiutil.AddPermissions(ctx, apiClient, writeName, writePermissions)
	if err != nil {
		r.Log.Error(err, "adding permission to role failed", "role", writeName, "permission", writePermissions)
		return ctrl.Result{}, err
	}

	// assign roles
	readUserRoles, err := apiClient.GetUserRoles(ctx, readName)
	if err != nil {
		r.Log.Error(err, "error getting user roles", "user", readName)
		return ctrl.Result{}, err
	}

	if len(readUserRoles) != 1 || readUserRoles[0] != readName {
		err = apiClient.SetUserRoles(ctx, readName, []string{readName})
		if err != nil {
			r.Log.Error(err, "error assigning roles to user", "user", readName, "roles", readUserRoles)
			return ctrl.Result{}, err
		}
	}

	writeUserRoles, err := apiClient.GetUserRoles(ctx, writeName)
	if err != nil {
		r.Log.Error(err, "error getting user roles", "user", writeName)
		return ctrl.Result{}, err
	}

	if len(writeUserRoles) != 2 || writeUserRoles[0] != readName || writeUserRoles[1] != writeName {
		err = apiClient.SetUserRoles(ctx, writeName, []string{readName, writeName})
		if err != nil {
			r.Log.Error(err, "error assigning roles to user", "user", readName, "roles", readUserRoles)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stardogv1beta1.Database{}).
		Complete(r)
}

// Generate and save credentials for Database
func (r *DatabaseReconciler) createCredentials(ctx context.Context, database *stardogv1beta1.Database, secretName string, readName string, writeName string) error {
	secret := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: database.Namespace}, secret)

	if err == nil {
		return nil
	} else if !errors.IsNotFound(err) {
		r.Log.Error(err, fmt.Sprintf("error getting Secret %s/%s", database.Namespace, secretName))
		return err
	}

	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: database.Namespace,
		}}
	err = controllerutil.SetControllerReference(database, secret, r.Scheme)
	if err != nil {
		return err
	}

	readPassword, err := password.Generate(20, 5, 0, false, false)
	if err != nil {
		r.Log.Error(err, "generation of password for read user failed", "db", database.Spec.DatabaseName)
		return err
	}

	writePassword, err := password.Generate(20, 5, 0, false, false)
	if err != nil {
		r.Log.Error(err, "generation of password for write user failed", "db", database.Spec.DatabaseName)
		return err
	}

	secret.StringData = map[string]string{}

	secret.StringData[readName] = readPassword
	secret.StringData[writeName] = writePassword

	err = r.Create(ctx, secret)
	if err != nil {
		r.Log.Error(err, fmt.Sprintf("error creating secret %s/%s", database.Namespace, secretName))
		return err
	} else {
		r.Log.Info("created credential secret", "namespace", secret.Namespace, "name", secret.Name)
	}

	return nil
}

// Handle deletion of all Stardog resources related to the database
func (r *DatabaseReconciler) handleDeletion(ctx context.Context, apiClient stardogapi.StardogAPI, database *stardogv1beta1.Database, users []string, roles []string) error {
	err := apiClient.DropDatabase(ctx, database.Spec.DatabaseName)
	if err != nil {
		return fmt.Errorf("error dropping database %s/%s: %w", database.Namespace, database.Name, err)
	}

	for _, user := range users {
		err = apiClient.DeleteUser(ctx, user)
		if err != nil {
			return fmt.Errorf("error deleting user %s: %w", user, err)
		}
	}

	for _, role := range roles {
		err = apiClient.DeleteRole(ctx, role)
		if err != nil {
			return fmt.Errorf("error deleting role %s: %w", role, err)
		}
	}

	return nil
}
