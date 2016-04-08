// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"net/http"
	"testing"
)

func stringInSlice(v string, slice []string) bool {
	for _, s := range slice {
		if v == s {
			return true
		}
	}

	return false
}

func TestParsePathString(t *testing.T) {
	tests := []struct {
		path, p, n string
		params     Params
		strict     bool
	}{
		// Strict vs non-strict
		{path: "/abc", p: "/abc"},
		{path: "/abc", p: "/abc/"},
		{path: "/abc", p: "/abc", n: "/abc/", strict: true},

		// Wildcards
		{path: "/abc/*/def", p: "/abc//def", n: "/abc//ktz"},
		{path: "/ab*", p: "/abcdef/khi", n: "/def"},
		{path: "/*", p: "/"},

		// Params
		{path: "/:id", p: "/tab", n: "/", params: Params{"id": "tab"}},
		{path: "/:i_d", p: "/tab", n: "/", params: Params{"i_d": "tab"}},
		{path: "/:i-d/abc", p: "/tab/abc", n: "/tab/adc", params: Params{"i-d": "tab"}},
		{path: "/:id1/abc/:id2", p: "/tab/abc/akad", n: "/tab/adc/akad", params: Params{"id1": "tab", "id2": "akad"}},
	}

	for _, test := range tests {
		p := NewFullPath(test.path, test.strict, nil)

		or, _ := http.NewRequest("", test.p, nil)
		r := newRequest(or)

		if !p.Match(test.p) {
			t.Errorf("Path did not match: %s != %s", test.p, test.path)
			continue
		}

		if p.Match(test.n) {
			t.Errorf("Path did match: %s == %s", test.n, test.path)
			continue
		}

		p.FillParams(r)
		for name, value := range test.params {
			if !stringInSlice(name, p.params) {
				t.Errorf("Missing param name: %s, %v", name, p.params)
				continue
			}

			v, ok := r.Params[name]
			if !ok {
				t.Errorf("Param not extracted: %s", name)
				continue
			}

			if v != value {
				t.Errorf("Wrong param value: %s != %s", v, value)
				continue
			}
		}
	}
}
