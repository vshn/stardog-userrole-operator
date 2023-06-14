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

// ListRoleUsersReader is a Reader for the ListRoleUsers structure.
type ListRoleUsersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListRoleUsersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListRoleUsersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListRoleUsersDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListRoleUsersOK creates a ListRoleUsersOK with default headers values
func NewListRoleUsersOK() *ListRoleUsersOK {
	return &ListRoleUsersOK{}
}

/*
ListRoleUsersOK describes a response with status code 200, with default header values.

The users assigned to the role
*/
type ListRoleUsersOK struct {
	Payload *models.Users
}

// IsSuccess returns true when this list role users o k response has a 2xx status code
func (o *ListRoleUsersOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list role users o k response has a 3xx status code
func (o *ListRoleUsersOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list role users o k response has a 4xx status code
func (o *ListRoleUsersOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list role users o k response has a 5xx status code
func (o *ListRoleUsersOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list role users o k response a status code equal to that given
func (o *ListRoleUsersOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the list role users o k response
func (o *ListRoleUsersOK) Code() int {
	return 200
}

func (o *ListRoleUsersOK) Error() string {
	return fmt.Sprintf("[GET /admin/roles/{role}/users][%d] listRoleUsersOK  %+v", 200, o.Payload)
}

func (o *ListRoleUsersOK) String() string {
	return fmt.Sprintf("[GET /admin/roles/{role}/users][%d] listRoleUsersOK  %+v", 200, o.Payload)
}

func (o *ListRoleUsersOK) GetPayload() *models.Users {
	return o.Payload
}

func (o *ListRoleUsersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Users)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListRoleUsersDefault creates a ListRoleUsersDefault with default headers values
func NewListRoleUsersDefault(code int) *ListRoleUsersDefault {
	return &ListRoleUsersDefault{
		_statusCode: code,
	}
}

/*
ListRoleUsersDefault describes a response with status code -1, with default header values.

unexpected error
*/
type ListRoleUsersDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this list role users default response has a 2xx status code
func (o *ListRoleUsersDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list role users default response has a 3xx status code
func (o *ListRoleUsersDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list role users default response has a 4xx status code
func (o *ListRoleUsersDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list role users default response has a 5xx status code
func (o *ListRoleUsersDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list role users default response a status code equal to that given
func (o *ListRoleUsersDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the list role users default response
func (o *ListRoleUsersDefault) Code() int {
	return o._statusCode
}

func (o *ListRoleUsersDefault) Error() string {
	return fmt.Sprintf("[GET /admin/roles/{role}/users][%d] listRoleUsers default  %+v", o._statusCode, o.Payload)
}

func (o *ListRoleUsersDefault) String() string {
	return fmt.Sprintf("[GET /admin/roles/{role}/users][%d] listRoleUsers default  %+v", o._statusCode, o.Payload)
}

func (o *ListRoleUsersDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *ListRoleUsersDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
