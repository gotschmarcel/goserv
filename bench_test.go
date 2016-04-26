// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"testing"
)

func BenchmarkServer(b *testing.B) {
	server := NewServer()
	server.ErrorHandler = nil

	server.Get("/v1/users/:id", func(ResponseWriter, *Request) {})

	req, _ := http.NewRequest(http.MethodGet, "/v1/users/123456", nil)
	for i := 0; i < b.N; i++ {
		server.ServeHTTP(nil, req)
	}
}

func BenchmarkServerManyParams(b *testing.B) {
	server := NewServer()
	server.ErrorHandler = nil

	server.Get("/v1/:p1/:p2/:p3/:p4/:p5", func(ResponseWriter, *Request) {})

	req, _ := http.NewRequest(http.MethodGet, "/v1/1/2/3/4/5", nil)
	for i := 0; i < b.N; i++ {
		server.ServeHTTP(nil, req)
	}
}

func BenchmarkNestedRouter(b *testing.B) {
	handler := func(ResponseWriter, *Request) {}
	s := NewServer()
	s.ErrorHandler = nil
	s.SubRouter("/v2").SubRouter("/v3").SubRouter("/v4").SubRouter("/v5").Get("/1", handler)

	req, _ := http.NewRequest(http.MethodGet, "/v2/v3/v4/v5/1", nil)
	for i := 0; i < b.N; i++ {
		s.ServeHTTP(nil, req)
	}
}

func BenchmarkNestedRouterWithParams(b *testing.B) {
	handler := func(ResponseWriter, *Request) {}
	s := NewServer()
	s.ErrorHandler = nil
	s.SubRouter("/v2").SubRouter("/v3").SubRouter("/v4").SubRouter("/v5").Get("/:id", handler)

	req, _ := http.NewRequest(http.MethodGet, "/v2/v3/v4/v5/1", nil)
	for i := 0; i < b.N; i++ {
		s.ServeHTTP(nil, req)
	}
}
