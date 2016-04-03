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

func TestRouter(t *testing.T) {
	h := &historyWriter{}

	router := NewRouter()
	router.UseFunc(func(res ResponseWriter, req *Request) {
		h.WriteString("1")
	}).AllFunc("/", func(res ResponseWriter, req *Request) {
		h.WriteString("2")
	}).GetFunc("/", func(res ResponseWriter, req *Request) {
		h.WriteString("3")
		res.Write([]byte("Done"))
	}).PostFunc("/", func(res ResponseWriter, req *Request) {
		h.WriteString("4")
		res.SetError(fmt.Errorf("Error in 4"))
	}).GetFunc("/dog", func(res ResponseWriter, req *Request) {
		h.WriteString("5")
		res.Write([]byte("Done"))
	}).GetFunc("/dog/:id", func(res ResponseWriter, req *Request) {
		h.WriteString("6")
	}).GetFunc("/dog/:id", func(res ResponseWriter, req *Request) {
		h.WriteString("7")
		res.Write([]byte(req.Params.Get("id")))
	}).ParamFunc("id", func(res ResponseWriter, req *Request, id string) {
		h.WriteString(id)
	}).Router("/cat").GetFunc("/ape", func(res ResponseWriter, req *Request) {
		h.WriteString("8")
		res.Write([]byte("Done"))
	}).GetFunc("/error", func(res ResponseWriter, req *Request) {
		h.WriteString("9")
		res.SetError(fmt.Errorf("Error in 9"))
	})

	tests := []struct {
		method string
		path   string
		writes []string
		body   string
		err    error
	}{
		{http.MethodGet, "/", []string{"1", "2", "3"}, "Done", nil},
		{http.MethodPost, "/", []string{"1", "2", "4"}, "", fmt.Errorf("Error in 4")},
		{http.MethodGet, "/dog", []string{"1", "5"}, "Done", nil},
		{http.MethodGet, "/dog/123456", []string{"1", "123456", "6", "7"}, "123456", nil},
		{http.MethodGet, "/cat/ape", []string{"1", "8"}, "Done", nil},
		{http.MethodGet, "/cat/error", []string{"1", "9"}, "", fmt.Errorf("Error in 7")},
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

		router.ServeHTTP(w, r)

		if test.err != nil && err == nil {
			t.Errorf("Expected error in ServeHTTP, but there is none (no. %d)", index)
		}

		if test.err == nil && err != nil {
			t.Errorf("Unexpected error in ServeHTTP: %v (no. %d)", err, index)
		}

		if w.Body.String() != test.body {
			t.Errorf("Wrong body: %s != %s (no. %d)", w.Body.String(), test.body, index)
		}

		if len(test.writes) != h.Len() {
			t.Errorf("Wrong write count %d != %d, %v (no. %d)", len(test.writes), h.Len(), h.writes, index)
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
