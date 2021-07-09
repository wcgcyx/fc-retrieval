package handlers

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"

	"github.com/wcgcyx/fc-retrieval/register/models"
	op "github.com/wcgcyx/fc-retrieval/register/restapi/operations/homepage"
)

// HomepageHandler handler
func HomepageHandler() middleware.Responder {
	serviceName := apiconfig.GetString("SERVICE_NAME")

	// Response
	payload := models.Ack{Status: "success", Message: strings.Join([]string{"Welcome to", serviceName, "api, please go to /docs."}, " ")}
	return op.NewHomepageOK().WithPayload(&payload)
}
