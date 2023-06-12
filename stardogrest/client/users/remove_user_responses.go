// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
)

// RemoveUserReader is a Reader for the RemoveUser structure.
type RemoveUserReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RemoveUserReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewRemoveUserNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRemoveUserDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRemoveUserNoContent creates a RemoveUserNoContent with default headers values
func NewRemoveUserNoContent() *RemoveUserNoContent {
	return &RemoveUserNoContent{}
}

/*
RemoveUserNoContent describes a response with status code 204, with default header values.

Null response
*/
type RemoveUserNoContent struct {
}

// IsSuccess returns true when this remove user no content response has a 2xx status code
func (o *RemoveUserNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this remove user no content response has a 3xx status code
func (o *RemoveUserNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this remove user no content response has a 4xx status code
func (o *RemoveUserNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this remove user no content response has a 5xx status code
func (o *RemoveUserNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this remove user no content response a status code equal to that given
func (o *RemoveUserNoContent) IsCode(code int) bool {
	return code == 204
}

func (o *RemoveUserNoContent) Error() string {
	return fmt.Sprintf("[DELETE /users/{user}][%d] removeUserNoContent ", 204)
}

func (o *RemoveUserNoContent) String() string {
	return fmt.Sprintf("[DELETE /users/{user}][%d] removeUserNoContent ", 204)
}

func (o *RemoveUserNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewRemoveUserDefault creates a RemoveUserDefault with default headers values
func NewRemoveUserDefault(code int) *RemoveUserDefault {
	return &RemoveUserDefault{
		_statusCode: code,
	}
}

/*
RemoveUserDefault describes a response with status code -1, with default header values.

unexpected error
*/
type RemoveUserDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the remove user default response
func (o *RemoveUserDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this remove user default response has a 2xx status code
func (o *RemoveUserDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this remove user default response has a 3xx status code
func (o *RemoveUserDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this remove user default response has a 4xx status code
func (o *RemoveUserDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this remove user default response has a 5xx status code
func (o *RemoveUserDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this remove user default response a status code equal to that given
func (o *RemoveUserDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *RemoveUserDefault) Error() string {
	return fmt.Sprintf("[DELETE /users/{user}][%d] removeUser default  %+v", o._statusCode, o.Payload)
}

func (o *RemoveUserDefault) String() string {
	return fmt.Sprintf("[DELETE /users/{user}][%d] removeUser default  %+v", o._statusCode, o.Payload)
}

func (o *RemoveUserDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *RemoveUserDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
