// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrNotFound is passed to the error handler if
	// no route matched the request path or none of the matching routes wrote
	// a response.
	ErrNotFound = errors.New(http.StatusText(http.StatusNotFound))

	// ErrDisallowedHost is passed to the error handler if a handler
	// created with .AllowedHosts() found a disallowed host.
	ErrDisallowedHost = errors.New("disallowed host")
)

// StdErrorHandler is the default ErrorHandler added to all Server instances
// created with NewServer().
var StdErrorHandler = func(w http.ResponseWriter, r *http.Request, err *ContextError) {
	w.WriteHeader(err.Code)
	fmt.Fprintf(w, err.Error())
}

// A ContextError stores an error along with a response code usually in the range
// 4xx or 5xx. The ContextError is passed to the ErrorHandler.
type ContextError struct {
	Err  error
	Code int
}

// Error returns the result of calling .Error() on the stored error.
func (c *ContextError) Error() string {
	return c.Err.Error()
}

// String returns a formatted string with this format: (<code>) <error>.
func (c *ContextError) String() string {
	return fmt.Sprintf("(%d) %s", c.Code, c.Err)
}
