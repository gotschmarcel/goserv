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
// through Request.
type Request struct {
	// Embedded http.Request.
	*http.Request

	sanitizedPath string
}

func newRequest(r *http.Request) *Request {
	return &Request{r, SanitizePath(r.URL.Path)}
}
