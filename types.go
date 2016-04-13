// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

// Params is a key-value store mapping parameter names
// of routes to their extracted values from the request path.
type Params map[string]string

// Get retrieves the value for key and returns an empty string
// if the key doesn't exist.
func (p Params) Get(key string) string { return p[key] }

type empty struct{}
