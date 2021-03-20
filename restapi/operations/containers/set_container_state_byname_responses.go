// Code generated by go-swagger; DO NOT EDIT.

package containers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/kraken-hpc/imageapi/models"
)

// SetContainerStateBynameOKCode is the HTTP code returned for type SetContainerStateBynameOK
const SetContainerStateBynameOKCode int = 200

/*SetContainerStateBynameOK Container state changed

swagger:response setContainerStateBynameOK
*/
type SetContainerStateBynameOK struct {

	/*
	  In: Body
	*/
	Payload *models.Container `json:"body,omitempty"`
}

// NewSetContainerStateBynameOK creates SetContainerStateBynameOK with default headers values
func NewSetContainerStateBynameOK() *SetContainerStateBynameOK {

	return &SetContainerStateBynameOK{}
}

// WithPayload adds the payload to the set container state byname o k response
func (o *SetContainerStateBynameOK) WithPayload(payload *models.Container) *SetContainerStateBynameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set container state byname o k response
func (o *SetContainerStateBynameOK) SetPayload(payload *models.Container) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetContainerStateBynameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*SetContainerStateBynameDefault error

swagger:response setContainerStateBynameDefault
*/
type SetContainerStateBynameDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSetContainerStateBynameDefault creates SetContainerStateBynameDefault with default headers values
func NewSetContainerStateBynameDefault(code int) *SetContainerStateBynameDefault {
	if code <= 0 {
		code = 500
	}

	return &SetContainerStateBynameDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the set container state byname default response
func (o *SetContainerStateBynameDefault) WithStatusCode(code int) *SetContainerStateBynameDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the set container state byname default response
func (o *SetContainerStateBynameDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the set container state byname default response
func (o *SetContainerStateBynameDefault) WithPayload(payload *models.Error) *SetContainerStateBynameDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set container state byname default response
func (o *SetContainerStateBynameDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetContainerStateBynameDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
