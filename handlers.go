// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"strings"
)

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

// AddHeaders returns a new HandlerFunc which adds the specified response headers.
func AddHeaders(headers map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()

		for name, value := range headers {
			h.Add(name, value)
		}
	}
}

// AllowedHosts returns a new HandlerFunc validating the HTTP Host header.
//
// Values can be fully qualified (e.g. "www.example.com"), in which case they
// must match the Host header exactly. Values starting with a period
// (e.g. ".example.com") will match example.com and all subdomains (e.g. www.example.com)
//
// If useXForwardedHost is true the X-Forwarded-Host header will be used in preference
// to the Host header. This is only useful if a proxy which sets the header is in use.
func AllowedHosts(hosts []string, useXForwardedHost bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host := r.Host

		if useXForwardedHost {
			host = r.Header.Get("X-Forwarded-Host")
		}

		for _, allowedHost := range hosts {
			if strings.HasPrefix(allowedHost, ".") && strings.HasSuffix(host, allowedHost) {
				return
			}

			if host == allowedHost {
				return
			}
		}

		Context(r).Error(ErrDisallowedHost, http.StatusBadRequest)
	}
}
