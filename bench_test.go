// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"testing"
)

func BenchmarkRouter(b *testing.B) {
	router := NewRouter()
	router.ErrorHandler = nil

	router.GetFunc("/v1/users/:id", func(ResponseWriter, *Request) {})

	req, _ := http.NewRequest(http.MethodGet, "/v1/users/123456", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, req)
	}
}

func BenchmarkRouterManyParams(b *testing.B) {
	router := NewRouter()
	router.ErrorHandler = nil

	router.GetFunc("/v1/:p1/:p2/:p3/:p4/:p5", func(ResponseWriter, *Request) {})

	req, _ := http.NewRequest(http.MethodGet, "/v1/1/2/3/4/5", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, req)
	}
}

func BenchmarkNestedRouter(b *testing.B) {
	handler := func(ResponseWriter, *Request) {}
	r1 := NewRouter()
	r1.ErrorHandler = nil
	r1.Router("/v2").Router("/v3").Router("/v4").Router("/v5").GetFunc("/1", handler)

	req, _ := http.NewRequest(http.MethodGet, "/v2/v3/v4/v5/1", nil)
	for i := 0; i < b.N; i++ {
		r1.ServeHTTP(nil, req)
	}
}

func BenchmarkNestedRouterWithParams(b *testing.B) {
	handler := func(ResponseWriter, *Request) {}
	r1 := NewRouter()
	r1.ErrorHandler = nil
	r1.Router("/v2").Router("/v3").Router("/v4").Router("/v5").GetFunc("/:id", handler)

	req, _ := http.NewRequest(http.MethodGet, "/v2/v3/v4/v5/1", nil)
	for i := 0; i < b.N; i++ {
		r1.ServeHTTP(nil, req)
	}
}
