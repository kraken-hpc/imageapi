// Code generated by go-swagger; DO NOT EDIT.

package attachments

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

// NewAttachParams creates a new AttachParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewAttachParams() *AttachParams {
	return &AttachParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewAttachParamsWithTimeout creates a new AttachParams object
// with the ability to set a timeout on a request.
func NewAttachParamsWithTimeout(timeout time.Duration) *AttachParams {
	return &AttachParams{
		timeout: timeout,
	}
}

// NewAttachParamsWithContext creates a new AttachParams object
// with the ability to set a context for a request.
func NewAttachParamsWithContext(ctx context.Context) *AttachParams {
	return &AttachParams{
		Context: ctx,
	}
}

// NewAttachParamsWithHTTPClient creates a new AttachParams object
// with the ability to set a custom HTTPClient for a request.
func NewAttachParamsWithHTTPClient(client *http.Client) *AttachParams {
	return &AttachParams{
		HTTPClient: client,
	}
}

/* AttachParams contains all the parameters to send to the API endpoint
   for the attach operation.

   Typically these are written to a http.Request.
*/
type AttachParams struct {

	// Attach.
	Attach *models.Attach

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the attach params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AttachParams) WithDefaults() *AttachParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the attach params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AttachParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the attach params
func (o *AttachParams) WithTimeout(timeout time.Duration) *AttachParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the attach params
func (o *AttachParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the attach params
func (o *AttachParams) WithContext(ctx context.Context) *AttachParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the attach params
func (o *AttachParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the attach params
func (o *AttachParams) WithHTTPClient(client *http.Client) *AttachParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the attach params
func (o *AttachParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAttach adds the attach to the attach params
func (o *AttachParams) WithAttach(attach *models.Attach) *AttachParams {
	o.SetAttach(attach)
	return o
}

// SetAttach adds the attach to the attach params
func (o *AttachParams) SetAttach(attach *models.Attach) {
	o.Attach = attach
}

// WriteToRequest writes these params to a swagger request
func (o *AttachParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Attach != nil {
		if err := r.SetBodyParam(o.Attach); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}