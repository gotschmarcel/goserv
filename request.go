// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"encoding/json"
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

	// Sanitized http.Request.URL.Path
	SanitizedPath string
}

// ReadJSON parses the request's body using the encoding/json Decoder. In case
// of a decoding error the error is returned.
//
// Note: The request's body is closed after calling this method.
func (r *Request) ReadJSON(v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	r.Body.Close()

	if err != nil {
		return err
	}

	return nil
}

func newRequest(r *http.Request) *Request {
	return &Request{r, newContext(), make(Params), SanitizePath(r.URL.Path)}
}
