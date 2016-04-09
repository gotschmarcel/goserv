// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	paramMatcher            = regexp.MustCompile(":([^:/]+)")
	allMatcher              = regexp.MustCompile(".*")
	maybeSlashSuffixMatcher = regexp.MustCompile("/?$")
)

const (
	paramValueCapture = "([^/]+)"
	nonStrictSlash    = "$1/?"
)

func pathComponents(pattern string, strict, prefixOnly bool) (*regexp.Regexp, []string) {
	pattern, params := replaceAndExtractParams(pattern)
	pattern = wildcardsToRegexp(pattern)

	if !strict {
		pattern = maybeSlashSuffixMatcher.ReplaceAllString(pattern, nonStrictSlash)
	}

	if prefixOnly {
		pattern = prefixPattern(pattern)
	} else {
		pattern = fullPattern(pattern)
	}

	return regexp.MustCompile(pattern), params
}

func replaceAndExtractParams(pattern string) (string, []string) {
	var params []string

	pattern = paramMatcher.ReplaceAllStringFunc(pattern, func(m string) string {
		name := m[1:] // Remove leading ':'
		params = append(params, name)
		return paramValueCapture
	})

	return pattern, params
}

func wildcardsToRegexp(path string) string {
	return strings.Replace(path, "*", ".*", -1)
}

func fullPattern(pattern string) string {
	return fmt.Sprintf("^%s$", pattern)
}

func prefixPattern(pattern string) string {
	return fmt.Sprintf("^%s", pattern)
}
