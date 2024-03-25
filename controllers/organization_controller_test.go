package controllers

import (
	"context"
	"github.com/go-logr/logr/testr"
	"github.com/go-openapi/runtime"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vshn/stardog-userrole-operator/api/v1alpha1"
	"github.com/vshn/stardog-userrole-operator/api/v1beta1"
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

func Test_createHiddenGraph(t *testing.T) {

	namespace := "namespace-test"
	stardogInstanceName := "instance-test"
	secretName := "secret-test"
	serverURL := "http://url-test.ch"
	username := "admin"
	password := "1234"
	action := "READ"
	resourceTypeNG := "named-graph"
	resourceTypeMeta := "metadata"
	resourceTypeDB := "db"
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
		name                string
		stardogInstance     v1alpha1.StardogInstance
		secret              v1.Secret
		stardogDB           *v1beta1.Database
		stardogOrg          *v1beta1.Organization
		or                  OrganizationReconciliation
		expectedPermissions []models.Permission
	}{
		{
			name:            "GivenReconciliationContext_WhenCreateOrgWithHiddenGraphs_ThenCreateHiddenGraphs",
			stardogInstance: *createStardogInstanceWithFinalizers(namespace, stardogInstanceName, secretName, serverURL),
			secret:          *createFullSecret(namespace, secretName, username, password),
			stardogDB:       db,
			stardogOrg:      org,
			or: OrganizationReconciliation{
				database: db,
				resource: org,
				reconciliationContext: &ReconciliationContext{
					context:       context.Background(),
					conditions:    make(map[v1alpha1.StardogConditionType]v1alpha1.StardogCondition),
					stardogClient: stardogClient,
				},
			},
			expectedPermissions: []models.Permission{
				{
					Action:       &action,
					Resource:     []string{"test-db"},
					ResourceType: &resourceTypeDB,
				},
				{
					Action:       &action,
					Resource:     []string{"test-db"},
					ResourceType: &resourceTypeMeta,
				},
				{
					Action:       &action,
					Resource:     []string{"test-db", "https://graph.ch/org-test/graph1"},
					ResourceType: &resourceTypeNG,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKubeClient, err := createKubeFakeClientWithSub(tt.stardogOrg, tt.stardogDB, &tt.stardogInstance, &tt.secret)
			assert.NoError(t, err)
			r := OrganizationReconciler{
				Log:    testr.New(t),
				Scheme: scheme.Scheme,
				Client: fakeKubeClient,
			}

			stardogMocked.EXPECT().
				SetTransport(gomock.Any()).
				AnyTimes()
			stardogMocked.EXPECT().
				ListUsers(gomock.Any(), gomock.Any()).
				Return(&users.ListUsersOK{Payload: &models.Users{}}, nil).
				Times(1)
			stardogMocked.EXPECT().
				CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
				Return(users.NewCreateUserCreated(), nil)
			stardogMocked.EXPECT().
				ListRoles(gomock.Any(), gomock.Any()).Times(1).
				Return(&roles.ListRolesOK{Payload: &models.Roles{}}, nil)
			stardogMocked.EXPECT().
				CreateRole(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
				Return(roles.NewCreateRoleCreated(), nil)
			stardogMocked.EXPECT().
				ListRolePermissions(gomock.Any(), gomock.Any()).AnyTimes().
				Return(&roles_permissions.ListRolePermissionsOK{Payload: &models.Permissions{}}, nil)
			stardogMocked.EXPECT().
				ListUserRoles(gomock.Any(), gomock.Any()).AnyTimes().
				Return(&users_roles.ListUserRolesOK{Payload: &models.Roles{}}, nil)
			stardogMocked.EXPECT().
				AddRole(gomock.Any(), gomock.Any()).AnyTimes().
				Return(users_roles.NewAddRoleNoContent(), nil)

			hiddenRoleExists := false
			var permHiddenRole []models.Permission
			stardogMocked.EXPECT().
				AddRolePermission(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(params *roles_permissions.AddRolePermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_roles.ClientOption) (*roles_permissions.AddRolePermissionCreated, error) {
					if params.Role == "hidden-user" {
						hiddenRoleExists = true
						permHiddenRole = append(permHiddenRole, *params.Permission)
					}
					return roles_permissions.NewAddRolePermissionCreated(), nil
				}).AnyTimes()

			_, err = r.reconcileOrganization(&tt.or)
			assert.True(t, hiddenRoleExists, "hidden graph exists")
			assert.Equal(t, tt.expectedPermissions, permHiddenRole)
		})
	}
}

func createOrg(name, dbRef string, ngs []v1beta1.NamedGraph) *v1beta1.Organization {
	return &v1beta1.Organization{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1beta1.OrganizationSpec{
			Name:        name,
			DisplayName: name,
			DatabaseRef: dbRef,
			NamedGraphs: ngs,
		},
	}
}
