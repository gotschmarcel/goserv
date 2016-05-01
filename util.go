// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"encoding/json"
	"fmt"
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
func WriteJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}

	return nil
}

// WriteString writes the s to the ResponseWriter utilizing io.WriteString. It also
// sets the Content-Type to "text/plain; charset=utf8".
// Any errors occured during Write are returned.
func WriteString(w http.ResponseWriter, s string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := io.WriteString(w, s); err != nil {
		return err
	}

	return nil
}

// WriteStringf writes a formatted string to the ResponseWriter utilizing fmt.Fprintf. It also
// sets the Content-Type to "text/plain; charset=utf8".
func WriteStringf(w http.ResponseWriter, format string, v ...interface{}) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, format, v...)
}

// ReadJSONBody decodes the request's body utilizing encoding/json. The body
// is closed after the decoding and any errors occured are returned.
func ReadJSONBody(r *http.Request, result interface{}) error {
	err := json.NewDecoder(r.Body).Decode(result)
	r.Body.Close()

	if err != nil {
		return err
	}

	return nil
}

// Returns true if either a response was written or a ContextError occured.
func doneProcessing(w *responseWriter, ctx *RequestContext) bool {
	return w.Written() || ctx.err != nil
}
