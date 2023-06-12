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

// Superuser superuser
//
// swagger:model Superuser
type Superuser struct {

	// superuser
	// Required: true
	Superuser *bool `json:"superuser"`
}

// Validate validates this superuser
func (m *Superuser) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSuperuser(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Superuser) validateSuperuser(formats strfmt.Registry) error {

	if err := validate.Required("superuser", "body", m.Superuser); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this superuser based on context it is used
func (m *Superuser) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Superuser) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Superuser) UnmarshalBinary(b []byte) error {
	var res Superuser
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
