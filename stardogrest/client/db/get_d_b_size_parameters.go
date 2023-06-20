// Code generated by go-swagger; DO NOT EDIT.

package db

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
	"github.com/go-openapi/swag"
)

// NewGetDBSizeParams creates a new GetDBSizeParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetDBSizeParams() *GetDBSizeParams {
	return &GetDBSizeParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetDBSizeParamsWithTimeout creates a new GetDBSizeParams object
// with the ability to set a timeout on a request.
func NewGetDBSizeParamsWithTimeout(timeout time.Duration) *GetDBSizeParams {
	return &GetDBSizeParams{
		timeout: timeout,
	}
}

// NewGetDBSizeParamsWithContext creates a new GetDBSizeParams object
// with the ability to set a context for a request.
func NewGetDBSizeParamsWithContext(ctx context.Context) *GetDBSizeParams {
	return &GetDBSizeParams{
		Context: ctx,
	}
}

// NewGetDBSizeParamsWithHTTPClient creates a new GetDBSizeParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetDBSizeParamsWithHTTPClient(client *http.Client) *GetDBSizeParams {
	return &GetDBSizeParams{
		HTTPClient: client,
	}
}

/*
GetDBSizeParams contains all the parameters to send to the API endpoint

	for the get d b size operation.

	Typically these are written to a http.Request.
*/
type GetDBSizeParams struct {

	/* Db.

	   Database name
	*/
	Db string

	/* Exact.

	   Whether to request that the database size be exact instead of approximate
	*/
	Exact *bool

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get d b size params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDBSizeParams) WithDefaults() *GetDBSizeParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get d b size params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDBSizeParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get d b size params
func (o *GetDBSizeParams) WithTimeout(timeout time.Duration) *GetDBSizeParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get d b size params
func (o *GetDBSizeParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get d b size params
func (o *GetDBSizeParams) WithContext(ctx context.Context) *GetDBSizeParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get d b size params
func (o *GetDBSizeParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get d b size params
func (o *GetDBSizeParams) WithHTTPClient(client *http.Client) *GetDBSizeParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get d b size params
func (o *GetDBSizeParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDb adds the db to the get d b size params
func (o *GetDBSizeParams) WithDb(db string) *GetDBSizeParams {
	o.SetDb(db)
	return o
}

// SetDb adds the db to the get d b size params
func (o *GetDBSizeParams) SetDb(db string) {
	o.Db = db
}

// WithExact adds the exact to the get d b size params
func (o *GetDBSizeParams) WithExact(exact *bool) *GetDBSizeParams {
	o.SetExact(exact)
	return o
}

// SetExact adds the exact to the get d b size params
func (o *GetDBSizeParams) SetExact(exact *bool) {
	o.Exact = exact
}

// WriteToRequest writes these params to a swagger request
func (o *GetDBSizeParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param db
	if err := r.SetPathParam("db", o.Db); err != nil {
		return err
	}

	if o.Exact != nil {

		// query param exact
		var qrExact bool

		if o.Exact != nil {
			qrExact = *o.Exact
		}
		qExact := swag.FormatBool(qrExact)
		if qExact != "" {

			if err := r.SetQueryParam("exact", qExact); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
