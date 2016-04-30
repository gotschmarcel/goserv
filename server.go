// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
)

// A TLS contains both the certificate and key file paths.
type TLS struct {
	CertFile, KeyFile string
}

// A Server is the main instance and entry point for all routing.
//
// It is compatible with the http package an can be used as a http.Handler.
// A Server is also a Router and provides the same fields and methods as the
// goserv.Router.
//
// Additionally to all routing methods a Server provides methods to register
// static file servers, short-hand methods for
// http.ListenAndServe as well as http.ListenAndServeTLS and the possibility
// to recover from panics.
//
type Server struct {
	// Embedded Router
	*Router

	// TCP address to listen on, set by .Listen or .ListenTLS
	Addr string

	// TLS information set by .ListenTLS or nil if .Listen was used
	TLS *TLS
}

// Listen is a convenience method that uses http.ListenAndServe.
func (s *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, s)
}

// ListenTLS is a convenience method that uses http.ListenAndServeTLS.
// The TLS informations used are stored in .TLS after calling this method.
func (s *Server) ListenTLS(addr, certFile, keyFile string) error {
	s.TLS = &TLS{certFile, keyFile}
	return http.ListenAndServeTLS(addr, certFile, keyFile, s)
}

// ServeHTTP dispatches the request to the internal Router.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	iw := newResponseWriter(w)
	createRequestContext(r)

	s.serveHTTP(iw, r)
}

// NewServer returns a newly allocated and initialized Server instance.
//
// By default the Server has no template engine, the template root is "" and
// panic recovery is disabled. The Router's ErrorHandler is set to the StdErrorHandler.
func NewServer() *Server {
	s := &Server{
		Router: newRouter(),
		Addr:   "",
		TLS:    nil,
	}

	s.ErrorHandler = StdErrorHandler

	return s
}
