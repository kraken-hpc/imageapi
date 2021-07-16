// Code generated by go-swagger; DO NOT EDIT.

package attachments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kraken-hpc/imageapi/models"
)

// DeleteAttachReader is a Reader for the DeleteAttach structure.
type DeleteAttachReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteAttachReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteAttachOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDeleteAttachDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteAttachOK creates a DeleteAttachOK with default headers values
func NewDeleteAttachOK() *DeleteAttachOK {
	return &DeleteAttachOK{}
}

/* DeleteAttachOK describes a response with status code 200, with default header values.

Detach succeed
*/
type DeleteAttachOK struct {
	Payload *models.Attach
}

func (o *DeleteAttachOK) Error() string {
	return fmt.Sprintf("[DELETE /attach][%d] deleteAttachOK  %+v", 200, o.Payload)
}
func (o *DeleteAttachOK) GetPayload() *models.Attach {
	return o.Payload
}

func (o *DeleteAttachOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Attach)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteAttachDefault creates a DeleteAttachDefault with default headers values
func NewDeleteAttachDefault(code int) *DeleteAttachDefault {
	return &DeleteAttachDefault{
		_statusCode: code,
	}
}

/* DeleteAttachDefault describes a response with status code -1, with default header values.

Detach failed
*/
type DeleteAttachDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the delete attach default response
func (o *DeleteAttachDefault) Code() int {
	return o._statusCode
}

func (o *DeleteAttachDefault) Error() string {
	return fmt.Sprintf("[DELETE /attach][%d] DeleteAttach default  %+v", o._statusCode, o.Payload)
}
func (o *DeleteAttachDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteAttachDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
