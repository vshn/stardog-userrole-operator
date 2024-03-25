package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-openapi/runtime"
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
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sort"
	"strings"
	"time"
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
		r.Log.Error(err, "Specification cannot be validated")
		rc.SetStatusCondition(createStatusConditionInvalid(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Specification cannot be validated"))
		return ctrl.Result{Requeue: false}, r.updateStatus(or)
	}
	rc.SetStatusIfExisting(stardogv1alpha1.StardogInvalid, v1.ConditionFalse)

	if err := r.getDatabaseRef(or); err != nil {
		r.Log.Error(err, "Cannot get StardogDatabase from reference")
		rc.SetStatusCondition(createStatusConditionTerminating(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Cannot get StardogDatabase from reference"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(or)
	}

	isStardogOrganizationMarkedToBeDeleted := organization.GetDeletionTimestamp() != nil
	if isStardogOrganizationMarkedToBeDeleted {
		if err := r.deleteOrganizations(or); err != nil {
			r.Log.Error(err, "Organization cannot be deleted")
			rc.SetStatusCondition(createStatusConditionTerminating(err))
			rc.SetStatusCondition(createStatusConditionReady(false, "Organization cannot be deleted"))
			return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(or)
		}
		return ctrl.Result{Requeue: false}, nil
	}

	if err := r.syncOrganization(or); err != nil {
		r.Log.Error(err, "Synchronization failed")
		rc.SetStatusCondition(createStatusConditionErrored(err))
		rc.SetStatusCondition(createStatusConditionReady(false, "Synchronization failed"))
		return ctrl.Result{Requeue: true, RequeueAfter: ReconFreqErr}, r.updateStatus(or)
	}
	rc.SetStatusIfExisting(stardogv1alpha1.StardogErrored, v1.ConditionFalse)

	r.Log.V(1).Info("adding Finalizer for the StardogDatabase")
	controllerutil.AddFinalizer(or.resource, orgFinalizer)

	if err := r.Update(or.reconciliationContext.context, or.resource); err != nil {
		r.Log.Error(err, "Cannot update organization")
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

	spec.Name = status.Name
	spec.DatabaseRef = status.DatabaseRef

	return nil
}

func (r *OrganizationReconciler) syncOrganization(or *OrganizationReconciliation) error {
	dbInstances := or.database.Spec.StardogInstanceRefs
	orgInstances := or.resource.Status.StardogInstanceRefs

	// Create an organization for each instance in spec.StardogInstanceRefs
	for _, instance := range dbInstances {
		if err := r.sync(or, instance); err != nil {
			return fmt.Errorf("unable to sync instance %s for organization %s: %v", instance.Name, or.resource.Name, err)
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

	auth, disabled, err := rc.initStardogClientFromRef(r.Client, instance)
	if err != nil {
		return fmt.Errorf("cannot initialize stardog client: %v", err)
	}
	if disabled {
		r.Log.Info("skipping resource from reconciliation", "instance", instance.Name, "resource", or.resource.Name)
		return nil
	}

	// Generate and save credentials in k8s
	secretName := getUsersCredentialSecret(dbName, orgName)
	userRoleName := getUserAndRoleName(dbName, orgName)
	if err != nil {
		return err
	}

	// create default write user for organization
	pass, err := generatePassword()
	if err != nil {
		return err
	}
	usr := models.User{
		Password: []string{pass}, Username: &userRoleName,
	}
	usrs, err := createDefaultUsersForDB(stardogClient, auth, []models.User{usr})
	if err != nil {
		r.Log.Error(err, "error creating users", "users", usrs)
		return err
	}
	if len(usrs) != 0 {
		err = r.createCredentials(or, secretName, usrs)
		if err != nil {
			r.Log.Error(err, "error creating secret credentials", "users", usrs)
			return err
		}
	}

	// create default read and write roles
	rolenames := []models.Rolename{
		{Rolename: &userRoleName},
	}
	err = createDefaultRolesForDB(stardogClient, auth, rolenames)
	if err != nil {
		r.Log.Error(err, "Cannot create roles", "roles", rolenames)
		return err
	}

	//create read and write permissions for user roles
	perms := getOrganizationPerms(database, org, true, false)
	err = createDefaultPermissions(stardogClient, auth, userRoleName, perms)
	if err != nil {
		r.Log.Error(err, "Adding permission to role failed", "role", userRoleName, "permission", perms)
		return err
	}

	// Remove permissions in case name graphs have been removed
	err = removePermissions(org, database, stardogClient, auth, userRoleName)
	if err != nil {
		r.Log.Error(err, "Cannot remove permissions")
		return err
	}

	// assign role to user
	err = assignDefaultRole(stardogClient, auth, userRoleName)
	if err != nil {
		r.Log.Error(err, "Cannot assign defaultRoles", "user", userRoleName)
		return err
	}

	//create read permission for public user for this organisation
	roleNameCustomUser := database.Spec.AddUserForNonHiddenGraphs
	permsCustomUser := getOrganizationPerms(database, org, false, true)
	err = createDefaultPermissions(stardogClient, auth, roleNameCustomUser, permsCustomUser)
	if err != nil {
		r.Log.Error(err, "Adding permission to role failed", "role", roleNameCustomUser, "permission", permsCustomUser)
		return err
	}

	err = adjustPermissionsForCustomUser(org, database, stardogClient, auth, roleNameCustomUser)
	if err != nil {
		r.Log.Error(err, "Cannot remove permissions")
		return err
	}

	return nil
}

func assignDefaultRole(stardogClient *stardog.Stardog, auth runtime.ClientAuthInfoWriter, userRoleName string) error {
	readUserRolesResp, err := stardogClient.UsersRoles.ListUserRoles(users_roles.NewListUserRolesParams().WithUser(userRoleName), auth)
	if err != nil || !readUserRolesResp.IsSuccess() {
		return fmt.Errorf("error getting user roles for user %s; %v", userRoleName, err)
	}

	if !slices.Contains(readUserRolesResp.Payload.Roles, userRoleName) {
		params := users_roles.NewAddRoleParams().
			WithUser(userRoleName).
			WithRole(&models.Rolename{Rolename: &userRoleName})
		roleResp, err := stardogClient.UsersRoles.AddRole(params, auth)
		if err != nil || !roleResp.IsSuccess() {
			return fmt.Errorf("error assigning role %s to user %s: %v", userRoleName, userRoleName, err)
		}
	}
	return nil
}

func adjustPermissionsForCustomUser(org *stardogv1beta1.Organization, database *stardogv1beta1.Database, stardogClient *stardog.Stardog, auth runtime.ClientAuthInfoWriter, userRoleName string) error {
	for _, statusGraph := range org.Status.NamedGraphs {
		statusGraphName := statusGraph.Name
		perm := roles_permissions.NewRemoveRolePermissionParams().WithRole(userRoleName)
		if !contains(stardogv1beta1.GetNamedGraphNames(org.Spec.NamedGraphs), statusGraph.Name) {
			return removePermissionForCustomUser(org, database, stardogClient, auth, perm, statusGraphName)
		}
	}
	return nil
}

func removePermissionForCustomUser(org *stardogv1beta1.Organization, database *stardogv1beta1.Database, stardogClient *stardog.Stardog, auth runtime.ClientAuthInfoWriter, perm *roles_permissions.RemoveRolePermissionParams, statusGraphName string) error {
	ng := getFullNamedGraph(org.Spec.Name, database.Spec.NamedGraphPrefix, statusGraphName, false)
	resources := []string{ng, database.Spec.DatabaseName}
	sort.Strings(resources)
	p := models.Permission{
		Action:       pointer.String("READ"),
		Resource:     resources,
		ResourceType: pointer.String("named-graph"),
	}
	pResp, err := stardogClient.RolesPermissions.RemoveRolePermission(perm.WithPermission(&p), auth)
	if err != nil || !pResp.IsSuccess() {
		return fmt.Errorf("cannot remove permission %+v for graph %s: %v", p, ng, err)
	}
	return nil
}

func removePermissions(org *stardogv1beta1.Organization, database *stardogv1beta1.Database, stardogClient *stardog.Stardog, auth runtime.ClientAuthInfoWriter, userRoleName string) error {
	for _, statusGraph := range org.Status.NamedGraphs {
		perm := roles_permissions.NewRemoveRolePermissionParams().WithRole(userRoleName)
		if !contains(stardogv1beta1.GetNamedGraphNames(org.Spec.NamedGraphs), statusGraph.Name) {
			ng := getFullNamedGraph(org.Spec.Name, database.Spec.NamedGraphPrefix, statusGraph.Name, false)
			for _, p := range getGraphPermissionForNameGraphs(ng, database.Spec.DatabaseName) {
				pResp, err := stardogClient.RolesPermissions.RemoveRolePermission(perm.WithPermission(&p), auth)
				if err != nil || !pResp.IsSuccess() {
					return fmt.Errorf("cannot remove permission %+v for graph %s of role %s: %v", p, ng, userRoleName, err)
				}
			}

			// If a NamedGraph has AddHidden then remove the permission from the hidden graph as well
			if statusGraph.AddHidden {
				ngh := getFullNamedGraph(org.Spec.Name, database.Spec.NamedGraphPrefix, statusGraph.Name, true)
				for _, p := range getGraphPermissionForNameGraphs(ngh, database.Spec.DatabaseName) {
					pResp, err := stardogClient.RolesPermissions.RemoveRolePermission(perm.WithPermission(&p), auth)
					if err != nil || !pResp.IsSuccess() {
						return fmt.Errorf("cannot remove permission %+v for graph %s of user %s: %v", p, ngh, userRoleName, err)
					}
				}
			}
		} else {
			// If the named graph still exists but has AddHidden true then delete the hidden graph
			specGraph, err := stardogv1beta1.FindNamedGraphByName(org.Spec.NamedGraphs, statusGraph.Name)
			if err != nil {
				return fmt.Errorf("cannot find spec graph by name %s: %v", statusGraph.Name, err)
			}
			if specGraph.AddHidden == false && statusGraph.AddHidden == true {
				ngh := getFullNamedGraph(org.Spec.Name, database.Spec.NamedGraphPrefix, statusGraph.Name, true)
				for _, p := range getGraphPermissionForNameGraphs(ngh, database.Spec.DatabaseName) {
					pResp, err := stardogClient.RolesPermissions.RemoveRolePermission(perm.WithPermission(&p), auth)
					if err != nil || !pResp.IsSuccess() {
						return fmt.Errorf("cannot remove permission %+v for graph %s of role %s: %v", p, ngh, userRoleName, err)
					}
				}
			}
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
	auth, disabled, err := or.reconciliationContext.initStardogClientFromRef(r.Client, instance)
	if err != nil {
		return fmt.Errorf("cannot initialize stardog client: %v", err)

	}
	if disabled {
		r.Log.Info("skipping resource from reconciliation", "instance", instance.Name, "resource", or.resource.Name)
		return nil
	}

	stardogClient := or.reconciliationContext.stardogClient
	dbName := database.Spec.DatabaseName

	userRoleName := getUserAndRoleName(dbName, orgName)

	// Remove all permissions
	permParam := roles_permissions.NewRemoveRolePermissionParams().WithRole(userRoleName)
	for _, p := range getOrganizationPerms(database, org, true, false) {
		_, err = stardogClient.RolesPermissions.RemoveRolePermission(permParam.WithPermission(&p), auth)
		if err != nil && !NotFound(err) {
			return fmt.Errorf("cannot remove permission %#v of role %s: %v", p, userRoleName, err)
		}
	}

	// Adjust permissions for custom user from the database if exists
	customUser := database.Spec.AddUserForNonHiddenGraphs
	if customUser != "" {
		for _, statusGraph := range org.Status.NamedGraphs {
			statusGraphName := statusGraph.Name
			perm := roles_permissions.NewRemoveRolePermissionParams().WithRole(customUser)
			return removePermissionForCustomUser(org, database, stardogClient, auth, perm, statusGraphName)
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

func (r *OrganizationReconciler) createCredentials(or *OrganizationReconciliation, secretName string, usrs []models.User) error {
	org := or.resource
	rc := or.reconciliationContext
	ctx := rc.context
	namespace := rc.namespace
	secret := &v1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: namespace}, secret)

	if err == nil {
		return nil
	}

	if !apierrors.IsNotFound(err) {
		r.Log.Error(err, fmt.Sprintf("error getting Secret %s/%s", namespace, secretName))
		return err
	}

	secret = &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		}}
	err = controllerutil.SetControllerReference(org, secret, r.Scheme)
	if err != nil {
		return err
	}

	secret.StringData = map[string]string{}
	for _, u := range usrs {
		secret.StringData[*u.Username] = u.Password[0]
	}

	err = r.Create(ctx, secret)
	if err != nil && apierrors.IsAlreadyExists(err) {
		err = r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: secret.Namespace}, secret)
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("error creating secret %s/%s", namespace, secretName))
			return fmt.Errorf("cannot get existing credentials from secret %s: %v", secretName, err)
		}
		return nil
	}
	if err != nil {
		r.Log.Error(err, fmt.Sprintf("error creating secret %s/%s", namespace, secretName))
		return err
	}

	r.Log.Info("created credential secret", "namespace", secret.Namespace, "name", secret.Name)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrganizationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	h := handler.EnqueueRequestsFromMapFunc(triggerOrgReconciliationFromDB(mgr.GetClient()))
	return ctrl.NewControllerManagedBy(mgr).
		For(&stardogv1beta1.Organization{}).
		Watches(&stardogv1beta1.Database{}, h).
		Complete(r)
}

// triggerOrgReconciliation Trigger a reconciliation of all organizations of this database
// This is required to update the custom user permissions from the organizations linked to this db
func triggerOrgReconciliationFromDB(c client.Client) handler.MapFunc {
	return func(ctx context.Context, db client.Object) []reconcile.Request {
		l := log.FromContext(ctx).WithName("triggerOrgReconciliation")
		var orgList stardogv1beta1.OrganizationList
		if err := c.List(ctx, &orgList); err != nil {
			l.Error(err, "failed to get organization list")
			return nil
		}

		reqs := make([]reconcile.Request, 0)
		for _, o := range orgList.Items {
			if db.GetName() == o.Spec.DatabaseRef {
				reqs = append(reqs, reconcile.Request{NamespacedName: types.NamespacedName{Name: o.GetName()}})
			}
		}
		// Wait until database reconciled
		time.Sleep(time.Second * 3)
		return reqs
	}
}

func getLoggingKeysAndValuesForOrganization(organization *stardogv1beta1.Organization) []interface{} {
	return []interface{}{
		"StardogOrganization", organization.Namespace + "/" + organization.Name,
	}
}

func getUserAndRoleName(dbName, orgName string) string {
	return fmt.Sprintf("%s-%s", dbName, orgName)
}

func getGraphPermissions(org *stardogv1beta1.Organization, namedGraphPrefix, dbName string, withHidden, readOnly bool) []models.Permission {
	perms := make([]models.Permission, 0)
	orgName := org.Spec.Name
	for _, ng := range org.Spec.NamedGraphs {
		fullNameNG := getFullNamedGraph(orgName, namedGraphPrefix, ng.Name, false)
		ngPerm := []models.Permission{
			{
				Action:       pointer.String("READ"),
				Resource:     []string{dbName, fullNameNG},
				ResourceType: pointer.String("named-graph"),
			},
		}
		if !readOnly {
			ngPerm = append(ngPerm, models.Permission{
				Action:       pointer.String("WRITE"),
				Resource:     []string{dbName, fullNameNG},
				ResourceType: pointer.String("named-graph"),
			})
		}
		if withHidden && ng.AddHidden {
			fullNameNGH := getFullNamedGraph(orgName, namedGraphPrefix, ng.Name, true)
			ngPermHidden := []models.Permission{
				{
					Action:       pointer.String("READ"),
					Resource:     []string{dbName, fullNameNGH},
					ResourceType: pointer.String("named-graph"),
				},
			}
			if !readOnly {
				ngPermHidden = append(ngPermHidden, models.Permission{
					Action:       pointer.String("WRITE"),
					Resource:     []string{dbName, fullNameNGH},
					ResourceType: pointer.String("named-graph"),
				})
			}
			perms = append(perms, ngPermHidden...)
		}
		perms = append(perms, ngPerm...)
	}
	return perms
}

func getFullNamedGraph(orgName, namedGraphPrefix, ng string, hidden bool) string {
	if hidden {
		return strings.TrimSuffix(namedGraphPrefix, "/") + "/" + orgName + "/" + ng + "/hidden"
	}
	return strings.TrimSuffix(namedGraphPrefix, "/") + "/" + orgName + "/" + ng
}

func getGraphPermissionForNameGraphs(namedGraph, dbName string) []models.Permission {
	resources := []string{namedGraph, dbName}
	sort.Strings(resources)
	return []models.Permission{
		{
			Action:       pointer.String("READ"),
			Resource:     resources,
			ResourceType: pointer.String("named-graph"),
		},
		{
			Action:       pointer.String("WRITE"),
			Resource:     resources,
			ResourceType: pointer.String("named-graph"),
		},
	}
}

func getOrganizationPerms(database *stardogv1beta1.Database, org *stardogv1beta1.Organization, withHidden, readOnly bool) []models.Permission {
	db := database.Spec.DatabaseName
	graphPerm := getGraphPermissions(org, database.Spec.NamedGraphPrefix, db, withHidden, readOnly)
	dbReadPerm := getDBReadPermissions(db)
	perm := append(dbReadPerm, graphPerm...)
	if readOnly {
		return perm
	}
	return append(perm, getDBWritePermissions(db)...)
}
