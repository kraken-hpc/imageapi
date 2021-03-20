// Code generated by go-swagger; DO NOT EDIT.

package mounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/kraken-hpc/imageapi/models"
)

// ListMountsRbdOKCode is the HTTP code returned for type ListMountsRbdOK
const ListMountsRbdOKCode int = 200

/*ListMountsRbdOK list of rbd mounts

swagger:response listMountsRbdOK
*/
type ListMountsRbdOK struct {

	/*
	  In: Body
	*/
	Payload []*models.MountRbd `json:"body,omitempty"`
}

// NewListMountsRbdOK creates ListMountsRbdOK with default headers values
func NewListMountsRbdOK() *ListMountsRbdOK {

	return &ListMountsRbdOK{}
}

// WithPayload adds the payload to the list mounts rbd o k response
func (o *ListMountsRbdOK) WithPayload(payload []*models.MountRbd) *ListMountsRbdOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list mounts rbd o k response
func (o *ListMountsRbdOK) SetPayload(payload []*models.MountRbd) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListMountsRbdOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.MountRbd, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*ListMountsRbdDefault error

swagger:response listMountsRbdDefault
*/
type ListMountsRbdDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewListMountsRbdDefault creates ListMountsRbdDefault with default headers values
func NewListMountsRbdDefault(code int) *ListMountsRbdDefault {
	if code <= 0 {
		code = 500
	}

	return &ListMountsRbdDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the list mounts rbd default response
func (o *ListMountsRbdDefault) WithStatusCode(code int) *ListMountsRbdDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the list mounts rbd default response
func (o *ListMountsRbdDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the list mounts rbd default response
func (o *ListMountsRbdDefault) WithPayload(payload *models.Error) *ListMountsRbdDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list mounts rbd default response
func (o *ListMountsRbdDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListMountsRbdDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
