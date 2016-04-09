// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
)

type TLS struct {
	CertFile, KeyFile string
}

type Server struct {
	*Router
	Addr string
	TLS  *TLS
}

func (s *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) ListenTLS(addr, certFile, keyFile string) error {
	s.TLS = &TLS{certFile, keyFile}
	return http.ListenAndServeTLS(addr, certFile, keyFile, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(newResponseWriter(w), newRequest(r))
}

func (s *Server) Static(prefix string, dir http.Dir) {
	s.Prefix(prefix, WrapHTTPHandler(http.StripPrefix(prefix, http.FileServer(dir))))
}

func NewServer() *Server {
	s := &Server{NewRouter(), "", nil}
	s.ErrorHandler = StdErrorHandler

	return s
}
