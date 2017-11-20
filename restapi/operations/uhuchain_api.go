// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	loads "github.com/go-openapi/loads"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	security "github.com/go-openapi/runtime/security"
	spec "github.com/go-openapi/spec"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/uhuchain/uhuchain-api/restapi/operations/car"
	"github.com/uhuchain/uhuchain-api/restapi/operations/claim"
	"github.com/uhuchain/uhuchain-api/restapi/operations/policy"
	"github.com/uhuchain/uhuchain-api/restapi/operations/status"
)

// NewUhuchainAPI creates a new Uhuchain instance
func NewUhuchainAPI(spec *loads.Document) *UhuchainAPI {
	return &UhuchainAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		ServerShutdown:      func() {},
		spec:                spec,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,
		JSONConsumer:        runtime.JSONConsumer(),
		JSONProducer:        runtime.JSONProducer(),
		CarAddCarHandler: car.AddCarHandlerFunc(func(params car.AddCarParams) middleware.Responder {
			return middleware.NotImplemented("operation CarAddCar has not yet been implemented")
		}),
		ClaimAddClaimHandler: claim.AddClaimHandlerFunc(func(params claim.AddClaimParams) middleware.Responder {
			return middleware.NotImplemented("operation ClaimAddClaim has not yet been implemented")
		}),
		PolicyAddPolicyHandler: policy.AddPolicyHandlerFunc(func(params policy.AddPolicyParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyAddPolicy has not yet been implemented")
		}),
		CarGetCarHandler: car.GetCarHandlerFunc(func(params car.GetCarParams) middleware.Responder {
			return middleware.NotImplemented("operation CarGetCar has not yet been implemented")
		}),
		StatusGetStatusHandler: status.GetStatusHandlerFunc(func(params status.GetStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation StatusGetStatus has not yet been implemented")
		}),
	}
}

/*UhuchainAPI Uhuchain simple REST API for cars */
type UhuchainAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for a "application/json" mime type
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for a "application/json" mime type
	JSONProducer runtime.Producer

	// CarAddCarHandler sets the operation handler for the add car operation
	CarAddCarHandler car.AddCarHandler
	// ClaimAddClaimHandler sets the operation handler for the add claim operation
	ClaimAddClaimHandler claim.AddClaimHandler
	// PolicyAddPolicyHandler sets the operation handler for the add policy operation
	PolicyAddPolicyHandler policy.AddPolicyHandler
	// CarGetCarHandler sets the operation handler for the get car operation
	CarGetCarHandler car.GetCarHandler
	// StatusGetStatusHandler sets the operation handler for the get status operation
	StatusGetStatusHandler status.GetStatusHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// SetDefaultProduces sets the default produces media type
func (o *UhuchainAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *UhuchainAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *UhuchainAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *UhuchainAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *UhuchainAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *UhuchainAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *UhuchainAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the UhuchainAPI
func (o *UhuchainAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.CarAddCarHandler == nil {
		unregistered = append(unregistered, "car.AddCarHandler")
	}

	if o.ClaimAddClaimHandler == nil {
		unregistered = append(unregistered, "claim.AddClaimHandler")
	}

	if o.PolicyAddPolicyHandler == nil {
		unregistered = append(unregistered, "policy.AddPolicyHandler")
	}

	if o.CarGetCarHandler == nil {
		unregistered = append(unregistered, "car.GetCarHandler")
	}

	if o.StatusGetStatusHandler == nil {
		unregistered = append(unregistered, "status.GetStatusHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *UhuchainAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *UhuchainAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {

	return nil

}

// Authorizer returns the registered authorizer
func (o *UhuchainAPI) Authorizer() runtime.Authorizer {

	return nil

}

// ConsumersFor gets the consumers for the specified media types
func (o *UhuchainAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {

	result := make(map[string]runtime.Consumer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONConsumer

		}
	}
	return result

}

// ProducersFor gets the producers for the specified media types
func (o *UhuchainAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {

	result := make(map[string]runtime.Producer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONProducer

		}
	}
	return result

}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *UhuchainAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the uhuchain API
func (o *UhuchainAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *UhuchainAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened

	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/cars"] = car.NewAddCar(o.context, o.CarAddCarHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/cars/{carId}/policies/{policyId}/claims"] = claim.NewAddClaim(o.context, o.ClaimAddClaimHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/cars/{id}/policies"] = policy.NewAddPolicy(o.context, o.PolicyAddPolicyHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/cars/{id}"] = car.NewGetCar(o.context, o.CarGetCarHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/status"] = status.NewGetStatus(o.context, o.StatusGetStatusHandler)

}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *UhuchainAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middelware as you see fit
func (o *UhuchainAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}
