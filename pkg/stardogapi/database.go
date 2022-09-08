package stardogapi

import (
	"context"
	"fmt"
	"net/http"
)

// TODO extend
type createDatabaseRequest struct {
	Root createDatabaseRequestRoot `json:"root"`
}

type createDatabaseRequestRoot struct {
	Name string `json:"dbname"`
}

type listDatabasesResponse struct {
	Databases []string `json:"databases"`
}

// Creates a database with the given name and options
func (c *Client) CreateDatabase(ctx context.Context, name string, options map[string]string) (err error) {
	return c.sendRequest(ctx,
		http.MethodPost,
		"/admin/databases",
		&createDatabaseRequest{
			Root: createDatabaseRequestRoot{Name: name}},
		nil,
	)
}

// Drops the given database
func (c *Client) DropDatabase(ctx context.Context, name string) (err error) {
	return c.sendRequest(ctx,
		http.MethodDelete,
		fmt.Sprintf("/admin/databases/%s", name),
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
