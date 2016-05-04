// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goserv

import (
	"fmt"
	"net/http"
)

type historyWriter struct{ writes []string }

func (h *historyWriter) Write(b []byte) (int, error) {
	h.WriteString(string(b))
	return len(b), nil
}
func (h *historyWriter) WriteString(s string) {
	h.writes = append(h.writes, s)
}
func (h *historyWriter) Contains(v string) (bool, int) {
	for index, b := range h.writes {
		if v == b {
			return true, index
		}
	}

	return false, -1
}
func (h *historyWriter) At(pos int) string { return h.writes[pos] }
func (h *historyWriter) Len() int          { return len(h.writes) }
func (h *historyWriter) Clear()            { h.writes = nil }

type historyHandler struct {
	*historyWriter
}

func (h historyHandler) Handler(id string) http.HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		h.WriteString(id)
	}
}

func (h historyHandler) WriteHandler(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.WriteString(id)
		w.Write([]byte(id))
	}
}

func (h historyHandler) ParamHandler() ParamHandlerFunc {
	return ParamHandlerFunc(func(w http.ResponseWriter, r *http.Request, value string) {
		h.WriteString(value)
	})
}

func (h historyHandler) HandlerWithError(v string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Context(r).Error(fmt.Errorf(v), 500)
	}
}

func (h historyHandler) SkipRouterHandler(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.WriteString(id)
		Context(r).SkipRouter()
	}
}

func newHistoryHandler() *historyHandler {
	return &historyHandler{&historyWriter{}}
}
