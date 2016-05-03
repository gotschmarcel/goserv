// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"fmt"
	"net/http"
	"reflect"
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

func TestPathComponents(t *testing.T) {
	tests := []struct {
		Path, TestPath string
		Match          bool
		Params         params
		Strict         bool
		Prefix         bool
		Err            error
		Regexp         string
		MatcherType    matcher
	}{
		// POSITIVE TESTS //
		{
			Path:     "/",
			TestPath: "/",
			Match:    true,
		},
		{
			Path:     "/abc//def",
			TestPath: "/abc//def",
			Match:    true,
		},

		// Strict vs non-strict
		{
			Path:     "/abc",
			TestPath: "/abc",
			Match:    true,
			Strict:   false,
		},
		{
			Path:     "/abc",
			TestPath: "/abc/",
			Match:    true,
			Strict:   false,
		},
		{
			Path:     "/abc",
			TestPath: "/abc",
			Match:    true,
			Strict:   true,
		},
		{
			Path:     "/abc",
			TestPath: "/abc/",
			Match:    false,
			Strict:   true,
		},

		// Wildcards
		{
			Path:     "/abc/*/def",
			TestPath: "/abc//def",
			Match:    true,
			Regexp:   "^/abc/(.*)/def/?$",
		},
		{
			Path:     "/abc/*/def",
			TestPath: "/hki//def",
			Match:    false,
		},
		{
			Path:     "/ab*",
			TestPath: "/abcdef/khi",
			Match:    true,
			Regexp:   "^/ab(.*)/?$",
		},
		{
			Path:     "/ab*",
			TestPath: "/khi",
			Match:    false,
		},

		// Groups
		{
			Path:     "/ab?c",
			TestPath: "/abc",
			Match:    true,
		},
		{
			Path:     "/ab?c",
			TestPath: "/ac",
			Match:    true,
		},
		{
			Path:     "/ab?c",
			TestPath: "/akc",
			Match:    false,
		},
		{
			Path:     "/a(bc)?d",
			TestPath: "/ad",
			Match:    true,
		},
		{
			Path:     "/a(bc)?d",
			TestPath: "/abcd",
			Match:    true,
		},
		{
			Path:     "/a(bc)?d",
			TestPath: "/abd",
			Match:    false,
		},
		{
			Path:     "/abc/(def)?",
			TestPath: "/abc",
			Match:    true,
			Strict:   true,
		},
		{
			Path:     "/abc/(def)?",
			TestPath: "/abc/def",
			Match:    true,
			Strict:   true,
		},
		{
			Path:     "/abc/(def)?/ghi",
			TestPath: "/abc/ghi",
			Match:    true,
		},
		{
			Path:     "/abc/(def)?/ghi",
			TestPath: "/abc/def/ghi",
			Match:    true,
		},
		{
			Path:     "/abc/(def)?/ghi",
			TestPath: "/abc/jkl/ghi",
			Match:    false,
		},
		{
			Path:     "/abc/(def",
			TestPath: "/abc/(def",
			Match:    true,
		},
		{
			Path:     "/abc/(def",
			TestPath: "/abc",
			Match:    false,
		},

		// Params
		{
			Path:     "/:id",
			TestPath: "/tab",
			Params:   params{"id": "tab"},
			Match:    true,
			Regexp:   "^/(?P<id>[^/]+)/?$",
		},
		{
			Path:     "/:id1/abc/:id2",
			TestPath: "/tab/abc/akad",
			Params:   params{"id1": "tab", "id2": "akad"},
			Match:    true,
			Regexp:   "^/(?P<id1>[^/]+)/abc/(?P<id2>[^/]+)/?$",
		},
		{
			Path:     "/:id1(\\d+)",
			TestPath: "/12345",
			Params:   params{"id1": "12345"},
			Match:    true,
			Regexp:   "^/(?P<id1>\\d+)/?$",
		},
		{
			Path:     "/:id1(\\d+)",
			TestPath: "/abc",
			Match:    false,
		},

		// Prefix
		{
			Path:     "/abc",
			TestPath: "/abcdef",
			Prefix:   true,
			Match:    true,
		},
		{
			Path:     "/abc",
			TestPath: "/def",
			Prefix:   false,
			Match:    false,
		},

		// Matcher Types
		{
			Path:        "/abc",
			MatcherType: &stringMatcher{},
		},
		{
			Path:        "/abc",
			Prefix:      true,
			MatcherType: &stringPrefixMatcher{},
		},
		{
			Path:        "/*",
			Match:       true,
			MatcherType: &allMatcher{},
		},
		{
			Path:        "/abc*",
			MatcherType: &regexpMatcher{},
		},

		// NEGATIVE TESTS //
		{
			Path: "",
			Err:  fmt.Errorf("Paths must not be empty"),
		},
		{
			Path: "abc",
			Err:  fmt.Errorf("Error at index 0, paths must start with '/'"),
		},
		{
			Path: "/abc(:)",
			Err:  fmt.Errorf("Error at index 6, invalid rune ')'"),
		},
	}

	for _, test := range tests {
		path, err := parsePath(test.Path, test.Strict, test.Prefix)

		// Test Parser Error
		if test.Err != nil {
			if err == nil {
				t.Error("Expected parser error")
				continue
			}

			if m1, m2 := test.Err.Error(), err.Error(); m1 != m2 {
				t.Errorf("Wrong error message, expected: %s, actual: %s", m1, m2)
			}

			continue

		} else if err != nil {
			t.Errorf("Unexpected parser error: %s", err)
			continue
		}

		// Test Matcher Type
		if testType, matcherType := reflect.TypeOf(test.MatcherType), reflect.TypeOf(path.matcher); test.MatcherType != nil && testType != matcherType {
			t.Errorf("Wrong matcher type, expected: %s, actual: %s", testType, matcherType)
		}

		// Test Regexp
		if len(test.Regexp) > 0 {
			rxMatcher, ok := path.matcher.(*regexpMatcher)
			if !ok {
				t.Error("Expected matcher to be of type *regexpMatcher")
				continue
			}

			if rxString := rxMatcher.rx.String(); rxString != test.Regexp {
				t.Errorf("Regexp did not match, expected: %s, actual: %s", test.Regexp, rxString)
			}
		}

		// Test Match
		if res := path.Match(test.TestPath); res != test.Match {
			t.Errorf("Path match error: %s == %s, expected: %t, actual: %t", test.Path, test.TestPath, test.Match, res)
			continue
		}

		// Test Params
		if len(test.Params) > 0 {
			req, _ := http.NewRequest(http.MethodGet, test.TestPath, nil)

			if !path.ContainsParams() {
				t.Error("Expected path to have params")
				continue
			}

			params := make(params)
			path.FillParams(req.URL.Path, params)

			for name, testValue := range test.Params {
				value, ok := params[name]

				if !ok {
					t.Errorf("Missing parameter '%s'", name)
				}

				if testValue != value {
					t.Errorf("Wrong value for parameter '%s', expected: %s, actual: %s", name, testValue, value)
				}
			}
		}
	}
}
