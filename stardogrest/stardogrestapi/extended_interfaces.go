package stardogrestapi

import "github.com/vshn/stardog-userrole-operator/stardogrest"

//go:generate mockgen -source extended_interfaces.go -destination ../mocks/mock_client.go -package stardogrestapi -aux_files=github.com/vshn/stardog-userrole-operator/stardogrest/stardogrestapi=interfaces.go
type ExtendedBaseClientAPI interface {
	SetConnection(url, username, password string)
	BaseClientAPI
}

var _ ExtendedBaseClientAPI = (*stardogrest.ExtendedBaseClient)(nil)
