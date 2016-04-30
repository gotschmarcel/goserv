// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovery(t *testing.T) {
	server := NewServer()
	server.PanicRecovery = true
	expectedErr := "Panic: I am panicked"

	server.Get("/", func(res ResponseWriter, req *Request) {
		panic("I am panicked")
	})

	var err error
	server.ErrorHandler = func(res ResponseWriter, req *Request, e error) {
		err = e
	}

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, r)

	if err == nil {
		t.Fatal("Error expected")
	}

	errMsg := err.Error()
	if !strings.HasPrefix(errMsg, "Panic") {
		t.Error("Expected error to have prefix 'Panic'")
	}

	if errMsg != expectedErr {
		t.Errorf("Expected '%s', not '%s'", expectedErr, errMsg)
	}

}

func TestStatic(t *testing.T) {
	tests := []struct {
		prefix string
		path   string
		code   int
	}{
		{"/", "/server.go", http.StatusOK},
		{"/", "/nonexisting.go", http.StatusNotFound},

		{"/public", "/public/server.go", http.StatusOK},
		{"/public", "/public/./server.go", http.StatusOK},
		{"/public", "/public/folder/../server.go", http.StatusOK},
		{"/public", "/public/nonexisting.go", http.StatusNotFound},
		{"/public", "/public/../server.go", http.StatusNotFound},
	}

	root := http.Dir(".")

	for idx, test := range tests {
		s := NewServer()
		s.Static(test.prefix, root)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, test.path, nil)

		s.ServeHTTP(w, r)

		if w.Code != test.code {
			t.Errorf("Expected code %d, is %d (test no. %d)", test.code, w.Code, idx)
		}

		if test.code == http.StatusOK && w.Body.Len() == 0 {
			t.Errorf("Expected non-empty body (test no. %d)", idx)
		}
	}
}

func TestServerContext(t *testing.T) {
	server := NewServer()

	server.Use(func(res ResponseWriter, req *Request) {
		ctx := Context(req)
		ctx.Set("test_key", "test_value")
	})

	server.Use(func(res ResponseWriter, req *Request) {
		ctx := Context(req)

		if !ctx.Exists("test_key") {
			t.Fatal("Missing key: test_key")
		}

		if v, ok := ctx.Get("test_key").(string); !ok || v != "test_value" {
			t.Errorf("Wrong key value, wanted: %q, got: %q", "test_value", v)
		}
	})
}
