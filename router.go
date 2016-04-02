// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"fmt"
	"net/http"
	"regexp"
)

type Router struct {
	ErrorHandler  ErrorHandler
	StrictSlash   bool
	path          string
	pathMatcher   *regexp.Regexp
	paramHandlers paramHandlerMap
	handlers      []internalHandler
}

func (r *Router) All(path string, handlers ...Handler) *Router {
	return r.addHandler(r.Route(path).All(handlers...))
}

func (r *Router) AllFunc(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.addHandler(r.Route(path).AllFunc(funcs...))
}

func (r *Router) AllNative(path string, handlers ...http.Handler) *Router {
	return r.addHandler(r.Route(path).AllNative(handlers...))
}

func (r *Router) Method(method, path string, handlers ...Handler) *Router {
	return r.addHandler(r.Route(path).Method(method, handlers...))
}

func (r *Router) MethodFunc(method, path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.addHandler(r.Route(path).MethodFunc(method, funcs...))
}

func (r *Router) Get(path string, handlers ...Handler) *Router {
	return r.Method(http.MethodGet, path, handlers...)
}

func (r *Router) GetFunc(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.MethodFunc(http.MethodGet, path, funcs...)
}

func (r *Router) Post(path string, handlers ...Handler) *Router {
	return r.Method(http.MethodPost, path, handlers...)
}

func (r *Router) PostFunc(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.MethodFunc(http.MethodPost, path, funcs...)
}

func (r *Router) Put(path string, handlers ...Handler) *Router {
	return r.Method(http.MethodPut, path, handlers...)
}

func (r *Router) PutFunc(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.MethodFunc(http.MethodPut, path, funcs...)
}

func (r *Router) Delete(path string, handlers ...Handler) *Router {
	return r.Method(http.MethodDelete, path, handlers...)
}

func (r *Router) DeleteFunc(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.MethodFunc(http.MethodDelete, path, funcs...)
}

func (r *Router) Patch(path string, handlers ...Handler) *Router {
	return r.Method(http.MethodPatch, path, handlers...)
}

func (r *Router) PatchFunc(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.MethodFunc(http.MethodPatch, path, funcs...)
}

func (r *Router) Param(name string, handler ParamHandler) *Router {
	r.paramHandlers[name] = append(r.paramHandlers[name], handler)
	return r
}

func (r *Router) Use(handlers ...Handler) *Router {
	for _, handler := range handlers {
		r.addHandler(&middlewareHandler{handler})
	}

	return r
}

func (r *Router) UseNative(handlers ...http.Handler) *Router {
	for _, handler := range handlers {
		r.addHandler(&middlewareHandler{nativeWrapper(handler)})
	}

	return r
}

func (r *Router) UseFunc(funcs ...func(ResponseWriter, *Request)) *Router {
	for _, fn := range funcs {
		r.addHandler(&middlewareHandler{HandlerFunc(fn)})
	}

	return r
}

func (r *Router) Mount(prefix string, router *Router) *Router {
	path := fmt.Sprintf("%s%s", r.path, prefix)

	matcher, err := prefixStringToRegexp(path)
	if err != nil {
		panic(err)
	}

	router.path = path
	router.pathMatcher = matcher

	return r.addHandler(router)
}

func (r *Router) Router(prefix string) *Router {
	child := NewRouter()
	r.Mount(prefix, child)
	return child
}

func (r *Router) Route(path string) *Route {
	route, err := newRoute(path, r.StrictSlash)
	if err != nil {
		panic(err)
	}
	return route
}

func (r *Router) Path() string {
	return r.path
}

func (r *Router) ServeHTTP(nativeRes http.ResponseWriter, nativeReq *http.Request) {
	res := &responseWriter{w: nativeRes}
	req := &Request{nativeReq, &Context{}, nil, nil}
	r.serveHTTP(res, req)

	if r.ErrorHandler == nil {
		return
	}

	if err := res.Error(); err != nil {
		r.ErrorHandler(res, req, err)
		return
	}

	if !res.Written() {
		r.ErrorHandler(res, req, errNotFound)
		return
	}
}

func (r *Router) serveHTTP(res ResponseWriter, req *Request) {
	paramHandlerInvoked := make(map[string]bool)

	path := req.URL.Path[len(r.path):] // Strip own prefix

	for _, handler := range r.handlers {
		if !handler.match(path) {
			continue
		}

		req.Params = handler.params(path)
		if req.Params != nil && !r.handleParams(res, req, paramHandlerInvoked) {
			return
		}

		handler.serveHTTP(res, req)

		if res.Error() != nil {
			return
		}

		if res.Written() {
			return
		}
	}
}

func (r *Router) handleParams(res ResponseWriter, req *Request, memory map[string]bool) bool {
	for name, value := range req.Params {
		if memory[name] {
			continue
		}

		for _, paramHandler := range r.paramHandlers[name] {
			paramHandler(res, req, value)

			if res.Error() != nil {
				return false
			}

			if res.Written() {
				return false
			}
		}

		memory[name] = true
	}

	return true
}

func (r *Router) match(path string) bool {
	return r.pathMatcher.MatchString(path)
}

func (r *Router) params(path string) Params {
	return nil
}

func (r *Router) addHandler(handler internalHandler) *Router {
	r.handlers = append(r.handlers, handler)
	return r
}

func NewRouter() *Router {
	return &Router{paramHandlers: make(paramHandlerMap)}
}

func NewServer() *Router {
	r := NewRouter()
	r.ErrorHandler = defaultErrorHandler
	return r
}

type paramHandlerMap map[string][]ParamHandler
