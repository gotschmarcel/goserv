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

type pathComponents struct {
	matcher *regexp.Regexp
	params  []string
}

func (p *pathComponents) match(path string) bool {
	return p.matcher.MatchString(path)
}

func (p *pathComponents) parseParams(path string) Params {
	if len(p.params) == 0 {
		return nil
	}

	matches := p.matcher.FindAllStringSubmatch(path, -1)
	if len(matches) == 0 {
		return nil
	}

	// Iterate group matches only
	params := Params{}
	for index, value := range matches[0][1:] {
		name := p.params[index]
		params[name] = value
	}

	return params
}

func parsePath(path string, strictSlash bool) (*pathComponents, error) {
	var err error
	c := &pathComponents{}

	c.matcher, err = pathStringToRegexp(path, strictSlash)
	if err != nil {
		return nil, err
	}

	c.params = pathStringParameters(path)

	return c, nil
}

func pathStringToRegexp(path string, strictSlash bool) (*regexp.Regexp, error) {
	// TODO: Validate path?
	pattern := paramMatcher.ReplaceAllLiteralString(path, "([^/]+)")
	pattern = wildcardsToRegexp(pattern)

	if strictSlash == false {
		pattern = maybeSlashSuffixMatcher.ReplaceAllString(pattern, "$1/?")
	}

	pattern = wholePattern(pattern)

	return regexp.Compile(pattern)
}

func pathStringParameters(path string) []string {
	var params []string

	for _, match := range paramMatcher.FindAllStringSubmatch(path, -1) {
		// Append param name without leading ':'
		params = append(params, match[1])
	}

	return params
}

func parsePrefixPath(prefix string) (*pathComponents, error) {
	var err error
	c := &pathComponents{}

	pattern := wildcardsToRegexp(prefix)
	pattern = prefixPattern(pattern)

	c.matcher, err = regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func wildcardsToRegexp(path string) string {
	return strings.Replace(path, "*", ".*", -1)
}

func wholePattern(pattern string) string {
	return fmt.Sprintf("^%s$", pattern)
}

func prefixPattern(pattern string) string {
	return fmt.Sprintf("^%s", pattern)
}
