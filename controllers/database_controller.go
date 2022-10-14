package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/sethvargo/go-password/password"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "k8s.io/api/core/v1"
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
	r.Log = r.Log.WithValues("database", req.NamespacedName)
	// TODO track connected external resources in status

	database := &stardogv1beta1.Database{}
	err := r.Get(ctx, req.NamespacedName, database)
	if err != nil {
		r.Log.Error(err, "error getting database")
		return ctrl.Result{}, err
	}

	for _, instanceRef := range database.Spec.Instances {
		r.Log = r.Log.WithValues("instance", instanceRef)

		instance := &stardogv1beta1.Instance{}

		r.Get(ctx, types.NamespacedName{Namespace: instanceRef.Namespace, Name: instanceRef.Name}, instance)

		apiClient, err := getStardogApiClient(ctx, r.Client, instance)
		if err != nil {
			r.Log.Error(err, "error creating API client")
			break
		}

		liveDatabases, err := apiClient.ListDatabases(ctx)
		if err != nil {
			r.Log.Error(err, "error listing databases")
			break
		}

		if !slices.Contains(liveDatabases, req.Name) {
			err = apiClient.CreateDatabase(ctx, req.Name, nil)
			if err != nil {
				r.Log.Error(err, "error creating database")
				break
			}
		}

		// Generate and save credentials
		secret := &v1.Secret{}

		secretName := fmt.Sprintf("%s-%s-credentials", req.Name, instance.Name)

		err = r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: secretName}, secret)

		if err != nil {
			err = client.IgnoreNotFound(err)
			if err != nil {
				r.Log.Error(err, fmt.Sprintf("error getting Secret %s/%s", req.Namespace, secretName))
				break
			} else {
				// TODO add managed by operator label
				secret = &v1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: req.Namespace, Name: secretName}}

				readPassword, err := password.Generate(20, 5, 0, false, false)
				if err != nil {
					r.Log.Error(err, "generation of password for read user failed", "db", req.Name)
					break
				}

				writePassword, err := password.Generate(20, 5, 0, false, false)
				if err != nil {
					r.Log.Error(err, "generation of password for write user failed", "db", req.Name)
					break
				}

				secret.StringData = map[string]string{}

				secret.StringData[fmt.Sprintf("%s-read", req.Name)] = readPassword
				secret.StringData[fmt.Sprintf("%s-write", req.Name)] = writePassword

				err = r.Create(ctx, secret)
				if err != nil {
					r.Log.Error(err, fmt.Sprintf("error creating Secret %s/%s", req.Namespace, secretName))
					break
				} else {
					r.Log.Info("Created credential secret", "db", req.Name)
				}
			}
		}

		readName := fmt.Sprintf("%s-read", req.Name)
		writeName := fmt.Sprintf("%s-write", req.Name)

		err = r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: secretName}, secret)
		if err != nil {
			r.Log.Error(err, "error getting secret", "secret", secretName)
			break
		}

		users := []User{{name: readName, password: string(secret.Data[readName])}, {name: writeName, password: string(secret.Data[writeName])}}

		err = r.createUsers(ctx, apiClient, users)
		if err != nil {
			r.Log.Error(err, "error creating users", "users", users)
			break
		}

		roles := []string{readName, writeName}
		err = r.createRoles(ctx, apiClient, roles)
		if err != nil {
			r.Log.Error(err, "error creating roles", "roles", roles)
			break
		}

		readPermissions := []stardogapi.Permission{
			{ResourceType: "db", Action: "READ", Resources: []string{req.Name}},
			{ResourceType: "metadata", Action: "READ", Resources: []string{req.Name}},
		}

		writePermissions := []stardogapi.Permission{
			{ResourceType: "db", Action: "WRITE", Resources: []string{req.Name}},
			{ResourceType: "metadata", Action: "WRITE", Resources: []string{req.Name}},
		}

		// TODO remove unwanted permissions?
		err = r.addPermissions(ctx, apiClient, readName, readPermissions)
		if err != nil {
			r.Log.Error(err, "adding permission to role failed", "role", readName, "permission", readPermissions)
			break
		}
		err = r.addPermissions(ctx, apiClient, writeName, writePermissions)
		if err != nil {
			r.Log.Error(err, "adding permission to role failed", "role", writeName, "permission", writePermissions)
			break
		}

		// assign roles
		readUserRoles, err := apiClient.GetUserRoles(ctx, readName)
		if err != nil {
			r.Log.Error(err, "error getting user roles", "user", readName)
			break
		}

		if len(readUserRoles) != 1 || readUserRoles[0] != readName {
			err = apiClient.SetUserRoles(ctx, readName, []string{readName})
			if err != nil {
				r.Log.Error(err, "error assigning roles to user", "user", readName, "roles", readUserRoles)
				break
			}
		}

		// TODO check if we actually need read when assigning write
		writeUserRoles, err := apiClient.GetUserRoles(ctx, writeName)
		if err != nil {
			r.Log.Error(err, "error getting user roles", "user", writeName)
			break
		}

		if len(writeUserRoles) != 2 || writeUserRoles[0] != readName || writeUserRoles[1] != writeName {
			err = apiClient.SetUserRoles(ctx, writeName, []string{readName, writeName})
			if err != nil {
				r.Log.Error(err, "error assigning roles to user", "user", readName, "roles", readUserRoles)
				break
			}
		}
	}

	return ctrl.Result{}, err
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
