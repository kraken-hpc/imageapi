// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ContainerScriptHook Describes a container script hook point with execution controls.
//
// Scripts will be executed in array order after any default scripts.
//
//
// swagger:model container_script_hook
type ContainerScriptHook struct {

	// Disable default script hooks.
	DisableDefaults *bool `json:"disable_defaults,omitempty"`

	// scripts
	Scripts []*ContainerScript `json:"scripts"`
}

// Validate validates this container script hook
func (m *ContainerScriptHook) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateScripts(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ContainerScriptHook) validateScripts(formats strfmt.Registry) error {
	if swag.IsZero(m.Scripts) { // not required
		return nil
	}

	for i := 0; i < len(m.Scripts); i++ {
		if swag.IsZero(m.Scripts[i]) { // not required
			continue
		}

		if m.Scripts[i] != nil {
			if err := m.Scripts[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("scripts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this container script hook based on the context it is used
func (m *ContainerScriptHook) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateScripts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ContainerScriptHook) contextValidateScripts(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Scripts); i++ {

		if m.Scripts[i] != nil {
			if err := m.Scripts[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("scripts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ContainerScriptHook) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ContainerScriptHook) UnmarshalBinary(b []byte) error {
	var res ContainerScriptHook
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
