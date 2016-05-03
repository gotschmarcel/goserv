// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteHandlerChain(t *testing.T) {
	w := httptest.NewRecorder()
	res := &responseWriter{w: w}
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	history := newHistoryHandler()

	createRequestContext(req)
	ctx := Context(req)

	route := newRoute("/", false, false)

	route.All(history.Handler("1"))
	route.All(history.Handler("2"))
	route.Get(history.WriteHandler("3"))
	route.Get(history.Handler("4"))
	route.Put(history.Handler("5"))

	route.serveHTTP(res, req)
	if err := ctx.err; err != nil {
		t.Errorf("Serve error: %v", err)
	}

	if history.Len() != 3 {
		t.Fatalf("Wrong write count: %d != 3", history.Len())
	}

	for index, value := range []string{"1", "2", "3"} {
		if history.At(index) != value {
			t.Errorf("Invalid write: %s != %s", history.At(index), value)
		}
	}

	if res.Code() != http.StatusOK {
		t.Errorf("Wrong status code: %d != %d", res.Code(), http.StatusOK)
	}

	if w.Body.String() != "3" {
		t.Errorf("Wrong body content: %s != %s", w.Body.String(), "3")
	}
}

func TestRouteRest(t *testing.T) {
	w := httptest.NewRecorder()
	res := &responseWriter{w: w}
	req, _ := http.NewRequest("", "/", nil)
	history := newHistoryHandler()

	createRequestContext(req)

	route := newRoute("/", false, false)

	route.Get(history.WriteHandler("get-handler"))
	route.Rest(history.WriteHandler("rest-handler"))

	for _, method := range methodNames {
		req.Method = method

		route.serveHTTP(res, req)

		wanted := "rest-handler"
		if method == http.MethodGet {
			wanted = "get-handler"
		}

		if first := history.At(0); first != wanted {
			t.Errorf("Wrong write value, wanted: %q, got: %q", wanted, first)
		}

		history.Clear()
	}
}
