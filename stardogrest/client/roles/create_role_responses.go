// Code generated by go-swagger; DO NOT EDIT.

package roles

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/vshn/stardog-userrole-operator/stardogrest/models"
)

// CreateRoleReader is a Reader for the CreateRole structure.
type CreateRoleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateRoleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateRoleCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewCreateRoleDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCreateRoleCreated creates a CreateRoleCreated with default headers values
func NewCreateRoleCreated() *CreateRoleCreated {
	return &CreateRoleCreated{}
}

/*
CreateRoleCreated describes a response with status code 201, with default header values.

Null response
*/
type CreateRoleCreated struct {
}

// IsSuccess returns true when this create role created response has a 2xx status code
func (o *CreateRoleCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create role created response has a 3xx status code
func (o *CreateRoleCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create role created response has a 4xx status code
func (o *CreateRoleCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this create role created response has a 5xx status code
func (o *CreateRoleCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this create role created response a status code equal to that given
func (o *CreateRoleCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the create role created response
func (o *CreateRoleCreated) Code() int {
	return 201
}

func (o *CreateRoleCreated) Error() string {
	return fmt.Sprintf("[POST /admin/roles][%d] createRoleCreated ", 201)
}

func (o *CreateRoleCreated) String() string {
	return fmt.Sprintf("[POST /admin/roles][%d] createRoleCreated ", 201)
}

func (o *CreateRoleCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateRoleDefault creates a CreateRoleDefault with default headers values
func NewCreateRoleDefault(code int) *CreateRoleDefault {
	return &CreateRoleDefault{
		_statusCode: code,
	}
}

/*
CreateRoleDefault describes a response with status code -1, with default header values.

unexpected error
*/
type CreateRoleDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this create role default response has a 2xx status code
func (o *CreateRoleDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this create role default response has a 3xx status code
func (o *CreateRoleDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this create role default response has a 4xx status code
func (o *CreateRoleDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this create role default response has a 5xx status code
func (o *CreateRoleDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this create role default response a status code equal to that given
func (o *CreateRoleDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the create role default response
func (o *CreateRoleDefault) Code() int {
	return o._statusCode
}

func (o *CreateRoleDefault) Error() string {
	return fmt.Sprintf("[POST /admin/roles][%d] createRole default  %+v", o._statusCode, o.Payload)
}

func (o *CreateRoleDefault) String() string {
	return fmt.Sprintf("[POST /admin/roles][%d] createRole default  %+v", o._statusCode, o.Payload)
}

func (o *CreateRoleDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *CreateRoleDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
