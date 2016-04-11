// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"html/template"
	"io"
	"path"
	"sync"
)

type TemplateEngine interface {
	RenderAndWrite(w io.Writer, filePath string, locals interface{}) error
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

	return tpl.ExecuteTemplate(w, path.Base(filePath), locals)
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

	tpl := s.tpl.Lookup(path.Base(filePath))
	if tpl == nil {
		tpl, err := s.tpl.ParseFiles(filePath)
		return tpl, err
	}

	return tpl, nil
}

func NewStdTemplateEngine(ext string, cacheTemplates bool) TemplateEngine {
	return &stdTemplateEngine{ext, cacheTemplates, template.New(""), sync.Mutex{}}
}
