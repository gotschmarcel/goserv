// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"regexp"
)

// A Route handles Requests by processing registered middleware and method handlers.
//
// Every route is build from a path specified in the call to a Server/Router's routing
// functions, like .Method(path, ...handlers) or .Get(path, ...handlers).
//
// The Path
//
// A path describes the request paths to which a Route should respond, e.g.
//	/mypath
//
// A Route created from "/mypath" will only match the request path if it is exactly
// "/mypath".
//
// In case that the route should match everything starting with "/mypath" a wildcard
// can be appended to the path, i.e.
//	/mypath*
//
// The wildcard can be at any position in the path. It also possible to just use the wildcard
// as path.
//
// Sometimes it is necessary to capture values from parts of the request path, so called parameters.
// To include parameters in a Route the path must contain a named parameter, e.g.
//	/users/:user_id
//
// Parameters always start with a ":" after a "/". The name (without the leading ":") can contain
// alphanumeric symbols as well as "_" and "-". A Route can have as many parameters as you like separated
// by at least a single "/".
//	/:value1/:value2
//
// When a Route processes a Request it automatically extracts the captured parameter values from the path
// and stores the values under their name in the .Param field of the Request.
//
// Strict vs non-strict
//
// Depending on the Server/Router's configuration routes are strict or non-strict. A strict Route matches
// the path only if the path matches exactly the Route's pattern including the last slash. That means that
// a strict Route created from "/path" matches "/path" but not "/path/".
//
// For non-strict routes the last slash is optional and the Route matches the path with or without a slash
// at the end.
//
// Request Processing
//
// The order in which Handlers are registered does matter in terms of the processing order. The only
// thing not relevant is wether a middleware handler (.All*) was registered before or
// after a method handler (.Method, .Get, ...).
//
// The processing ends as soon as a handler wrote a response, an error was set on the ResponseWriter
// or all handlers were invoked.
//
// Note that all handler functions return the Route itself to allow method chaining, e.g.
//	route.All(middleware).Get(getHandler).Put(putHandler)
type Route struct {
	middleware []Handler
	methods    map[string][]Handler
	matcher    *regexp.Regexp
	params     []string
}

// All registers the specified handlers as middleware in the order of appearance.
func (r *Route) All(handlers ...Handler) *Route {
	r.middleware = append(r.middleware, handlers...)
	return r
}

// AllFunc is an adapter for .All to register ordinary functions as middleware.
func (r *Route) AllFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	for _, fn := range funcs {
		r.middleware = append(r.middleware, HandlerFunc(fn))
	}
	return r
}

// Method registers the handlers for the specified HTTP method in the order of appearance.
func (r *Route) Method(method string, handlers ...Handler) *Route {
	r.addMethodHandlers(method, handlers...)
	return r
}

// MethodFunc is an adapater for .Method to register ordinary functions as method handlers.
func (r *Route) MethodFunc(method string, funcs ...func(ResponseWriter, *Request)) *Route {
	for _, fn := range funcs {
		r.addMethodHandlers(method, HandlerFunc(fn))
	}
	return r
}

// Methods is an adapter for .Method to register handlers
// for multiple methods in one call.
func (r *Route) Methods(methods []string, handlers ...Handler) *Route {
	for _, method := range methods {
		r.Method(method, handlers...)
	}
	return r
}

// MethodsFunc is an adapter for .Methods to register ordinary functions
// for multiple methods in one call.
func (r *Route) MethodsFunc(methods []string, funcs ...func(ResponseWriter, *Request)) *Route {
	for _, method := range methods {
		r.MethodFunc(method, funcs...)
	}
	return r
}

// Get is an adapter for .Method and registers the handlers for the "GET" method.
func (r *Route) Get(handlers ...Handler) *Route {
	return r.Method(http.MethodGet, handlers...)
}

// GetFunc is an adapater for .Get to register ordinary functions as Handlers
// for the "GET" method.
func (r *Route) GetFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodGet, funcs...)
}

// Post is an adapter for .Method and registers the handlers for the "POST" method.
func (r *Route) Post(handlers ...Handler) *Route {
	return r.Method(http.MethodPost, handlers...)
}

// PostFunc is an adapater for .Post to register ordinary functions as Handlers
// for the "POST" method.
func (r *Route) PostFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodPost, funcs...)
}

// Put is an adapter for .Method and registers the handlers for the "PUT" method.
func (r *Route) Put(handlers ...Handler) *Route {
	return r.Method(http.MethodPut, handlers...)
}

// PutFunc is an adapater for .Put to register ordinary functions as Handlers
// for the "PUT" method.
func (r *Route) PutFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodPut, funcs...)
}

// Delete is an adapter for .Method and registers the handlers for the "DELETE" method.
func (r *Route) Delete(handlers ...Handler) *Route {
	return r.Method(http.MethodDelete, handlers...)
}

// DeleteFunc is an adapater for .Delete to register ordinary functions as Handlers
// for the "DELETE" method.
func (r *Route) DeleteFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodDelete, funcs...)
}

// Patch is an adapter for .Method and registers the handlers for the "PATCH" method.
func (r *Route) Patch(handlers ...Handler) *Route {
	return r.Method(http.MethodPatch, handlers...)
}

// PatchFunc is an adapater for .Patch to register ordinary functions as Handlers
// for the "PATCH" method.
func (r *Route) PatchFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodPatch, funcs...)
}

// ServeHTTP processes the Request by invoking all middleware and all method handlers for the
// corresponding method of the Request in the order they were registered.
//
// The processing stops as soon as a handler writes a response or set's an error
// on the ResponseWriter.
func (r *Route) ServeHTTP(res ResponseWriter, req *Request) {
	if r.ContainsParams() && len(req.Params) == 0 {
		r.fillParams(req)
	}

	for _, handler := range append(r.middleware, r.methods[req.Method]...) {
		handler.ServeHTTP(res, req)

		if res.Error() != nil {
			return
		}

		if res.Written() {
			return
		}
	}
}

// Match returns true if the given path fulfills the routes matching pattern.
func (r *Route) Match(path string) bool {
	return r.matcher.MatchString(path)
}

// ContainsParams returns true if the route has any parameters registered.
func (r *Route) ContainsParams() bool {
	return len(r.params) > 0
}

// Params returns a list containing the names of all registered parameters.
func (r *Route) Params() []string {
	return append([]string{}, r.params...) // return a copy
}

func (r *Route) fillParams(req *Request) {
	if !r.ContainsParams() {
		return
	}

	matches := r.matcher.FindAllStringSubmatch(req.SanitizedPath(), -1)
	if len(matches) == 0 {
		return
	}

	// Iterate group matches only
	for index, value := range matches[0][1:] {
		name := r.params[index]
		req.Params[name] = value
	}
}

func (r *Route) addMethodHandlers(method string, handlers ...Handler) {
	r.methods[method] = append(r.methods[method], handlers...)
}

func newRoute(pattern string, strict, prefixOnly bool) *Route {
	matcher, params := pathComponents(pattern, strict, prefixOnly)

	return &Route{
		methods: make(map[string][]Handler),
		matcher: matcher,
		params:  params,
	}
}
