package stardogapi

import (
	"context"
	"net/http"
	"path"
	"strings"
)

type addUserRequest struct {
	Username string   `json:"username"`
	Password []string `json:"password"`
}

type userRolesRequestAndResponse struct {
	Roles []string `json:"roles"`
}

// Add a new user
func (c *Client) AddUser(ctx context.Context, name, password string) (err error) {
	return c.sendRequest(ctx,
		http.MethodPost,
		"/admin/users",
		&addUserRequest{
			Username: name,
			Password: strings.Split(password, ""), // The API expects the password split into an array of single characters
		},
		nil,
	)
}

// Delete a user
func (c *Client) DeleteUser(ctx context.Context, name string) (err error) {
	return c.sendRequest(ctx,
		http.MethodDelete,
		path.Join("/admin/users/", sanitizePathValue(name)),
		nil,
		nil,
	)
}

// Get a user
func (c *Client) GetUser(ctx context.Context, name string) (user User, err error) {
	return user, c.sendRequest(ctx,
		http.MethodGet,
		path.Join("/admin/users/", sanitizePathValue(name)),
		nil,
		&user,
	)
}

// Set the roles of a user
func (c *Client) SetUserRoles(ctx context.Context, name string, roles []string) (err error) {
	return c.sendRequest(ctx,
		http.MethodPut,
		path.Join("/admin/users/", sanitizePathValue(name), "/roles"),
		&userRolesRequestAndResponse{Roles: roles},
		nil,
	)
}

// Get the roles of a user
func (c *Client) GetUserRoles(ctx context.Context, name string) (roles []string, err error) {
	var response userRolesRequestAndResponse

	return response.Roles, c.sendRequest(ctx,
		http.MethodGet,
		path.Join("/admin/users/", sanitizePathValue(name), "/roles"),
		nil,
		&response,
	)
}
