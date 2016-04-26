// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type matcher interface {
	Match(path string) bool
}

type allMatcher struct{}

func (a *allMatcher) Match(string) bool {
	return true
}

type stringMatcher struct {
	path   string
	strict bool
}

func (s *stringMatcher) Match(path string) bool {
	if path != "/" && !s.strict && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	return s.path == path
}

type stringPrefixMatcher struct {
	path string
}

func (s *stringPrefixMatcher) Match(path string) bool {
	return strings.HasPrefix(path, s.path)
}

type regexpMatcher struct {
	rx *regexp.Regexp
}

func (r *regexpMatcher) Match(path string) bool {
	return r.rx.MatchString(path)
}

type path struct {
	matcher
	params *regexp.Regexp
}

func (p *path) ContainsParams() bool {
	return p.params != nil
}

func (p *path) Params() []string {
	if !p.ContainsParams() {
		return nil
	}

	subexpNames := p.params.SubexpNames()
	names := make([]string, 0, p.params.NumSubexp())

	// Filter unnamed subexpressions.
	for _, name := range subexpNames {
		if len(name) == 0 {
			continue
		}

		names = append(names, name)
	}

	return names
}

func (p *path) FillParams(req *Request) {
	if !p.ContainsParams() {
		return
	}

	matches := p.params.FindAllStringSubmatch(req.SanitizedPath, -1)
	if len(matches) == 0 {
		return
	}

	// Iterate group matches only
	names := p.params.SubexpNames()[1:]
	for index, value := range matches[0][1:] {
		name := names[index]

		// Skip unnamed groups
		if len(name) == 0 {
			continue
		}

		req.Params[name] = value
	}
}

type runeStream struct {
	data []rune
	idx  int
}

func (r *runeStream) Err(msg string) error {
	return fmt.Errorf("Error at index %d, %s", r.idx-1, msg)
}

func (r *runeStream) Peek() rune {
	if r.AtEnd() {
		panic("cannot peek after end")
	}

	return r.data[r.idx+1]
}

func (r *runeStream) Back() {
	if r.idx == 0 {
		return
	}

	r.idx--
}

func (r *runeStream) Next() rune {
	if r.AtEnd() {
		panic("cannot read after end")
	}

	v := r.data[r.idx]
	r.idx++
	return v
}

func (r *runeStream) AtEnd() bool {
	return r.idx == len(r.data)
}

func (r *runeStream) String() string {
	return string(r.data)
}

type pathParser struct {
	p      *runeStream
	rxBuf  bytes.Buffer
	pBuf   bytes.Buffer
	simple bool // Contains no regexp expressions
}

func (p *pathParser) Parse(pattern string, strict, prefix bool) (*path, error) {
	p.Reset()
	p.p = &runeStream{data: []rune(pattern)}

	// Write start.
	p.rxBuf.WriteByte('^')

	// Check correct start.
	if len(pattern) == 0 {
		return nil, fmt.Errorf("Paths must not be empty")
	}

	slash := p.p.Next()
	if slash != '/' {
		return nil, p.p.Err("paths must start with '/'")
	}

	p.startPart(slash)

	// Iterate over runes.
	for !p.p.AtEnd() {
		r := p.p.Next()

		var err error

		switch r {
		case '/':
			p.startPart(r)
		case '*':
			p.wildcard()
		case '(':
			err = p.group()
		case ')':
			err = p.p.Err("unmatched ')'")
		case '?':
			p.simple = false
			p.flushPart()
			_, err = p.rxBuf.WriteRune(r)
		case ':':
			p.simple = false
			err = p.param()
		default:
			_, err = p.pBuf.WriteRune(r)
		}

		if err != nil {
			return nil, err
		}
	}

	// Flush the last part
	p.flushPart()

	// Check all matcher
	if p.rxBuf.String() == "^/(.*)" {
		return &path{&allMatcher{}, nil}, nil
	}

	if p.simple {
		if prefix {
			return &path{&stringPrefixMatcher{p.p.String()}, nil}, nil
		}

		return &path{&stringMatcher{p.p.String(), strict}, nil}, nil
	}

	if !strict && !prefix {
		if p.rxBuf.Bytes()[p.rxBuf.Len()-1] != '/' {
			p.rxBuf.WriteByte('/')
		}

		p.rxBuf.WriteByte('?')
	}

	if !prefix {
		p.rxBuf.WriteByte('$')
	}

	regexpPattern, err := regexp.Compile(p.rxBuf.String())
	if err != nil {
		return nil, err
	}

	return &path{&regexpMatcher{regexpPattern}, regexpPattern}, nil
}

func (p *pathParser) Reset() {
	p.rxBuf.Reset()
	p.pBuf.Reset()
	p.simple = true
}

func (p *pathParser) flushPart() {
	if p.pBuf.Len() == 0 {
		return
	}

	safePattern := regexp.QuoteMeta(p.pBuf.String())
	p.rxBuf.WriteString(safePattern)
	p.pBuf.Reset()
}

func (p *pathParser) startPart(r rune) {
	p.flushPart()
	p.pBuf.WriteRune(r)
}

func (p *pathParser) wildcard() {
	p.flushPart()
	p.simple = false
	p.rxBuf.WriteString("(.*)")
}

func (p *pathParser) group() error {
	p.flushPart()

	// Walk over runes until end or until the group is complete.
	quote := true
	level := 1
Loop:
	for !p.p.AtEnd() {
		r := p.p.Next()

		switch {
		case r == '(':
			// Increase group level
			level++
		case r == ')':
			// Decrease group level
			level--
		case r == ':' || r == '/':
			p.p.Back()
			break Loop
		case r == '?':
			quote = false
			break Loop
		default:
			// Copy runes until the closing brace is found.
			p.pBuf.WriteRune(r)
		}
	}

	if level != 0 {
		return p.p.Err("unmatched '(' or ')'")
	}

	part := p.pBuf.String()

	if quote {
		part = fmt.Sprintf("(%s)", part)
		part = regexp.QuoteMeta(part)
	} else {
		part = fmt.Sprintf("(%s)?", regexp.QuoteMeta(part))
		p.simple = false
	}

	p.rxBuf.WriteString(part)
	p.pBuf.Reset()

	return nil
}

func (p *pathParser) param() error {
	var name bytes.Buffer
	var pattern bytes.Buffer

Loop:
	for !p.p.AtEnd() {
		r := p.p.Next()

		switch {
		case isAlphaNumDash(r):
			name.WriteRune(r)
		case r == '(':
			if name.Len() == 0 {
				return p.p.Err("missing parameter name")
			}

			var err error
			pattern, err = p.paramPattern()
			if err != nil {
				return err
			}

			break Loop
		case r == ':' || r == '/':
			p.p.Back()
			break Loop
		default:
			return p.p.Err("invalid rune '" + string(r) + "'")
		}
	}

	if pattern.Len() == 0 {
		// Use the default pattern, which captures everything until
		// the next '/'.
		pattern.WriteString("[^/]+")
	}

	part := fmt.Sprintf("(?P<%s>%s)", name.String(), pattern.String())

	p.flushPart()
	p.rxBuf.WriteString(part)

	return nil
}

func (p *pathParser) paramPattern() (bytes.Buffer, error) {
	var pattern bytes.Buffer
	level := 1

Loop:
	for !p.p.AtEnd() {
		r := p.p.Next()

		switch r {
		case '(':
			level++
		case ')':
			level--

			if level > 0 {
				break
			}

			// End found
			break Loop
		default:
			pattern.WriteRune(r)
		}
	}

	return pattern, nil
}

func parsePath(pattern string, strict, prefixOnly bool) (*path, error) {
	parser := &pathParser{}

	path, err := parser.Parse(pattern, strict, prefixOnly)
	if err != nil {
		return nil, err
	}

	return path, nil
}

func isAlphaNum(r rune) bool {
	return unicode.In(r, unicode.Digit, unicode.Letter)
}

func isAlphaNumDash(r rune) bool {
	return isAlphaNum(r) || r == '_' || r == '-'
}
