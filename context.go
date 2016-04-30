// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

// Context stores key-value pairs supporting all data types.
type Context struct {
	store map[string]interface{}
}

// Set sets the value for the specified the key. It replaces any existing values.
func (c *Context) Set(key string, value interface{}) {
	c.store[key] = value
}

// Get retrieves the value for key. If the key doesn't exist in the Context,
// Get returns nil.
func (c *Context) Get(key string) interface{} {
	return c.store[key]
}

// Delete deletes the value associated with key. If the key doesn't exist nothing happens.
func (c *Context) Delete(key string) {
	delete(c.store, key)
}

// Exists returns true if the specified key exists in the Context, otherwise false is returned.
func (c *Context) Exists(key string) bool {
	_, exists := c.store[key]
	return exists
}

func newContext() *Context {
	return &Context{make(map[string]interface{})}
}

// Stores a Context for each Request.
var requestContextMap = make(map[*Request]*Context)

// RequestContext returns the corresponding Context for the given Request.
func RequestContext(r *Request) *Context {
	return requestContextMap[r]
}

// Stores a new Context for the specified Request in the requestContextMap.
// This may overwrite an existing Context!
func createRequestContext(r *Request) {
	requestContextMap[r] = newContext()
}
