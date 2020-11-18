package controllers

import (
	"context"
	"errors"
	"github.com/Azure/go-autorest/autorest"
	testing2 "github.com/go-logr/logr/testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	"github.com/vshn/stardog-userrole-operator/stardogrest"
	stardogrestapi "github.com/vshn/stardog-userrole-operator/stardogrest/mocks"
	stardogrestapi2 "github.com/vshn/stardog-userrole-operator/stardogrest/stardogrestapi"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"testing"
	"time"
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
	permission1 := stardogrest.Permission{
		Action:       &action1,
		ResourceType: &resourceType1,
		Resource:     &resources1,
	}
	permission2 := stardogrest.Permission{
		Action:       &action2,
		ResourceType: &resourceType2,
		Resource:     &resources2,
	}
	permission3 := stardogrest.Permission{
		Action:       &action1,
		ResourceType: &resourceType2,
		Resource:     &resources2,
	}
	ctx := context.Background()

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name            string
		stardogRole     v1alpha1.StardogRole
		stardogInstance v1alpha1.StardogInstance
		secret          v1.Secret
		srr             StardogRoleReconciliation
		expectations    []func(stardogrestapi2.ExtendedBaseClientAPI)
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
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						SetConnection(serverURL, username, password).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListRoles(ctx).
						Return(stardogrest.Roles{Roles: &[]string{}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						CreateRole(ctx, gomock.Eq(&stardogrest.Rolename{Rolename: &stardogRoleName})).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListRolePermissions(ctx, stardogRoleName).
						Return(stardogrest.Permissions{Permissions: &[]stardogrest.Permission{permission1, permission2}}, nil).
						Times(0)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						AddRolePermission(ctx, stardogRoleName, gomock.Eq(permission1)).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						AddRolePermission(ctx, stardogRoleName, gomock.Eq(permission2)).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						RemoveRolePermission(ctx, stardogRoleName, gomock.Any()).
						Times(0)
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
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						SetConnection(serverURL, username, password).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListRoles(ctx).
						Return(stardogrest.Roles{Roles: &[]string{stardogRoleName}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						CreateRole(ctx, gomock.Any()).
						Times(0)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListRolePermissions(ctx, stardogRoleName).
						Return(stardogrest.Permissions{Permissions: &[]stardogrest.Permission{permission2, permission3}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						AddRolePermission(ctx, stardogRoleName, gomock.Eq(permission1)).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						RemoveRolePermission(ctx, stardogRoleName, gomock.Eq(permission3)).
						Times(1)
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
				Log:               testing2.TestLogger{},
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			for _, addExpectation := range tt.expectations {
				addExpectation(tt.srr.reconciliationContext.stardogClient)
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
				Log:               testing2.TestLogger{},
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
	serverURL := "server"
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
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name               string
		stardogRole        v1alpha1.StardogRole
		stardogInstance    v1alpha1.StardogInstance
		secret             v1.Secret
		srr                StardogRoleReconciliation
		condition          func(stardogrestapi2.ExtendedBaseClientAPI)
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
			condition: func(stardogrestapi2.ExtendedBaseClientAPI) {
				stardogClient.EXPECT().
					ListRoleUsers(gomock.Any(), gomock.Any()).
					Return(stardogrest.Users{Users: &[]string{}}, nil)
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
			condition: func(stardogrestapi2.ExtendedBaseClientAPI) {
				stardogClient.EXPECT().
					ListRoleUsers(gomock.Any(), gomock.Any()).
					Return(stardogrest.Users{Users: &[]string{"user1"}}, nil)
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
				Log:               testing2.TestLogger{},
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			tt.condition(stardogClient)
			stardogClient.EXPECT().SetConnection(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			stardogClient.EXPECT().RemoveRole1(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
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
	serverURL := "server"
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
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name            string
		stardogRole     v1alpha1.StardogRole
		stardogInstance v1alpha1.StardogInstance
		secret          v1.Secret
		srr             StardogRoleReconciliation
		conditions      []func(stardogrestapi2.ExtendedBaseClientAPI)
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
			conditions: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						SetConnection(serverURL, username, password).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						ListRoleUsers(ctx, stardogRoleName).
						Return(stardogrest.Users{Users: &[]string{}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						RemoveRole1(ctx, stardogRoleName, &[]bool{false}[0]).
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
			conditions: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						SetConnection(serverURL, username, password).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						ListRoleUsers(ctx, stardogRoleName).
						Return(stardogrest.Users{Users: &[]string{"user1"}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						RemoveRole1(ctx, stardogRoleName, &[]bool{false}[0]).
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
			conditions: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						SetConnection(serverURL, username, password).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						ListRoleUsers(ctx, stardogRoleName).
						Return(stardogrest.Users{Users: &[]string{}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						RemoveRole1(ctx, stardogRoleName, &[]bool{false}[0]).
						Return(autorest.Response{}, errors.New("cannot update role")).
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
				Log:               testing2.TestLogger{},
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			for _, addCondition := range tt.conditions {
				addCondition(stardogClient)
			}

			err = r.finalize(&tt.srr, roleFinalizer)

			assert.Equal(t, tt.err, err)
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

func createStardogRoleWithFinalizer(namespace, stardogRoleName, stardogInstanceRef string, permissions []v1alpha1.StardogPermissionSpec) *v1alpha1.StardogRole {
	role := createStardogRole(namespace, stardogRoleName, stardogInstanceRef, permissions)
	role.SetFinalizers([]string{roleFinalizer})
	return role
}
