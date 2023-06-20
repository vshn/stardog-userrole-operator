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

// RemoveRoleOfUserReader is a Reader for the RemoveRoleOfUser structure.
type RemoveRoleOfUserReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RemoveRoleOfUserReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewRemoveRoleOfUserNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewRemoveRoleOfUserNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewRemoveRoleOfUserDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRemoveRoleOfUserNoContent creates a RemoveRoleOfUserNoContent with default headers values
func NewRemoveRoleOfUserNoContent() *RemoveRoleOfUserNoContent {
	return &RemoveRoleOfUserNoContent{}
}

/*
RemoveRoleOfUserNoContent describes a response with status code 204, with default header values.

Null response
*/
type RemoveRoleOfUserNoContent struct {
}

// IsSuccess returns true when this remove role of user no content response has a 2xx status code
func (o *RemoveRoleOfUserNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this remove role of user no content response has a 3xx status code
func (o *RemoveRoleOfUserNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this remove role of user no content response has a 4xx status code
func (o *RemoveRoleOfUserNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this remove role of user no content response has a 5xx status code
func (o *RemoveRoleOfUserNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this remove role of user no content response a status code equal to that given
func (o *RemoveRoleOfUserNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the remove role of user no content response
func (o *RemoveRoleOfUserNoContent) Code() int {
	return 204
}

func (o *RemoveRoleOfUserNoContent) Error() string {
	return fmt.Sprintf("[DELETE /admin/users/{user}/roles/{role}][%d] removeRoleOfUserNoContent ", 204)
}

func (o *RemoveRoleOfUserNoContent) String() string {
	return fmt.Sprintf("[DELETE /admin/users/{user}/roles/{role}][%d] removeRoleOfUserNoContent ", 204)
}

func (o *RemoveRoleOfUserNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewRemoveRoleOfUserNotFound creates a RemoveRoleOfUserNotFound with default headers values
func NewRemoveRoleOfUserNotFound() *RemoveRoleOfUserNotFound {
	return &RemoveRoleOfUserNotFound{}
}

/*
RemoveRoleOfUserNotFound describes a response with status code 404, with default header values.

Role user does not exist
*/
type RemoveRoleOfUserNotFound struct {
	Payload *models.NotExists
}

// IsSuccess returns true when this remove role of user not found response has a 2xx status code
func (o *RemoveRoleOfUserNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this remove role of user not found response has a 3xx status code
func (o *RemoveRoleOfUserNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this remove role of user not found response has a 4xx status code
func (o *RemoveRoleOfUserNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this remove role of user not found response has a 5xx status code
func (o *RemoveRoleOfUserNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this remove role of user not found response a status code equal to that given
func (o *RemoveRoleOfUserNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the remove role of user not found response
func (o *RemoveRoleOfUserNotFound) Code() int {
	return 404
}

func (o *RemoveRoleOfUserNotFound) Error() string {
	return fmt.Sprintf("[DELETE /admin/users/{user}/roles/{role}][%d] removeRoleOfUserNotFound  %+v", 404, o.Payload)
}

func (o *RemoveRoleOfUserNotFound) String() string {
	return fmt.Sprintf("[DELETE /admin/users/{user}/roles/{role}][%d] removeRoleOfUserNotFound  %+v", 404, o.Payload)
}

func (o *RemoveRoleOfUserNotFound) GetPayload() *models.NotExists {
	return o.Payload
}

func (o *RemoveRoleOfUserNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.NotExists)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRemoveRoleOfUserDefault creates a RemoveRoleOfUserDefault with default headers values
func NewRemoveRoleOfUserDefault(code int) *RemoveRoleOfUserDefault {
	return &RemoveRoleOfUserDefault{
		_statusCode: code,
	}
}

/*
RemoveRoleOfUserDefault describes a response with status code -1, with default header values.

unexpected error
*/
type RemoveRoleOfUserDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this remove role of user default response has a 2xx status code
func (o *RemoveRoleOfUserDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this remove role of user default response has a 3xx status code
func (o *RemoveRoleOfUserDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this remove role of user default response has a 4xx status code
func (o *RemoveRoleOfUserDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this remove role of user default response has a 5xx status code
func (o *RemoveRoleOfUserDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this remove role of user default response a status code equal to that given
func (o *RemoveRoleOfUserDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the remove role of user default response
func (o *RemoveRoleOfUserDefault) Code() int {
	return o._statusCode
}

func (o *RemoveRoleOfUserDefault) Error() string {
	return fmt.Sprintf("[DELETE /admin/users/{user}/roles/{role}][%d] removeRoleOfUser default  %+v", o._statusCode, o.Payload)
}

func (o *RemoveRoleOfUserDefault) String() string {
	return fmt.Sprintf("[DELETE /admin/users/{user}/roles/{role}][%d] removeRoleOfUser default  %+v", o._statusCode, o.Payload)
}

func (o *RemoveRoleOfUserDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *RemoveRoleOfUserDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
