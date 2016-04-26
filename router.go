// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
)

// A Router dispatches incoming requests to matching routes and routers.
//
// Note that most methods return the Router itself to allow method chaining.
// Some methods like .Route or .SubRouter return the created instances instead.
type Router struct {
	// Handles errors set on the ResponseWriter with .SetError(err), not found errors
	// and recovered panics.
	ErrorHandler func(ResponseWriter, *Request, error)

	// Defines how Routes treat the trailing slash in a path.
	//
	// When enabled routes with a trailing slash are considered to be different routes
	// than routes without a trailing slash.
	StrictSlash bool

	path          string
	paramHandlers paramHandlerMap
	routes        []*Route
}

// All registers the specified functions in the order of appearance for the given path for
// all http methods.
func (r *Router) All(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	r.Route(path).All(funcs...)
	return r
}

// Method registers the specified functions in the order of appearance for the given path
// and method.
func (r *Router) Method(method, path string, funcs ...func(ResponseWriter, *Request)) *Router {
	r.Route(path).Method(method, funcs...)
	return r
}

// Methods is an adapter for Method for registering functions on a path for multiple methods
// in a single call.
func (r *Router) Methods(methods []string, path string, funcs ...func(ResponseWriter, *Request)) *Router {
	r.Route(path).Methods(methods, funcs...)
	return r
}

// Get is an adapter for Method registering the functions for the "GET" method on path.
func (r *Router) Get(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.Method(http.MethodGet, path, funcs...)
}

// Post is an adapter for Method registering the functions for the "POST" method on path.
func (r *Router) Post(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.Method(http.MethodPost, path, funcs...)
}

// Put is an adapter for Method registering the functions for the "PUT" method on path.
func (r *Router) Put(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.Method(http.MethodPut, path, funcs...)
}

// Delete is an adapter for Method registering the functions for the "DELETE" method on path.
func (r *Router) Delete(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.Method(http.MethodDelete, path, funcs...)
}

// Patch is an adapter for Method registering the functions for the "PATCH" method on path.
func (r *Router) Patch(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.Method(http.MethodPatch, path, funcs...)
}

// Use registers the specified functions in the order of appearance as middleware.
// Middleware is always processed before any dispatching happens.
func (r *Router) Use(funcs ...func(ResponseWriter, *Request)) *Router {
	r.Route("/*").All(funcs...)
	return r
}

// UseHandler is an adapter for Use to register handlers as middleware.
func (r *Router) UseHandler(handlers ...Handler) *Router {
	route := r.Route("/*")
	for _, handler := range handlers {
		route.All(handler.ServeHTTP)
	}

	return r
}

// Param registers a handler for the specified parameter name (without the leading ":").
//
// Parameter handlers are invoked with the extracted value before any route is processed.
// All handlers are only invoked once per request, even though the request may be dispatched
// to multiple routes.
func (r *Router) Param(name string, fn func(ResponseWriter, *Request, string)) *Router {
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

	r.addRoute(newRoute(prefix, r.StrictSlash, true).All(router.ServeHTTP))

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

// ServeHTTP dispatches the request to middleware and matching handlers.
//
// Errors are only processed if an ErrorHandler was configured.
func (r *Router) ServeHTTP(res ResponseWriter, req *Request) {
	r.invokeHandlers(res, req)

	if res.Written() || r.ErrorHandler == nil {
		return
	}

	err := res.Error()
	if err == nil {
		err = ErrNotFound
	}

	r.ErrorHandler(res, req, err)
}

func (r *Router) invokeHandlers(res ResponseWriter, req *Request) {
	path := req.sanitizedPath[len(r.path):] // Strip own prefix

	paramInvokedMem := make(map[string]bool)

	for _, route := range r.routes {
		if !route.match(path) {
			continue
		}

		route.fillParams(req)
		if !r.handleParams(res, req, route.params(), paramInvokedMem) {
			return
		}

		route.ServeHTTP(res, req)

		if res.Error() != nil {
			return
		}

		if res.Written() {
			return
		}
	}
}

func (r *Router) handleParams(res ResponseWriter, req *Request, orderedParams []string, invoked map[string]bool) bool {
	// Call param handlers in the same order in which the parameters appear in the path.
	for _, name := range orderedParams {
		if invoked[name] {
			continue
		}

		value := req.Params.Get(name)

		for _, paramHandler := range r.paramHandlers[name] {
			paramHandler(res, req, value)

			if res.Error() != nil {
				return false
			}

			if res.Written() {
				return false
			}
		}

		invoked[name] = true
	}

	return true
}

func (r *Router) addRoute(route *Route) *Router {
	r.routes = append(r.routes, route)
	return r
}

func newRouter() *Router {
	return &Router{paramHandlers: make(paramHandlerMap)}
}

type paramHandlerMap map[string][]ParamHandlerFunc
