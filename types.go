// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

type Locals interface{}

type Params map[string]string

func (p Params) Get(key string) string { return p[key] }

type ErrorHandler func(ResponseWriter, *Request, error)

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (h HandlerFunc) ServeHTTP(res ResponseWriter, req *Request) {
	h(res, req)
}

type ParamHandler func(ResponseWriter, *Request, string)

type internalHandler interface {
	serveHTTP(ResponseWriter, *Request)
	match(path string) bool
	params(path string) Params
}

type middlewareHandler struct {
	handler Handler
}

func (m *middlewareHandler) serveHTTP(res ResponseWriter, req *Request) {
	m.handler.ServeHTTP(res, req)
}

func (m *middlewareHandler) match(path string) bool {
	return true
}

func (m *middlewareHandler) params(path string) Params {
	return nil
}
