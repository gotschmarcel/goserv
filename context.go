// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

type Context struct {
	store storage
}

func (c *Context) Set(key string, value interface{}) {
	c.assureStorage()
	c.store[key] = value
}

func (c *Context) Get(key string) interface{} {
	c.assureStorage()
	return c.store[key]
}

func (c *Context) assureStorage() {
	if c.store != nil {
		return
	}

	c.store = make(storage)
}

type storage map[string]interface{}
