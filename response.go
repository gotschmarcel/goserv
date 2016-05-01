// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
)

type responseWriter struct {
	w      http.ResponseWriter
	status int
}

func (r *responseWriter) Header() http.Header {
	return r.w.Header()
}

func (r *responseWriter) Write(b []byte) (int, error) {
	if !r.Written() {
		r.WriteHeader(http.StatusOK)
	}

	return r.w.Write(b)
}

func (r *responseWriter) WriteHeader(status int) {
	r.status = status
	r.w.WriteHeader(status)
}

func (r *responseWriter) Written() bool {
	return r.status != 0
}

func (r *responseWriter) Code() int {
	return r.status
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w: w}
}
