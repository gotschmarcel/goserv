// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"fmt"
	"io"
	"net/http"
	"path"
)

type TLS struct {
	CertFile, KeyFile string
}

type Server struct {
	*Router
	Addr           string
	TLS            *TLS
	ViewRoot       string
	TemplateEngine TemplateEngine
	PanicRecovery  bool
}

func (s *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) ListenTLS(addr, certFile, keyFile string) error {
	s.TLS = &TLS{certFile, keyFile}
	return http.ListenAndServeTLS(addr, certFile, keyFile, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := newResponseWriter(w, s)
	req := newRequest(r)

	if s.PanicRecovery {
		defer s.handleRecovery(res, req)
	}

	s.Router.ServeHTTP(res, req)
}

func (s *Server) Static(prefix string, dir http.Dir) {
	s.Prefix(prefix, WrapHTTPHandler(http.StripPrefix(prefix, http.FileServer(dir))))
}

func (s *Server) renderView(w io.Writer, name string, locals interface{}) error {
	if s.TemplateEngine == nil {
		panic("template engine not set")
	}

	filePath := path.Join(s.ViewRoot, name) + s.TemplateEngine.Ext()
	return s.TemplateEngine.RenderAndWrite(w, filePath, locals)
}

func (s *Server) handleRecovery(res ResponseWriter, req *Request) {
	if r := recover(); r != nil {
		s.ErrorHandler.ServeHTTP(res, req, fmt.Errorf("Panic: %v", r))
	}
}

func NewServer() *Server {
	s := &Server{
		Router:         newRouter(),
		Addr:           "",
		TLS:            nil,
		ViewRoot:       "",
		TemplateEngine: nil,
		PanicRecovery:  false,
	}

	s.ErrorHandler = StdErrorHandler

	return s
}
