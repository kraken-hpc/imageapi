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

// MountReader is a Reader for the Mount structure.
type MountReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MountReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewMountCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewMountDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewMountCreated creates a MountCreated with default headers values
func NewMountCreated() *MountCreated {
	return &MountCreated{}
}

/* MountCreated describes a response with status code 201, with default header values.

mount succeed
*/
type MountCreated struct {
	Payload *models.Mount
}

func (o *MountCreated) Error() string {
	return fmt.Sprintf("[POST /mount][%d] mountCreated  %+v", 201, o.Payload)
}
func (o *MountCreated) GetPayload() *models.Mount {
	return o.Payload
}

func (o *MountCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Mount)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewMountDefault creates a MountDefault with default headers values
func NewMountDefault(code int) *MountDefault {
	return &MountDefault{
		_statusCode: code,
	}
}

/* MountDefault describes a response with status code -1, with default header values.

error
*/
type MountDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the mount default response
func (o *MountDefault) Code() int {
	return o._statusCode
}

func (o *MountDefault) Error() string {
	return fmt.Sprintf("[POST /mount][%d] mount default  %+v", o._statusCode, o.Payload)
}
func (o *MountDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *MountDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
