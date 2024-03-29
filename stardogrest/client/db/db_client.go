// Code generated by go-swagger; DO NOT EDIT.

package db

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new db API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for db API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	CreateNewDatabase(params *CreateNewDatabaseParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateNewDatabaseCreated, error)

	DropDatabase(params *DropDatabaseParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DropDatabaseOK, error)

	GetDBSize(params *GetDBSizeParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetDBSizeOK, error)

	ListDatabases(params *ListDatabasesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListDatabasesOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
CreateNewDatabase creates database

Add a new database to the server, optionally with RDF bulk-loaded
*/
func (a *Client) CreateNewDatabase(params *CreateNewDatabaseParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateNewDatabaseCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateNewDatabaseParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "createNewDatabase",
		Method:             "POST",
		PathPattern:        "/admin/databases",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"multipart/form-data"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateNewDatabaseReader{formats: a.formats},
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
	success, ok := result.(*CreateNewDatabaseCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for createNewDatabase: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
DropDatabase drops database

Delete the database
*/
func (a *Client) DropDatabase(params *DropDatabaseParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DropDatabaseOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDropDatabaseParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "dropDatabase",
		Method:             "DELETE",
		PathPattern:        "/admin/databases/{db}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DropDatabaseReader{formats: a.formats},
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
	success, ok := result.(*DropDatabaseOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*DropDatabaseDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
GetDBSize gets d b size

Retrieve the size of the db. Size is approximate unless the exact parameter is set to true
*/
func (a *Client) GetDBSize(params *GetDBSizeParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetDBSizeOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetDBSizeParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getDBSize",
		Method:             "GET",
		PathPattern:        "/{db}/size",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetDBSizeReader{formats: a.formats},
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
	success, ok := result.(*GetDBSizeOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getDBSize: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
ListDatabases lists databases

List all the databases in the server
*/
func (a *Client) ListDatabases(params *ListDatabasesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListDatabasesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListDatabasesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listDatabases",
		Method:             "GET",
		PathPattern:        "/admin/databases",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListDatabasesReader{formats: a.formats},
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
	success, ok := result.(*ListDatabasesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for listDatabases: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
