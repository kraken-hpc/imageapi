// Code generated by go-swagger; DO NOT EDIT.

package containers

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

	"github.com/kraken-hpc/imageapi/models"
)

// NewCreateContainerParams creates a new CreateContainerParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateContainerParams() *CreateContainerParams {
	return &CreateContainerParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateContainerParamsWithTimeout creates a new CreateContainerParams object
// with the ability to set a timeout on a request.
func NewCreateContainerParamsWithTimeout(timeout time.Duration) *CreateContainerParams {
	return &CreateContainerParams{
		timeout: timeout,
	}
}

// NewCreateContainerParamsWithContext creates a new CreateContainerParams object
// with the ability to set a context for a request.
func NewCreateContainerParamsWithContext(ctx context.Context) *CreateContainerParams {
	return &CreateContainerParams{
		Context: ctx,
	}
}

// NewCreateContainerParamsWithHTTPClient creates a new CreateContainerParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateContainerParamsWithHTTPClient(client *http.Client) *CreateContainerParams {
	return &CreateContainerParams{
		HTTPClient: client,
	}
}

/* CreateContainerParams contains all the parameters to send to the API endpoint
   for the create container operation.

   Typically these are written to a http.Request.
*/
type CreateContainerParams struct {

	// Container.
	Container *models.Container

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create container params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateContainerParams) WithDefaults() *CreateContainerParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create container params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateContainerParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create container params
func (o *CreateContainerParams) WithTimeout(timeout time.Duration) *CreateContainerParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create container params
func (o *CreateContainerParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create container params
func (o *CreateContainerParams) WithContext(ctx context.Context) *CreateContainerParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create container params
func (o *CreateContainerParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create container params
func (o *CreateContainerParams) WithHTTPClient(client *http.Client) *CreateContainerParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create container params
func (o *CreateContainerParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithContainer adds the container to the create container params
func (o *CreateContainerParams) WithContainer(container *models.Container) *CreateContainerParams {
	o.SetContainer(container)
	return o
}

// SetContainer adds the container to the create container params
func (o *CreateContainerParams) SetContainer(container *models.Container) {
	o.Container = container
}

// WriteToRequest writes these params to a swagger request
func (o *CreateContainerParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Container != nil {
		if err := r.SetBodyParam(o.Container); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
