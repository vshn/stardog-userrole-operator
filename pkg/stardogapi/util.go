package stardogapi

import (
	"context"
	"strings"

	"golang.org/x/exp/slices"
)

// Compares permissions
func ComparePermission(x, y Permission) bool {
	if x.ResourceType == y.ResourceType {
		if x.Action == y.Action {
			return slices.Equal(x.Resources, y.Resources)
		}
	}
	return false
}

// Adds permissions to a role
func AddPermissions(ctx context.Context, stardogAPI StardogAPI, role string, permissions []Permission) error {
	rolePermissions, err := stardogAPI.GetRolePermissions(ctx, role)
	if err != nil {
		return err
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
				return err
			}
		}
	}
	return nil
}

// Create multiple roles
func CreateRoles(ctx context.Context, stardogAPI StardogAPI, roles []string) error {
	activeRoles, err := stardogAPI.GetRoles(ctx)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if !slices.Contains(activeRoles, role) {
			err = stardogAPI.AddRole(ctx, role)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Create multiple users
func CreateUsers(ctx context.Context, stardogAPI StardogAPI, users []UserCredentials) error {
	for _, user := range users {
		_, err := stardogAPI.GetUser(ctx, user.Name)
		if err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				err = stardogAPI.AddUser(ctx, user.Name, user.Password)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}
