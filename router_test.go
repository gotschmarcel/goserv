// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type historyHandler struct {
	*historyWriter
}

func (h historyHandler) Handler(id string) Handler {
	return HandlerFunc(func(ResponseWriter, *Request) {
		h.WriteString(id)
	})
}

func (h historyHandler) WriteHandler(id string) Handler {
	return HandlerFunc(func(w ResponseWriter, r *Request) {
		h.WriteString(id)
		w.Write([]byte(id))
	})
}

func (h historyHandler) ParamHandler() ParamHandler {
	return ParamHandlerFunc(func(w ResponseWriter, r *Request, value string) {
		h.WriteString(value)
	})
}

func (h historyHandler) HandlerWithError(v string) Handler {
	return HandlerFunc(func(w ResponseWriter, r *Request) {
		w.SetError(fmt.Errorf(v))
	})
}

func newHistoryHandler() *historyHandler {
	return &historyHandler{&historyWriter{}}
}

func TestRouter(t *testing.T) {
	h := newHistoryHandler()

	router := newRouter()

	// Register middleware
	router.Use(h.Handler("middleware"))

	// Register all handlers
	router.All("/", h.Handler("all-handler"))

	// Register simple method path
	router.Get("/", h.WriteHandler("get-handler"))

	// Register multi-method path
	router.Methods([]string{http.MethodGet, http.MethodDelete}, "/multi", h.WriteHandler("multi-handler"))

	// Register route
	router.NewRoute("/route").Get(h.WriteHandler("route-handler"))

	// Register parameter routes
	router.Get("/param1/:value1/param2/:value2", h.Handler("param-route1"))
	router.Get("/param1/:value1/param2/:value2", h.WriteHandler("param-route2"))
	router.Param("value1", h.ParamHandler())
	router.Param("value2", h.ParamHandler())

	// Register sub routers
	sr1 := router.NewRouter("/srouter1").Get("/get", h.WriteHandler("srouter1-handler"))
	sr1.NewRouter("/srouter2").Get("/get", h.WriteHandler("srouter2-handler")).Get("/error", h.HandlerWithError("srouter2-error"))

	// Register prefix route
	router.Prefix("/prefix", h.WriteHandler("prefix-route"))

	tests := []struct {
		method string
		path   string
		writes []string
		body   string
		err    error
	}{
		{http.MethodGet, "/", []string{"middleware", "all-handler", "get-handler"}, "get-handler", nil},
		{http.MethodPost, "/", []string{"middleware", "all-handler"}, "", ErrNotFound},

		{http.MethodGet, "/multi", []string{"middleware", "multi-handler"}, "multi-handler", nil},
		{http.MethodDelete, "/multi", []string{"middleware", "multi-handler"}, "multi-handler", nil},
		{http.MethodPost, "/multi", []string{"middleware"}, "", ErrNotFound},

		{http.MethodGet, "/route", []string{"middleware", "route-handler"}, "route-handler", nil},

		{http.MethodGet, "/param1/123/param2/456", []string{"middleware", "123", "456", "param-route1", "param-route2"}, "param-route2", nil},

		{http.MethodGet, "/srouter1/get", []string{"middleware", "srouter1-handler"}, "srouter1-handler", nil},
		{http.MethodGet, "/srouter1/srouter2/get", []string{"middleware", "srouter2-handler"}, "srouter2-handler", nil},
		{http.MethodGet, "/srouter1/srouter2/error", []string{"middleware"}, "", fmt.Errorf("srouter2-error")},

		{http.MethodGet, "/prefix", []string{"middleware", "prefix-route"}, "prefix-route", nil},
		{http.MethodGet, "/prefix/more", []string{"middleware", "prefix-route"}, "prefix-route", nil},
	}

	for index, test := range tests {
		w := httptest.NewRecorder()
		r := &http.Request{Method: test.method}
		var err error

		r.URL, _ = url.Parse(test.path)
		h.Clear()

		router.ErrorHandler = ErrorHandlerFunc(func(res ResponseWriter, req *Request, e error) {
			err = e
		})

		router.ServeHTTP(newResponseWriter(w, nil), newRequest(r))

		if test.err != nil {
			if err == nil {
				t.Errorf("Expected error in ServeHTTP, but there is none (no. %d)", index)
			} else if test.err.Error() != err.Error() {
				t.Errorf("Wrong error message: %s != %s", err.Error(), test.err.Error())
			}
		}

		if test.err == nil && err != nil {
			t.Errorf("Unexpected error in ServeHTTP: %v (no. %d)", err, index)
		}

		if w.Body.String() != test.body {
			t.Errorf("Wrong body: %s != %s (no. %d)", w.Body.String(), test.body, index)
		}

		if len(test.writes) != h.Len() {
			t.Errorf("Wrong write count %d != %d, %v (no. %d)", h.Len(), len(test.writes), h.writes, index)
			continue
		}

		for index, value := range test.writes {
			writeValue := h.At(index)
			if value != writeValue {
				t.Errorf("Wrong write value at %d: %s != %s (no. %d)", index, writeValue, value, index)
			}
		}
	}
}
