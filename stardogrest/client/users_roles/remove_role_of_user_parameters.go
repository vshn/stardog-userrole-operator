// Code generated by go-swagger; DO NOT EDIT.

package users_roles

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewRemoveRoleOfUserParams creates a new RemoveRoleOfUserParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewRemoveRoleOfUserParams() *RemoveRoleOfUserParams {
	return &RemoveRoleOfUserParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewRemoveRoleOfUserParamsWithTimeout creates a new RemoveRoleOfUserParams object
// with the ability to set a timeout on a request.
func NewRemoveRoleOfUserParamsWithTimeout(timeout time.Duration) *RemoveRoleOfUserParams {
	return &RemoveRoleOfUserParams{
		timeout: timeout,
	}
}

// NewRemoveRoleOfUserParamsWithContext creates a new RemoveRoleOfUserParams object
// with the ability to set a context for a request.
func NewRemoveRoleOfUserParamsWithContext(ctx context.Context) *RemoveRoleOfUserParams {
	return &RemoveRoleOfUserParams{
		Context: ctx,
	}
}

// NewRemoveRoleOfUserParamsWithHTTPClient creates a new RemoveRoleOfUserParams object
// with the ability to set a custom HTTPClient for a request.
func NewRemoveRoleOfUserParamsWithHTTPClient(client *http.Client) *RemoveRoleOfUserParams {
	return &RemoveRoleOfUserParams{
		HTTPClient: client,
	}
}

/*
RemoveRoleOfUserParams contains all the parameters to send to the API endpoint

	for the remove role of user operation.

	Typically these are written to a http.Request.
*/
type RemoveRoleOfUserParams struct {

	/* Role.

	   The name of the role to remove
	*/
	Role string

	/* User.

	   The username of the user whose role should be removed
	*/
	User string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the remove role of user params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RemoveRoleOfUserParams) WithDefaults() *RemoveRoleOfUserParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the remove role of user params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RemoveRoleOfUserParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the remove role of user params
func (o *RemoveRoleOfUserParams) WithTimeout(timeout time.Duration) *RemoveRoleOfUserParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the remove role of user params
func (o *RemoveRoleOfUserParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the remove role of user params
func (o *RemoveRoleOfUserParams) WithContext(ctx context.Context) *RemoveRoleOfUserParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the remove role of user params
func (o *RemoveRoleOfUserParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the remove role of user params
func (o *RemoveRoleOfUserParams) WithHTTPClient(client *http.Client) *RemoveRoleOfUserParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the remove role of user params
func (o *RemoveRoleOfUserParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRole adds the role to the remove role of user params
func (o *RemoveRoleOfUserParams) WithRole(role string) *RemoveRoleOfUserParams {
	o.SetRole(role)
	return o
}

// SetRole adds the role to the remove role of user params
func (o *RemoveRoleOfUserParams) SetRole(role string) {
	o.Role = role
}

// WithUser adds the user to the remove role of user params
func (o *RemoveRoleOfUserParams) WithUser(user string) *RemoveRoleOfUserParams {
	o.SetUser(user)
	return o
}

// SetUser adds the user to the remove role of user params
func (o *RemoveRoleOfUserParams) SetUser(user string) {
	o.User = user
}

// WriteToRequest writes these params to a swagger request
func (o *RemoveRoleOfUserParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param role
	if err := r.SetPathParam("role", o.Role); err != nil {
		return err
	}

	// path param user
	if err := r.SetPathParam("user", o.User); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
