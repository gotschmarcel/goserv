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
	req := &Request{&http.Request{Method: http.MethodGet}, "/"}
	history := &historyWriter{}

	route := newRoute("/", false, false)

	route.All(func(ResponseWriter, *Request) {
		history.WriteString("1")
	})

	route.All(func(ResponseWriter, *Request) {
		history.WriteString("2")
	})

	route.Get(func(ResponseWriter, *Request) {
		history.WriteString("3")
		res.Write([]byte("Done"))
	})

	route.Get(func(ResponseWriter, *Request) {
		history.WriteString("4")
	})

	route.Put(func(ResponseWriter, *Request) {
		history.WriteString("5")
	})

	route.ServeHTTP(res, req)
	if err := res.Error(); err != nil {
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

	if w.Body.String() != "Done" {
		t.Errorf("Wrong body content: %s != %s", w.Body.String(), "Done")
	}
}
