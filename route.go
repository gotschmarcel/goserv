// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
)

// A Route handles requests by processing method handlers.
//
// Note that all handler functions return the Route itself to allow method chaining, e.g.
//	route.All(middleware).Get(getHandler).Put(putHandler)
type Route struct {
	methods map[string][]Handler
	path    *path
}

// All registers the specified handlers for all methods in the order of appearance.
func (r *Route) All(handlers ...Handler) *Route {
	for _, method := range methodNames {
		r.addMethodHandlers(method, handlers...)
	}
	return r
}

// AllFunc is an adapter for .All to register ordinary functions as middleware.
func (r *Route) AllFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	for _, fn := range funcs {
		r.All(HandlerFunc(fn))
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

	for _, handler := range r.methods[req.Method] {
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
	return r.path.Match(path)
}

// ContainsParams returns true if the route has any parameters registered.
func (r *Route) ContainsParams() bool {
	return r.path.ContainsParams()
}

func (r *Route) params() []string {
	return r.path.Params()
}

func (r *Route) fillParams(req *Request) {
	r.path.FillParams(req)
}

func (r *Route) addMethodHandlers(method string, handlers ...Handler) {
	r.methods[method] = append(r.methods[method], handlers...)
}

func newRoute(pattern string, strict, prefixOnly bool) *Route {
	path, err := parsePath(pattern, strict, prefixOnly)

	if err != nil {
		panic(err)
	}

	return &Route{
		methods: make(map[string][]Handler),
		path:    path,
	}
}
