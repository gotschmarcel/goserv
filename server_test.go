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

	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		panic("I am panicked")
	})

	var err error
	server.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e *ContextError) {
		err = e.Err
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

func TestServerContext(t *testing.T) {
	server := NewServer()

	server.Use(func(w http.ResponseWriter, r *http.Request) {
		ctx := Context(r)
		ctx.Set("test_key", "test_value")
	})

	server.Use(func(w http.ResponseWriter, r *http.Request) {
		ctx := Context(r)

		if !ctx.Exists("test_key") {
			t.Fatal("Missing key: test_key")
		}

		if v, ok := ctx.Get("test_key").(string); !ok || v != "test_value" {
			t.Errorf("Wrong key value, wanted: %q, got: %q", "test_value", v)
		}
	})
}
