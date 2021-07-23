// Code generated by go-swagger; DO NOT EDIT.

package containers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kraken-hpc/imageapi/models"
)

// DeleteContainerReader is a Reader for the DeleteContainer structure.
type DeleteContainerReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteContainerReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteContainerOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDeleteContainerDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteContainerOK creates a DeleteContainerOK with default headers values
func NewDeleteContainerOK() *DeleteContainerOK {
	return &DeleteContainerOK{}
}

/* DeleteContainerOK describes a response with status code 200, with default header values.

Container deleted
*/
type DeleteContainerOK struct {
	Payload *models.Container
}

func (o *DeleteContainerOK) Error() string {
	return fmt.Sprintf("[DELETE /container][%d] deleteContainerOK  %+v", 200, o.Payload)
}
func (o *DeleteContainerOK) GetPayload() *models.Container {
	return o.Payload
}

func (o *DeleteContainerOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Container)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteContainerDefault creates a DeleteContainerDefault with default headers values
func NewDeleteContainerDefault(code int) *DeleteContainerDefault {
	return &DeleteContainerDefault{
		_statusCode: code,
	}
}

/* DeleteContainerDefault describes a response with status code -1, with default header values.

error
*/
type DeleteContainerDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the delete container default response
func (o *DeleteContainerDefault) Code() int {
	return o._statusCode
}

func (o *DeleteContainerDefault) Error() string {
	return fmt.Sprintf("[DELETE /container][%d] delete_container default  %+v", o._statusCode, o.Payload)
}
func (o *DeleteContainerDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteContainerDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
