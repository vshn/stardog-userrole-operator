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

	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
)

// NewPutRolesParams creates a new PutRolesParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPutRolesParams() *PutRolesParams {
	return &PutRolesParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPutRolesParamsWithTimeout creates a new PutRolesParams object
// with the ability to set a timeout on a request.
func NewPutRolesParamsWithTimeout(timeout time.Duration) *PutRolesParams {
	return &PutRolesParams{
		timeout: timeout,
	}
}

// NewPutRolesParamsWithContext creates a new PutRolesParams object
// with the ability to set a context for a request.
func NewPutRolesParamsWithContext(ctx context.Context) *PutRolesParams {
	return &PutRolesParams{
		Context: ctx,
	}
}

// NewPutRolesParamsWithHTTPClient creates a new PutRolesParams object
// with the ability to set a custom HTTPClient for a request.
func NewPutRolesParamsWithHTTPClient(client *http.Client) *PutRolesParams {
	return &PutRolesParams{
		HTTPClient: client,
	}
}

/*
PutRolesParams contains all the parameters to send to the API endpoint

	for the put roles operation.

	Typically these are written to a http.Request.
*/
type PutRolesParams struct {

	/* Roles.

	   The new set of roles
	*/
	Roles *models.Roles

	/* User.

	   The username of the user whose roles should be changed
	*/
	User string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the put roles params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutRolesParams) WithDefaults() *PutRolesParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the put roles params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutRolesParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the put roles params
func (o *PutRolesParams) WithTimeout(timeout time.Duration) *PutRolesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put roles params
func (o *PutRolesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put roles params
func (o *PutRolesParams) WithContext(ctx context.Context) *PutRolesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put roles params
func (o *PutRolesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put roles params
func (o *PutRolesParams) WithHTTPClient(client *http.Client) *PutRolesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put roles params
func (o *PutRolesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRoles adds the roles to the put roles params
func (o *PutRolesParams) WithRoles(roles *models.Roles) *PutRolesParams {
	o.SetRoles(roles)
	return o
}

// SetRoles adds the roles to the put roles params
func (o *PutRolesParams) SetRoles(roles *models.Roles) {
	o.Roles = roles
}

// WithUser adds the user to the put roles params
func (o *PutRolesParams) WithUser(user string) *PutRolesParams {
	o.SetUser(user)
	return o
}

// SetUser adds the user to the put roles params
func (o *PutRolesParams) SetUser(user string) {
	o.User = user
}

// WriteToRequest writes these params to a swagger request
func (o *PutRolesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Roles != nil {
		if err := r.SetBodyParam(o.Roles); err != nil {
			return err
		}
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
