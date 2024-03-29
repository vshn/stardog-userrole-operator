// Code generated by go-swagger; DO NOT EDIT.

package users_permissions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
)

// ListEffectivePermissionsReader is a Reader for the ListEffectivePermissions structure.
type ListEffectivePermissionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListEffectivePermissionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListEffectivePermissionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListEffectivePermissionsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListEffectivePermissionsOK creates a ListEffectivePermissionsOK with default headers values
func NewListEffectivePermissionsOK() *ListEffectivePermissionsOK {
	return &ListEffectivePermissionsOK{}
}

/*
ListEffectivePermissionsOK describes a response with status code 200, with default header values.

The user's permissions
*/
type ListEffectivePermissionsOK struct {
	Payload *models.Permissions
}

// IsSuccess returns true when this list effective permissions o k response has a 2xx status code
func (o *ListEffectivePermissionsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list effective permissions o k response has a 3xx status code
func (o *ListEffectivePermissionsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list effective permissions o k response has a 4xx status code
func (o *ListEffectivePermissionsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list effective permissions o k response has a 5xx status code
func (o *ListEffectivePermissionsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list effective permissions o k response a status code equal to that given
func (o *ListEffectivePermissionsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the list effective permissions o k response
func (o *ListEffectivePermissionsOK) Code() int {
	return 200
}

func (o *ListEffectivePermissionsOK) Error() string {
	return fmt.Sprintf("[GET /admin/permissions/effective/user/{user}][%d] listEffectivePermissionsOK  %+v", 200, o.Payload)
}

func (o *ListEffectivePermissionsOK) String() string {
	return fmt.Sprintf("[GET /admin/permissions/effective/user/{user}][%d] listEffectivePermissionsOK  %+v", 200, o.Payload)
}

func (o *ListEffectivePermissionsOK) GetPayload() *models.Permissions {
	return o.Payload
}

func (o *ListEffectivePermissionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Permissions)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListEffectivePermissionsDefault creates a ListEffectivePermissionsDefault with default headers values
func NewListEffectivePermissionsDefault(code int) *ListEffectivePermissionsDefault {
	return &ListEffectivePermissionsDefault{
		_statusCode: code,
	}
}

/*
ListEffectivePermissionsDefault describes a response with status code -1, with default header values.

unexpected error
*/
type ListEffectivePermissionsDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this list effective permissions default response has a 2xx status code
func (o *ListEffectivePermissionsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list effective permissions default response has a 3xx status code
func (o *ListEffectivePermissionsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list effective permissions default response has a 4xx status code
func (o *ListEffectivePermissionsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list effective permissions default response has a 5xx status code
func (o *ListEffectivePermissionsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list effective permissions default response a status code equal to that given
func (o *ListEffectivePermissionsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the list effective permissions default response
func (o *ListEffectivePermissionsDefault) Code() int {
	return o._statusCode
}

func (o *ListEffectivePermissionsDefault) Error() string {
	return fmt.Sprintf("[GET /admin/permissions/effective/user/{user}][%d] listEffectivePermissions default  %+v", o._statusCode, o.Payload)
}

func (o *ListEffectivePermissionsDefault) String() string {
	return fmt.Sprintf("[GET /admin/permissions/effective/user/{user}][%d] listEffectivePermissions default  %+v", o._statusCode, o.Payload)
}

func (o *ListEffectivePermissionsDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *ListEffectivePermissionsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
