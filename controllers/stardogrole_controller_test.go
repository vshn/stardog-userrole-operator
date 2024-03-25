package controllers

import (
	"context"
	"errors"
	stardog_client "github.com/vshn/stardog-userrole-operator/stardogrest/client"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles_permissions"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users_roles"
	stardogmock "github.com/vshn/stardog-userrole-operator/stardogrest/mocks"
	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
	"k8s.io/utils/pointer"
	"os"
	"testing"
	"time"

	testr "github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

func Test_syncRole(t *testing.T) {
	namespace := "namespace-test"
	stardogInstanceRef := "instance-test"
	stardogRoleName := "role-test"
	secretName := "secret-test"
	username := "admin"
	password := "1234"
	action1 := "READ"
	action2 := "WRITE"
	resourceType1 := "DB"
	resourceType2 := "*"
	serverURL := "https://stardog-test.com"
	resources1 := []string{"Database1", "Database2"}
	resources2 := []string{"Graph1", "Database2"}

	permissionSpec1 := v1alpha1.StardogPermissionSpec{
		Action:       action1,
		ResourceType: resourceType1,
		Resources:    resources1,
	}
	permissionSpec2 := v1alpha1.StardogPermissionSpec{
		Action:       action2,
		ResourceType: resourceType2,
		Resources:    resources2,
	}
	permission1 := models.Permission{
		Action:       &action1,
		ResourceType: &resourceType1,
		Resource:     resources1,
	}
	permission2 := models.Permission{
		Action:       &action2,
		ResourceType: &resourceType2,
		Resource:     resources2,
	}
	permission3 := models.Permission{
		Action:       &action1,
		ResourceType: &resourceType2,
		Resource:     resources2,
	}

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogMocked := stardogmock.NewMockStardogTestClient(mockCtrl)
	stardogClient := createStardogClientFromMock(stardogMocked)

	tests := []struct {
		name            string
		stardogRole     v1alpha1.StardogRole
		stardogInstance v1alpha1.StardogInstance
		secret          v1.Secret
		srr             StardogRoleReconciliation
		expectations    []func(stardog_client.Stardog)
		err             error
	}{
		{
			name:            "GivenReconciliationContext_WhenThereIsNoRole_ThenCreateRole",
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1, permissionSpec2}),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRole(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1, permissionSpec2}),
			},
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						ListRoles(gomock.Any(), gomock.Any()).
						Return(&roles.ListRolesOK{Payload: &models.Roles{}}, nil).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						CreateRole(roles.NewCreateRoleParams().WithRole(&models.Rolename{Rolename: &stardogRoleName}), gomock.Any()).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						ListRolePermissions(roles_permissions.NewListRolePermissionsParams().WithRole(stardogRoleName), gomock.Any()).
						Return(&roles_permissions.ListRolePermissionsOK{Payload: &models.Permissions{Permissions: []*models.Permission{&permission1, &permission2}}}, nil).
						Times(0)
				},
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						AddRolePermission(roles_permissions.NewAddRolePermissionParams().WithRole(stardogRoleName).WithPermission(&permission1), gomock.Any()).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						AddRolePermission(roles_permissions.NewAddRolePermissionParams().WithRole(stardogRoleName).WithPermission(&permission2), gomock.Any()).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						RemoveRolePermission(roles_permissions.NewRemoveRolePermissionParams().WithRole(stardogRoleName), gomock.Any()).
						Times(0)
				},
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						SetTransport(gomock.Any()).
						AnyTimes()
				},
			},
			err: nil,
		},
		{
			name:            "GivenReconciliationContext_WhenThereIsRole_ThenUpdateRole",
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1, permissionSpec2}),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRole(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1, permissionSpec2}),
			},
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						ListRoles(gomock.Any(), gomock.Any()).
						Return(&roles.ListRolesOK{Payload: &models.Roles{Roles: []string{stardogRoleName}}}, nil).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						CreateRole(gomock.Any(), gomock.Any()).
						Times(0)
				},
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						ListRolePermissions(
							roles_permissions.NewListRolePermissionsParams().
								WithRole(stardogRoleName),
							gomock.Any()).
						Return(&roles_permissions.ListRolePermissionsOK{Payload: &models.Permissions{Permissions: []*models.Permission{&permission2, &permission3}}}, nil).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						AddRolePermission(
							roles_permissions.NewAddRolePermissionParams().
								WithRole(stardogRoleName).
								WithPermission(&permission1),
							gomock.Any()).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.
						EXPECT().
						RemoveRolePermission(
							roles_permissions.NewRemoveRolePermissionParams().
								WithRole(stardogRoleName).
								WithPermission(&permission3),
							gomock.Any()).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						SetTransport(gomock.Any()).
						AnyTimes()
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.stardogInstance, &tt.secret)
			assert.NoError(t, err)
			r := StardogRoleReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			for _, addExpectation := range tt.expectations {
				addExpectation(*tt.srr.reconciliationContext.stardogClient)
			}
			err = r.syncRole(&tt.srr)

			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_validateSpecification(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceRef := "instance-test"
	stardogRoleName := "role-test"
	action1 := "READ"
	resourceType1 := "DB"
	resources1 := []string{"Database1", "Database2"}
	permissionSpec1 := v1alpha1.StardogPermissionSpec{
		Action:       action1,
		ResourceType: resourceType1,
		Resources:    resources1,
	}
	permissionSpec2 := v1alpha1.StardogPermissionSpec{
		ResourceType: resourceType1,
		Resources:    resources1,
	}
	permissionSpec3 := v1alpha1.StardogPermissionSpec{
		Action:    action1,
		Resources: resources1,
	}
	permissionSpec4 := v1alpha1.StardogPermissionSpec{
		Action:       action1,
		ResourceType: resourceType1,
	}

	tests := []struct {
		name        string
		stardogRole v1alpha1.StardogRole
		err         error
	}{
		{
			name:        "GivenReconciliationContext_WhenStardogRoleIsValid_ThenReturnNoError",
			stardogRole: *createStardogRole(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			err:         nil,
		},
		{
			name:        "GivenReconciliationContext_WhenStardogInstanceIsMissing_ThenRaiseError",
			stardogRole: *createStardogRole(namespace, stardogRoleName, "", []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			err:         errors.New(".spec.StardogInstanceRef is required"),
		},
		{
			name:        "GivenReconciliationContext_WhenActionIsMissing_ThenRaiseError",
			stardogRole: *createStardogRole(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec2}),
			err:         errors.New(".spec.Permissions[0].Action is required"),
		},
		{
			name:        "GivenReconciliationContext_WhenResourceTypeIsMissing_ThenRaiseError",
			stardogRole: *createStardogRole(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec3}),
			err:         errors.New(".spec.Permissions[0].ResourceType is required"),
		},
		{
			name:        "GivenReconciliationContext_WhenResourcesIsMissing_ThenRaiseError",
			stardogRole: *createStardogRole(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec4}),
			err:         errors.New(".spec.Permissions[0].Resources at least one resource is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := StardogRoleReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            nil,
			}

			err := r.validateSpecification(&tt.stardogRole)

			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_deleteStardogRole(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceRef := "instance-test"
	stardogRoleName := "role-test"
	action1 := "READ"
	resourceType1 := "DB"
	secretName := "secret-test"
	username := "admin"
	password := "1234"
	serverURL := "http://server:8080"
	resources1 := []string{"Database1", "Database2"}
	permissionSpec1 := v1alpha1.StardogPermissionSpec{
		Action:       action1,
		ResourceType: resourceType1,
		Resources:    resources1,
	}
	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogMocked := stardogmock.NewMockStardogTestClient(mockCtrl)
	stardogClient := createStardogClientFromMock(stardogMocked)

	tests := []struct {
		name               string
		stardogRole        v1alpha1.StardogRole
		stardogInstance    v1alpha1.StardogInstance
		secret             v1.Secret
		srr                StardogRoleReconciliation
		condition          func(stardog_client.Stardog)
		expectedFinalizers []string
		err                error
	}{
		{
			name:            "GivenReconciliationContext_WhenStardogRoleCanBeDeleted_ThenDeleteIt",
			stardogRole:     *createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			},
			condition: func(stardog_client.Stardog) {
				stardogMocked.EXPECT().
					ListRoleUsers(gomock.Any(), gomock.Any()).
					Return(&users_roles.ListRoleUsersOK{Payload: &models.Users{Users: []string{}}}, nil)
			},
			expectedFinalizers: nil,
			err:                nil,
		},
		{
			name:            "GivenReconciliationContext_WhenThereAreUsersWithRole_ThenRaiseError",
			stardogRole:     *createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			},
			condition: func(stardog_client.Stardog) {
				stardogMocked.EXPECT().
					ListRoleUsers(gomock.Any(), gomock.Any()).
					Return(&users_roles.ListRoleUsersOK{Payload: &models.Users{Users: []string{"user1"}}}, nil)
			},
			expectedFinalizers: []string{roleFinalizer},
			err:                errors.New("cannot delete role role-test as it is used by user1 users in namespace-test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.stardogRole, &tt.stardogInstance, &tt.secret)
			assert.NoError(t, err)
			r := StardogRoleReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			tt.condition(*stardogClient)
			stardogMocked.EXPECT().RemoveRole(gomock.Any(), gomock.Any()).AnyTimes()
			stardogMocked.EXPECT().SetTransport(gomock.Any()).AnyTimes()

			err = fakeKubeClient.Get(context.Background(), types.NamespacedName{
				Namespace: namespace,
				Name:      stardogRoleName,
			}, tt.srr.resource)
			assert.NoError(t, err)

			err = r.deleteStardogRole(&tt.srr)
			actualRole := v1alpha1.StardogRole{}
			_ = fakeKubeClient.Get(context.Background(), types.NamespacedName{
				Namespace: namespace,
				Name:      stardogRoleName,
			}, &actualRole)

			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expectedFinalizers, actualRole.GetFinalizers())
		})
	}
}

func Test_finalizeRole(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceRef := "instance-test"
	stardogRoleName := "role-test"
	action1 := "READ"
	resourceType1 := "DB"
	secretName := "secret-test"
	username := "admin"
	password := "1234"
	serverURL := "http://server:8080"
	resources1 := []string{"Database1", "Database2"}
	permissionSpec1 := v1alpha1.StardogPermissionSpec{
		Action:       action1,
		ResourceType: resourceType1,
		Resources:    resources1,
	}
	ctx := context.Background()
	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogMocked := stardogmock.NewMockStardogTestClient(mockCtrl)
	stardogClient := createStardogClientFromMock(stardogMocked)

	tests := []struct {
		name            string
		stardogRole     v1alpha1.StardogRole
		stardogInstance v1alpha1.StardogInstance
		secret          v1.Secret
		srr             StardogRoleReconciliation
		conditions      []func(stardog_client.Stardog)
		err             error
	}{
		{
			name:            "GivenReconciliationContext_WhenThereAreUsersWithRole_ThenRaiseError",
			stardogRole:     *createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       ctx,
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			},
			conditions: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						ListRoleUsers(gomock.Any(), gomock.Any()).
						Return(&users_roles.ListRoleUsersOK{Payload: &models.Users{Users: []string{}}}, nil).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						RemoveRole(roles.NewRemoveRoleParams().WithRole(stardogRoleName).WithForce(pointer.Bool(false)), gomock.Any()).
						Times(1)
				},
			},
			err: nil,
		},
		{
			name:            "GivenReconciliationContext_WhenThereAreUsersWithRole_ThenRaiseError",
			stardogRole:     *createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       ctx,
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			},
			conditions: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						ListRoleUsers(users_roles.NewListRoleUsersParams().WithRole(stardogRoleName), gomock.Any()).
						Return(&users_roles.ListRoleUsersOK{Payload: &models.Users{Users: []string{"user1"}}}, nil).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						RemoveRole(gomock.Any(), gomock.Any()).
						Times(0)
				},
			},
			err: errors.New("cannot delete role role-test as it is used by user1 users in namespace-test"),
		},
		{
			name:            "GivenReconciliationContext_WhenCannotUpdateRole_ThenRaiseError",
			stardogRole:     *createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       ctx,
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
			},
			conditions: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						ListRoleUsers(gomock.Any(), gomock.Any()).
						Return(&users_roles.ListRoleUsersOK{Payload: &models.Users{Users: []string{}}}, nil).
						Times(1)
				},
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						RemoveRole(roles.NewRemoveRoleParams().WithRole(stardogRoleName).WithForce(pointer.Bool(false)), gomock.Any()).
						Return(&roles.RemoveRoleNoContent{}, errors.New("cannot update role")).
						Times(1)
				},
			},
			err: errors.New("cannot remove Stardog Role namespace-test/role-test: cannot update role"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.stardogRole, &tt.stardogInstance, &tt.secret)
			assert.NoError(t, err)
			r := StardogRoleReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			for _, addCondition := range tt.conditions {
				addCondition(*stardogClient)
			}
			stardogMocked.EXPECT().SetTransport(gomock.Any()).AnyTimes()

			err = r.finalize(&tt.srr)

			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_ReconcileRole(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	instance := "instance"

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{Name: stardogInstanceName, Namespace: namespace},
	}

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		namespace      v1.Namespace
		stardogRole    v1alpha1.StardogRole
		srr            StardogRoleReconciliation
		expectedResult ctrl.Result
	}{
		{
			name:        "GivenReconciliation_WhenStardogRoleMissing_ThenReturnNoRequeue",
			namespace:   *createNamespace(namespace),
			stardogRole: *createStardogRole(namespace, "nonExistingRole", instance, []v1alpha1.StardogPermissionSpec{}),
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.namespace, &tt.stardogRole)
			assert.NoError(t, err)
			r := StardogRoleReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}

			result, err := r.Reconcile(context.Background(), req)

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func Test_ReconcileStardogRole(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	secretName := "secret-test"
	stardogRoleName := "stardog-role"
	stardogUserName := "stardog-user"
	serverURL := "http://localhost:5820/"
	username := "admin"
	password := "1234"

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	os.Setenv("RECONCILIATION_FREQUENCY_ON_ERROR", "1s")
	os.Setenv("RECONCILIATION_FREQUENCY", "1m")
	InitEnv()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogMocked := stardogmock.NewMockStardogTestClient(mockCtrl)
	stardogClient := createStardogClientFromMock(stardogMocked)

	tests := []struct {
		name            string
		namespace       v1.Namespace
		stardogInstance v1alpha1.StardogInstance
		stardogRole     v1alpha1.StardogRole
		stardogUser     v1alpha1.StardogUser
		srr             StardogRoleReconciliation
		secret          v1.Secret
		expectations    []func(stardog_client.Stardog)
		expectedResult  ctrl.Result
	}{
		{
			name:            "GivenReconciliation_WhenStardogRoleIsDeleted_ThenReturnNoRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createDeletedStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			},
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						ListRoleUsers(gomock.Any(), gomock.Any()).
						Return(&users_roles.ListRoleUsersOK{Payload: &models.Users{Users: []string{}}}, nil)
				},
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						RemoveRole(gomock.Any(), gomock.Any())
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogRoleIsDeletedWithFinalizers_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createDeletedStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createDeletedStardogRoleWithFinalizers(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			},
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						ListRoleUsers(gomock.Any(), gomock.Any()).
						Return(&users_roles.ListRoleUsersOK{Payload: &models.Users{Users: []string{}}}, nil)
				},
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						RemoveRole(gomock.Any(), gomock.Any())
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogRoleIsDeletedWithFinalizersAndUser_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createDeletedStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createDeletedStardogRoleWithFinalizers(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			stardogUser:     *createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{stardogRoleName}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			},
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						ListRoleUsers(gomock.Any(), gomock.Any()).
						Return(&users_roles.ListRoleUsersOK{Payload: &models.Users{Users: []string{stardogUserName}}}, nil)
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      true,
				RequeueAfter: ReconFreqErr,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogRoleIsInvalid_ThenReturnNoRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createStardogRole(namespace, stardogRoleName, "", []v1alpha1.StardogPermissionSpec{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			},
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogRoleNoSync_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			},
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						ListRoles(gomock.Any(), gomock.Any()).
						Return(&roles.ListRolesOK{Payload: &models.Roles{}}, errors.New("cannot list roles"))
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      true,
				RequeueAfter: ReconFreqErr,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogRoleCannotBeUpdated_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createStardogRoleWithWrongVersion(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			},
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						ListRoles(gomock.Any(), gomock.Any()).
						Return(&roles.ListRolesOK{Payload: &models.Roles{}}, nil)
				},
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						CreateRole(gomock.Any(), gomock.Any())
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      true,
				RequeueAfter: ReconFreqErr,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogRoleIsReconciled_ThenReturnNoRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			srr: StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			},
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						ListRoles(gomock.Any(), gomock.Any()).
						Return(&roles.ListRolesOK{Payload: &models.Roles{}}, nil)
				},
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						CreateRole(gomock.Any(), gomock.Any())
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      true,
				RequeueAfter: ReconFreq,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.namespace, &tt.stardogInstance, &tt.secret, &tt.stardogRole, &tt.stardogUser)
			assert.NoError(t, err)
			r := StardogRoleReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			srr := &StardogRoleReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(map[v1alpha1.StardogConditionType]v1alpha1.StardogCondition),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: &tt.stardogRole,
			}
			for _, addExpectation := range tt.expectations {
				addExpectation(*tt.srr.reconciliationContext.stardogClient)
			}
			stardogMocked.EXPECT().SetTransport(gomock.Any()).AnyTimes()

			result, err := r.ReconcileStardogRole(srr)

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func createStardogRole(namespace, stardogRoleName, stardogInstanceRef string, permissions []v1alpha1.StardogPermissionSpec) *v1alpha1.StardogRole {
	return &v1alpha1.StardogRole{
		TypeMeta:   metav1.TypeMeta{Kind: "StardogRole", APIVersion: "v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{Name: stardogRoleName, Namespace: namespace},
		Spec: v1alpha1.StardogRoleSpec{
			RoleName:           stardogRoleName,
			StardogInstanceRef: stardogInstanceRef,
			Permissions:        permissions,
		},
	}
}

func createStardogRoleWithWrongVersion(namespace, stardogRoleName, stardogInstanceRef string, permissions []v1alpha1.StardogPermissionSpec) *v1alpha1.StardogRole {
	stardogRole := createStardogRole(namespace, stardogRoleName, stardogInstanceRef, permissions)
	stardogRole.SetResourceVersion("wrongVersion")
	return stardogRole
}

func createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef string, permissions []v1alpha1.StardogPermissionSpec) *v1alpha1.StardogRole {
	role := createStardogRole(namespace, stardogRoleName, stardogInstanceRef, permissions)
	role.SetFinalizers([]string{roleFinalizer})
	return role
}

func createDeletedStardogRole(namespace, stardogRoleName, stardogInstanceRef string, permissions []v1alpha1.StardogPermissionSpec) *v1alpha1.StardogRole {
	stardogRole := createStardogRole(namespace, stardogRoleName, stardogInstanceRef, permissions)
	newTime := metav1.NewTime(time.Now())
	stardogRole.SetDeletionTimestamp(&newTime)
	stardogRole.SetFinalizers([]string{"finalizer"})
	return stardogRole
}

func createDeletedStardogRoleWithFinalizers(namespace, stardogRoleName, stardogInstanceRef string, permissions []v1alpha1.StardogPermissionSpec) *v1alpha1.StardogRole {
	stardogRole := createStardogRole(namespace, stardogRoleName, stardogInstanceRef, permissions)
	newTime := metav1.NewTime(time.Now())
	stardogRole.SetDeletionTimestamp(&newTime)
	stardogRole.SetFinalizers([]string{instanceRoleFinalizer})
	return stardogRole
}
