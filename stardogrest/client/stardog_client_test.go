package client

import (
	"github.com/go-openapi/runtime"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/db"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/roles_permissions"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users_permissions"
	"github.com/vshn/stardog-userrole-operator/stardogrest/client/users_roles"
)

//go:generate mockgen -source . -destination ../mocks/mock_client.go
type StardogTestClient interface {
	db.ClientService
	roles.ClientService
	roles_permissions.ClientService
	users.ClientService
	users_permissions.ClientService
	users_roles.ClientService
	runtime.ClientTransport
}
