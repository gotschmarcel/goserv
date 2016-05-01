// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"strings"
)

type fileServer struct {
	prefix string
	index  string
	fs     http.Handler
}

func (f *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Strip prefix
	p := strings.TrimPrefix(r.URL.Path, f.prefix)
	if len(p) >= len(r.URL.Path) {
		return
	}

	if p == "" {
		p += "/"
	}

	// Fallback to index file
	if p == "/" {
		p += f.index
	}

	f.fs.ServeHTTP(w, r)
}

// FileServer returns a thin wrapper around http.FileServer which
// strips the specified prefix from the request path and
// delivers the index file on "/".
func FileServer(root http.Dir, prefix, index string) http.Handler {
	return &fileServer{prefix, index, http.FileServer(root)}
}
