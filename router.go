// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
)

type Router struct {
	ErrorHandler ErrorHandler
	StrictSlash  bool

	*pathComponents
	path          string
	paramHandlers paramHandlerMap
	handlers      []pathHandler
}

func (r *Router) All(path string, handlers ...Handler) *Router {
	r.NewRoute(path).All(handlers...)
	return r
}

func (r *Router) AllFunc(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	r.NewRoute(path).AllFunc(funcs...)
	return r
}

func (r *Router) AllNative(path string, handlers ...http.Handler) *Router {
	r.NewRoute(path).AllNative(handlers...)
	return r
}

func (r *Router) Method(method, path string, handlers ...Handler) *Router {
	r.NewRoute(path).Method(method, handlers...)
	return r
}

func (r *Router) MethodFunc(method, path string, funcs ...func(ResponseWriter, *Request)) *Router {
	r.NewRoute(path).MethodFunc(method, funcs...)
	return r
}

func (r *Router) Methods(methods []string, path string, handlers ...Handler) *Router {
	r.NewRoute(path).Methods(methods, handlers...)
	return r
}

func (r *Router) MethodsFunc(methods []string, path string, funcs ...func(ResponseWriter, *Request)) *Router {
	r.NewRoute(path).MethodsFunc(methods, funcs...)
	return r
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

func (r *Router) ParamFunc(name string, fn func(ResponseWriter, *Request, string)) *Router {
	return r.Param(name, ParamHandlerFunc(fn))
}

func (r *Router) Use(handlers ...Handler) *Router {
	r.NewRoute("*").All(handlers...)
	return r
}

func (r *Router) UseNative(handlers ...http.Handler) *Router {
	r.NewRoute("*").AllNative(handlers...)
	return r
}

func (r *Router) UseFunc(funcs ...func(ResponseWriter, *Request)) *Router {
	r.NewRoute("*").AllFunc(funcs...)
	return r
}

func (r *Router) Mount(prefix string, router *Router) *Router {
	var err error
	router.pathComponents, err = parsePrefixPath(prefix)
	if err != nil {
		panic(err)
	}

	router.path = r.path + prefix

	return r.addHandler(router)
}

func (r *Router) NewRouter(prefix string) *Router {
	child := NewRouter()
	child.ErrorHandler = nil
	child.StrictSlash = r.StrictSlash

	r.Mount(prefix, child)

	return child
}

func (r *Router) NewRoute(path string) *Route {
	route, err := newRoute(path, r.StrictSlash)
	if err != nil {
		panic(err)
	}

	r.addHandler(route)

	return route
}

func (r *Router) Path() string {
	return r.path
}

func (r *Router) serveHTTP(res ResponseWriter, req *Request) {
	r.invokeHandlers(res, req)

	if r.ErrorHandler == nil {
		return
	}

	err := res.Error()
	if err == nil && !res.Written() {
		err = ErrNotFound
	}

	r.ErrorHandler.ServeHTTP(res, req, err)
}

func (r *Router) invokeHandlers(res ResponseWriter, req *Request) {
	paramHandlerInvoked := make(map[string]bool)
	path := req.SanitizedPath[len(r.path):] // Strip own prefix

	for _, handler := range r.handlers {
		if !handler.match(path) {
			continue
		}

		req.Params = handler.parseParams(path)
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
			paramHandler.ServeHTTP(res, req, value)

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

func (r *Router) addHandler(handler pathHandler) *Router {
	r.handlers = append(r.handlers, handler)
	return r
}

func NewRouter() *Router {
	return &Router{paramHandlers: make(paramHandlerMap)}
}

type paramHandlerMap map[string][]ParamHandler
