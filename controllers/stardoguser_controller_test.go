package controllers

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest"
	testr "github.com/go-logr/logr/testr"
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
	ctrl "sigs.k8s.io/controller-runtime"
)

func Test_deleteStardogUser(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceRef := "instance-test"
	stardogUserName := "user-test"
	secretNameAdmin := "secret-test"
	usernameAdmin := "admin"
	passwordAdmin := "1234"
	secretNameUser := "user-secret-test"
	usernameUser := "user"
	passwordUser := "1234"
	serverURL := "server"
	roles := []string{"role1", "role2"}
	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name               string
		stardogUser        v1alpha1.StardogUser
		stardogInstance    v1alpha1.StardogInstance
		secretAdmin        v1.Secret
		secretUser         v1.Secret
		sur                StardogUserReconciliation
		condition          func(stardogrestapi2.ExtendedBaseClientAPI)
		expectedFinalizers []string
		err                error
	}{
		{
			name:            "GivenReconciliationContext_WhenStardogUserCanBeDeleted_ThenDeleteIt",
			stardogUser:     *createStardogUserWithFinalizer(namespace, stardogUserName, stardogInstanceRef, secretNameUser, roles),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretNameAdmin, serverURL),
			secretAdmin:     *createFullSecret(namespace, secretNameAdmin, usernameAdmin, passwordAdmin),
			secretUser:      *createFullSecret(namespace, secretNameUser, usernameUser, passwordUser),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUserWithFinalizer(namespace, stardogUserName, stardogInstanceRef, secretNameAdmin, roles),
			},
			condition: func(stardogrestapi2.ExtendedBaseClientAPI) {
				stardogClient.EXPECT().
					RemoveUser(gomock.Any(), gomock.Any())
			},
			expectedFinalizers: nil,
			err:                nil,
		},
		{
			name:            "GivenReconciliationContext_WhenStardogUserCannotSetConnection_ThenRaiseError",
			stardogUser:     *createStardogUserWithFinalizer(namespace, stardogUserName, stardogInstanceRef, secretNameUser, roles),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretNameAdmin, serverURL),
			secretAdmin:     *createFullSecret(namespace, secretNameAdmin, usernameAdmin, passwordAdmin),
			secretUser:      *createFullSecret(namespace, secretNameUser, usernameUser, passwordUser),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUserWithFinalizer(namespace, stardogUserName, stardogInstanceRef, secretNameUser, roles),
			},
			condition: func(stardogrestapi2.ExtendedBaseClientAPI) {
				stardogClient.EXPECT().
					RemoveUser(gomock.Any(), gomock.Any()).
					Return(autorest.Response{}, errors.New("cannot remove user"))
			},
			expectedFinalizers: []string{userFinalizer},
			err:                errors.New("cannot remove Stardog user namespace-test/dXNlcg==: cannot remove user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.stardogUser, &tt.stardogInstance, &tt.secretAdmin, &tt.secretUser)
			assert.NoError(t, err)
			r := StardogUserReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			tt.condition(stardogClient)
			stardogClient.EXPECT().SetConnection(gomock.Any(), gomock.Any(), gomock.Any())
			err = r.deleteStardogUser(&tt.sur)
			actualUser := v1alpha1.StardogUser{}
			_ = fakeKubeClient.Get(context.Background(), types.NamespacedName{
				Namespace: namespace,
				Name:      stardogUserName,
			}, &actualUser)

			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expectedFinalizers, actualUser.GetFinalizers())
		})
	}
}

func Test_validateSpecificationUser(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceRef := "instance-test"
	stardogUser := "user-test"
	secretNameUser := "user-secret-test"
	roles := []string{"role1", "role2"}

	tests := []struct {
		name        string
		stardogUser v1alpha1.StardogUser
		err         error
	}{
		{
			name:        "GivenReconciliationContext_WhenStardogUserIsValid_ThenReturnNoError",
			stardogUser: *createStardogUser(namespace, stardogUser, stardogInstanceRef, secretNameUser, roles),
			err:         nil,
		},
		{
			name:        "GivenReconciliationContext_WhenStardogUserIsValid_ThenReturnNoError",
			stardogUser: *createStardogUser(namespace, stardogUser, stardogInstanceRef, secretNameUser, []string{}),
			err:         nil,
		},
		{
			name:        "GivenReconciliationContext_WhenStardogUserIsInvalid_ThenRaiseError",
			stardogUser: *createStardogUser(namespace, stardogUser, "", secretNameUser, roles),
			err:         errors.New(".spec.StardogInstanceRef is required"),
		},
		{
			name:        "GivenReconciliationContext_WhenStardogUserIsInvalid_ThenRaiseError",
			stardogUser: *createStardogUser(namespace, stardogUser, stardogInstanceRef, "", roles),
			err:         errors.New(".spec.Credentials.SecretRef is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := StardogUserReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            nil,
			}

			err := r.validateSpecification(&tt.stardogUser.Spec)

			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_syncUser(t *testing.T) {
	namespace := "namespace-test"
	stardogInstanceRef := "instance-test"
	stardogUser := "user-test"
	secretNameAdmin := "secret-test"
	usernameAdmin := "admin"
	passwordAdmin := "1234"
	secretNameUser := "user-secret-test"
	usernameUser := "user"
	passwordUser := "1234"
	serverURL := "https://stardog-test.com"
	role1 := "roleA"
	role2 := "roleB"
	role3 := "roleC"
	roles1 := []string{role1, role2}
	roles2 := []string{role3, role2}
	ctx := context.Background()

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)
	encodedUser := base64.StdEncoding.EncodeToString([]byte(usernameUser))

	tests := []struct {
		name            string
		stardogUser     v1alpha1.StardogUser
		stardogInstance v1alpha1.StardogInstance
		secretAdmin     v1.Secret
		secretUser      v1.Secret
		sur             StardogUserReconciliation
		expectations    []func(stardogrestapi2.ExtendedBaseClientAPI)
		err             error
	}{
		{
			name:            "GivenReconciliationContext_WhenCreateUser1_ThenUpdateStardogDB",
			stardogUser:     *createStardogUser(namespace, stardogUser, stardogInstanceRef, secretNameUser, roles1),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretNameAdmin, serverURL),
			secretAdmin:     *createFullSecret(namespace, secretNameAdmin, usernameAdmin, passwordAdmin),
			secretUser:      *createFullSecret(namespace, secretNameUser, usernameUser, passwordUser),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUser, stardogInstanceRef, secretNameUser, roles1),
			},
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						SetConnection(
							serverURL,
							base64.StdEncoding.EncodeToString([]byte(usernameAdmin)),
							base64.StdEncoding.EncodeToString([]byte(passwordAdmin))).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListUsers(ctx).
						Return(stardogrest.Users{Users: &[]string{}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						CreateUser(ctx, gomock.Any()).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListUserRoles(ctx, encodedUser).
						Return(stardogrest.Roles{Roles: &[]string{}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						AddRole(ctx, encodedUser, stardogrest.Rolename{Rolename: &role1}).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						AddRole(ctx, encodedUser, stardogrest.Rolename{Rolename: &role2}).
						Times(1)
				},
			},
			err: nil,
		},
		{
			name:            "GivenReconciliationContext_WhenCreateUser2_ThenUpdateStardogDB",
			stardogUser:     *createStardogUser(namespace, stardogUser, stardogInstanceRef, secretNameUser, roles1),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretNameAdmin, serverURL),
			secretAdmin:     *createFullSecret(namespace, secretNameAdmin, usernameAdmin, passwordAdmin),
			secretUser:      *createFullSecret(namespace, secretNameUser, usernameUser, passwordUser),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUser, stardogInstanceRef, secretNameUser, roles1),
			},
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						SetConnection(
							serverURL,
							base64.StdEncoding.EncodeToString([]byte(usernameAdmin)),
							base64.StdEncoding.EncodeToString([]byte(passwordAdmin))).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListUsers(ctx).
						Return(stardogrest.Users{Users: &[]string{"random-user"}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						CreateUser(ctx, gomock.Any()).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListUserRoles(ctx, encodedUser).
						Return(stardogrest.Roles{Roles: &[]string{}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						AddRole(ctx, encodedUser, stardogrest.Rolename{Rolename: &role1}).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						AddRole(ctx, encodedUser, stardogrest.Rolename{Rolename: &role2}).
						Times(1)
				},
			},
			err: nil,
		},
		{
			name:            "GivenReconciliationContext_WhenUpdateUser_ThenUpdateStardogDB",
			stardogUser:     *createStardogUser(namespace, stardogUser, stardogInstanceRef, secretNameUser, roles1),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceRef, secretNameAdmin, serverURL),
			secretAdmin:     *createFullSecret(namespace, secretNameAdmin, usernameAdmin, passwordAdmin),
			secretUser:      *createFullSecret(namespace, secretNameUser, usernameUser, passwordUser),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUser, stardogInstanceRef, secretNameUser, roles2),
			},
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						SetConnection(
							serverURL,
							base64.StdEncoding.EncodeToString([]byte(usernameAdmin)),
							base64.StdEncoding.EncodeToString([]byte(passwordAdmin))).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListUsers(ctx).
						Return(stardogrest.Users{Users: &[]string{encodedUser}}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ChangePassword(ctx, encodedUser, gomock.Any()).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						ListUserRoles(ctx, encodedUser).
						Return(stardogrest.Roles{Roles: &roles1}, nil).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						AddRole(ctx, encodedUser, stardogrest.Rolename{Rolename: &role3}).
						Times(1)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.
						EXPECT().
						RemoveRole(ctx, encodedUser, role1).
						Times(1)
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.stardogInstance, &tt.secretAdmin, &tt.secretUser, &tt.stardogUser)
			assert.NoError(t, err)
			r := StardogUserReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			for _, addExpectation := range tt.expectations {
				addExpectation(tt.sur.reconciliationContext.stardogClient)
			}
			err = r.syncUser(&tt.sur)

			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_ReconcileUser(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	instance := "instance"
	secretName := "secret-name"

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{Name: stardogInstanceName, Namespace: namespace},
	}

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		namespace      v1.Namespace
		stardogUser    v1alpha1.StardogUser
		sur            StardogUserReconciliation
		expectedResult ctrl.Result
	}{
		{
			name:        "GivenReconciliation_WhenStardogUserMissing_ThenReturnNoRequeue",
			namespace:   *createNamespace(namespace),
			stardogUser: *createStardogUser(namespace, "nonExistingUser", instance, secretName, []string{}),
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.namespace, &tt.stardogUser)
			assert.NoError(t, err)
			r := StardogUserReconciler{
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

func Test_ReconcileStardogUser(t *testing.T) {

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
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name            string
		namespace       v1.Namespace
		stardogInstance v1alpha1.StardogInstance
		stardogRole     v1alpha1.StardogRole
		stardogUser     v1alpha1.StardogUser
		sur             StardogUserReconciliation
		secret          v1.Secret
		expectations    []func(stardogrestapi2.ExtendedBaseClientAPI)
		expectedResult  ctrl.Result
	}{
		{
			name:            "GivenReconciliation_WhenStardogUserIsDeleted_ThenReturnNoRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			stardogUser:     *createDeletedStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			},
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						SetConnection(gomock.Any(), gomock.Any(), gomock.Any())
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						RemoveUser(gomock.Any(), gomock.Any())
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogUserIsDeletedWithFinalizers_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createDeletedStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
			stardogUser:     *createDeletedStardogUserWithFinalizers(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			stardogRole:     *createDeletedStardogRoleWithFinalizers(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			},
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						SetConnection(gomock.Any(), gomock.Any(), gomock.Any())
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						RemoveUser(gomock.Any(), gomock.Any())
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogUserIsDeletedWithFinalizersAndUser_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createDeletedStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			stardogUser:     *createDeletedStardogUserWithFinalizers(namespace, stardogUserName, stardogInstanceName, secretName, []string{stardogRoleName}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			},
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						SetConnection(gomock.Any(), gomock.Any(), gomock.Any())
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						RemoveUser(gomock.Any(), gomock.Any()).
						Return(autorest.Response{}, errors.New("cannot delete user"))
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      true,
				RequeueAfter: ReconFreqErr,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogUserIsInvalid_ThenReturnNoRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogUser:     *createStardogUser(namespace, stardogUserName, "", secretName, []string{}),
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogRoleName, []v1alpha1.StardogPermissionSpec{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			},
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogUserNoSync_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			stardogUser:     *createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			},
			expectedResult: ctrl.Result{
				Requeue:      true,
				RequeueAfter: ReconFreqErr,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogUserCannotBeUpdated_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
			stardogUser:     *createStardogUserWithWrongVersion(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			},
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						SetConnection(gomock.Any(), gomock.Any(), gomock.Any())
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						ListUsers(gomock.Any()).
						Return(stardogrest.Users{Users: &[]string{}}, nil)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						CreateUser(gomock.Any(), gomock.Any())
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						ListUserRoles(gomock.Any(), gomock.Any()).
						Return(stardogrest.Roles{Roles: &[]string{}}, nil)
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
			stardogUser:     *createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			secret:          *createFullSecret(namespace, secretName, username, password),
			sur: StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogUser(namespace, stardogUserName, stardogInstanceName, secretName, []string{}),
			},
			expectations: []func(stardogrestapi2.ExtendedBaseClientAPI){
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						SetConnection(gomock.Any(), gomock.Any(), gomock.Any())
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						ListUsers(gomock.Any()).
						Return(stardogrest.Users{Users: &[]string{}}, nil)
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						CreateUser(gomock.Any(), gomock.Any())
				},
				func(stardogrestapi2.ExtendedBaseClientAPI) {
					stardogClient.EXPECT().
						ListUserRoles(gomock.Any(), gomock.Any()).
						Return(stardogrest.Roles{Roles: &[]string{}}, nil)
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
			r := StardogUserReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			sur := &StardogUserReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(map[v1alpha1.StardogConditionType]v1alpha1.StardogCondition),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: &tt.stardogUser,
			}
			for _, addExpectation := range tt.expectations {
				addExpectation(tt.sur.reconciliationContext.stardogClient)
			}

			result, err := r.ReconcileStardogUser(sur)

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func createStardogUser(namespace, stardogUserName, stardogInstanceRef, secretRef string, roles []string) *v1alpha1.StardogUser {
	return &v1alpha1.StardogUser{
		TypeMeta:   metav1.TypeMeta{Kind: "StardogUser", APIVersion: "v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{Name: stardogUserName, Namespace: namespace},
		Spec: v1alpha1.StardogUserSpec{
			StardogInstanceRef: stardogInstanceRef,
			Credentials: v1alpha1.StardogUserCredentialsSpec{
				Namespace: namespace,
				SecretRef: secretRef,
			},
			Roles: roles,
		},
	}
}

func createStardogUserWithFinalizer(namespace, stardogUserName, stardogInstanceRef, secretRef string, roles []string) *v1alpha1.StardogUser {
	user := createStardogUser(namespace, stardogUserName, stardogInstanceRef, secretRef, roles)
	user.SetFinalizers([]string{userFinalizer})
	return user
}

func createDeletedStardogUser(namespace, stardogUserName, stardogInstanceRef, secretRef string, roles []string) *v1alpha1.StardogUser {
	user := createStardogUser(namespace, stardogUserName, stardogInstanceRef, secretRef, roles)
	newTime := metav1.NewTime(time.Now())
	user.SetDeletionTimestamp(&newTime)
	return user
}

func createDeletedStardogUserWithFinalizers(namespace, stardogUserName, stardogInstanceRef, secretRef string, roles []string) *v1alpha1.StardogUser {
	user := createStardogUser(namespace, stardogUserName, stardogInstanceRef, secretRef, roles)
	newTime := metav1.NewTime(time.Now())
	user.SetDeletionTimestamp(&newTime)
	user.SetFinalizers([]string{instanceUserFinalizer})
	return user
}

func createStardogUserWithWrongVersion(namespace, stardogUserName, stardogInstanceRef, secretRef string, roles []string) *v1alpha1.StardogUser {
	stardogUser := createStardogUser(namespace, stardogUserName, stardogInstanceRef, secretRef, roles)
	stardogUser.SetResourceVersion("wrongVersion")
	return stardogUser
}
