// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrNotFound gets passed to a Router's ErrorHandler if
// no route matched the request path or none of the matching routes wrote
// a response.
var ErrNotFound = errors.New(http.StatusText(http.StatusNotFound))

// StdErrorHandler is the default ErrorHandler added to all Server instances
// created with NewServer().
//
// All errors, except ErrNotFound, passed to it result in an internal server error (500) including
// the message in the response body. The ErrNotFound error results
// in a "not found" (404) response.
var StdErrorHandler = func(res ResponseWriter, req *Request, err error) {
	status := http.StatusInternalServerError

	if err == ErrNotFound {
		status = http.StatusNotFound
	}

	res.WriteHeader(status)
	fmt.Fprintf(res, err.Error())
}
