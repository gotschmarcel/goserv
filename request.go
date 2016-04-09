// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
)

type Request struct {
	*http.Request
	Context       *Context
	Locals        Locals
	Params        Params
	sanitizedPath string
}

func (r *Request) SanitizedPath() string {
	return r.sanitizedPath
}

func newRequest(r *http.Request) *Request {
	return &Request{r, newContext(), nil, make(Params), SanitizePath(r.URL.Path)}
}
