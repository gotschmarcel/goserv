// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

type ErrorHandler interface {
	ServeHTTP(ResponseWriter, *Request, error)
}

type ErrorHandlerFunc func(ResponseWriter, *Request, error)

func (e ErrorHandlerFunc) ServeHTTP(res ResponseWriter, req *Request, err error) {
	e(res, req, err)
}

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (h HandlerFunc) ServeHTTP(res ResponseWriter, req *Request) {
	h(res, req)
}

type ParamHandler interface {
	ServeHTTP(ResponseWriter, *Request, string)
}

type ParamHandlerFunc func(ResponseWriter, *Request, string)

func (p ParamHandlerFunc) ServeHTTP(res ResponseWriter, req *Request, value string) {
	p(res, req, value)
}
