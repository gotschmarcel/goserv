// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	gopath "path"
	"strings"
)

var methodNames = []string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

// WrapHTTPHandler wraps a native http.Handler in a Handler.
func WrapHTTPHandler(handler http.Handler) Handler {
	return HandlerFunc(func(res ResponseWriter, req *Request) {
		handler.ServeHTTP(res, req.Request)
	})
}

// WrapHTTPHandlerFunc wraps ordinary functions with the http.HandlerFunc
// format in a Handler.
func WrapHTTPHandlerFunc(fn func(w http.ResponseWriter, r *http.Request)) Handler {
	return HandlerFunc(func(res ResponseWriter, req *Request) {
		fn(res, req.Request)
	})
}

// SanitizePath returns the clean version of the specified path.
//
// It prepends a "/" to the path if none was found, uses path.Clean to resolve
// any "." and ".." and adds back any trailing slashes.
func SanitizePath(p string) string {
	if len(p) == 0 {
		return "/"
	}

	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	trailingSlash := strings.HasSuffix(p, "/")
	p = gopath.Clean(p)

	if p != "/" && trailingSlash {
		p += "/"
	}

	return p
}
