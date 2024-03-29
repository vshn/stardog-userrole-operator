package stardogapi

import (
	"context"
	"net/http"
	"path"
)

// TODO extend
type createDatabaseRequest struct {
	Name string `json:"dbname"`
}

type listDatabasesResponse struct {
	Databases []string `json:"databases"`
}

// Creates a database with the given name and options
func (c *Client) CreateDatabase(ctx context.Context, name string, options map[string]string) (err error) {
	return c.sendMultipartJsonRequest(ctx,
		http.MethodPost,
		"/admin/databases",
		map[string]any{"root": &createDatabaseRequest{Name: name}},
		nil,
	)
}

// Drops the given database
func (c *Client) DropDatabase(ctx context.Context, name string) (err error) {
	return c.sendRequest(ctx,
		http.MethodDelete,
		path.Join("/admin/databases/", sanitizePathValue(name)),
		nil,
		nil,
	)
}

// Returns the list of databases
func (c *Client) ListDatabases(ctx context.Context) (databases []string, err error) {
	var response listDatabasesResponse

	return response.Databases, c.sendRequest(ctx,
		http.MethodGet,
		"/admin/databases",
		nil,
		&response,
	)
}
