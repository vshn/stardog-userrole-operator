package controllers

import (
	"context"
	"fmt"
	"strings"

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

	types "k8s.io/apimachinery/pkg/types"
)

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type User struct {
	name     string
	password string
}

//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Database object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.6.4/pkg/reconcile
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	// TODO track connected external resources in status

	database := &stardogv1beta1.Database{}
	err := r.Get(ctx, req.NamespacedName, database)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Generate and save credentials
	secret := &corev1.Secret{}
	secretName := fmt.Sprintf("%s-%s-credentials", database.Spec.DatabaseName, database.Spec.InstanceRef.Name)
	err = r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: database.Namespace}, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			secret = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: database.Namespace,
				}}
			err = controllerutil.SetControllerReference(database, secret, r.Scheme)
			if err != nil {
				return ctrl.Result{}, err
			}

			readPassword, err := password.Generate(20, 5, 0, false, false)
			if err != nil {
				r.Log.Error(err, "generation of password for read user failed", "db", database.Spec.DatabaseName)
				return ctrl.Result{}, err
			}

			writePassword, err := password.Generate(20, 5, 0, false, false)
			if err != nil {
				r.Log.Error(err, "generation of password for write user failed", "db", database.Spec.DatabaseName)
				return ctrl.Result{}, err
			}

			secret.StringData = map[string]string{}

			secret.StringData[fmt.Sprintf("%s-read", database.Spec.DatabaseName)] = readPassword
			secret.StringData[fmt.Sprintf("%s-write", database.Spec.DatabaseName)] = writePassword

			err = r.Create(ctx, secret)
			if err != nil {
				r.Log.Error(err, fmt.Sprintf("error creating Secret %s/%s", database.Namespace, secretName))
				return ctrl.Result{}, err
			} else {
				r.Log.Info("created credential secret", "db", database.Spec.DatabaseName)
			}
		} else {
			r.Log.Error(err, fmt.Sprintf("error getting Secret %s/%s", database.Namespace, secretName))
			return ctrl.Result{}, err
		}
	}

	readName := fmt.Sprintf("%s-read", database.Spec.DatabaseName)
	writeName := fmt.Sprintf("%s-write", database.Spec.DatabaseName)

	err = r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: database.Namespace}, secret)
	if err != nil {
		return ctrl.Result{}, err
	}

	users := []User{{name: readName, password: string(secret.Data[readName])}, {name: writeName, password: string(secret.Data[writeName])}}

	// Create Stardog resources
	instance := &stardogv1beta1.Instance{}
	err = r.Get(ctx, types.NamespacedName{Name: database.Spec.InstanceRef.Name, Namespace: database.Namespace}, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	credentialSecret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Spec.AdminCredentialRef.Name, Namespace: database.Namespace}, credentialSecret)
	if err != nil {
		return ctrl.Result{}, err
	}

	apiClient := stardogapi.NewClient(instance.Spec.AdminCredentialRef.Key, string(credentialSecret.Data[instance.Spec.AdminCredentialRef.Key]), instance.Spec.URL)

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

	err = r.createUsers(ctx, apiClient, users)
	if err != nil {
		r.Log.Error(err, "error creating users", "users", users)
		return ctrl.Result{}, err
	}

	roles := []string{readName, writeName}
	err = r.createRoles(ctx, apiClient, roles)
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

	// TODO remove unwanted permissions?
	err = r.addPermissions(ctx, apiClient, readName, readPermissions)
	if err != nil {
		r.Log.Error(err, "adding permission to role failed", "role", readName, "permission", readPermissions)
		return ctrl.Result{}, err
	}
	err = r.addPermissions(ctx, apiClient, writeName, writePermissions)
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

func comparePermission(x, y stardogapi.Permission) bool {
	if x.ResourceType == y.ResourceType {
		if x.Action == y.Action {
			return slices.Equal(x.Resources, y.Resources)
		}
	}
	return false
}

func (r *DatabaseReconciler) addPermissions(ctx context.Context, apiClient *stardogapi.Client, role string, permissions []stardogapi.Permission) error {
	for _, permission := range permissions {
		exists := false
		rolePermissions, err := apiClient.GetRolePermissions(ctx, role)
		if err != nil {
			return err
		}
		for _, rolePermission := range rolePermissions {
			if comparePermission(rolePermission, permission) {
				exists = true
				break
			}
		}

		if !exists {
			err = apiClient.AddRolePermission(ctx, role, permission)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *DatabaseReconciler) createRoles(ctx context.Context, apiClient *stardogapi.Client, roles []string) error {
	for _, role := range roles {
		activeRoles, err := apiClient.GetRoles(ctx)
		if err != nil {
			return err
		}

		if !slices.Contains(activeRoles, role) {
			err = apiClient.AddRole(ctx, role)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *DatabaseReconciler) createUsers(ctx context.Context, apiClient *stardogapi.Client, users []User) error {
	for _, user := range users {
		_, err := apiClient.GetUser(ctx, user.name)
		if err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				err = apiClient.AddUser(ctx, user.name, user.password)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}
