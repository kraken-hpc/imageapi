// Code generated by go-swagger; DO NOT EDIT.

package attachments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/kraken-hpc/imageapi/models"
)

// AttachCreatedCode is the HTTP code returned for type AttachCreated
const AttachCreatedCode int = 201

/*AttachCreated attach succeed

swagger:response attachCreated
*/
type AttachCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Attach `json:"body,omitempty"`
}

// NewAttachCreated creates AttachCreated with default headers values
func NewAttachCreated() *AttachCreated {

	return &AttachCreated{}
}

// WithPayload adds the payload to the attach created response
func (o *AttachCreated) WithPayload(payload *models.Attach) *AttachCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the attach created response
func (o *AttachCreated) SetPayload(payload *models.Attach) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AttachCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*AttachDefault error

swagger:response attachDefault
*/
type AttachDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAttachDefault creates AttachDefault with default headers values
func NewAttachDefault(code int) *AttachDefault {
	if code <= 0 {
		code = 500
	}

	return &AttachDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the attach default response
func (o *AttachDefault) WithStatusCode(code int) *AttachDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the attach default response
func (o *AttachDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the attach default response
func (o *AttachDefault) WithPayload(payload *models.Error) *AttachDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the attach default response
func (o *AttachDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AttachDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
