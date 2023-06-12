// Code generated by go-swagger; DO NOT EDIT.

package roles_permissions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
)

// AddRolePermissionReader is a Reader for the AddRolePermission structure.
type AddRolePermissionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddRolePermissionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewAddRolePermissionCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewAddRolePermissionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAddRolePermissionCreated creates a AddRolePermissionCreated with default headers values
func NewAddRolePermissionCreated() *AddRolePermissionCreated {
	return &AddRolePermissionCreated{}
}

/*
AddRolePermissionCreated describes a response with status code 201, with default header values.

Null response
*/
type AddRolePermissionCreated struct {
}

// IsSuccess returns true when this add role permission created response has a 2xx status code
func (o *AddRolePermissionCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this add role permission created response has a 3xx status code
func (o *AddRolePermissionCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this add role permission created response has a 4xx status code
func (o *AddRolePermissionCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this add role permission created response has a 5xx status code
func (o *AddRolePermissionCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this add role permission created response a status code equal to that given
func (o *AddRolePermissionCreated) IsCode(code int) bool {
	return code == 201
}

func (o *AddRolePermissionCreated) Error() string {
	return fmt.Sprintf("[PUT /permissions/role/{role}][%d] addRolePermissionCreated ", 201)
}

func (o *AddRolePermissionCreated) String() string {
	return fmt.Sprintf("[PUT /permissions/role/{role}][%d] addRolePermissionCreated ", 201)
}

func (o *AddRolePermissionCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewAddRolePermissionDefault creates a AddRolePermissionDefault with default headers values
func NewAddRolePermissionDefault(code int) *AddRolePermissionDefault {
	return &AddRolePermissionDefault{
		_statusCode: code,
	}
}

/*
AddRolePermissionDefault describes a response with status code -1, with default header values.

unexpected error
*/
type AddRolePermissionDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the add role permission default response
func (o *AddRolePermissionDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this add role permission default response has a 2xx status code
func (o *AddRolePermissionDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this add role permission default response has a 3xx status code
func (o *AddRolePermissionDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this add role permission default response has a 4xx status code
func (o *AddRolePermissionDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this add role permission default response has a 5xx status code
func (o *AddRolePermissionDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this add role permission default response a status code equal to that given
func (o *AddRolePermissionDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *AddRolePermissionDefault) Error() string {
	return fmt.Sprintf("[PUT /permissions/role/{role}][%d] addRolePermission default  %+v", o._statusCode, o.Payload)
}

func (o *AddRolePermissionDefault) String() string {
	return fmt.Sprintf("[PUT /permissions/role/{role}][%d] addRolePermission default  %+v", o._statusCode, o.Payload)
}

func (o *AddRolePermissionDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *AddRolePermissionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
