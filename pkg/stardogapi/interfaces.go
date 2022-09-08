package stardogapi

import "context"

// StardogAPI provides an interface to interact with a subset of the Stardog API
//
//go:generate mockgen -source interfaces.go -destination mock/mock_client.go -package mock -aux_files=github.com/vshn/stardog-userrole-operator/pkg/stardogapi=interfaces.go
type StardogAPI interface {
	// DB
	CreateDatabase(ctx context.Context, name string, options map[string]string) (err error)
	DropDatabase(ctx context.Context, name string) (err error)
	ListDatabases(ctx context.Context) (databases []string, err error)

	// User
	AddUser(ctx context.Context, name, password string) (err error)
	DeleteUser(ctx context.Context, name string) (err error)
	GetUser(ctx context.Context, name string) (user User, err error)
	SetUserRoles(ctx context.Context, name string, roles []string) (err error)
	GetUserRoles(ctx context.Context, name string) (roles []string, err error)

	// Roles
	AddRole(ctx context.Context, name string) (err error)
	DeleteRole(ctx context.Context, name string) (err error)
	GetRoles(ctx context.Context) (roles []string, err error)

	// Permissions
	AddRolePermission(ctx context.Context, name string, permission Permission) (err error)
	DeleteRolePermission(ctx context.Context, name string, permission Permission) (err error)
	GetRolePermissions(ctx context.Context, name string) (permissions []Permission, err error)
}

var _ StardogAPI = (*Client)(nil)
