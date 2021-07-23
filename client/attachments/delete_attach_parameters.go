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
	"github.com/go-openapi/swag"
)

// NewDeleteAttachParams creates a new DeleteAttachParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteAttachParams() *DeleteAttachParams {
	return &DeleteAttachParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteAttachParamsWithTimeout creates a new DeleteAttachParams object
// with the ability to set a timeout on a request.
func NewDeleteAttachParamsWithTimeout(timeout time.Duration) *DeleteAttachParams {
	return &DeleteAttachParams{
		timeout: timeout,
	}
}

// NewDeleteAttachParamsWithContext creates a new DeleteAttachParams object
// with the ability to set a context for a request.
func NewDeleteAttachParamsWithContext(ctx context.Context) *DeleteAttachParams {
	return &DeleteAttachParams{
		Context: ctx,
	}
}

// NewDeleteAttachParamsWithHTTPClient creates a new DeleteAttachParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteAttachParamsWithHTTPClient(client *http.Client) *DeleteAttachParams {
	return &DeleteAttachParams{
		HTTPClient: client,
	}
}

/* DeleteAttachParams contains all the parameters to send to the API endpoint
   for the delete attach operation.

   Typically these are written to a http.Request.
*/
type DeleteAttachParams struct {

	/* Force.

	   Force deletion
	*/
	Force *bool

	// ID.
	//
	// Format: int64
	ID int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete attach params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteAttachParams) WithDefaults() *DeleteAttachParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete attach params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteAttachParams) SetDefaults() {
	var (
		forceDefault = bool(false)
	)

	val := DeleteAttachParams{
		Force: &forceDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the delete attach params
func (o *DeleteAttachParams) WithTimeout(timeout time.Duration) *DeleteAttachParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete attach params
func (o *DeleteAttachParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete attach params
func (o *DeleteAttachParams) WithContext(ctx context.Context) *DeleteAttachParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete attach params
func (o *DeleteAttachParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete attach params
func (o *DeleteAttachParams) WithHTTPClient(client *http.Client) *DeleteAttachParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete attach params
func (o *DeleteAttachParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithForce adds the force to the delete attach params
func (o *DeleteAttachParams) WithForce(force *bool) *DeleteAttachParams {
	o.SetForce(force)
	return o
}

// SetForce adds the force to the delete attach params
func (o *DeleteAttachParams) SetForce(force *bool) {
	o.Force = force
}

// WithID adds the id to the delete attach params
func (o *DeleteAttachParams) WithID(id int64) *DeleteAttachParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete attach params
func (o *DeleteAttachParams) SetID(id int64) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteAttachParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Force != nil {

		// query param force
		var qrForce bool

		if o.Force != nil {
			qrForce = *o.Force
		}
		qForce := swag.FormatBool(qrForce)
		if qForce != "" {

			if err := r.SetQueryParam("force", qForce); err != nil {
				return err
			}
		}
	}

	// query param id
	qrID := o.ID
	qID := swag.FormatInt64(qrID)
	if qID != "" {

		if err := r.SetQueryParam("id", qID); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
