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

// SetEnabledReader is a Reader for the SetEnabled structure.
type SetEnabledReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetEnabledReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSetEnabledOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewSetEnabledDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewSetEnabledOK creates a SetEnabledOK with default headers values
func NewSetEnabledOK() *SetEnabledOK {
	return &SetEnabledOK{}
}

/*
SetEnabledOK describes a response with status code 200, with default header values.

Null response
*/
type SetEnabledOK struct {
}

// IsSuccess returns true when this set enabled o k response has a 2xx status code
func (o *SetEnabledOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this set enabled o k response has a 3xx status code
func (o *SetEnabledOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this set enabled o k response has a 4xx status code
func (o *SetEnabledOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this set enabled o k response has a 5xx status code
func (o *SetEnabledOK) IsServerError() bool {
	return false
}

// IsCode returns true when this set enabled o k response a status code equal to that given
func (o *SetEnabledOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the set enabled o k response
func (o *SetEnabledOK) Code() int {
	return 200
}

func (o *SetEnabledOK) Error() string {
	return fmt.Sprintf("[PUT /users/{user}/enabled][%d] setEnabledOK ", 200)
}

func (o *SetEnabledOK) String() string {
	return fmt.Sprintf("[PUT /users/{user}/enabled][%d] setEnabledOK ", 200)
}

func (o *SetEnabledOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewSetEnabledDefault creates a SetEnabledDefault with default headers values
func NewSetEnabledDefault(code int) *SetEnabledDefault {
	return &SetEnabledDefault{
		_statusCode: code,
	}
}

/*
SetEnabledDefault describes a response with status code -1, with default header values.

unexpected error
*/
type SetEnabledDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this set enabled default response has a 2xx status code
func (o *SetEnabledDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this set enabled default response has a 3xx status code
func (o *SetEnabledDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this set enabled default response has a 4xx status code
func (o *SetEnabledDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this set enabled default response has a 5xx status code
func (o *SetEnabledDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this set enabled default response a status code equal to that given
func (o *SetEnabledDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the set enabled default response
func (o *SetEnabledDefault) Code() int {
	return o._statusCode
}

func (o *SetEnabledDefault) Error() string {
	return fmt.Sprintf("[PUT /users/{user}/enabled][%d] setEnabled default  %+v", o._statusCode, o.Payload)
}

func (o *SetEnabledDefault) String() string {
	return fmt.Sprintf("[PUT /users/{user}/enabled][%d] setEnabled default  %+v", o._statusCode, o.Payload)
}

func (o *SetEnabledDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *SetEnabledDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
