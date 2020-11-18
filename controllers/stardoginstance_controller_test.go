package controllers

import (
	"context"
	"errors"
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
			stardogClient.EXPECT().SetConnection(serverURL, username, password)
			stardogClient.EXPECT().ListUsers(context.Background()).Return(stardogrest.Users{Users: &[]string{}}, nil)
			stardogClient.EXPECT().ListRoles(context.Background()).Return(stardogrest.Roles{Roles: &[]string{}}, nil)
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

	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name          string
		conditionUser func(stardogrestapi2.ExtendedBaseClientAPI)
		err           error
	}{
		{
			name: "GivenUserFinalizer_WhenThereAreNoUsers_ThenDoNotReturnAnyError",
			conditionUser: func(stardogrestapi2.ExtendedBaseClientAPI) {
				stardogClient.EXPECT().ListUsers(ctx).Return(stardogrest.Users{Users: &[]string{}}, nil)
			},
			err: nil,
		},
		{
			name: "GivenUserFinalizer_WhenThereAreUsers_ThenRaiseError",
			conditionUser: func(stardogrestapi2.ExtendedBaseClientAPI) {
				stardogClient.EXPECT().ListUsers(ctx).Return(stardogrest.Users{Users: &[]string{"pippo"}}, nil)
			},
			err: errors.New("cannot delete StardogInstance, found 1 users"),
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
			tt.conditionUser(stardogClient)
			err := r.userFinalizer(stardogClient, ctx)

			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_roleFinalize(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogClient := stardogrestapi.NewMockExtendedBaseClientAPI(mockCtrl)

	tests := []struct {
		name          string
		conditionRole func(stardogrestapi2.ExtendedBaseClientAPI)
		err           error
	}{
		{
			name: "GivenRoleFinalizer_WhenThereAreNoRoles_ThenDoNotReturnAnyError",
			conditionRole: func(stardogrestapi2.ExtendedBaseClientAPI) {
				stardogClient.EXPECT().ListRoles(ctx).Return(stardogrest.Roles{Roles: &[]string{}}, nil)
			},
			err: nil,
		},
		{
			name: "GivenRoleFinalizer_WhenThereAreRoles_ThenRaiseError",
			conditionRole: func(stardogrestapi2.ExtendedBaseClientAPI) {
				stardogClient.EXPECT().ListRoles(ctx).Return(stardogrest.Roles{Roles: &[]string{"role1"}}, nil)
			},
			err: errors.New("cannot delete StardogInstance, found 1 remaining role(s)"),
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
			tt.conditionRole(stardogClient)
			err := r.roleFinalizer(stardogClient, ctx)

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
			stardogClient.EXPECT().SetConnection(serverURL, username, password)
			stardogClient.EXPECT().IsEnabled(context.Background(), username).Return(stardogrest.Enabled{}, tt.err)
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
