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
	methods map[string][]func(ResponseWriter, *Request)
	path    *path
}

// All registers the specified functions for all methods in the order of appearance.
func (r *Route) All(funcs ...func(ResponseWriter, *Request)) *Route {
	for _, method := range methodNames {
		r.addMethodHandlers(method, funcs...)
	}
	return r
}

// Method registers the functions for the specified HTTP method in the order of appearance.
func (r *Route) Method(method string, funcs ...func(ResponseWriter, *Request)) *Route {
	r.addMethodHandlers(method, funcs...)
	return r
}

// Methods is an adapter for .Method to register functions
// for multiple methods in one call.
func (r *Route) Methods(methods []string, funcs ...func(ResponseWriter, *Request)) *Route {
	for _, method := range methods {
		r.Method(method, funcs...)
	}
	return r
}

// Get is an adapter for .Method and registers the functions for the "GET" method.
func (r *Route) Get(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.Method(http.MethodGet, funcs...)
}

// Post is an adapter for .Method and registers the functions for the "POST" method.
func (r *Route) Post(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.Method(http.MethodPost, funcs...)
}

// Put is an adapter for .Method and registers the functions for the "PUT" method.
func (r *Route) Put(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.Method(http.MethodPut, funcs...)
}

// Delete is an adapter for .Method and registers the functions for the "DELETE" method.
func (r *Route) Delete(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.Method(http.MethodDelete, funcs...)
}

// Patch is an adapter for .Method and registers the functions for the "PATCH" method.
func (r *Route) Patch(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.Method(http.MethodPatch, funcs...)
}

// ServeHTTP processes the Request by invoking all middleware and all method handlers for the
// corresponding method of the Request in the order they were registered.
//
// The processing stops as soon as a handler writes a response or set's an error
// on the ResponseWriter.
func (r *Route) ServeHTTP(res ResponseWriter, req *Request) {
	if r.containsParams() && len(req.Params) == 0 {
		r.fillParams(req)
	}

	for _, handler := range r.methods[req.Method] {
		handler(res, req)

		if res.Error() != nil {
			return
		}

		if res.Written() {
			return
		}
	}
}

func (r *Route) match(path string) bool {
	return r.path.Match(path)
}

func (r *Route) containsParams() bool {
	return r.path.ContainsParams()
}

func (r *Route) params() []string {
	return r.path.Params()
}

func (r *Route) fillParams(req *Request) {
	r.path.FillParams(req)
}

func (r *Route) addMethodHandlers(method string, funcs ...func(ResponseWriter, *Request)) {
	r.methods[method] = append(r.methods[method], funcs...)
}

func newRoute(pattern string, strict, prefixOnly bool) *Route {
	path, err := parsePath(pattern, strict, prefixOnly)

	if err != nil {
		panic(err)
	}

	return &Route{
		methods: make(map[string][]func(ResponseWriter, *Request)),
		path:    path,
	}
}
