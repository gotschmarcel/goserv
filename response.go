// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
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
	Code() int

	// Redirect replies to the request with a redirect url. The specified code should
	// be in the 3xx range.
	Redirect(req *Request, url string, code int)
}

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

func (r *responseWriter) Redirect(req *Request, url string, code int) {
	http.Redirect(r, req.Request, url, code)
}

func newResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{w: w}
}
