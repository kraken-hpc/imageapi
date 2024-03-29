// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AttachLocal `attach_local` describes a block device that is locally present.
// This can be used to get a reference to a local disk, for instance.
//
// Local only supports finding device files on the local (root) system.
// It only takes one parameter: the path to the device file.
//
//
// swagger:model attach_local
type AttachLocal struct {

	// A unix-formatted filesystem path pointing to a block device file.
	// Required: true
	Path *string `json:"path"`
}

// Validate validates this attach local
func (m *AttachLocal) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePath(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AttachLocal) validatePath(formats strfmt.Registry) error {

	if err := validate.Required("path", "body", m.Path); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this attach local based on context it is used
func (m *AttachLocal) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AttachLocal) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AttachLocal) UnmarshalBinary(b []byte) error {
	var res AttachLocal
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
