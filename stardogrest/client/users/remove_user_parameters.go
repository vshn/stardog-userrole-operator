// Code generated by go-swagger; DO NOT EDIT.

package users

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

// NewRemoveUserParams creates a new RemoveUserParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewRemoveUserParams() *RemoveUserParams {
	return &RemoveUserParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewRemoveUserParamsWithTimeout creates a new RemoveUserParams object
// with the ability to set a timeout on a request.
func NewRemoveUserParamsWithTimeout(timeout time.Duration) *RemoveUserParams {
	return &RemoveUserParams{
		timeout: timeout,
	}
}

// NewRemoveUserParamsWithContext creates a new RemoveUserParams object
// with the ability to set a context for a request.
func NewRemoveUserParamsWithContext(ctx context.Context) *RemoveUserParams {
	return &RemoveUserParams{
		Context: ctx,
	}
}

// NewRemoveUserParamsWithHTTPClient creates a new RemoveUserParams object
// with the ability to set a custom HTTPClient for a request.
func NewRemoveUserParamsWithHTTPClient(client *http.Client) *RemoveUserParams {
	return &RemoveUserParams{
		HTTPClient: client,
	}
}

/*
RemoveUserParams contains all the parameters to send to the API endpoint

	for the remove user operation.

	Typically these are written to a http.Request.
*/
type RemoveUserParams struct {

	/* User.

	   The username of the user to delete
	*/
	User string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the remove user params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RemoveUserParams) WithDefaults() *RemoveUserParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the remove user params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RemoveUserParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the remove user params
func (o *RemoveUserParams) WithTimeout(timeout time.Duration) *RemoveUserParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the remove user params
func (o *RemoveUserParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the remove user params
func (o *RemoveUserParams) WithContext(ctx context.Context) *RemoveUserParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the remove user params
func (o *RemoveUserParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the remove user params
func (o *RemoveUserParams) WithHTTPClient(client *http.Client) *RemoveUserParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the remove user params
func (o *RemoveUserParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUser adds the user to the remove user params
func (o *RemoveUserParams) WithUser(user string) *RemoveUserParams {
	o.SetUser(user)
	return o
}

// SetUser adds the user to the remove user params
func (o *RemoveUserParams) SetUser(user string) {
	o.User = user
}

// WriteToRequest writes these params to a swagger request
func (o *RemoveUserParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param user
	if err := r.SetPathParam("user", o.User); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
