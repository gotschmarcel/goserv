// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import "net/http"

type Route struct {
	components *pathComponents
	middleware []Handler
	methods    map[string][]Handler
}

func (r *Route) All(handlers ...Handler) *Route {
	r.middleware = append(r.middleware, handlers...)
	return r
}

func (r *Route) AllNative(handlers ...http.Handler) *Route {
	for _, handler := range handlers {
		r.middleware = append(r.middleware, nativeWrapper(handler))
	}
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

func (r *Route) match(path string) bool {
	return r.components.MatchString(path)
}

func (r *Route) params(path string) Params {
	return r.components.Params(path)
}

func (r *Route) serveHTTP(res ResponseWriter, req *Request) {
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

func (r *Route) addMethodHandlers(method string, handlers ...Handler) {
	r.methods[method] = append(r.methods[method], handlers...)
}

func newRoute(path string, strictSlash bool) (*Route, error) {
	r := &Route{methods: make(map[string][]Handler)}

	pathComponents, err := parsePathString(path, strictSlash)
	if err != nil {
		return nil, err
	}

	r.components = pathComponents
	return r, nil
}
