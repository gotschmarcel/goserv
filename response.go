// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	Written() bool
	Status() int
	Error() error
	SetError(error)
}

type responseWriter struct {
	w      http.ResponseWriter
	status int
	err    error
}

func (r *responseWriter) Header() http.Header {
	return r.w.Header()
}

func (r *responseWriter) Write(b []byte) (int, error) {
	r.status = http.StatusOK
	return r.w.Write(b)
}

func (r *responseWriter) WriteHeader(status int) {
	r.status = status
	r.w.WriteHeader(status)
}

func (r *responseWriter) Written() bool {
	return r.status != 0
}

func (r *responseWriter) Status() int {
	return r.status
}

func (r *responseWriter) Error() error {
	return r.err
}

func (r *responseWriter) SetError(err error) {
	if r.err != nil {
		panic("error set twice")
	}

	r.err = err
}

func newResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{w: w}
}
