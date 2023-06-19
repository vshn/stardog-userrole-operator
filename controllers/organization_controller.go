package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/sethvargo/go-password/password"
	stardogv1alpha1 "github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	stardogv1beta1 "github.com/vshn/stardog-userrole-operator/api/v1beta1"
	stardog "github.com/vshn/stardog-userrole-operator/stardogrest/client"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles_permissions"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users_roles"
	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	scheme "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
)

const orgFinalizer = "finalizer.stardog.organizations"

// OrganizationReconciler reconciles a Organization object
type OrganizationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *scheme.Scheme
}

//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=stardog.vshn.ch,resources=databases/status,verbs=get;update;patch

// Reconcile manages the Stardog resources for a Database object
func (r *OrganizationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	organization := &stardogv1beta1.Organization{}
	err := r.Get(ctx, req.NamespacedName, organization)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("Organization not found, ignoring reconcile.")
			return ctrl.Result{Requeue: false}, nil
		}
		r.Log.Error(err, "Could not retrieve Organization.")
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, err
	}

	or := &OrganizationReconciliation{
		reconciliationContext: &ReconciliationContext{
			context:       ctx,
			conditions:    make(map[stardogv1alpha1.StardogConditionType]stardogv1alpha1.StardogCondition),
			stardogClient: stardog.NewHTTPClient(nil),
		},
		resource: organization,
	}

	return r.reconcileOrganization(or)
}

func (r *OrganizationReconciler) reconcileOrganization(or *OrganizationReconciliation) (ctrl.Result, error) {
	rc := or.reconciliationContext
	organization := or.resource

	r.Log.Info("reconciling", getLoggingKeysAndValuesForOrganization(organization)...)

	if err := r.validateSpecification(or.reconciliationContext.context, organization); err != nil {
		rc.SetStatusCondition(createStatusConditionInvalid(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Specification cannot be validated"))
		return ctrl.Result{Requeue: false}, r.updateStatus(or)
	}
	rc.SetStatusIfExisting(stardogv1alpha1.StardogInvalid, v1.ConditionFalse)

	if err := r.getDatabaseRef(or); err != nil {
		rc.SetStatusCondition(createStatusConditionTerminating(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Cannot get StardogDatabase from reference"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(or)
	}

	isStardogOrganizationMarkedToBeDeleted := organization.GetDeletionTimestamp() != nil
	if isStardogOrganizationMarkedToBeDeleted {
		if err := r.deleteOrganizations(or); err != nil {
			rc.SetStatusCondition(createStatusConditionTerminating(err))
			rc.SetStatusCondition(createStatusConditionReady(false, "Organization cannot be deleted"))
			return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(or)
		}
		return ctrl.Result{Requeue: false}, nil
	}

	if err := r.syncOrganization(or); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Synchronization failed"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(or)
	}
	rc.SetStatusIfExisting(stardogv1alpha1.StardogErrored, v1.ConditionFalse)

	r.Log.V(1).Info("adding Finalizer for the StardogDatabase")
	controllerutil.AddFinalizer(or.resource, orgFinalizer)

	if err := r.Update(or.reconciliationContext.context, or.resource); err != nil {
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Cannot update organization"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(or)
	}

	or.resource.Status.StardogInstanceRefs = or.database.Status.StardogInstanceRefs
	or.resource.Status.NamedGraphs = or.resource.Spec.NamedGraphs
	rc.SetStatusCondition(createStatusConditionReady(true, "Synchronized"))
	return ctrl.Result{Requeue: true, RequeueAfter: ReconFreq}, r.updateStatus(or)
}

func (r *OrganizationReconciler) updateStatus(or *OrganizationReconciliation) error {
	res := or.resource
	status := res.Status
	// Once we are on Kubernetes 0.19, we can use metav1.Conditions, but for now, we have to implement our helpers on
	// our own.
	status.Conditions = mergeWithExistingConditions(status.Conditions, or.reconciliationContext.conditions)
	res.Status = status

	err := r.Client.Status().Update(or.reconciliationContext.context, res)
	if err != nil {
		r.Log.Error(err, "could not update Organization", getLoggingKeysAndValuesForOrganization(res)...)
		return err
	}
	r.Log.Info("updated Organization status", getLoggingKeysAndValuesForOrganization(res)...)
	return nil
}

func (r *OrganizationReconciler) validateSpecification(ctx context.Context, organization *stardogv1beta1.Organization) error {
	r.Log.V(1).Info("validating StardogRoleSpec")
	spec := &organization.Spec
	status := &organization.Status

	if len(spec.NamedGraphs) == 0 {
		return fmt.Errorf(".spec.NamedGraphs is required to have at least one graph")
	}
	if spec.Name == "" {
		return fmt.Errorf(".spec.Name is required")
	}
	if spec.DatabaseRef == "" {
		return fmt.Errorf(".spec.DatabaseRef is required")
	}
	if spec.DisplayName == "" {
		return fmt.Errorf(".spec.DisplayName is required")
	}

	// If status is not set for organization name then we treat it as a creation (first object reconciliation)
	if status.Name == "" {
		status.Name = spec.Name
		status.NamedGraphs = spec.NamedGraphs
		status.DatabaseRef = spec.DatabaseRef

		err := r.Client.Status().Update(ctx, organization)
		if err != nil {
			r.Log.Error(err, "Could not update Organization Status", getLoggingKeysAndValuesForOrganization(organization)...)
			return err
		}
	}

	// If status is set for database name then we treat it as an update (2 - n object reconciliation)
	if status.Name != "" {
		spec.Name = status.Name
		spec.DatabaseRef = status.DatabaseRef
	}

	return nil
}

func (r *OrganizationReconciler) syncOrganization(or *OrganizationReconciliation) error {
	dbInstances := or.database.Status.StardogInstanceRefs
	orgInstances := or.resource.Status.StardogInstanceRefs

	// Create an organization for each instance in spec.StardogInstanceRefs
	for _, instance := range dbInstances {
		if err := r.sync(or, instance); err != nil {
			return fmt.Errorf("unable to sync instance %s for organization %s", instance.Name, or.resource.Name)
		}
	}

	// Remove an organization for any removed instance from spec.StardogInstanceRefs
	for _, instance := range getRemovedInstances(dbInstances, orgInstances) {
		if err := r.deleteOrganization(or, instance); err != nil {
			return fmt.Errorf("cannot delete organization %s for instance %s: %v", or.resource.Name, instance.Name, err)
		}
	}

	return nil
}

func (r *OrganizationReconciler) sync(or *OrganizationReconciliation, instance stardogv1beta1.StardogInstanceRef) error {
	rc := or.reconciliationContext
	database := or.database
	dbName := database.Spec.DatabaseName
	org := or.resource
	orgName := org.Spec.Name
	stardogClient := or.reconciliationContext.stardogClient

	auth, err := rc.initStardogClientFromRef(r.Client, instance)
	if err != nil {
		return fmt.Errorf("cannot initialize stardog client: %v", err)
	}

	// Generate and save credentials in k8s
	secretName := getUsersCredentialSecret(dbName, orgName)
	userRoleName := getUserAndRoleName(dbName, orgName)
	credDBSecret, err := r.createCredentials(or, secretName, userRoleName)
	if err != nil {
		return err
	}

	// create default write user for organization
	usrs := []models.User{
		{Password: []string{string(credDBSecret.Data[userRoleName])}, Username: &userRoleName},
	}
	err = createDefaultUsersForDB(stardogClient, auth, usrs)
	if err != nil {
		r.Log.Error(err, "error creating users", "users", usrs)
		return err
	}

	// create default read and write roles
	rolenames := []models.Rolename{
		{Rolename: &userRoleName},
	}
	err = createDefaultRolesForDB(stardogClient, auth, rolenames)
	if err != nil {
		r.Log.Error(err, "error creating roles", "roles", rolenames)
		return err
	}

	//create read and write permissions for user roles
	perms := getOrganizationPerms(database, org, dbName)
	err = createDefaultPermissions(stardogClient, auth, userRoleName, perms)
	if err != nil {
		r.Log.Error(err, "adding permission to role failed", "role", userRoleName, "permission", perms)
		return err
	}

	// Remove permissions in case name graphs have been removed
	for _, graph := range org.Status.NamedGraphs {
		if !contains(org.Spec.NamedGraphs, graph) {
			perm := roles_permissions.NewRemoveRolePermissionParams().WithRole(userRoleName)
			ng := getFullNamedGraph(org.Spec.Name, database.Spec.NamedGraphPrefix, graph)
			for _, p := range getGraphPermissionForNameGraphs(ng, dbName) {
				pResp, err := stardogClient.RolesPermissions.RemoveRolePermission(perm.WithPermission(&p), auth)
				if err != nil || !pResp.IsSuccess() {
					return fmt.Errorf("cannot remove permission %+v for graph %s: %v", p, ng, err)
				}
			}
		}
	}

	// assign roles to users
	readUserRolesResp, err := stardogClient.UsersRoles.ListUserRoles(users_roles.NewListUserRolesParams().WithUser(userRoleName), auth)
	if err != nil || !readUserRolesResp.IsSuccess() {
		r.Log.Error(err, "error getting user roles", "user", userRoleName)
		return err
	}

	if !slices.Contains(readUserRolesResp.Payload.Roles, userRoleName) {
		params := users_roles.NewAddRoleParams().
			WithUser(userRoleName).
			WithRole(&models.Rolename{Rolename: &userRoleName})
		roleResp, err := stardogClient.UsersRoles.AddRole(params, auth)
		if err != nil || !roleResp.IsSuccess() {
			r.Log.Error(err, "error assigning role to user", "user", userRoleName, "role", userRoleName)
			return err
		}
	}
	return nil
}

func (r *OrganizationReconciler) getDatabaseRef(or *OrganizationReconciliation) error {
	org := or.resource
	ctx := or.reconciliationContext.context
	r.Log.V(1).Info("Getting Database from reference", "db", org.Spec.DatabaseRef, "organization", org.Name)

	or.database = &stardogv1beta1.Database{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: org.Spec.DatabaseRef}, or.database)
	if err != nil {
		return fmt.Errorf("cannot get database %s for organization %s: %v", org.Spec.DatabaseRef, org.Name, err)
	}

	return nil
}

func (r *OrganizationReconciler) deleteOrganizations(or *OrganizationReconciliation) error {
	instances := or.database.Status.StardogInstanceRefs
	org := or.resource
	r.Log.Info(fmt.Sprintf("deleting organization %s for each Stardog instances %s", or.resource.Name, instances))

	for _, instance := range instances {
		if err := r.deleteOrganization(or, instance); err != nil {
			return fmt.Errorf("cannot delete organization %s: %v", org.Spec.Name, err)
		}
	}

	controllerutil.RemoveFinalizer(org, orgFinalizer)
	err := r.Update(or.reconciliationContext.context, org)
	if err != nil {
		return fmt.Errorf("cannot update organization: %v", err)
	}

	return nil
}

func (r *OrganizationReconciler) deleteOrganization(or *OrganizationReconciliation, instance stardogv1beta1.StardogInstanceRef) error {
	org := or.resource
	database := or.database
	orgName := org.Spec.Name

	r.Log.V(1).Info("setup Stardog Client from ", "ref", instance)
	auth, err := or.reconciliationContext.initStardogClientFromRef(r.Client, instance)
	if err != nil {
		return err
	}

	stardogClient := or.reconciliationContext.stardogClient
	dbName := database.Spec.DatabaseName

	userRoleName := getUserAndRoleName(dbName, orgName)

	// Remove all permissions
	permParam := roles_permissions.NewRemoveRolePermissionParams().WithRole(userRoleName)
	for _, p := range getOrganizationPerms(database, org, dbName) {
		_, err = stardogClient.RolesPermissions.RemoveRolePermission(permParam.WithPermission(&p), auth)
		if err != nil && !NotFound(err) {
			return fmt.Errorf("cannot remove permission %#v of role %s: %v", p, userRoleName, err)
		}
	}

	// Remove assigned role to user
	param := users_roles.NewRemoveRoleOfUserParams()
	_, err = stardogClient.UsersRoles.RemoveRoleOfUser(param.WithUser(userRoleName).WithRole(userRoleName), auth)
	if err != nil && !NotFound(err) {
		return fmt.Errorf("cannot remove assigned role %s from user %s: %v", userRoleName, userRoleName, err)
	}

	// Remove role
	roleParam := roles.NewRemoveRoleParams()
	_, err = stardogClient.Roles.RemoveRole(roleParam.WithRole(userRoleName), auth)
	if err != nil && !NotFound(err) {
		return fmt.Errorf("cannot remove role %s: %v", userRoleName, err)
	}

	// Remove read and write users
	userParam := users.NewRemoveUserParams()
	_, err = stardogClient.Users.RemoveUser(userParam.WithUser(userRoleName), auth)
	if err != nil && !NotFound(err) {
		return fmt.Errorf("cannot remove user %s: %v", userRoleName, err)
	}

	return nil
}

func (r *OrganizationReconciler) createCredentials(or *OrganizationReconciliation, secretName, userName string) (*v1.Secret, error) {

	org := or.resource
	rc := or.reconciliationContext
	ctx := rc.context
	namespace := rc.namespace
	secret := &v1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: namespace}, secret)

	if err == nil {
		return secret, nil
	}

	if !apierrors.IsNotFound(err) {
		r.Log.Error(err, fmt.Sprintf("error getting Secret %s/%s", namespace, secretName))
		return nil, err
	}

	secret = &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		}}
	err = controllerutil.SetControllerReference(org, secret, r.Scheme)
	if err != nil {
		return nil, err
	}

	pwd, err := password.Generate(20, 5, 0, false, false)
	if err != nil {
		r.Log.Error(err, "generation of password for read user failed", "organization", org.Spec.Name)
		return nil, err
	}

	secret.StringData = map[string]string{
		userName: pwd,
	}

	err = r.Create(ctx, secret)
	if err != nil && apierrors.IsAlreadyExists(err) {
		err = r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: secret.Namespace}, secret)
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("error creating secret %s/%s", namespace, secretName))
			return nil, fmt.Errorf("cannot get existing credentials from secret %s: %v", secretName, err)
		}
		return secret, nil
	}
	if err != nil {
		r.Log.Error(err, fmt.Sprintf("error creating secret %s/%s", namespace, secretName))
		return nil, err
	}

	r.Log.Info("created credential secret", "namespace", secret.Namespace, "name", secret.Name)
	return secret, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrganizationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stardogv1beta1.Organization{}).
		Complete(r)
}

func getLoggingKeysAndValuesForOrganization(organization *stardogv1beta1.Organization) []interface{} {
	return []interface{}{
		"StardogOrganization", organization.Namespace + "/" + organization.Name,
	}
}

func getUserAndRoleName(dbName, orgName string) string {
	return fmt.Sprintf("%s-%s", dbName, orgName)
}

func getGraphPermissions(org *stardogv1beta1.Organization, namedGraphPrefix, dbName string) []models.Permission {
	perms := make([]models.Permission, 0)
	for _, ng := range org.Spec.NamedGraphs {
		fullNameNG := getFullNamedGraph(org.Spec.Name, namedGraphPrefix, ng)
		ngPerm := []models.Permission{
			{
				Action:       pointer.String("READ"),
				Resource:     []string{fullNameNG, dbName},
				ResourceType: pointer.String("named-graph"),
			},
			{
				Action:       pointer.String("WRITE"),
				Resource:     []string{fullNameNG, dbName},
				ResourceType: pointer.String("named-graph"),
			},
		}
		perms = append(perms, ngPerm...)
	}
	return perms
}

func getFullNamedGraph(orgName, namedGraphPrefix string, ng string) string {
	return strings.TrimSuffix(namedGraphPrefix, "/") + "/" + orgName + "/" + ng
}

func getGraphPermissionForNameGraphs(namedGraph, dbName string) []models.Permission {
	return []models.Permission{
		{
			Action:       pointer.String("READ"),
			Resource:     []string{namedGraph, dbName},
			ResourceType: pointer.String("named-graph"),
		},
		{
			Action:       pointer.String("WRITE"),
			Resource:     []string{namedGraph, dbName},
			ResourceType: pointer.String("named-graph"),
		},
	}
}

func getOrganizationPerms(database *stardogv1beta1.Database, org *stardogv1beta1.Organization, dbName string) []models.Permission {
	dbPerms := append(getDBReadPermissions(database), getDBWritePermissions(database)...)
	return append(dbPerms, getGraphPermissions(org, database.Spec.NamedGraphPrefix, dbName)...)
}
