// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"sync"
)

// RequestContext allows sharing of data between handlers by storing
// key-value pairs of arbitrary types. It also provides the captured
// URL parameter values depending on the current route.
//
// Any occuring errors during the processing of handlers can be
// set on the RequestContext using .Error. By setting an error
// all Routers and Routes will stop processing immediately and the
// error is passed to the next error handler.
type RequestContext struct {
	mutex  sync.RWMutex
	store  anyMap
	params params
	err    *ContextError
	skip   bool
}

// Set sets the value for the specified the key. It replaces any existing values.
func (r *RequestContext) Set(key string, value interface{}) {
	r.mutex.Lock()
	r.store[key] = value
	r.mutex.Unlock()
}

// Get retrieves the value for key. If the key doesn't exist in the RequestContext,
// Get returns nil.
func (r *RequestContext) Get(key string) interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.store[key]
}

// Delete deletes the value associated with key. If the key doesn't exist nothing happens.
func (r *RequestContext) Delete(key string) {
	r.mutex.Lock()
	delete(r.store, key)
	r.mutex.Unlock()
}

// Exists returns true if the specified key exists in the RequestContext, otherwise false is returned.
func (r *RequestContext) Exists(key string) bool {
	r.mutex.RLock()
	_, exists := r.store[key]
	r.mutex.RUnlock()
	return exists
}

// Param returns the capture URL parameter value for the given parameter name. The name is
// the one specified in one of the routing functions without the leading ":".
func (r *RequestContext) Param(name string) string {
	return r.params[name]
}

// Error sets a ContextError which will be passed to the next error handler and
// forces all Routers and Routes to stop processing.
//
// Note: Calling Error from different threads can cause race conditions. Also
// calling Error more than once causes a runtime panic!
func (r *RequestContext) Error(err error, code int) {
	if r.err != nil {
		panic("RequestContext: called .Error() twice")
	}
	r.err = &ContextError{err, code}
}

// SkipRouter tells the current router to end processing, which means that the
// parent router will continue processing.
//
// Calling SkipRouter in a top level route causes a "not found" error.
func (r *RequestContext) SkipRouter() {
	r.skip = true
}

func (r *RequestContext) skipped() {
	r.skip = false
}

func newRequestContext() *RequestContext {
	return &RequestContext{
		store:  make(anyMap),
		params: make(params),
		err:    nil,
		skip:   false,
	}
}

// Stores a RequestContext for each Request.
var contextMutex sync.RWMutex
var requestContextMap = make(map[*http.Request]*RequestContext)

// Context returns the corresponding RequestContext for the given Request.
func Context(r *http.Request) *RequestContext {
	contextMutex.RLock()
	defer contextMutex.RUnlock()
	return requestContextMap[r]
}

// Stores a new RequestContext for the specified Request in the requestContextMap.
// This may overwrite an existing RequestContext!
func createRequestContext(r *http.Request) {
	contextMutex.Lock()
	requestContextMap[r] = newRequestContext()
	contextMutex.Unlock()
}

// Removes the RequestContext for the given Request from the requestContextMap.
func deleteRequestContext(r *http.Request) {
	contextMutex.Lock()
	delete(requestContextMap, r)
	contextMutex.Unlock()
}

type params map[string]string
type anyMap map[string]interface{}
