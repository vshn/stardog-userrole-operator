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

// ChangePasswordReader is a Reader for the ChangePassword structure.
type ChangePasswordReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ChangePasswordReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewChangePasswordOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewChangePasswordDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewChangePasswordOK creates a ChangePasswordOK with default headers values
func NewChangePasswordOK() *ChangePasswordOK {
	return &ChangePasswordOK{}
}

/*
ChangePasswordOK describes a response with status code 200, with default header values.

Null response
*/
type ChangePasswordOK struct {
}

// IsSuccess returns true when this change password o k response has a 2xx status code
func (o *ChangePasswordOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this change password o k response has a 3xx status code
func (o *ChangePasswordOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this change password o k response has a 4xx status code
func (o *ChangePasswordOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this change password o k response has a 5xx status code
func (o *ChangePasswordOK) IsServerError() bool {
	return false
}

// IsCode returns true when this change password o k response a status code equal to that given
func (o *ChangePasswordOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the change password o k response
func (o *ChangePasswordOK) Code() int {
	return 200
}

func (o *ChangePasswordOK) Error() string {
	return fmt.Sprintf("[PUT /users/{user}/pwd][%d] changePasswordOK ", 200)
}

func (o *ChangePasswordOK) String() string {
	return fmt.Sprintf("[PUT /users/{user}/pwd][%d] changePasswordOK ", 200)
}

func (o *ChangePasswordOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewChangePasswordDefault creates a ChangePasswordDefault with default headers values
func NewChangePasswordDefault(code int) *ChangePasswordDefault {
	return &ChangePasswordDefault{
		_statusCode: code,
	}
}

/*
ChangePasswordDefault describes a response with status code -1, with default header values.

unexpected error
*/
type ChangePasswordDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this change password default response has a 2xx status code
func (o *ChangePasswordDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this change password default response has a 3xx status code
func (o *ChangePasswordDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this change password default response has a 4xx status code
func (o *ChangePasswordDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this change password default response has a 5xx status code
func (o *ChangePasswordDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this change password default response a status code equal to that given
func (o *ChangePasswordDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the change password default response
func (o *ChangePasswordDefault) Code() int {
	return o._statusCode
}

func (o *ChangePasswordDefault) Error() string {
	return fmt.Sprintf("[PUT /users/{user}/pwd][%d] changePassword default  %+v", o._statusCode, o.Payload)
}

func (o *ChangePasswordDefault) String() string {
	return fmt.Sprintf("[PUT /users/{user}/pwd][%d] changePassword default  %+v", o._statusCode, o.Payload)
}

func (o *ChangePasswordDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *ChangePasswordDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
