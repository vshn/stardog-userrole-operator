package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-openapi/runtime"
	"github.com/sethvargo/go-password/password"

	stardog "github.com/vshn/stardog-userrole-operator/stardogrest/client"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/db"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles_permissions"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users_roles"
	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	scheme "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	stardogv1alpha1 "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	stardogv1beta1 "github.com/vshn/stardog-userrole-operator/api/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/types"
)

const databaseFinalizer = "finalizer.stardog.databases"

var defaultDBOptions = map[string]interface{}{
	"transaction.write.conflict.strategy": "abort_on_conflict",
	"index.aggregate":                     "On",
	"spatial.enabled":                     "true",
	"transaction.logging":                 "true",
	"query.all.graphs":                    "true",
	"preserve.bnode.ids":                  "false",
}

type stardogDatabaseCreate struct {
	Name    string                 `json:"dbname,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *scheme.Scheme
}

//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases/status,verbs=get;update;patch

// Reconcile manages the Stardog resources for a Database object
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	database := &stardogv1beta1.Database{}
	err := r.Get(ctx, req.NamespacedName, database)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("Database not found, ignoring reconcile.")
			return ctrl.Result{Requeue: false}, nil
		}
		r.Log.Error(err, "Could not retrieve Database.")
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, err
	}

	dr := &DatabaseReconciliation{
		reconciliationContext: &ReconciliationContext{
			context:       ctx,
			conditions:    make(map[stardogv1alpha1.StardogConditionType]stardogv1alpha1.StardogCondition),
			stardogClient: stardog.NewHTTPClient(nil),
		},
		resource: database,
	}

	return r.reconcileDatabase(dr)
}

func (r *DatabaseReconciler) reconcileDatabase(dr *DatabaseReconciliation) (ctrl.Result, error) {
	rc := dr.reconciliationContext
	database := dr.resource

	r.Log.Info("reconciling", getLoggingKeysAndValuesForDatabase(database)...)

	isStardogDatabaseMarkedToBeDeleted := database.GetDeletionTimestamp() != nil
	if isStardogDatabaseMarkedToBeDeleted {
		if err := r.deleteDatabases(dr); err != nil {
			rc.SetStatusCondition(createStatusConditionTerminating(err))
			rc.SetStatusCondition(createStatusConditionReady(false, "StardogDatabase cannot be deleted"))
			return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(dr)
		}
		return ctrl.Result{Requeue: false}, nil
	}

	if err := r.validateSpecification(dr.reconciliationContext.context, dr.resource); err != nil {
		rc.SetStatusCondition(createStatusConditionInvalid(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Specification cannot be validated"))
		return ctrl.Result{Requeue: false}, r.updateStatus(dr)
	}
	rc.SetStatusIfExisting(stardogv1alpha1.StardogInvalid, v1.ConditionFalse)

	if err := r.syncDB(dr); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Synchronization failed"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(dr)
	}
	rc.SetStatusIfExisting(stardogv1alpha1.StardogErrored, v1.ConditionFalse)

	r.Log.V(1).Info("adding Finalizer for the StardogDatabase")
	controllerutil.AddFinalizer(dr.resource, databaseFinalizer)

	if err := r.Update(dr.reconciliationContext.context, dr.resource); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Cannot update database"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(dr)
	}

	rc.SetStatusCondition(createStatusConditionReady(true, "Synchronized"))
	return ctrl.Result{Requeue: true, RequeueAfter: ReconFreq}, r.updateStatus(dr)
}

func (r *DatabaseReconciler) updateStatus(dr *DatabaseReconciliation) error {
	res := dr.resource
	status := res.Status
	// Once we are on Kubernetes 0.19, we can use metav1.Conditions, but for now, we have to implement our helpers on
	// our own.
	status.Conditions = mergeWithExistingConditions(status.Conditions, dr.reconciliationContext.conditions)
	res.Status = status

	err := r.Client.Status().Update(dr.reconciliationContext.context, res)
	if err != nil {
		r.Log.Error(err, "could not update Database", getLoggingKeysAndValuesForDatabase(res)...)
		return err
	}
	r.Log.Info("updated Database status", getLoggingKeysAndValuesForDatabase(res)...)
	return nil
}

func (r *DatabaseReconciler) deleteDatabases(dr *DatabaseReconciliation) error {
	instances := dr.resource.Spec.StardogInstanceRefs
	database := dr.resource
	r.Log.Info(fmt.Sprintf("checking if Stardog Database %s is deletable for each instance %s", dr.resource.Name, instances))

	// Do not delete the database unless there are no organizations
	orgs := &stardogv1beta1.OrganizationList{}
	err := r.Client.List(dr.reconciliationContext.context, orgs)
	if err != nil {
		return fmt.Errorf("cannot get organization list for database %s: %v", database.Spec.DatabaseName, err)
	}
	if len(orgs.Items) > 0 {
		return fmt.Errorf("cannot delete database while having %d organizations", len(orgs.Items))
	}

	for _, instance := range instances {
		if err := r.deleteDatabase(dr, instance); err != nil {
			return fmt.Errorf("cannot delete database: %v", err)
		}

		database.Spec.StardogInstanceRefs = removeIndex(database.Spec.StardogInstanceRefs, instance)
		err := r.Update(dr.reconciliationContext.context, database)
		if err != nil {
			return fmt.Errorf("cannot update Database CRD: %v", err)
		}
	}

	controllerutil.RemoveFinalizer(database, databaseFinalizer)
	err = r.Update(dr.reconciliationContext.context, database)
	if err != nil {
		return fmt.Errorf("cannot update database: %v", err)
	}

	return nil
}

func (r *DatabaseReconciler) deleteDatabase(dr *DatabaseReconciliation, instance stardogv1beta1.StardogInstanceRef) error {
	database := dr.resource

	r.Log.V(1).Info("setup Stardog Client from ", "ref", instance)
	auth, err := dr.reconciliationContext.initStardogClientFromRef(r.Client, instance)
	if err != nil {
		return err
	}

	stardogClient := dr.reconciliationContext.stardogClient
	dbName := database.Spec.DatabaseName

	// Do not delete the database unless it's empty
	sizeParams := db.NewGetDBSizeParams().WithDb(dbName).WithExact(pointer.Bool(false))
	dbSize, err := stardogClient.Db.GetDBSize(sizeParams, auth)
	if err != nil || !dbSize.IsSuccess() {
		return fmt.Errorf("cannot determine the size of the database %s: %v", dbName, err)
	}
	if dbSize.Payload != "0" {
		return fmt.Errorf("cannot delete non empty database %s", dbName)
	}

	read, write := getUserRoleNames(database.Spec.DatabaseName)

	// Remove permissions
	permReadParam := roles_permissions.NewRemoveRolePermissionParams().WithRole(read)
	for _, p := range getDBReadPermissions(database) {
		permResp, err := stardogClient.RolesPermissions.RemoveRolePermission(permReadParam.WithPermission(&p), auth)
		if err != nil || !permResp.IsSuccess() {
			return fmt.Errorf("cannot remove permission %#v of role %s: %v", p, read, err)
		}
	}
	permWriteParam := roles_permissions.NewRemoveRolePermissionParams().WithRole(write)
	for _, p := range getDBWritePermissions(database) {
		permResp, err := stardogClient.RolesPermissions.RemoveRolePermission(permWriteParam.WithPermission(&p), auth)
		if err != nil || !permResp.IsSuccess() {
			return fmt.Errorf("cannot remove permission %#v of role %s: %v", p, write, err)
		}
	}

	// Remove assigned roles to users
	param := users_roles.NewRemoveRoleOfUserParams()
	roleUserResp, err := stardogClient.UsersRoles.RemoveRoleOfUser(param.WithUser(read).WithRole(read), auth)
	if err != nil || !roleUserResp.IsSuccess() {
		return fmt.Errorf("cannot remove assigned role %s from user %s: %v", read, read, err)
	}
	roleUserResp, err = stardogClient.UsersRoles.RemoveRoleOfUser(param.WithUser(write).WithRole(read), auth)
	if err != nil || !roleUserResp.IsSuccess() {
		return fmt.Errorf("cannot remove assigned role %s from user %s: %v", read, write, err)
	}
	roleUserResp, err = stardogClient.UsersRoles.RemoveRoleOfUser(param.WithUser(write).WithRole(write), auth)
	if err != nil || !roleUserResp.IsSuccess() {
		return fmt.Errorf("cannot remove assigned role %s from user %s: %v", write, write, err)
	}

	// Remove read and write roles
	roleParam := roles.NewRemoveRoleParams()
	roleRead, err := stardogClient.Roles.RemoveRole(roleParam.WithRole(read), auth)
	if err != nil || !roleRead.IsSuccess() {
		return fmt.Errorf("cannot remove read role %s: %v", read, err)
	}
	roleWrite, err := stardogClient.Roles.RemoveRole(roleParam.WithRole(write), auth)
	if err != nil || !roleWrite.IsSuccess() {
		return fmt.Errorf("cannot remove write role %s: %v", write, err)
	}

	// Remove read and write users
	userParam := users.NewRemoveUserParams()
	userRead, err := stardogClient.Users.RemoveUser(userParam.WithUser(read), auth)
	if err != nil || !userRead.IsSuccess() {
		return fmt.Errorf("cannot remove read user %s: %v", read, err)
	}
	userWrite, err := stardogClient.Users.RemoveUser(userParam.WithUser(write), auth)
	if err != nil || !userWrite.IsSuccess() {
		return fmt.Errorf("cannot remove write user %s: %v", write, err)
	}

	// Remove database
	params := db.NewDropDatabaseParams().WithDb(database.Spec.DatabaseName)
	dbResponse, err := stardogClient.Db.DropDatabase(params, auth)
	if err != nil || !dbResponse.IsSuccess() {
		return fmt.Errorf("error dropping database %s: %w", database.Name, err)
	}

	return nil
}

func (r *DatabaseReconciler) validateSpecification(ctx context.Context, database *stardogv1beta1.Database) error {
	r.Log.V(1).Info("validating StardogRoleSpec")
	spec := &database.Spec
	status := &database.Status

	if len(spec.StardogInstanceRefs) == 0 {
		return fmt.Errorf(".spec.StardogInstanceRefs is required to have at least one instance")
	}
	if spec.DatabaseName == "" {
		return fmt.Errorf(".spec.DatabaseName is required")
	}
	if spec.NamedGraphPrefix == "" {
		return fmt.Errorf(".spec.NamedGraphPrefix is required")
	}

	// If status is not set for database name then we treat it as a creation (first object reconciliation)
	if status.DatabaseName == "" {
		status.NamedGraphPrefix = spec.NamedGraphPrefix
		status.StardogInstanceRefs = spec.StardogInstanceRefs
		status.DatabaseName = spec.DatabaseName
		status.Options = spec.Options

		err := r.Client.Status().Update(ctx, database)
		if err != nil {
			r.Log.Error(err, "could not update Database Status", getLoggingKeysAndValuesForDatabase(database)...)
			return err
		}
	}

	// If status is set for database name then we treat it as an update (2 - n object reconciliation)
	if status.DatabaseName != "" {
		spec.NamedGraphPrefix = status.NamedGraphPrefix
		spec.DatabaseName = status.DatabaseName
		spec.StardogInstanceRefs = status.StardogInstanceRefs
		spec.Options = status.Options
	}

	return nil
}

func (r *DatabaseReconciler) syncDB(dr *DatabaseReconciliation) error {
	rc := dr.reconciliationContext
	database := dr.resource
	stardogClient := dr.reconciliationContext.stardogClient

	for _, instance := range database.Spec.StardogInstanceRefs {
		auth, err := rc.initStardogClientFromRef(r.Client, instance)
		if err != nil {
			return fmt.Errorf("cannot initialize stardog client: %v", err)
		}

		// Generate and save credentials in k8s
		secretName := getUsersCredentialSecret(database.Spec.DatabaseName, instance.Name)
		readName, writeName := getUserRoleNames(database.Spec.DatabaseName)
		credDBSecret, err := r.createCredentials(dr, secretName, readName, writeName)
		if err != nil {
			return err
		}

		// Create database in Stardog if it does not exist
		liveDatabases, err := stardogClient.Db.ListDatabases(nil, auth)
		if err != nil {
			r.Log.Error(err, "error listing databases")
			return err
		}

		if !slices.Contains(liveDatabases.Payload.Databases, database.Spec.DatabaseName) {
			err = createDatabase(database, stardogClient, auth)
			if err != nil {
				return fmt.Errorf("failed to create database %v", err)
			}
			r.Log.Info("created Stardog database", "name", database.Spec.DatabaseName)
		}

		// create default read and write users
		usrs := []models.User{
			{Password: []string{string(credDBSecret.Data[readName])}, Username: &readName},
			{Password: []string{string(credDBSecret.Data[writeName])}, Username: &writeName},
		}
		err = createDefaultUsersForDB(stardogClient, auth, usrs)
		if err != nil {
			r.Log.Error(err, "error creating users", "users", usrs)
			return err
		}

		// create default read and write roles
		rolenames := []models.Rolename{
			{Rolename: &readName},
			{Rolename: &writeName},
		}
		err = createDefaultRolesForDB(stardogClient, auth, rolenames)
		if err != nil {
			r.Log.Error(err, "error creating roles", "roles", rolenames)
			return err
		}

		//create read and write permissions for user roles
		readPerms := getDBReadPermissions(database)
		err = createDefaultPermissions(stardogClient, auth, readName, readPerms)
		if err != nil {
			r.Log.Error(err, "adding permission to role failed", "role", readName, "permission", readPerms)
			return err
		}

		writePerms := getDBWritePermissions(database)
		err = createDefaultPermissions(stardogClient, auth, writeName, writePerms)
		if err != nil {
			r.Log.Error(err, "adding permission to role failed", "role", writeName, "permission", writePerms)
			return err
		}

		// assign roles to users
		readUserRolesResp, err := stardogClient.UsersRoles.ListUserRoles(users_roles.NewListUserRolesParams().WithUser(readName), auth)
		if err != nil || !readUserRolesResp.IsSuccess() {
			r.Log.Error(err, "error getting user roles", "user", readName)
			return err
		}

		if !slices.Contains(readUserRolesResp.Payload.Roles, readName) {
			params := users_roles.NewAddRoleParams().
				WithUser(readName).
				WithRole(&models.Rolename{Rolename: &readName})
			roleResp, err := stardogClient.UsersRoles.AddRole(params, auth)
			if err != nil || !roleResp.IsSuccess() {
				r.Log.Error(err, "error assigning role to user", "user", readName, "role", readName)
				return err
			}
		}

		writeUserRolesResp, err := stardogClient.UsersRoles.ListUserRoles(users_roles.NewListUserRolesParams().WithUser(writeName), auth)
		if err != nil || !writeUserRolesResp.IsSuccess() {
			r.Log.Error(err, "error getting user roles", "user", writeName)
			return err
		}

		if !slices.Contains(writeUserRolesResp.Payload.Roles, readName) {
			params := users_roles.NewAddRoleParams().
				WithUser(writeName).
				WithRole(&models.Rolename{Rolename: &readName})
			roleResp, err := stardogClient.UsersRoles.AddRole(params, auth)
			if err != nil || !roleResp.IsSuccess() {
				r.Log.Error(err, "error assigning role to user", "user", writeName, "role", readName)
				return err
			}
		}

		if !slices.Contains(writeUserRolesResp.Payload.Roles, writeName) {
			params := users_roles.NewAddRoleParams().
				WithUser(writeName).
				WithRole(&models.Rolename{Rolename: &writeName})
			roleResp, err := stardogClient.UsersRoles.AddRole(params, auth)
			if err != nil || !roleResp.IsSuccess() {
				r.Log.Error(err, "error assigning role to user", "user", writeName, "role", writeName)
				return err
			}
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stardogv1beta1.Database{}).
		Complete(r)
}

// Generate and save credentials for Database
func (r *DatabaseReconciler) createCredentials(dr *DatabaseReconciliation, secretName string, readName string, writeName string) (*corev1.Secret, error) {
	database := dr.resource
	rc := dr.reconciliationContext
	ctx := dr.reconciliationContext.context
	secret := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: rc.namespace}, secret)

	if err == nil {
		return secret, nil
	} else if !apierrors.IsNotFound(err) {
		r.Log.Error(err, fmt.Sprintf("error getting Secret %s/%s", rc.namespace, secretName))
		return nil, err
	}

	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: dr.reconciliationContext.namespace,
		}}
	err = controllerutil.SetControllerReference(database, secret, r.Scheme)
	if err != nil {
		return nil, err
	}

	readPassword, err := password.Generate(20, 5, 0, false, false)
	if err != nil {
		r.Log.Error(err, "generation of password for read user failed", "db", database.Spec.DatabaseName)
		return nil, err
	}

	writePassword, err := password.Generate(20, 5, 0, false, false)
	if err != nil {
		r.Log.Error(err, "generation of password for write user failed", "db", database.Spec.DatabaseName)
		return nil, err
	}

	secret.StringData = map[string]string{}

	secret.StringData[readName] = readPassword
	secret.StringData[writeName] = writePassword

	err = r.Create(ctx, secret)
	if err != nil && apierrors.IsAlreadyExists(err) {
		err = r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: secret.Namespace}, secret)
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("error creating secret %s/%s", rc.namespace, secretName))
			return nil, fmt.Errorf("cannot get existing credentials from secret %s: %v", secretName, err)
		}
		return secret, nil
	}
	if err != nil {
		r.Log.Error(err, fmt.Sprintf("error creating secret %s/%s", rc.namespace, secretName))
		return nil, err
	}

	r.Log.Info("created credential secret", "namespace", secret.Namespace, "name", secret.Name)
	return secret, nil
}

func createDatabase(database *stardogv1beta1.Database, stardogClient *stardog.Stardog, auth runtime.ClientAuthInfoWriter) error {
	dbName := database.Spec.DatabaseName
	options := database.Spec.Options

	if options != "" {
		err := json.Unmarshal([]byte(options), &defaultDBOptions)
		if err != nil {
			return fmt.Errorf("cannot unmarshal options of json type: %v", err)
		}
	}

	payload, err := json.Marshal(stardogDatabaseCreate{
		Name:    dbName,
		Options: defaultDBOptions,
	})
	if err != nil {
		return fmt.Errorf("cannot marshal root parameter to create database %s: %v", dbName, err)
	}

	params := db.NewCreateNewDatabaseParams().WithRoot(string(payload))
	newDatabaseResp, err := stardogClient.Db.CreateNewDatabase(params, auth)
	if err != nil || !newDatabaseResp.IsSuccess() {
		return fmt.Errorf("error creating database %s: %v", dbName, err)
	}
	return nil
}

func createDefaultUsersForDB(stardogClient *stardog.Stardog, auth runtime.ClientAuthInfoWriter, usrs []models.User) error {
	existingUsers, err := stardogClient.Users.ListUsers(nil, auth)
	if err != nil || !existingUsers.IsSuccess() {
		return fmt.Errorf("error listing users: %w", err)
	}

	for _, user := range usrs {
		if !slices.Contains(existingUsers.Payload.Users, *user.Username) {
			createUserResp, err := stardogClient.Users.CreateUser(users.NewCreateUserParams().WithUser(&user), auth)
			if err != nil || !createUserResp.IsSuccess() {
				return fmt.Errorf("error create database user %s: %w", *user.Username, err)
			}
		}
	}
	return nil
}

func createDefaultPermissions(stardogClient *stardog.Stardog, auth runtime.ClientAuthInfoWriter, role string, perms []models.Permission) error {
	listParams := roles_permissions.NewListRolePermissionsParams().WithRole(role)
	existingPermissionsResp, err := stardogClient.RolesPermissions.ListRolePermissions(listParams, auth)
	if err != nil || !existingPermissionsResp.IsSuccess() {
		return fmt.Errorf("error listing role permissions: %w", err)
	}

	for _, perm := range perms {
		if !containsPermission(existingPermissionsResp.Payload.Permissions, perm) {
			params := roles_permissions.NewAddRolePermissionParams().WithRole(role).WithPermission(&perm)
			permissionResp, err := stardogClient.RolesPermissions.AddRolePermission(params, auth)
			if err != nil || !permissionResp.IsSuccess() {
				return fmt.Errorf("error create permission %#v for role %s: %w", perm, role, err)
			}
		}
	}
	return nil
}

func createDefaultRolesForDB(stardogClient *stardog.Stardog, auth runtime.ClientAuthInfoWriter, rolenames []models.Rolename) error {
	existingRoles, err := stardogClient.Roles.ListRoles(nil, auth)
	if err != nil || !existingRoles.IsSuccess() {
		return fmt.Errorf("error getting roles: %w", err)
	}

	for _, role := range rolenames {
		if !slices.Contains(existingRoles.Payload.Roles, *role.Rolename) {
			createRoleResp, err := stardogClient.Roles.CreateRole(roles.NewCreateRoleParams().WithRole(&role), auth)
			if err != nil || !createRoleResp.IsSuccess() {
				return fmt.Errorf("error creating role %s: %w", *role.Rolename, err)
			}
		}
	}
	return nil
}

func getDBWritePermissions(database *stardogv1beta1.Database) []models.Permission {
	return []models.Permission{
		{
			Action:       pointer.String("WRITE"),
			Resource:     []string{database.Spec.DatabaseName},
			ResourceType: pointer.String("db"),
		},
		{
			Action:       pointer.String("WRITE"),
			Resource:     []string{database.Spec.DatabaseName},
			ResourceType: pointer.String("metadata"),
		},
	}
}

func getDBReadPermissions(database *stardogv1beta1.Database) []models.Permission {
	return []models.Permission{
		{
			Action:       pointer.String("READ"),
			Resource:     []string{database.Spec.DatabaseName},
			ResourceType: pointer.String("db"),
		},
		{
			Action:       pointer.String("READ"),
			Resource:     []string{database.Spec.DatabaseName},
			ResourceType: pointer.String("metadata"),
		},
	}
}

func getLoggingKeysAndValuesForDatabase(stardogDatabase *stardogv1beta1.Database) []interface{} {
	return []interface{}{
		"StardogDatabase", stardogDatabase.Name,
	}
}

func getUserRoleNames(dbName string) (read, write string) {
	return fmt.Sprintf("%s-read", dbName), fmt.Sprintf("%s-write", dbName)
}

func getUsersCredentialSecret(dbName, instance string) string {
	return fmt.Sprintf("%s-%s-credentials", dbName, instance)
}
