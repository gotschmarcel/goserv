// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

// A Handler processes an HTTP request and may respond to it.
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

// HandlerFunc is an adapter to allow ordinary functions to be Handlers.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls the actual Handler function.
func (h HandlerFunc) ServeHTTP(res ResponseWriter, req *Request) {
	h(res, req)
}

// A ParamHandlerFunc can be registered to a Router using a parameter's name.
// It gets invoked with the corresponding value extracted from the request's
// path.
//
// Parameters are part of a Route's path. To learn more about parameters take
// a look at the documentation of Route.
type ParamHandlerFunc func(ResponseWriter, *Request, string)
