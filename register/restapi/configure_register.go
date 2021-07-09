// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/wcgcyx/fc-retrieval/register/restapi/operations"
	"github.com/wcgcyx/fc-retrieval/register/restapi/operations/gateway"
	"github.com/wcgcyx/fc-retrieval/register/restapi/operations/homepage"
	"github.com/wcgcyx/fc-retrieval/register/restapi/operations/provider"

	"github.com/rs/cors"

	"github.com/wcgcyx/fc-retrieval/register/internal/handlers"
)

//go:generate swagger generate server --target ../../fc-retrieval/register --name Register --spec ../docs/swagger.yml --principal interface{}

func configureFlags(_ *operations.RegisterAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.RegisterAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Homepage
	api.HomepageHomepageHandler = homepage.HomepageHandlerFunc(func(params homepage.HomepageParams) middleware.Responder {
		return handlers.HomepageHandler()
	})

	// Gateway
	api.GatewayAddGatewayRegisterHandler = gateway.AddGatewayRegisterHandlerFunc(func(params gateway.AddGatewayRegisterParams) middleware.Responder {
		return handlers.AddGatewayRegister(params)
	})
	api.GatewayGetGatewayRegistersHandler = gateway.GetGatewayRegistersHandlerFunc(func(params gateway.GetGatewayRegistersParams) middleware.Responder {
		return handlers.GetGatewayRegisters(params)
	})
	api.GatewayGetGatewayRegistersByIDHandler = gateway.GetGatewayRegistersByIDHandlerFunc(func(params gateway.GetGatewayRegistersByIDParams) middleware.Responder {
		return handlers.GetGatewayRegisterByID(params)
	})
	api.GatewayDeleteGatewayRegisterHandler = gateway.DeleteGatewayRegisterHandlerFunc(func(params gateway.DeleteGatewayRegisterParams) middleware.Responder {
		return handlers.DeleteGatewayRegisters(params)
	})

	// Register
	api.ProviderAddProviderRegisterHandler = provider.AddProviderRegisterHandlerFunc(func(params provider.AddProviderRegisterParams) middleware.Responder {
		return handlers.AddProviderRegister(params)
	})
	api.ProviderGetProviderRegistersHandler = provider.GetProviderRegistersHandlerFunc(func(params provider.GetProviderRegistersParams) middleware.Responder {
		return handlers.GetProviderRegisters(params)
	})
	api.ProviderGetProviderRegistersByIDHandler = provider.GetProviderRegistersByIDHandlerFunc(func(params provider.GetProviderRegistersByIDParams) middleware.Responder {
		return handlers.GetProviderRegisterByID(params)
	})
	api.ProviderDeleteProviderRegisterHandler = provider.DeleteProviderRegisterHandlerFunc(func(params provider.DeleteProviderRegisterParams) middleware.Responder {
		return handlers.DeleteProviderRegisters(params)
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(_ *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(_ *http.Server, _, _ string) {
	// params: server, scheme, addr
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"*"}, //TODO(magicking) restrict headers list
	})

	return c.Handler(handler)
}
