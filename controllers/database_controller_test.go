package controllers

import (
	"context"
	"github.com/go-logr/logr/testr"
	"github.com/go-openapi/runtime"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	"github.com/vshn/stardog-userrole-operator/api/v1beta1"
	db2 "github.com/vshn/stardog-userrole-operator/stardogrest/client/db"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles_permissions"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users_roles"
	stardogmock "github.com/vshn/stardog-userrole-operator/stardogrest/mocks"
	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"testing"
)

func Test_createCustomUser(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	secretName := "secret-test"
	serverURL := "http://url-test.ch"
	username := "admin"
	password := "1234"
	err := v1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)
	err = v1beta1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stardogMocked := stardogmock.NewMockStardogTestClient(mockCtrl)
	stardogClient := createStardogClientFromMock(stardogMocked)
	db := createStardogDB("test-db", "hidden-user", v1beta1.StardogInstanceRef{
		Name:      stardogInstanceName,
		Namespace: namespace,
	})
	org := createOrg("org-test", db.Name, []v1beta1.NamedGraph{{
		Name:      "graph1",
		AddHidden: true,
	}})

	tests := []struct {
		name                     string
		stardogInstance          v1alpha1.StardogInstance
		secret                   v1.Secret
		stardogDB                *v1beta1.Database
		stardogOrg               *v1beta1.Organization
		dr                       DatabaseReconciliation
		expectedCustomUser       bool
		expectedCustomRole       bool
		expectedWritePermissions bool
	}{
		{
			name:            "GivenReconciliationContext_WhenCreateDatabaseWithCustomUser_ThenCreateCustomUser",
			stardogInstance: *createStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			stardogDB:       db,
			stardogOrg:      org,
			dr: DatabaseReconciliation{
				resource: db,
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(map[v1alpha1.StardogConditionType]v1alpha1.StardogCondition),
					stardogClient: stardogClient,
				},
			},
			expectedCustomUser:       true,
			expectedCustomRole:       true,
			expectedWritePermissions: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClientWithSub(tt.stardogDB, &tt.stardogInstance, &tt.secret)
			assert.NoError(t, err)
			r := DatabaseReconciler{
				Log:    testr.New(t),
				Scheme: scheme.Scheme,
				Client: fakeKubeClient,
			}

			customUserExists := false
			writePermissions := false
			stardogMocked.EXPECT().
				SetTransport(gomock.Any()).
				AnyTimes()
			stardogMocked.EXPECT().
				ListDatabases(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&db2.ListDatabasesOK{Payload: &models.Databases{}}, nil).
				Times(1)
			stardogMocked.EXPECT().
				CreateNewDatabase(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(db2.NewCreateNewDatabaseCreated(), nil).
				Times(1)
			stardogMocked.EXPECT().
				ListUsers(gomock.Any(), gomock.Any()).
				Return(&users.ListUsersOK{Payload: &models.Users{}}, nil).
				Times(1)
			stardogMocked.EXPECT().
				CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(params *users.CreateUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.CreateUserCreated, error) {
					if *params.User.Username == "hidden-user" {
						customUserExists = true
					}
					return users.NewCreateUserCreated(), nil
				}).
				Times(3)
			stardogMocked.EXPECT().
				ListRoles(gomock.Any(), gomock.Any()).
				Return(&roles.ListRolesOK{Payload: &models.Roles{}}, nil).
				Times(1)
			stardogMocked.EXPECT().
				CreateRole(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(roles.NewCreateRoleCreated(), nil).
				Times(3)
			stardogMocked.EXPECT().
				ListRolePermissions(gomock.Any(), gomock.Any()).
				Return(&roles_permissions.ListRolePermissionsOK{Payload: &models.Permissions{}}, nil).
				Times(3)
			stardogMocked.EXPECT().
				ListUserRoles(gomock.Any(), gomock.Any()).
				Return(&users_roles.ListUserRolesOK{Payload: &models.Roles{}}, nil).
				Times(3)
			stardogMocked.EXPECT().
				AddRole(gomock.Any(), gomock.Any()).
				Return(users_roles.NewAddRoleNoContent(), nil).
				Times(3)
			stardogMocked.EXPECT().
				AddRolePermission(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(params *roles_permissions.AddRolePermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_roles.ClientOption) (*roles_permissions.AddRolePermissionCreated, error) {
					if *params.Permission.Action == "WRITE" && params.Role == "hidden-user" {
						writePermissions = true
					}
					return roles_permissions.NewAddRolePermissionCreated(), nil
				}).AnyTimes()

			_, err = r.reconcileDatabase(&tt.dr)
			assert.Equal(t, tt.expectedCustomUser, customUserExists)
			assert.Equal(t, tt.expectedWritePermissions, writePermissions)
		})
	}
}

func createStardogDB(name, hiddenUser string, instanceRef v1beta1.StardogInstanceRef) *v1beta1.Database {
	return &v1beta1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1beta1.DatabaseSpec{
			DatabaseName:              name,
			AddUserForNonHiddenGraphs: hiddenUser,
			StardogInstanceRefs:       []v1beta1.StardogInstanceRef{instanceRef},
			NamedGraphPrefix:          "https://graph.ch",
		},
	}
}
