// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

// A ErrorHandler is the last handler in the request chain and
// is responsible for handling errors that occur during the
// request processing.
//
// A ErrorHandler should always write a response!
type ErrorHandler interface {
	ServeHTTP(ResponseWriter, *Request, error)
}

// ErrorHandlerFunc is an adapter to allow ordinary functions
// to be ErrorHandlers.
type ErrorHandlerFunc func(ResponseWriter, *Request, error)

// ServeHTTP calls the actual function.
func (e ErrorHandlerFunc) ServeHTTP(res ResponseWriter, req *Request, err error) {
	e(res, req, err)
}

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

// A ParamHandler can be registered to a Router using the parameter's name.
// It gets invoked with the corresponding value extracted from the request's
// path.
//
// Parameters are part of a Route's path. To learn more about parameters take
// a look at the documentation of Route.
type ParamHandler interface {
	ServeHTTP(ResponseWriter, *Request, string)
}

// ParamHandlerFunc is an adapter to allow ordinary functions to be ParamHandlers.
type ParamHandlerFunc func(ResponseWriter, *Request, string)

// ServeHTTP calls the actual function.
func (p ParamHandlerFunc) ServeHTTP(res ResponseWriter, req *Request, value string) {
	p(res, req, value)
}
