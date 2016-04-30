// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"encoding/json"
	"io"
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
func WrapHTTPHandler(handler http.Handler) HandlerFunc {
	return HandlerFunc(func(res ResponseWriter, req *Request) {
		handler.ServeHTTP(res, req.Request)
	})
}

// WrapHTTPHandlerFunc wraps ordinary functions with the http.HandlerFunc
// format in a Handler.
func WrapHTTPHandlerFunc(fn func(w http.ResponseWriter, r *http.Request)) HandlerFunc {
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

// WriteJSON writes the passed value as JSON to the ResponseWriter utilizing the
// encoding/json package. It also sets the Content-Type header to "application/json".
// Any errors occured during encoding are returned.
func WriteJSON(w ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}

	return nil
}

// WriteString writes the s to the ResponseWriter utilizing io.WriteString. It also
// sets the Content-Type to "text/plain; charset=utf8".
// Any errors occured during Write are returned.
func WriteString(w ResponseWriter, s string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := io.WriteString(w, s); err != nil {
		return err
	}

	return nil
}

// Returns true if either a response was written or a ContextError occured.
func doneProcessing(w ResponseWriter, ctx *RequestContext) bool {
	return w.Written() || ctx.err != nil
}
