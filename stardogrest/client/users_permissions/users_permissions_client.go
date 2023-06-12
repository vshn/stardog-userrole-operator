// Code generated by go-swagger; DO NOT EDIT.

package users_permissions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new users permissions API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for users permissions API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	AddUserPermission(params *AddUserPermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*AddUserPermissionCreated, error)

	ListEffectivePermissions(params *ListEffectivePermissionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEffectivePermissionsOK, error)

	ListUserPermissions(params *ListUserPermissionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListUserPermissionsOK, error)

	RemoveUserPermission(params *RemoveUserPermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*RemoveUserPermissionCreated, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
AddUserPermission adds a permission to a user
*/
func (a *Client) AddUserPermission(params *AddUserPermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*AddUserPermissionCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewAddUserPermissionParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "addUserPermission",
		Method:             "PUT",
		PathPattern:        "/permissions/user/{user}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &AddUserPermissionReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*AddUserPermissionCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*AddUserPermissionDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEffectivePermissions lists the user s effective permissions all permissions
*/
func (a *Client) ListEffectivePermissions(params *ListEffectivePermissionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEffectivePermissionsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEffectivePermissionsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEffectivePermissions",
		Method:             "GET",
		PathPattern:        "/permissions/effective/user/{user}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListEffectivePermissionsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEffectivePermissionsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEffectivePermissionsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListUserPermissions lists the user s direct permissions not via roles
*/
func (a *Client) ListUserPermissions(params *ListUserPermissionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListUserPermissionsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListUserPermissionsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listUserPermissions",
		Method:             "GET",
		PathPattern:        "/permissions/user/{user}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListUserPermissionsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListUserPermissionsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListUserPermissionsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
RemoveUserPermission removes a permission from a user
*/
func (a *Client) RemoveUserPermission(params *RemoveUserPermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*RemoveUserPermissionCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewRemoveUserPermissionParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "removeUserPermission",
		Method:             "POST",
		PathPattern:        "/permissions/user/{user}/delete",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &RemoveUserPermissionReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*RemoveUserPermissionCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*RemoveUserPermissionDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
