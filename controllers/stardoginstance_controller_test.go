package controllers

import (
	"context"
	"errors"
	stardog_client "github.com/vshn/stardog-userrole-operator/stardogrest/client"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users"
	stardogmock "github.com/vshn/stardog-userrole-operator/stardogrest/mocks"
	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
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
	stardogMocked := stardogmock.NewMockStardogTestClient(mockCtrl)
	stardogClient := createStardogClientFromMock(stardogMocked)

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
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}

			err = fakeKubeClient.Get(context.Background(), types.NamespacedName{
				Namespace: namespace,
				Name:      stardogInstanceName,
			}, tt.sir.resource)
			assert.NoError(t, err)

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
	stardogMocked := stardogmock.NewMockStardogTestClient(mockCtrl)
	stardogClient := createStardogClientFromMock(stardogMocked)

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
				Log:               testr.New(t),
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
	stardogMocked := stardogmock.NewMockStardogTestClient(mockCtrl)
	stardogClient := createStardogClientFromMock(stardogMocked)

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
				Log:               testr.New(t),
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
				Log:               testr.New(t),
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
	serverURL := "http://url-test:5820"
	username := "admin"
	password := "1234"

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogMocked := stardogmock.NewMockStardogTestClient(mockCtrl)
	stardogClient := createStardogClientFromMock(stardogMocked)

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
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}

			stardogMocked.EXPECT().
				SetTransport(gomock.Any()).
				AnyTimes()

			stardogMocked.EXPECT().
				IsEnabled(gomock.Any(), gomock.Any()).
				Return(&users.IsEnabledOK{Payload: &models.Enabled{Enabled: false}}, tt.err).
				Times(1)

			err = r.validateConnection(&tt.sir)

			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_ReconcileInstance(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	secretName := "secret-test"
	serverURL := "http://localhost:5820/"

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{Name: stardogInstanceName, Namespace: namespace},
	}

	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	tests := []struct {
		name            string
		namespace       v1.Namespace
		stardogInstance v1alpha1.StardogInstance
		sir             StardogInstanceReconciliation
		expectedResult  ctrl.Result
	}{
		{
			name:            "GivenReconciliation_WhenStardogInstanceMissing_ThenReturnNoRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, "nonExistingInstance", secretName, serverURL),
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClient(&tt.namespace, &tt.stardogInstance)
			assert.NoError(t, err)
			r := StardogInstanceReconciler{
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

func Test_ReconcileStardogInstance(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	secretName := "secret-test"
	stardogRoleName := "stardog-role"
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
		sir             StardogInstanceReconciliation
		secret          v1.Secret
		expectations    []func(stardog_client.Stardog)
		expectedResult  ctrl.Result
	}{
		{
			name:            "GivenReconciliation_WhenStardogInstanceIsDeleted_ThenReturnNoRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createDeletedStardogInstance(namespace, stardogInstanceName, secretName, serverURL),
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
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogInstanceIsDeletedWithFinalizers_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createDeletedStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
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
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogInstanceIsDeletedWithFinalizersAndRole_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createDeletedStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
			stardogRole:     *createStardogRole(namespace, stardogRoleName, stardogInstanceName, []v1alpha1.StardogPermissionSpec{}),
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
			expectedResult: ctrl.Result{
				Requeue:      true,
				RequeueAfter: ReconFreqErr,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogInstanceIsInvalid_ThenReturnNoRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, stardogInstanceName, "", serverURL),
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
			expectedResult: ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogInstanceUnreachable_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
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
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						SetTransport(gomock.Any()).
						AnyTimes()

					stardogMocked.EXPECT().
						IsEnabled(gomock.Any(), gomock.Any()).
						Return(&users.IsEnabledOK{}, errors.New("cannot connect to Stardog"))
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      true,
				RequeueAfter: ReconFreqErr,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogInstanceCannotBeUpdated_ThenReturnRequeue",
			namespace:       *createNamespace(namespace),
			stardogInstance: *createStardogInstance(namespace, "", secretName, serverURL),
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
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						SetTransport(gomock.Any()).
						AnyTimes()
					stardogMocked.EXPECT().
						IsEnabled(gomock.Any(), gomock.Any())
				},
			},
			expectedResult: ctrl.Result{
				Requeue:      true,
				RequeueAfter: ReconFreqErr,
			},
		},
		{
			name:            "GivenReconciliation_WhenStardogInstanceIsReconciled_ThenReturnNoRequeue",
			namespace:       *createNamespace(namespace),
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
			expectations: []func(stardog_client.Stardog){
				func(stardog_client.Stardog) {
					stardogMocked.EXPECT().
						SetTransport(gomock.Any()).
						AnyTimes()
					stardogMocked.EXPECT().
						IsEnabled(gomock.Any(), gomock.Any())
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
			fakeKubeClient, err := createKubeFakeClient(&tt.namespace, &tt.stardogInstance, &tt.secret, &tt.stardogRole)
			assert.NoError(t, err)
			r := StardogInstanceReconciler{
				Log:               testr.New(t),
				ReconcileInterval: time.Duration(1),
				Scheme:            scheme.Scheme,
				Client:            fakeKubeClient,
			}
			sir := &StardogInstanceReconciliation{
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(map[v1alpha1.StardogConditionType]v1alpha1.StardogCondition),
					namespace:     namespace,
					stardogClient: stardogClient,
				},
				resource: &tt.stardogInstance,
			}
			for _, addExpectation := range tt.expectations {
				addExpectation(*tt.sir.reconciliationContext.stardogClient)
			}

			result, err := r.ReconcileStardogInstance(sir)

			assert.Equal(t, tt.expectedResult, result)
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

func createDeletedStardogInstance(namespace, name, secretName, serverURL string) *v1alpha1.StardogInstance {
	stardogInstace := createStardogInstance(namespace, name, secretName, serverURL)
	newTime := metav1.NewTime(time.Now())
	stardogInstace.SetDeletionTimestamp(&newTime)
	return stardogInstace
}

func createDeletedStardogInstanceWithFinalizers(namespace, name, secretName, serverURL string) *v1alpha1.StardogInstance {
	stardogInstance := createStardogInstance(namespace, name, secretName, serverURL)
	newTime := metav1.NewTime(time.Now())
	stardogInstance.SetDeletionTimestamp(&newTime)
	stardogInstance.SetFinalizers([]string{instanceRoleFinalizer, instanceUserFinalizer})
	return stardogInstance
}

func createStardogInstanceWithFinalizers(namespace, name, secretName, serverURL string) *v1alpha1.StardogInstance {
	stardogInstance := createStardogInstance(namespace, name, secretName, serverURL)
	stardogInstance.SetFinalizers([]string{instanceRoleFinalizer, instanceUserFinalizer})
	return stardogInstance
}
