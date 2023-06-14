// Code generated by go-swagger; DO NOT EDIT.

package users_roles

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
)

// AddRoleReader is a Reader for the AddRole structure.
type AddRoleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddRoleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewAddRoleNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewAddRoleDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAddRoleNoContent creates a AddRoleNoContent with default headers values
func NewAddRoleNoContent() *AddRoleNoContent {
	return &AddRoleNoContent{}
}

/*
AddRoleNoContent describes a response with status code 204, with default header values.

Null response
*/
type AddRoleNoContent struct {
}

// IsSuccess returns true when this add role no content response has a 2xx status code
func (o *AddRoleNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this add role no content response has a 3xx status code
func (o *AddRoleNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this add role no content response has a 4xx status code
func (o *AddRoleNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this add role no content response has a 5xx status code
func (o *AddRoleNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this add role no content response a status code equal to that given
func (o *AddRoleNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the add role no content response
func (o *AddRoleNoContent) Code() int {
	return 204
}

func (o *AddRoleNoContent) Error() string {
	return fmt.Sprintf("[POST /admin/users/{user}/roles][%d] addRoleNoContent ", 204)
}

func (o *AddRoleNoContent) String() string {
	return fmt.Sprintf("[POST /admin/users/{user}/roles][%d] addRoleNoContent ", 204)
}

func (o *AddRoleNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewAddRoleDefault creates a AddRoleDefault with default headers values
func NewAddRoleDefault(code int) *AddRoleDefault {
	return &AddRoleDefault{
		_statusCode: code,
	}
}

/*
AddRoleDefault describes a response with status code -1, with default header values.

unexpected error
*/
type AddRoleDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this add role default response has a 2xx status code
func (o *AddRoleDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this add role default response has a 3xx status code
func (o *AddRoleDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this add role default response has a 4xx status code
func (o *AddRoleDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this add role default response has a 5xx status code
func (o *AddRoleDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this add role default response a status code equal to that given
func (o *AddRoleDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the add role default response
func (o *AddRoleDefault) Code() int {
	return o._statusCode
}

func (o *AddRoleDefault) Error() string {
	return fmt.Sprintf("[POST /admin/users/{user}/roles][%d] addRole default  %+v", o._statusCode, o.Payload)
}

func (o *AddRoleDefault) String() string {
	return fmt.Sprintf("[POST /admin/users/{user}/roles][%d] addRole default  %+v", o._statusCode, o.Payload)
}

func (o *AddRoleDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *AddRoleDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
