// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import "testing"

func stringInSlice(v string, slice []string) bool {
	for _, s := range slice {
		if v == s {
			return true
		}
	}

	return false
}

func TestPathComponents(t *testing.T) {
	tests := []struct {
		path, p, n string
		params     []string
		strict     bool
		prefix     bool
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
		{path: "/:id", p: "/tab", n: "/", params: []string{"id"}},
		{path: "/:i_d", p: "/tab", n: "/", params: []string{"i_d"}},
		{path: "/:i-d/abc", p: "/tab/abc", n: "/tab/adc", params: []string{"i-d"}},
		{path: "/:id1/abc/:id2", p: "/tab/abc/akad", n: "/tab/adc/akad", params: []string{"id1", "id2"}},

		// Prefix
		{path: "/abc", p: "/abcdef", prefix: true},
		{path: "/abc", p: "/abc", n: "/abcdef", prefix: false},
	}

	for _, test := range tests {
		matcher, params := pathComponents(test.path, test.strict, test.prefix)

		if !matcher.MatchString(test.p) {
			t.Errorf("Path did not match: %s != %s", test.p, test.path)
			continue
		}

		if matcher.MatchString(test.n) {
			t.Errorf("Path did match: %s == %s", test.n, test.path)
			continue
		}

		for _, name := range test.params {
			if !stringInSlice(name, params) {
				t.Errorf("Missing param name: %s, %v", name, params)
				continue
			}
		}
	}
}
