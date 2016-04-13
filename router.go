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

	path          string
	paramHandlers paramHandlerMap
	routes        []*Route
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
	return r.addRoute(newRoute(prefix, r.StrictSlash, true).All(handlers...))
}

func (r *Router) PrefixFunc(prefix string, funcs ...func(ResponseWriter, *Request)) *Router {
	return r.addRoute(newRoute(prefix, r.StrictSlash, true).AllFunc(funcs...))
}

func (r *Router) NewRouter(prefix string) *Router {
	router := NewRouter()
	router.StrictSlash = r.StrictSlash
	router.path = r.path + prefix

	r.addRoute(newRoute(prefix, r.StrictSlash, true).All(router))

	return router
}

func (r *Router) NewRoute(pattern string) *Route {
	route := newRoute(pattern, r.StrictSlash, false)
	r.addRoute(route)
	return route
}

func (r *Router) Path() string {
	return r.path
}

func (r *Router) ServeHTTP(res ResponseWriter, req *Request) {
	r.invokeHandlers(res, req)

	if res.Written() || r.ErrorHandler == nil {
		return
	}

	err := ErrNotFound
	if e := res.Error(); e != nil {
		err = e
	}

	r.ErrorHandler.ServeHTTP(res, req, err)
}

func (r *Router) invokeHandlers(res ResponseWriter, req *Request) {
	path := req.SanitizedPath()[len(r.path):] // Strip own prefix

	paramInvokedMem := make(emptyKeyMap)

	for _, route := range r.routes {
		if !route.Match(path) {
			continue
		}

		route.FillParams(req)
		if !r.handleParams(res, req, route.params, paramInvokedMem) {
			return
		}

		route.ServeHTTP(res, req)

		if res.Error() != nil {
			return
		}

		if res.Written() {
			return
		}
	}
}

func (r *Router) handleParams(res ResponseWriter, req *Request, orderedParams []string, invoked emptyKeyMap) bool {
	// Call param handlers in the same order in which the parameters appear in the path.
	for _, name := range orderedParams {
		if _, ok := invoked[name]; ok {
			continue
		}

		value := req.Params.Get(name)

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

func (r *Router) addRoute(route *Route) *Router {
	r.routes = append(r.routes, route)
	return r
}

func newRouter() *Router {
	return &Router{paramHandlers: make(paramHandlerMap)}
}

type paramHandlerMap map[string][]ParamHandler
type emptyKeyMap map[string]empty
