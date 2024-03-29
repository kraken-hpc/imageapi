// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// MountBind `mount_bind` describes a local bind mount.
// Bind mounts can be relative to another mount, or to /, allowing a way to access local data.
//
//
// swagger:model mount_bind
type MountBind struct {

	// base determines the relative root for the path.  There are two options:
	// `root` means to use the current root (`/`) as the base path.
	// `mount` means to use a mount as the base path. If this is specified, `mount` must be specified as well.
	//
	// Required: true
	// Enum: [root mount]
	Base *string `json:"base"`

	// mount
	Mount *Mount `json:"mount,omitempty"`

	// A unix-formatted filesystem path with `/` relative to the respective base.
	// Required: true
	Path *string `json:"path"`

	// perform a recursive bind mount
	Recursive *bool `json:"recursive,omitempty"`

	// mount read-only
	Ro *bool `json:"ro,omitempty"`
}

// Validate validates this mount bind
func (m *MountBind) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBase(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMount(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePath(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var mountBindTypeBasePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["root","mount"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		mountBindTypeBasePropEnum = append(mountBindTypeBasePropEnum, v)
	}
}

const (

	// MountBindBaseRoot captures enum value "root"
	MountBindBaseRoot string = "root"

	// MountBindBaseMount captures enum value "mount"
	MountBindBaseMount string = "mount"
)

// prop value enum
func (m *MountBind) validateBaseEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, mountBindTypeBasePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *MountBind) validateBase(formats strfmt.Registry) error {

	if err := validate.Required("base", "body", m.Base); err != nil {
		return err
	}

	// value enum
	if err := m.validateBaseEnum("base", "body", *m.Base); err != nil {
		return err
	}

	return nil
}

func (m *MountBind) validateMount(formats strfmt.Registry) error {
	if swag.IsZero(m.Mount) { // not required
		return nil
	}

	if m.Mount != nil {
		if err := m.Mount.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("mount")
			}
			return err
		}
	}

	return nil
}

func (m *MountBind) validatePath(formats strfmt.Registry) error {

	if err := validate.Required("path", "body", m.Path); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this mount bind based on the context it is used
func (m *MountBind) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMount(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MountBind) contextValidateMount(ctx context.Context, formats strfmt.Registry) error {

	if m.Mount != nil {
		if err := m.Mount.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("mount")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *MountBind) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MountBind) UnmarshalBinary(b []byte) error {
	var res MountBind
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
