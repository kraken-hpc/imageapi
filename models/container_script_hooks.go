// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ContainerScriptHooks Container script execution hooks.
//
// We currently provide 4 hook points:
// 1. `create` is executed on container creation (root namespaces).
// 2. `init` is executed in the container namespaces before the provided `init` is called (container namespaces).
// 3. `exit` is executed on container exit (root namespaces).
// 4. `delete` is executed on container deletion (root namespaces).
//
//
// swagger:model container_script_hooks
type ContainerScriptHooks struct {

	// create
	Create *ContainerScriptHook `json:"create,omitempty"`

	// delete
	Delete *ContainerScriptHook `json:"delete,omitempty"`

	// exit
	Exit *ContainerScriptHook `json:"exit,omitempty"`

	// init
	Init *ContainerScriptHook `json:"init,omitempty"`
}

// Validate validates this container script hooks
func (m *ContainerScriptHooks) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDelete(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExit(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInit(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ContainerScriptHooks) validateCreate(formats strfmt.Registry) error {
	if swag.IsZero(m.Create) { // not required
		return nil
	}

	if m.Create != nil {
		if err := m.Create.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("create")
			}
			return err
		}
	}

	return nil
}

func (m *ContainerScriptHooks) validateDelete(formats strfmt.Registry) error {
	if swag.IsZero(m.Delete) { // not required
		return nil
	}

	if m.Delete != nil {
		if err := m.Delete.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("delete")
			}
			return err
		}
	}

	return nil
}

func (m *ContainerScriptHooks) validateExit(formats strfmt.Registry) error {
	if swag.IsZero(m.Exit) { // not required
		return nil
	}

	if m.Exit != nil {
		if err := m.Exit.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("exit")
			}
			return err
		}
	}

	return nil
}

func (m *ContainerScriptHooks) validateInit(formats strfmt.Registry) error {
	if swag.IsZero(m.Init) { // not required
		return nil
	}

	if m.Init != nil {
		if err := m.Init.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("init")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this container script hooks based on the context it is used
func (m *ContainerScriptHooks) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCreate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDelete(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateExit(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateInit(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ContainerScriptHooks) contextValidateCreate(ctx context.Context, formats strfmt.Registry) error {

	if m.Create != nil {
		if err := m.Create.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("create")
			}
			return err
		}
	}

	return nil
}

func (m *ContainerScriptHooks) contextValidateDelete(ctx context.Context, formats strfmt.Registry) error {

	if m.Delete != nil {
		if err := m.Delete.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("delete")
			}
			return err
		}
	}

	return nil
}

func (m *ContainerScriptHooks) contextValidateExit(ctx context.Context, formats strfmt.Registry) error {

	if m.Exit != nil {
		if err := m.Exit.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("exit")
			}
			return err
		}
	}

	return nil
}

func (m *ContainerScriptHooks) contextValidateInit(ctx context.Context, formats strfmt.Registry) error {

	if m.Init != nil {
		if err := m.Init.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("init")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ContainerScriptHooks) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ContainerScriptHooks) UnmarshalBinary(b []byte) error {
	var res ContainerScriptHooks
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
