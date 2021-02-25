// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
)

// ID An ID is a unique numeric ID that references an object.
// IDs are not necessarily unique across object types.
// IDs are generall readOnly and generated internally.
//
//
// swagger:model id
type ID int64

// Validate validates this id
func (m ID) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this id based on context it is used
func (m ID) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}