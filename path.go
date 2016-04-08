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

type Path struct {
	matcher *regexp.Regexp
	params  []string
	handler Handler
}

func (p *Path) Match(path string) bool {
	return p.matcher.MatchString(path)
}

func (p *Path) FillParams(r *Request) {
	if len(p.params) == 0 {
		return
	}

	matches := p.matcher.FindAllStringSubmatch(r.SanitizedPath(), -1)
	if len(matches) == 0 {
		return
	}

	// Iterate group matches only
	for index, value := range matches[0][1:] {
		name := p.params[index]
		r.Params[name] = value
	}
}

func (p *Path) ServeHTTP(res ResponseWriter, req *Request) {
	p.handler.ServeHTTP(res, req)
}

func NewPrefixPath(pattern string, handler Handler) *Path {
	return newPath(pattern, true, prefixPattern, handler)
}

func NewFullPath(pattern string, strict bool, handler Handler) *Path {
	return newPath(pattern, strict, fullPattern, handler)
}

func newPath(pattern string, strict bool, wrapper func(string) string, handler Handler) *Path {
	p := &Path{handler: handler}

	pattern = replaceAndExtractParams(pattern, &p.params)
	pattern = wildcardsToRegexp(pattern)

	if strict == false {
		pattern = maybeSlashSuffixMatcher.ReplaceAllString(pattern, nonStrictSlash)
	}

	pattern = wrapper(pattern)

	p.matcher = regexp.MustCompile(pattern)

	return p
}

func replaceAndExtractParams(pattern string, params *[]string) string {
	return paramMatcher.ReplaceAllStringFunc(pattern, func(m string) string {
		name := m[1:] // Remove leading ':'
		*params = append(*params, name)
		return paramValueCapture
	})
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
