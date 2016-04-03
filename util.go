// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"path"
	"strings"
)

func nativeWrapper(handler http.Handler) HandlerFunc {
	return HandlerFunc(func(res ResponseWriter, req *Request) {
		handler.ServeHTTP(http.ResponseWriter(res), req.Request)
	})
}

func sanitizePath(p string) string {
	if len(p) == 0 {
		return "/"
	}

	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	trailingSlash := strings.HasSuffix(p, "/")
	p = path.Clean(p)

	if p != "/" && trailingSlash {
		p += "/"
	}

	return p
}
