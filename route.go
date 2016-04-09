// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"regexp"
)

type Route struct {
	middleware []Handler
	methods    map[string][]Handler
	matcher    *regexp.Regexp
	params     []string
}

func (r *Route) All(handlers ...Handler) *Route {
	r.middleware = append(r.middleware, handlers...)
	return r
}

func (r *Route) AllFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	for _, fn := range funcs {
		r.middleware = append(r.middleware, HandlerFunc(fn))
	}
	return r
}

func (r *Route) Method(method string, handlers ...Handler) *Route {
	r.addMethodHandlers(method, handlers...)
	return r
}

func (r *Route) MethodFunc(method string, funcs ...func(ResponseWriter, *Request)) *Route {
	for _, fn := range funcs {
		r.addMethodHandlers(method, HandlerFunc(fn))
	}
	return r
}

func (r *Route) Methods(methods []string, handlers ...Handler) *Route {
	for _, method := range methods {
		r.Method(method, handlers...)
	}
	return r
}

func (r *Route) MethodsFunc(methods []string, funcs ...func(ResponseWriter, *Request)) *Route {
	for _, method := range methods {
		r.MethodFunc(method, funcs...)
	}
	return r
}

func (r *Route) Get(handlers ...Handler) *Route {
	return r.Method(http.MethodGet, handlers...)
}

func (r *Route) GetFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodGet, funcs...)
}

func (r *Route) Post(handlers ...Handler) *Route {
	return r.Method(http.MethodPost, handlers...)
}

func (r *Route) PostFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodPost, funcs...)
}

func (r *Route) Put(handlers ...Handler) *Route {
	return r.Method(http.MethodPut, handlers...)
}

func (r *Route) PutFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodPut, funcs...)
}

func (r *Route) Delete(handlers ...Handler) *Route {
	return r.Method(http.MethodDelete, handlers...)
}

func (r *Route) DeleteFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodDelete, funcs...)
}

func (r *Route) Patch(handlers ...Handler) *Route {
	return r.Method(http.MethodPatch, handlers...)
}

func (r *Route) PatchFunc(funcs ...func(ResponseWriter, *Request)) *Route {
	return r.MethodFunc(http.MethodPatch, funcs...)
}

func (r *Route) ServeHTTP(res ResponseWriter, req *Request) {
	if r.ContainsParams() && len(req.Params) == 0 {
		r.FillParams(req)
	}

	for _, handler := range append(r.middleware, r.methods[req.Method]...) {
		handler.ServeHTTP(res, req)

		if res.Error() != nil {
			return
		}

		if res.Written() {
			return
		}
	}
}

func (r *Route) Match(path string) bool {
	return r.matcher.MatchString(path)
}

func (r *Route) FillParams(req *Request) {
	if !r.ContainsParams() {
		return
	}

	matches := r.matcher.FindAllStringSubmatch(req.SanitizedPath(), -1)
	if len(matches) == 0 {
		return
	}

	// Iterate group matches only
	for index, value := range matches[0][1:] {
		name := r.params[index]
		req.Params[name] = value
	}
}

func (r *Route) ContainsParams() bool {
	return len(r.params) > 0
}

func (r *Route) addMethodHandlers(method string, handlers ...Handler) {
	r.methods[method] = append(r.methods[method], handlers...)
}

func newRoute(pattern string, strict, prefixOnly bool) *Route {
	matcher, params := pathComponents(pattern, strict, prefixOnly)

	return &Route{
		methods: make(map[string][]Handler),
		matcher: matcher,
		params:  params,
	}
}
