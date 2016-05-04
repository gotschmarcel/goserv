// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a BSD-style
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
	methods map[string][]http.HandlerFunc
	path    *path
}

// All registers the specified functions for all methods in the order of appearance.
func (r *Route) All(fn http.HandlerFunc) *Route {
	for _, method := range methodNames {
		r.addMethodHandlerFunc(method, fn)
	}
	return r
}

// Method registers the functions for the specified HTTP method in the order of appearance.
func (r *Route) Method(method string, fn http.HandlerFunc) *Route {
	r.addMethodHandlerFunc(method, fn)
	return r
}

// Methods is an adapter for .Method to register functions
// for multiple methods in one call.
func (r *Route) Methods(methods []string, fn http.HandlerFunc) *Route {
	for _, method := range methods {
		r.Method(method, fn)
	}
	return r
}

// Get is an adapter for .Method and registers the functions for the "GET" method.
func (r *Route) Get(fn http.HandlerFunc) *Route {
	return r.Method(http.MethodGet, fn)
}

// Post is an adapter for .Method and registers the functions for the "POST" method.
func (r *Route) Post(fn http.HandlerFunc) *Route {
	return r.Method(http.MethodPost, fn)
}

// Put is an adapter for .Method and registers the functions for the "PUT" method.
func (r *Route) Put(fn http.HandlerFunc) *Route {
	return r.Method(http.MethodPut, fn)
}

// Delete is an adapter for .Method and registers the functions for the "DELETE" method.
func (r *Route) Delete(fn http.HandlerFunc) *Route {
	return r.Method(http.MethodDelete, fn)
}

// Patch is an adapter for .Method and registers the functions for the "PATCH" method.
func (r *Route) Patch(fn http.HandlerFunc) *Route {
	return r.Method(http.MethodPatch, fn)
}

// Rest registers the given handler on all methods without a handler.
func (r *Route) Rest(fn http.HandlerFunc) *Route {
	for _, method := range methodNames {
		if len(r.methods[method]) > 0 {
			continue
		}

		r.addMethodHandlerFunc(method, fn)
	}

	return r
}

// serveHTTP processes the Request by invoking all middleware and all method handlers for the
// corresponding method of the Request in the order they were registered.
//
// The processing stops as soon as a handler writes a response or set's an error
// on the RequestContext.
func (r *Route) serveHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := Context(req)

	for _, handler := range r.methods[req.Method] {
		handler(res, req)

		if doneProcessing(res.(*responseWriter), ctx) {
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

func (r *Route) fillParams(req *http.Request, params map[string]string) {
	if !r.path.ContainsParams() {
		return
	}

	r.path.FillParams(SanitizePath(req.URL.Path), params)
}

func (r *Route) addMethodHandlerFunc(method string, fn http.HandlerFunc) {
	r.methods[method] = append(r.methods[method], fn)
}

func newRoute(pattern string, strict, prefixOnly bool) *Route {
	path, err := parsePath(pattern, strict, prefixOnly)

	if err != nil {
		panic(err)
	}

	return &Route{
		methods: make(map[string][]http.HandlerFunc),
		path:    path,
	}
}
