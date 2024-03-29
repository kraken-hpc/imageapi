// Code generated by go-swagger; DO NOT EDIT.

package attachments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewListAttachmentsParams creates a new ListAttachmentsParams object
//
// There are no default values defined in the spec.
func NewListAttachmentsParams() ListAttachmentsParams {

	return ListAttachmentsParams{}
}

// ListAttachmentsParams contains all the bound params for the list attachments operation
// typically these are obtained from a http.Request
//
// swagger:parameters list_attachments
type ListAttachmentsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*ID of a single attachment to query.
	  In: query
	*/
	ID *int64
	/*Kind of attachments to query.
	  In: query
	*/
	Kind *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewListAttachmentsParams() beforehand.
func (o *ListAttachmentsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qID, qhkID, _ := qs.GetOK("id")
	if err := o.bindID(qID, qhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	qKind, qhkKind, _ := qs.GetOK("kind")
	if err := o.bindKind(qKind, qhkKind, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindID binds and validates parameter ID from query.
func (o *ListAttachmentsParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("id", "query", "int64", raw)
	}
	o.ID = &value

	return nil
}

// bindKind binds and validates parameter Kind from query.
func (o *ListAttachmentsParams) bindKind(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Kind = &raw

	if err := o.validateKind(formats); err != nil {
		return err
	}

	return nil
}

// validateKind carries on validations for parameter Kind
func (o *ListAttachmentsParams) validateKind(formats strfmt.Registry) error {

	if err := validate.EnumCase("kind", "query", *o.Kind, []interface{}{"iscsi", "local", "loopback", "rbd"}, true); err != nil {
		return err
	}

	return nil
}
