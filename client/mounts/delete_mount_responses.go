// Code generated by go-swagger; DO NOT EDIT.

package mounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kraken-hpc/imageapi/models"
)

// DeleteMountReader is a Reader for the DeleteMount structure.
type DeleteMountReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteMountReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteMountOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDeleteMountDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteMountOK creates a DeleteMountOK with default headers values
func NewDeleteMountOK() *DeleteMountOK {
	return &DeleteMountOK{}
}

/* DeleteMountOK describes a response with status code 200, with default header values.

Unmount succeeded
*/
type DeleteMountOK struct {
	Payload *models.Mount
}

func (o *DeleteMountOK) Error() string {
	return fmt.Sprintf("[DELETE /mount][%d] deleteMountOK  %+v", 200, o.Payload)
}
func (o *DeleteMountOK) GetPayload() *models.Mount {
	return o.Payload
}

func (o *DeleteMountOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Mount)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteMountDefault creates a DeleteMountDefault with default headers values
func NewDeleteMountDefault(code int) *DeleteMountDefault {
	return &DeleteMountDefault{
		_statusCode: code,
	}
}

/* DeleteMountDefault describes a response with status code -1, with default header values.

Unmount failed
*/
type DeleteMountDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the delete mount default response
func (o *DeleteMountDefault) Code() int {
	return o._statusCode
}

func (o *DeleteMountDefault) Error() string {
	return fmt.Sprintf("[DELETE /mount][%d] DeleteMount default  %+v", o._statusCode, o.Payload)
}
func (o *DeleteMountDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteMountDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
