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

// NewIsSuperuserParams creates a new IsSuperuserParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewIsSuperuserParams() *IsSuperuserParams {
	return &IsSuperuserParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewIsSuperuserParamsWithTimeout creates a new IsSuperuserParams object
// with the ability to set a timeout on a request.
func NewIsSuperuserParamsWithTimeout(timeout time.Duration) *IsSuperuserParams {
	return &IsSuperuserParams{
		timeout: timeout,
	}
}

// NewIsSuperuserParamsWithContext creates a new IsSuperuserParams object
// with the ability to set a context for a request.
func NewIsSuperuserParamsWithContext(ctx context.Context) *IsSuperuserParams {
	return &IsSuperuserParams{
		Context: ctx,
	}
}

// NewIsSuperuserParamsWithHTTPClient creates a new IsSuperuserParams object
// with the ability to set a custom HTTPClient for a request.
func NewIsSuperuserParamsWithHTTPClient(client *http.Client) *IsSuperuserParams {
	return &IsSuperuserParams{
		HTTPClient: client,
	}
}

/*
IsSuperuserParams contains all the parameters to send to the API endpoint

	for the is superuser operation.

	Typically these are written to a http.Request.
*/
type IsSuperuserParams struct {

	/* User.

	   The username of the user whose status should be queried
	*/
	User string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the is superuser params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *IsSuperuserParams) WithDefaults() *IsSuperuserParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the is superuser params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *IsSuperuserParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the is superuser params
func (o *IsSuperuserParams) WithTimeout(timeout time.Duration) *IsSuperuserParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the is superuser params
func (o *IsSuperuserParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the is superuser params
func (o *IsSuperuserParams) WithContext(ctx context.Context) *IsSuperuserParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the is superuser params
func (o *IsSuperuserParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the is superuser params
func (o *IsSuperuserParams) WithHTTPClient(client *http.Client) *IsSuperuserParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the is superuser params
func (o *IsSuperuserParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUser adds the user to the is superuser params
func (o *IsSuperuserParams) WithUser(user string) *IsSuperuserParams {
	o.SetUser(user)
	return o
}

// SetUser adds the user to the is superuser params
func (o *IsSuperuserParams) SetUser(user string) {
	o.User = user
}

// WriteToRequest writes these params to a swagger request
func (o *IsSuperuserParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
