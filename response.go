// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	Written() bool
	Status() int
	Error() error
	SetError(error)
	Render(string, interface{})
	JSON(interface{})
}

type responseWriter struct {
	w      http.ResponseWriter
	s      *Server
	status int
	err    error
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

func (r *responseWriter) Render(name string, locals interface{}) {
	if err := r.s.renderView(r, name, locals); err != nil {
		r.SetError(err)
	}
}

func (r *responseWriter) JSON(v interface{}) {
	enc := json.NewEncoder(r)

	if err := enc.Encode(v); err != nil {
		r.SetError(err)
	}
}

func newResponseWriter(w http.ResponseWriter, server *Server) ResponseWriter {
	return &responseWriter{w: w, s: server}
}
