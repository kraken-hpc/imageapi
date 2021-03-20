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

// UnmountOverlayReader is a Reader for the UnmountOverlay structure.
type UnmountOverlayReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UnmountOverlayReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUnmountOverlayOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewUnmountOverlayDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUnmountOverlayOK creates a UnmountOverlayOK with default headers values
func NewUnmountOverlayOK() *UnmountOverlayOK {
	return &UnmountOverlayOK{}
}

/*UnmountOverlayOK handles this case with default header values.

Unmounted
*/
type UnmountOverlayOK struct {
	Payload *models.MountOverlay
}

func (o *UnmountOverlayOK) Error() string {
	return fmt.Sprintf("[DELETE /mount/overlay/{id}][%d] unmountOverlayOK  %+v", 200, o.Payload)
}

func (o *UnmountOverlayOK) GetPayload() *models.MountOverlay {
	return o.Payload
}

func (o *UnmountOverlayOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.MountOverlay)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUnmountOverlayDefault creates a UnmountOverlayDefault with default headers values
func NewUnmountOverlayDefault(code int) *UnmountOverlayDefault {
	return &UnmountOverlayDefault{
		_statusCode: code,
	}
}

/*UnmountOverlayDefault handles this case with default header values.

error
*/
type UnmountOverlayDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the unmount overlay default response
func (o *UnmountOverlayDefault) Code() int {
	return o._statusCode
}

func (o *UnmountOverlayDefault) Error() string {
	return fmt.Sprintf("[DELETE /mount/overlay/{id}][%d] unmount_overlay default  %+v", o._statusCode, o.Payload)
}

func (o *UnmountOverlayDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *UnmountOverlayDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
