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

// ListRolePermissionsReader is a Reader for the ListRolePermissions structure.
type ListRolePermissionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListRolePermissionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListRolePermissionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListRolePermissionsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListRolePermissionsOK creates a ListRolePermissionsOK with default headers values
func NewListRolePermissionsOK() *ListRolePermissionsOK {
	return &ListRolePermissionsOK{}
}

/*
ListRolePermissionsOK describes a response with status code 200, with default header values.

The roles's permissions
*/
type ListRolePermissionsOK struct {
	Payload *models.Permissions
}

// IsSuccess returns true when this list role permissions o k response has a 2xx status code
func (o *ListRolePermissionsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list role permissions o k response has a 3xx status code
func (o *ListRolePermissionsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list role permissions o k response has a 4xx status code
func (o *ListRolePermissionsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list role permissions o k response has a 5xx status code
func (o *ListRolePermissionsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list role permissions o k response a status code equal to that given
func (o *ListRolePermissionsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the list role permissions o k response
func (o *ListRolePermissionsOK) Code() int {
	return 200
}

func (o *ListRolePermissionsOK) Error() string {
	return fmt.Sprintf("[GET /permissions/role/{role}][%d] listRolePermissionsOK  %+v", 200, o.Payload)
}

func (o *ListRolePermissionsOK) String() string {
	return fmt.Sprintf("[GET /permissions/role/{role}][%d] listRolePermissionsOK  %+v", 200, o.Payload)
}

func (o *ListRolePermissionsOK) GetPayload() *models.Permissions {
	return o.Payload
}

func (o *ListRolePermissionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Permissions)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListRolePermissionsDefault creates a ListRolePermissionsDefault with default headers values
func NewListRolePermissionsDefault(code int) *ListRolePermissionsDefault {
	return &ListRolePermissionsDefault{
		_statusCode: code,
	}
}

/*
ListRolePermissionsDefault describes a response with status code -1, with default header values.

unexpected error
*/
type ListRolePermissionsDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this list role permissions default response has a 2xx status code
func (o *ListRolePermissionsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list role permissions default response has a 3xx status code
func (o *ListRolePermissionsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list role permissions default response has a 4xx status code
func (o *ListRolePermissionsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list role permissions default response has a 5xx status code
func (o *ListRolePermissionsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list role permissions default response a status code equal to that given
func (o *ListRolePermissionsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the list role permissions default response
func (o *ListRolePermissionsDefault) Code() int {
	return o._statusCode
}

func (o *ListRolePermissionsDefault) Error() string {
	return fmt.Sprintf("[GET /permissions/role/{role}][%d] listRolePermissions default  %+v", o._statusCode, o.Payload)
}

func (o *ListRolePermissionsDefault) String() string {
	return fmt.Sprintf("[GET /permissions/role/{role}][%d] listRolePermissions default  %+v", o._statusCode, o.Payload)
}

func (o *ListRolePermissionsDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *ListRolePermissionsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}