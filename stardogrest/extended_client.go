package stardogrest

import (
	"github.com/Azure/go-autorest/autorest"
)

type ExtendedBaseClient struct {
	BaseClient
}

func (client *ExtendedBaseClient) SetConnection(url, username, password string) {
	client.BaseClient.BaseURI = url
	client.BaseClient.Authorizer = autorest.NewBasicAuthorizer(username, password)
}

func NewExtendedBaseClient() *ExtendedBaseClient {
	return &ExtendedBaseClient{
		BaseClient: New(),
	}
}
