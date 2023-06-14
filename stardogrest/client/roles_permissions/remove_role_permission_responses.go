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

// RemoveRolePermissionReader is a Reader for the RemoveRolePermission structure.
type RemoveRolePermissionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RemoveRolePermissionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewRemoveRolePermissionCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRemoveRolePermissionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRemoveRolePermissionCreated creates a RemoveRolePermissionCreated with default headers values
func NewRemoveRolePermissionCreated() *RemoveRolePermissionCreated {
	return &RemoveRolePermissionCreated{}
}

/*
RemoveRolePermissionCreated describes a response with status code 201, with default header values.

Null response
*/
type RemoveRolePermissionCreated struct {
}

// IsSuccess returns true when this remove role permission created response has a 2xx status code
func (o *RemoveRolePermissionCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this remove role permission created response has a 3xx status code
func (o *RemoveRolePermissionCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this remove role permission created response has a 4xx status code
func (o *RemoveRolePermissionCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this remove role permission created response has a 5xx status code
func (o *RemoveRolePermissionCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this remove role permission created response a status code equal to that given
func (o *RemoveRolePermissionCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the remove role permission created response
func (o *RemoveRolePermissionCreated) Code() int {
	return 201
}

func (o *RemoveRolePermissionCreated) Error() string {
	return fmt.Sprintf("[POST /admin/permissions/role/{role}/delete][%d] removeRolePermissionCreated ", 201)
}

func (o *RemoveRolePermissionCreated) String() string {
	return fmt.Sprintf("[POST /admin/permissions/role/{role}/delete][%d] removeRolePermissionCreated ", 201)
}

func (o *RemoveRolePermissionCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewRemoveRolePermissionDefault creates a RemoveRolePermissionDefault with default headers values
func NewRemoveRolePermissionDefault(code int) *RemoveRolePermissionDefault {
	return &RemoveRolePermissionDefault{
		_statusCode: code,
	}
}

/*
RemoveRolePermissionDefault describes a response with status code -1, with default header values.

unexpected error
*/
type RemoveRolePermissionDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this remove role permission default response has a 2xx status code
func (o *RemoveRolePermissionDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this remove role permission default response has a 3xx status code
func (o *RemoveRolePermissionDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this remove role permission default response has a 4xx status code
func (o *RemoveRolePermissionDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this remove role permission default response has a 5xx status code
func (o *RemoveRolePermissionDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this remove role permission default response a status code equal to that given
func (o *RemoveRolePermissionDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the remove role permission default response
func (o *RemoveRolePermissionDefault) Code() int {
	return o._statusCode
}

func (o *RemoveRolePermissionDefault) Error() string {
	return fmt.Sprintf("[POST /admin/permissions/role/{role}/delete][%d] removeRolePermission default  %+v", o._statusCode, o.Payload)
}

func (o *RemoveRolePermissionDefault) String() string {
	return fmt.Sprintf("[POST /admin/permissions/role/{role}/delete][%d] removeRolePermission default  %+v", o._statusCode, o.Payload)
}

func (o *RemoveRolePermissionDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *RemoveRolePermissionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
