// Code generated by go-swagger; DO NOT EDIT.

package policy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// AddPolicyHandlerFunc turns a function with the right signature into a add policy handler
type AddPolicyHandlerFunc func(AddPolicyParams) middleware.Responder

// Handle executing the request and returning a response
func (fn AddPolicyHandlerFunc) Handle(params AddPolicyParams) middleware.Responder {
	return fn(params)
}

// AddPolicyHandler interface for that can handle valid add policy params
type AddPolicyHandler interface {
	Handle(AddPolicyParams) middleware.Responder
}

// NewAddPolicy creates a new http.Handler for the add policy operation
func NewAddPolicy(ctx *middleware.Context, handler AddPolicyHandler) *AddPolicy {
	return &AddPolicy{Context: ctx, Handler: handler}
}

/*AddPolicy swagger:route POST /cars/{id}/policies policy addPolicy

Add a new policy to a car

*/
type AddPolicy struct {
	Context *middleware.Context
	Handler AddPolicyHandler
}

func (o *AddPolicy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewAddPolicyParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}