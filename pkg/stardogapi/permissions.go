package stardogapi

import (
	"context"
	"net/http"
	"path"
)

type getRolePermissionsRolesResponse struct {
	Permissions []Permission `json:"permissions"`
}

// Get Permissions of a role
func (c *Client) GetRolePermissions(ctx context.Context, name string) ([]Permission, error) {
	var response getRolePermissionsRolesResponse

	return response.Permissions, c.sendRequest(ctx,
		http.MethodGet,
		path.Join("/admin/permissions/role/", sanitizePathValue(name)),
		nil,
		&response,
	)
}

// Add the Permission to a role
func (c *Client) AddRolePermission(ctx context.Context, name string, permission Permission) (err error) {
	return c.sendRequest(ctx,
		http.MethodPut,
		path.Join("/admin/permissions/role/", sanitizePathValue(name)),
		&permission,
		nil,
	)
}

// Delete the Permission from a role
func (c *Client) DeleteRolePermission(ctx context.Context, name string, permission Permission) (err error) {
	return c.sendRequest(ctx,
		http.MethodPost,
		path.Join("/admin/permissions/role/", sanitizePathValue(name), "/delete"),
		&permission,
		nil,
	)
}
