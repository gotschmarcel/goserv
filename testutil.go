// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

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
