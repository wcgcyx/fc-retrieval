// Code generated by go-swagger; DO NOT EDIT.

package gateway

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/wcgcyx/fc-retrieval/register/models"
)

// AddGatewayRegisterOKCode is the HTTP code returned for type AddGatewayRegisterOK
const AddGatewayRegisterOKCode int = 200

/*AddGatewayRegisterOK Gateway register added

swagger:response addGatewayRegisterOK
*/
type AddGatewayRegisterOK struct {

	/*
	  In: Body
	*/
	Payload *models.GatewayRegister `json:"body,omitempty"`
}

// NewAddGatewayRegisterOK creates AddGatewayRegisterOK with default headers values
func NewAddGatewayRegisterOK() *AddGatewayRegisterOK {

	return &AddGatewayRegisterOK{}
}

// WithPayload adds the payload to the add gateway register o k response
func (o *AddGatewayRegisterOK) WithPayload(payload *models.GatewayRegister) *AddGatewayRegisterOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add gateway register o k response
func (o *AddGatewayRegisterOK) SetPayload(payload *models.GatewayRegister) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddGatewayRegisterOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*AddGatewayRegisterDefault Internal error

swagger:response addGatewayRegisterDefault
*/
type AddGatewayRegisterDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddGatewayRegisterDefault creates AddGatewayRegisterDefault with default headers values
func NewAddGatewayRegisterDefault(code int) *AddGatewayRegisterDefault {
	if code <= 0 {
		code = 500
	}

	return &AddGatewayRegisterDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the add gateway register default response
func (o *AddGatewayRegisterDefault) WithStatusCode(code int) *AddGatewayRegisterDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the add gateway register default response
func (o *AddGatewayRegisterDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the add gateway register default response
func (o *AddGatewayRegisterDefault) WithPayload(payload *models.Error) *AddGatewayRegisterDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add gateway register default response
func (o *AddGatewayRegisterDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddGatewayRegisterDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
