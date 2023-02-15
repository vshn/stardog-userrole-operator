package stardogapi

import (
	"context"
	"net/http"
	"path"
)

type addRoleRequest struct {
	Rolename string `json:"rolename"`
}

type rolesResponse struct {
	Roles []string `json:"roles"`
}

// Get the available roles
func (c *Client) GetRoles(ctx context.Context) ([]string, error) {
	var rolesResponse rolesResponse

	return rolesResponse.Roles, c.sendRequest(ctx,
		http.MethodGet,
		"/admin/roles/",
		nil,
		&rolesResponse,
	)
}

// Add a new role
func (c *Client) AddRole(ctx context.Context, name string) (err error) {
	return c.sendRequest(ctx,
		http.MethodPost,
		"/admin/roles",
		&addRoleRequest{Rolename: name},
		nil,
	)
}

// Delete a role
func (c *Client) DeleteRole(ctx context.Context, name string) (err error) {
	return c.sendRequest(ctx,
		http.MethodDelete,
		path.Join("/admin/roles/", sanitizePathValue(name)),
		nil,
		nil,
	)
}
