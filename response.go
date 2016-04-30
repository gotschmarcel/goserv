// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"encoding/json"
	"io"
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

	// WriteJSON sends a JSON response by encoding the input using json.Encoder from encoding/json.
	WriteJSON(interface{}) error

	// WriteString sends a simple plain text response.
	WriteString(string) error

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

func (r *responseWriter) WriteJSON(v interface{}) error {
	r.w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(r).Encode(v); err != nil {
		return err
	}

	return nil
}

func (r *responseWriter) WriteString(data string) error {
	if _, err := io.WriteString(r, data); err != nil {
		return err
	}

	return nil
}

func (r *responseWriter) Redirect(req *Request, url string, code int) {
	http.Redirect(r, req.Request, url, code)
}

func newResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{w: w}
}
