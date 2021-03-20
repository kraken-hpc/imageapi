// Code generated by go-swagger; DO NOT EDIT.

package attach

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kraken-hpc/imageapi/models"
)

// MapRbdReader is a Reader for the MapRbd structure.
type MapRbdReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MapRbdReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewMapRbdCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewMapRbdDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewMapRbdCreated creates a MapRbdCreated with default headers values
func NewMapRbdCreated() *MapRbdCreated {
	return &MapRbdCreated{}
}

/*MapRbdCreated handles this case with default header values.

RBD attach succeed
*/
type MapRbdCreated struct {
	Payload *models.Rbd
}

func (o *MapRbdCreated) Error() string {
	return fmt.Sprintf("[POST /attach/rbd][%d] mapRbdCreated  %+v", 201, o.Payload)
}

func (o *MapRbdCreated) GetPayload() *models.Rbd {
	return o.Payload
}

func (o *MapRbdCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Rbd)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewMapRbdDefault creates a MapRbdDefault with default headers values
func NewMapRbdDefault(code int) *MapRbdDefault {
	return &MapRbdDefault{
		_statusCode: code,
	}
}

/*MapRbdDefault handles this case with default header values.

error
*/
type MapRbdDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the map rbd default response
func (o *MapRbdDefault) Code() int {
	return o._statusCode
}

func (o *MapRbdDefault) Error() string {
	return fmt.Sprintf("[POST /attach/rbd][%d] map_rbd default  %+v", o._statusCode, o.Payload)
}

func (o *MapRbdDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *MapRbdDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
