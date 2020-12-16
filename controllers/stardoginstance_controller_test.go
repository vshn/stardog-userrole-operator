package controllers

import (
	"context"
	"encoding/base64"
	"errors"
	testing2 "github.com/go-logr/logr/testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	"github.com/vshn/stardog-userrole-operator/stardogrest"
	stardogrestapi "github.com/vshn/stardog-userrole-operator/stardogrest/mocks"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"testing"
	"time"
)

func Test_deleteStardogInstance(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	secretName := "secret-test"
	serverURL := "url-test"
	username := "admin"
	password := "1234"

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name               string
		stardogInstance    v1alpha1.StardogInstance
		sir                StardogInstanceReconciliation
		secret             v1.Secret
		expectedFinalizers []string
		err                error
	}{
		{
			name:            "GivenReconciliationContext_WhenFinalizersArePresent_ThenDeleteStardogInstance",
			stardogInstance: *createStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			sir: StardogInstanceReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
			},
			expectedFinalizers: []string{},
			err:                nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.stardogInstance, &tt.secret)
			assert.NoError(t, err)
			r := StardogInstanceReconciler{
				Log:               testing2.TestLogger{},
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			stardogClient.EXPECT().SetConnection(
				serverURL,
				base64.StdEncoding.EncodeToString([]byte(username)),
				base64.StdEncoding.EncodeToString([]byte(password)))
			err = r.deleteStardogInstance(&tt.sir)

			actualInstance := v1alpha1.StardogInstance{}
			_ = fakeKubeClient.Get(context.Background(), types.NamespacedName{
				Namespace: namespace,
				Name:      stardogInstanceName,
			}, &actualInstance)

			assert.Equal(t, tt.err, err)
			assert.ElementsMatch(t, actualInstance.GetFinalizers(), tt.expectedFinalizers)
		})
	}
}

func Test_userFinalize(t *testing.T) {

	stardogInstanceRef := "instance-test"
	stardogUser := "user-test"
	stardogUser_2 := "user-test-2"
	secretNameUser := "user-secret-test"
	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	secretName := "secret-test"
	serverURL := "url-test"
	ctx := context.Background()
	roles := []string{"role1", "role2"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name     string
		userList v1alpha1.StardogUserList
		sir      StardogInstanceReconciliation
		err      error
	}{
		{
			name: "GivenUserFinalizer_WhenThereAreNoUsers_ThenDoNotReturnAnyError",
			userList: v1alpha1.StardogUserList{
				Items: []v1alpha1.StardogUser{},
			},
			sir: StardogInstanceReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       ctx,
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
			},
			err: nil,
		},
		{
			name: "GivenUserFinalizer_WhenThereAreUsers_ThenRaiseError",
			userList: v1alpha1.StardogUserList{
				Items: []v1alpha1.StardogUser{
					*createStardogUser(namespace, stardogUser, stardogInstanceRef, secretNameUser, roles),
					*createStardogUser(namespace, stardogUser_2, stardogInstanceRef, secretNameUser, roles),
				},
			},
			sir: StardogInstanceReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       ctx,
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
			},
			err: errors.New("cannot delete StardogInstance, found [  user-test user-test-2] user CRDs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.userList)
			assert.NoError(t, err)
			r := StardogInstanceReconciler{
				Log:               testing2.TestLogger{},
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			err = r.userFinalizer(&tt.sir)
			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_roleFinalize(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceRef := "instance-test"
	stardogRoleName := "role-test"
	stardogRoleName2 := "role-test-2"
	secretName := "secret-test"
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

	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name     string
		userList v1alpha1.StardogRoleList
		sir      StardogInstanceReconciliation
		err      error
	}{
		{
			name: "GivenRoleFinalizer_WhenThereAreNoRoles_ThenDoNotReturnAnyError",
			userList: v1alpha1.StardogRoleList{
				Items: []v1alpha1.StardogRole{},
			},
			sir: StardogInstanceReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       ctx,
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogInstanceWithFinalizers(namespace, stardogInstanceRef, secretName, serverURL),
			},
			err: nil,
		},
		{
			name: "GivenRoleFinalizer_WhenThereAreRoles_ThenRaiseError",
			userList: v1alpha1.StardogRoleList{
				Items: []v1alpha1.StardogRole{
					*createStardogRole(namespace, stardogRoleName, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1, permissionSpec2}),
					*createStardogRole(namespace, stardogRoleName2, stardogInstanceRef, []v1alpha1.StardogPermissionSpec{permissionSpec1}),
				},
			},
			sir: StardogInstanceReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       ctx,
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogInstanceWithFinalizers(namespace, stardogInstanceRef, secretName, serverURL),
			},
			err: errors.New("cannot delete StardogInstance, found [  role-test role-test-2] role CRDs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.userList)
			assert.NoError(t, err)
			r := StardogInstanceReconciler{
				Log:               testing2.TestLogger{},
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			err = r.roleFinalizer(&tt.sir)
			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_validateSpec(t *testing.T) {

	namespace := "namespace-test"
	secretName := "secret-test"
	serverURL := "https://test-stardog.ch/"

	tests := []struct {
		name                string
		stardogInstanceSpec v1alpha1.StardogInstanceSpec
		err                 error
	}{
		{
			name: "GivenStardogInstanceSpec_WhenThereIsUrlAndAdminCredentials_ThenDoNotReturnAnyError",
			stardogInstanceSpec: v1alpha1.StardogInstanceSpec{
				ServerUrl: serverURL,
				AdminCredentials: v1alpha1.StardogUserCredentialsSpec{
					SecretRef: secretName,
				},
			},
			err: nil,
		},
		{
			name: "GivenStardogInstanceSpec_WhenThereIsUrlAndNamespaceAndAdminCredentials_ThenDoNotReturnAnyError",
			stardogInstanceSpec: v1alpha1.StardogInstanceSpec{
				ServerUrl: serverURL,
				AdminCredentials: v1alpha1.StardogUserCredentialsSpec{
					Namespace: namespace,
					SecretRef: secretName,
				},
			},
			err: nil,
		},
		{
			name: "GivenStardogInstanceSpec_WhenThereIsNoUrl_ThenRaiseError",
			stardogInstanceSpec: v1alpha1.StardogInstanceSpec{
				AdminCredentials: v1alpha1.StardogUserCredentialsSpec{
					Namespace: namespace,
					SecretRef: secretName,
				},
			},
			err: errors.New(".spec.ServerUrl is required"),
		},
		{
			name: "GivenStardogInstanceSpec_WhenThereIsCredentials_ThenRaiseError",
			stardogInstanceSpec: v1alpha1.StardogInstanceSpec{
				ServerUrl: serverURL,
				AdminCredentials: v1alpha1.StardogUserCredentialsSpec{
					Namespace: namespace,
				},
			},
			err: errors.New(".spec.AdminCredentials.SecretRef is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := StardogInstanceReconciler{
				Log:               testing2.TestLogger{},
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            nil,
			}
			err := r.validateSpecification(tt.stardogInstanceSpec)

			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_validateConnection(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	secretName := "secret-test"
	serverURL := "url-test"
	username := "admin"
	password := "1234"

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)
	encodedUser := base64.StdEncoding.EncodeToString([]byte(username))

	tests := []struct {
		name            string
		stardogInstance v1alpha1.StardogInstance
		sir             StardogInstanceReconciliation
		secret          v1.Secret
		err             error
	}{
		{
			name:            "GivenReconciliationContext_WhenConnectionToStardogIsOK_ThenDoNotReturnAnyError",
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			sir: StardogInstanceReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			},
			err: nil,
		},
		{
			name:            "GivenReconciliationContext_WhenConnectionToStardogIsKO_ThenRaiseError",
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			sir: StardogInstanceReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(v1alpha1.StardogConditionMap),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: createStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
			},
			err: errors.New("error connection"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.stardogInstance, &tt.secret)
			assert.NoError(t, err)
			r := StardogInstanceReconciler{
				Log:               testing2.TestLogger{},
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			stardogClient.EXPECT().SetConnection(
				serverURL,
				encodedUser,
				base64.StdEncoding.EncodeToString([]byte(password)))
			stardogClient.EXPECT().IsEnabled(context.Background(), encodedUser).Return(stardogrest.Enabled{}, tt.err)
			err = r.validateConnection(&tt.sir)

			assert.Equal(t, tt.err, err)
		})
	}
}

func createStardogInstance(namespace, name, secretName, serverURL string) *v1alpha1.StardogInstance {
	return &v1alpha1.StardogInstance{
		TypeMeta:   metav1.TypeMeta{Kind: "StardogInstance", APIVersion: "v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: v1alpha1.StardogInstanceSpec{
			AdminCredentials: v1alpha1.StardogUserCredentialsSpec{
				Namespace: namespace,
				SecretRef: secretName,
			},
			ServerUrl: serverURL,
		},
	}
}

func createStardogInstanceWithFinalizers(namespace, name, secretName, serverURL string) *v1alpha1.StardogInstance {
	instance := createStardogInstance(namespace, name, secretName, serverURL)
	instance.SetFinalizers([]string{instanceRoleFinalizer, instanceUserFinalizer})
	return instance
}
