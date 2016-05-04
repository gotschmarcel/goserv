// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goserv

import (
	"fmt"
	"net/http"
	"strings"
)

// A Router dispatches incoming requests to matching routes and routers.
//
// Note that most methods return the Router itself to allow method chaining.
// Some methods like .Route or .SubRouter return the created instances instead.
type Router struct {
	// Handles errors set on the RequestContext with .Error, not found errors
	// and recovered panics.
	ErrorHandler ErrorHandlerFunc

	// Defines how Routes treat the trailing slash in a path.
	//
	// When enabled routes with a trailing slash are considered to be different routes
	// than routes without a trailing slash.
	StrictSlash bool

	// Enables/Disables panic recovery
	PanicRecovery bool

	path          string
	paramHandlers paramHandlerMap
	routes        []*Route
}

// All registers the specified HandlerFunc for the given path for
// all http methods.
func (r *Router) All(path string, fn http.HandlerFunc) *Router {
	r.Route(path).All(fn)
	return r
}

// Method registers the specified HandlerFunc for the given path
// and method.
func (r *Router) Method(method, path string, fn http.HandlerFunc) *Router {
	r.Route(path).Method(method, fn)
	return r
}

// Methods is an adapter for Method for registering a HandlerFunc on a path for multiple methods
// in a single call.
func (r *Router) Methods(methods []string, path string, fn http.HandlerFunc) *Router {
	r.Route(path).Methods(methods, fn)
	return r
}

// Get is an adapter for Method registering a HandlerFunc for the "GET" method on path.
func (r *Router) Get(path string, fn http.HandlerFunc) *Router {
	return r.Method(http.MethodGet, path, fn)
}

// Post is an adapter for Method registering a HandlerFunc for the "POST" method on path.
func (r *Router) Post(path string, fn http.HandlerFunc) *Router {
	return r.Method(http.MethodPost, path, fn)
}

// Put is an adapter for Method registering a HandlerFunc for the "PUT" method on path.
func (r *Router) Put(path string, fn http.HandlerFunc) *Router {
	return r.Method(http.MethodPut, path, fn)
}

// Delete is an adapter for Method registering a HandlerFunc for the "DELETE" method on path.
func (r *Router) Delete(path string, fn http.HandlerFunc) *Router {
	return r.Method(http.MethodDelete, path, fn)
}

// Patch is an adapter for Method registering a HandlerFunc for the "PATCH" method on path.
func (r *Router) Patch(path string, fn http.HandlerFunc) *Router {
	return r.Method(http.MethodPatch, path, fn)
}

// Use registers the specified function as middleware.
// Middleware is always processed before any dispatching happens.
func (r *Router) Use(fn http.HandlerFunc) *Router {
	r.Route("/*").All(fn)
	return r
}

// UseHandler is an adapter for Use to register a Handler as middleware.
func (r *Router) UseHandler(handler http.Handler) *Router {
	r.Route("/*").All(handler.ServeHTTP)
	return r
}

// Param registers a handler for the specified parameter name (without the leading ":").
//
// Parameter handlers are invoked with the extracted value before any route is processed.
// All handlers are only invoked once per request, even though the request may be dispatched
// to multiple routes.
func (r *Router) Param(name string, fn ParamHandlerFunc) *Router {
	r.paramHandlers[name] = append(r.paramHandlers[name], fn)
	return r
}

// SubRouter returns a new sub router mounted on the specified prefix.
//
// All sub routers automatically inherit their StrictSlash behaviour,
// have the full mount path and no error handler. It is possible though
// to set a custom error handler for a sub router.
//
// Note that this function returns the new sub router instead of the
// parent router!
func (r *Router) SubRouter(prefix string) *Router {
	router := newRouter()
	router.StrictSlash = r.StrictSlash
	router.path = r.path + prefix

	r.addRoute(newRoute(prefix, r.StrictSlash, true).All(router.serveHTTP))

	return router
}

// Route returns a new Route for the given path.
func (r *Router) Route(path string) *Route {
	route := newRoute(path, r.StrictSlash, false)
	r.addRoute(route)
	return route
}

// Path returns the routers full mount path.
func (r *Router) Path() string {
	return r.path
}

func (r *Router) serveHTTP(res http.ResponseWriter, req *http.Request) {
	if r.PanicRecovery {
		defer r.handleRecovery(res, req)
	}

	ctx := Context(req)

	r.invokeHandlers(res, req, ctx)

	if res.(*responseWriter).Written() || r.ErrorHandler == nil {
		return
	}

	if ctx.err == nil {
		ctx.Error(ErrNotFound, http.StatusNotFound)
	}

	r.ErrorHandler(res, req, ctx.err)
}

func (r *Router) invokeHandlers(res http.ResponseWriter, req *http.Request, ctx *RequestContext) {
	defer ctx.skipped() // Clear any skip events

	path := strings.TrimPrefix(SanitizePath(req.URL.Path), r.path)

	paramInvoked := make(map[string]bool)

	for _, route := range r.routes {
		if !route.match(path) {
			continue
		}

		// Call param handlers in the same order in which the parameters appear in the path.
		route.fillParams(req, ctx.params)
		for _, name := range route.params() {
			if paramInvoked[name] {
				continue
			}

			value := ctx.Param(name)

			for _, paramHandler := range r.paramHandlers[name] {
				paramHandler(res, req, value)

				if doneProcessing(res.(*responseWriter), ctx) {
					return
				}
			}

			paramInvoked[name] = true
		}

		route.serveHTTP(res, req)

		if doneProcessing(res.(*responseWriter), ctx) {
			return
		}
	}
}

func (r *Router) handleRecovery(res http.ResponseWriter, req *http.Request) {
	if err := recover(); err != nil && r.ErrorHandler != nil {
		r.ErrorHandler(res, req, &ContextError{fmt.Errorf("Panic: %v", err), http.StatusInternalServerError})
	}
}

func (r *Router) addRoute(route *Route) *Router {
	r.routes = append(r.routes, route)
	return r
}

func newRouter() *Router {
	return &Router{paramHandlers: make(paramHandlerMap)}
}

type paramHandlerMap map[string][]ParamHandlerFunc
