// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
)

// A Request represents an HTTP request received by the Server.
//
// It embeds the native http.Request, thus all native fields are still available
// through Request. Every Request has it's own Context providing a key-value store to share
// data between multiple Handlers. In case that the Route handling the Request has parameters, the parameter
// values are extracted from the Request's path and stored in .Params.
type Request struct {
	// Embedded http.Request.
	*http.Request

	// Request specific key-value store to share data between Handlers
	Context *Context

	// Key-value store containing named parameter values extracted from
	// the Request's path. See Route.
	Params Params

	sanitizedPath string
}

// SanitizedPath returns the Request's path sanitized with SanitizePath().
func (r *Request) SanitizedPath() string {
	return r.sanitizedPath
}

func newRequest(r *http.Request) *Request {
	return &Request{r, newContext(), make(Params), SanitizePath(r.URL.Path)}
}
