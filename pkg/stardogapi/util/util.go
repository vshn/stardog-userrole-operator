package util

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/vshn/stardog-userrole-operator/pkg/stardogapi"
)

// Compares permissions
func ComparePermission(x, y stardogapi.Permission) bool {
	if x.ResourceType == y.ResourceType {
		if x.Action == y.Action {
			return slices.Equal(x.Resources, y.Resources)
		}
	}
	return false
}

// Adds permissions to a role
func AddPermissions(ctx context.Context, stardogAPI stardogapi.StardogAPI, role string, permissions []stardogapi.Permission) error {
	rolePermissions, err := stardogAPI.GetRolePermissions(ctx, role)
	if err != nil {
		return fmt.Errorf("error getting role permissions %s: %w", role, err)
	}

	for _, permission := range permissions {
		exists := false
		for _, rolePermission := range rolePermissions {
			if ComparePermission(rolePermission, permission) {
				exists = true
				break
			}
		}

		if !exists {
			err = stardogAPI.AddRolePermission(ctx, role, permission)
			if err != nil {
				return fmt.Errorf("error adding permission %s to role %s: %w", permission, role, err)
			}
		}
	}
	return nil
}

// Create multiple roles
func CreateRoles(ctx context.Context, stardogAPI stardogapi.StardogAPI, roles []string) error {
	activeRoles, err := stardogAPI.GetRoles(ctx)
	if err != nil {
		return fmt.Errorf("error getting roles: %w", err)
	}

	for _, role := range roles {
		if !slices.Contains(activeRoles, role) {
			err = stardogAPI.AddRole(ctx, role)
			if err != nil {
				return fmt.Errorf("error creating role %s: %w", role, err)
			}
		}
	}
	return nil
}

// Create multiple users
func CreateUsers(ctx context.Context, stardogAPI stardogapi.StardogAPI, users []stardogapi.UserCredentials) error {
	for _, user := range users {
		_, err := stardogAPI.GetUser(ctx, user.Name)
		if err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				err = stardogAPI.AddUser(ctx, user.Name, user.Password)
				if err != nil {
					return fmt.Errorf("error creating user %s: %w", user, err)
				}
			} else {
				return fmt.Errorf("error getting user %s: %w", user, err)
			}
		}
	}
	return nil
}
