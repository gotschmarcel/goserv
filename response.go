// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"encoding/json"
	"net/http"
)

// A ResponseWriter is used to build an HTTP response.
//
// The ResponseWriter embeds the native http.ResponseWriter, thus
// all fields are still available through ResponseWriter. Additionally
// to the native http.ResponseWriter a ResponseWriter implements
// functions to check if the body was written, the current status or
// if an error was set.
//
// It also provides helper functions to make JSON responses or the
// use of templates even easier.
type ResponseWriter interface {
	// Embedded http.ResponseWriter interface.
	http.ResponseWriter

	// Written returns true if the response header was written.
	Written() bool

	// Returns the status code written to the response header.
	// If no status was written 0 is returned.
	Status() int

	// Returns any error set with SetError() or nil if none was set.
	Error() error

	// Sets a response error which will be passed to the ErrorHandler
	// by the Server/Router.
	SetError(error)

	// Render renders the template with the given name and locals.
	// The result is written to the body.
	//
	// Render works only if a TemplateEngine is registered to the
	// main Server and results in a panic otherwise.
	Render(string, interface{})

	// JSON sends a JSON response by encoding the input using json.Encoder from encoding/json.
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
