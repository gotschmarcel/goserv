// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import "net/http"

// A ErrorHandlerFunc is the last handler in the request chain and
// is responsible for handling errors that occur during the
// request processing.
//
// A ErrorHandlerFunc should always write a response!
type ErrorHandlerFunc func(http.ResponseWriter, *http.Request, *ContextError)

// A ParamHandlerFunc can be registered to a Router using a parameter's name.
// It gets invoked with the corresponding value extracted from the request's
// path.
//
// Parameters are part of a Route's path. To learn more about parameters take
// a look at the documentation of Route.
type ParamHandlerFunc func(http.ResponseWriter, *http.Request, string)
