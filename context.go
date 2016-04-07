// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

type Context struct {
	store map[string]interface{}
}

func (c *Context) Set(key string, value interface{}) {
	c.assureStorage()
	c.store[key] = value
}

func (c *Context) Get(key string) interface{} {
	if c.store == nil {
		return nil
	}

	return c.store[key]
}

func (c *Context) Delete(key string) {
	if c.store == nil {
		return
	}

	delete(c.store, key)
}

func (c *Context) Exists(key string) bool {
	if c.store == nil {
		return false
	}

	_, exists := c.store[key]
	return exists
}

func (c *Context) assureStorage() {
	if c.store != nil {
		return
	}

	c.store = make(map[string]interface{})
}
