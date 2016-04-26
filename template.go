// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"html/template"
	"io"
	gopath "path"
	"sync"
)

// A TemplateEngine renders templates files using the specified locals and writes
// the result to an io.Writer.
//
// Note that RenderAndWrite must be thread-safe!
type TemplateEngine interface {
	// RenderAndWrite renders the template at filePath using the specified
	// locals and writes the result to w.
	RenderAndWrite(w io.Writer, filePath string, locals interface{}) error

	// Ext returns the file extension (including the leading ".") for the
	// files supported by the TemplateEngine.
	Ext() string
}

type stdTemplateEngine struct {
	ext            string
	cacheTemplates bool
	tpl            *template.Template
	tplMutex       sync.Mutex
}

func (s *stdTemplateEngine) RenderAndWrite(w io.Writer, filePath string, locals interface{}) error {
	tpl, err := s.template(filePath)
	if err != nil {
		return err
	}

	return tpl.ExecuteTemplate(w, gopath.Base(filePath), locals)
}

func (s *stdTemplateEngine) Ext() string {
	return s.ext
}

func (s *stdTemplateEngine) template(filePath string) (*template.Template, error) {
	if !s.cacheTemplates {
		return template.ParseFiles(filePath)
	}

	s.tplMutex.Lock()
	defer s.tplMutex.Unlock()

	tpl := s.tpl.Lookup(gopath.Base(filePath))
	if tpl == nil {
		tpl, err := s.tpl.ParseFiles(filePath)
		return tpl, err
	}

	return tpl, nil
}

// NewStdTemplateEngine returns a new TemplateEngine using Go's html/template package for
// the specified file extension.
//
// Also caching of templates can be enabled. This will cache every template after
// its first use.
func NewStdTemplateEngine(ext string, cacheTemplates bool) TemplateEngine {
	return &stdTemplateEngine{ext, cacheTemplates, template.New(""), sync.Mutex{}}
}
