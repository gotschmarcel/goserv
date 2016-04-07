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

func (s *Server) Listen(addr string, tls *TLS) error {
	if tls != nil {
		s.TLS = tls
		return http.ListenAndServeTLS(addr, tls.CertFile, tls.KeyFile, s)
	}

	return http.ListenAndServe(addr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := &responseWriter{w: w}
	req := &Request{r, &Context{}, nil, nil, sanitizePath(r.URL.Path)}

	s.Router.serveHTTP(res, req)
}

func NewServer() *Server {
	s := &Server{NewRouter(), "", nil}
	s.ErrorHandler = StdErrorHandler

	return s
}
