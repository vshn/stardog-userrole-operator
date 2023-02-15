package util_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vshn/stardog-userrole-operator/pkg/stardogapi"
	"github.com/vshn/stardog-userrole-operator/pkg/stardogapi/mock"
	"github.com/vshn/stardog-userrole-operator/pkg/stardogapi/util"
)

func TestComparePermission(t *testing.T) {
	x := stardogapi.Permission{
		Action:       "write",
		ResourceType: "db",
		Resources:    []string{"foobar"},
	}
	y := stardogapi.Permission{
		Action:       "write",
		ResourceType: "db",
		Resources:    []string{"foobar"},
	}

	assert.True(t, util.ComparePermission(x, y))

	y.Action = "read"
	assert.False(t, util.ComparePermission(x, y))

	x.Action = "read"
	y.ResourceType = "role"
	assert.False(t, util.ComparePermission(x, y))

	x.ResourceType = "role"
	y.Resources = []string{"foobar", "foo"}
	assert.False(t, util.ComparePermission(x, y))
}

func TestAddPermissions(t *testing.T) {
	ctrl, stardog := setupTestMock(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	readPermissions := []stardogapi.Permission{
		{ResourceType: "db", Action: "READ", Resources: []string{"foo"}},
		{ResourceType: "metadata", Action: "READ", Resources: []string{"foo"}},
	}

	stardog.EXPECT().
		GetRolePermissions(ctx, gomock.Any()).
		Return(nil, nil).
		Times(1)

	stardog.EXPECT().
		AddRolePermission(ctx, gomock.Any(), gomock.Any()).
		Return(nil).
		Times(2)

	err := util.AddPermissions(ctx, stardog, "foo-read", readPermissions)
	assert.NoError(t, err)
}

func TestCreateRoles(t *testing.T) {
	ctrl, stardog := setupTestMock(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	roles := []string{"foo-read", "foo-write"}

	stardog.EXPECT().
		GetRoles(ctx).
		Return(nil, nil).
		Times(1)

	stardog.EXPECT().
		AddRole(ctx, gomock.Any()).
		Return(nil).
		Times(2)

	err := util.CreateRoles(ctx, stardog, roles)
	assert.NoError(t, err)

	// Don't create roles that already exist
	stardog.EXPECT().
		GetRoles(ctx).
		Return([]string{"foo-read"}, nil).
		Times(1)

	stardog.EXPECT().
		AddRole(ctx, "foo-write").
		Return(nil).
		Times(1)

	err = util.CreateRoles(ctx, stardog, roles)
	assert.NoError(t, err)
}

func TestCreateUsers(t *testing.T) {
	ctrl, stardog := setupTestMock(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	users := []stardogapi.UserCredentials{
		{
			Name:     "foo-read",
			Password: "abcdef",
		},
		{
			Name:     "foo-write",
			Password: "abcdef",
		},
	}

	stardog.EXPECT().
		GetUser(ctx, gomock.Any()).
		Return(stardogapi.User{}, errors.New("user does not exist")).
		Times(2)

	stardog.EXPECT().
		AddUser(ctx, gomock.Any(), gomock.Any()).
		Return(nil).
		Times(2)

	err := util.CreateUsers(ctx, stardog, users)
	assert.NoError(t, err)

	// Don't create users that already exist
	stardog.EXPECT().
		GetUser(ctx, gomock.Any()).
		Return(stardogapi.User{Name: "foo-read"}, nil).
		Times(1)
	stardog.EXPECT().
		GetUser(ctx, gomock.Any()).
		Return(stardogapi.User{}, errors.New("user does not exist")).
		Times(1)

	stardog.EXPECT().
		AddUser(ctx, "foo-write", gomock.Any()).
		Return(nil).
		Times(1)

	err = util.CreateUsers(ctx, stardog, users)
	assert.NoError(t, err)
}

func setupTestMock(t *testing.T) (*gomock.Controller, *mock.MockStardogAPI) {
	t.Helper()
	ctrl := gomock.NewController(t)
	stardog := mock.NewMockStardogAPI(ctrl)

	return ctrl, stardog
}
