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

type Renderer interface {
	RenderAndWrite(w io.Writer, filePath string, locals interface{}) error
	Ext() string
}

type stdRenderer struct {
	ext            string
	cacheTemplates bool
	tpl            *template.Template
	tplMutex       sync.Mutex
}

func (s *stdRenderer) RenderAndWrite(w io.Writer, filePath string, locals interface{}) error {
	tpl, err := s.template(filePath)
	if err != nil {
		return err
	}

	return tpl.ExecuteTemplate(w, path.Base(filePath), locals)
}

func (s *stdRenderer) Ext() string {
	return s.ext
}

func (s *stdRenderer) template(filePath string) (*template.Template, error) {
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

func NewStdRenderer(ext string, cacheTemplates bool) Renderer {
	return &stdRenderer{ext, cacheTemplates, template.New(""), sync.Mutex{}}
}
