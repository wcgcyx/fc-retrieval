// Code generated by go-swagger; DO NOT EDIT.

package gateway

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/wcgcyx/fc-retrieval/register/models"
)

// GetGatewayRegistersByIDOKCode is the HTTP code returned for type GetGatewayRegistersByIDOK
const GetGatewayRegistersByIDOKCode int = 200

/*GetGatewayRegistersByIDOK Get a registered gateway by Id

swagger:response getGatewayRegistersByIdOK
*/
type GetGatewayRegistersByIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.GatewayRegister `json:"body,omitempty"`
}

// NewGetGatewayRegistersByIDOK creates GetGatewayRegistersByIDOK with default headers values
func NewGetGatewayRegistersByIDOK() *GetGatewayRegistersByIDOK {

	return &GetGatewayRegistersByIDOK{}
}

// WithPayload adds the payload to the get gateway registers by Id o k response
func (o *GetGatewayRegistersByIDOK) WithPayload(payload *models.GatewayRegister) *GetGatewayRegistersByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get gateway registers by Id o k response
func (o *GetGatewayRegistersByIDOK) SetPayload(payload *models.GatewayRegister) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetGatewayRegistersByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetGatewayRegistersByIDDefault Internal error

swagger:response getGatewayRegistersByIdDefault
*/
type GetGatewayRegistersByIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetGatewayRegistersByIDDefault creates GetGatewayRegistersByIDDefault with default headers values
func NewGetGatewayRegistersByIDDefault(code int) *GetGatewayRegistersByIDDefault {
	if code <= 0 {
		code = 500
	}

	return &GetGatewayRegistersByIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get gateway registers by Id default response
func (o *GetGatewayRegistersByIDDefault) WithStatusCode(code int) *GetGatewayRegistersByIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get gateway registers by Id default response
func (o *GetGatewayRegistersByIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get gateway registers by Id default response
func (o *GetGatewayRegistersByIDDefault) WithPayload(payload *models.Error) *GetGatewayRegistersByIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get gateway registers by Id default response
func (o *GetGatewayRegistersByIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetGatewayRegistersByIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
