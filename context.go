// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

type Context struct {
	store map[string]interface{}
}

func (c *Context) Set(key string, value interface{}) {
	c.store[key] = value
}

func (c *Context) Get(key string) interface{} {
	return c.store[key]
}

func (c *Context) Delete(key string) {
	delete(c.store, key)
}

func (c *Context) Exists(key string) bool {
	_, exists := c.store[key]
	return exists
}

func newContext() *Context {
	return &Context{make(map[string]interface{})}
}
