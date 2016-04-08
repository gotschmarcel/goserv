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

	path                      string
	paramHandlers             paramHandlerMap
	paramHandlerInvokedMemory map[*Request]emptyNameMap
	paths                     []*Path
}

func (r *Router) All(path string, handlers ...Handler) *Router {
	r.NewRoute(path).All(handlers...)
	return r
}

func (r *Router) AllFunc(path string, funcs ...func(ResponseWriter, *Request)) *Router {
	r.NewRoute(path).AllFunc(funcs...)
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

func (r *Router) UseFunc(funcs ...func(ResponseWriter, *Request)) *Router {
	r.NewRoute("*").AllFunc(funcs...)
	return r
}

func (r *Router) Prefix(prefix string, handlers ...Handler) *Router {
	path := NewPrefixPath(prefix, NewRoute())
	path.handler.(*Route).All(handlers...)
	return r.addPath(path)
}

func (r *Router) PrefixFunc(prefix string, funcs ...func(ResponseWriter, *Request)) *Router {
	path := NewPrefixPath(prefix, NewRoute())
	path.handler.(*Route).AllFunc(funcs...)
	return r.addPath(path)
}

func (r *Router) NewRouter(prefix string) *Router {
	child := NewRouter()
	child.StrictSlash = r.StrictSlash
	child.path = r.path + prefix

	r.addPath(NewPrefixPath(prefix, child))

	return child
}

func (r *Router) NewRoute(pattern string) *Route {
	path := NewFullPath(pattern, r.StrictSlash, NewRoute())
	r.addPath(path)
	return path.handler.(*Route)
}

func (r *Router) Path() string {
	return r.path
}

func (r *Router) ServeHTTP(res ResponseWriter, req *Request) {
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
	pathString := req.SanitizedPath()[len(r.path):] // Strip own prefix

	defer r.deleteParamHandlerInvokedMemory(req)

	for _, path := range r.paths {
		if !path.Match(pathString) {
			continue
		}

		path.FillParams(req)
		if !r.handleParams(res, req) {
			return
		}

		path.ServeHTTP(res, req)

		if res.Error() != nil {
			return
		}

		if res.Written() {
			return
		}
	}
}

func (r *Router) handleParams(res ResponseWriter, req *Request) bool {
	invoked := r.getParamHandlerInvokedMemory(req)

	for name, value := range req.Params {
		if _, ok := invoked[name]; ok {
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

		invoked[name] = empty{}
	}

	return true
}

func (r *Router) getParamHandlerInvokedMemory(req *Request) emptyNameMap {
	if memory, ok := r.paramHandlerInvokedMemory[req]; ok {
		return memory
	}

	memory := make(emptyNameMap)
	r.paramHandlerInvokedMemory[req] = memory

	return memory
}

func (r *Router) deleteParamHandlerInvokedMemory(req *Request) {
	delete(r.paramHandlerInvokedMemory, req)
}

func (r *Router) addPath(path *Path) *Router {
	r.paths = append(r.paths, path)
	return r
}

func NewRouter() *Router {
	return &Router{
		paramHandlers:             make(paramHandlerMap),
		paramHandlerInvokedMemory: make(map[*Request]emptyNameMap),
	}
}

type paramHandlerMap map[string][]ParamHandler
type emptyNameMap map[string]empty
