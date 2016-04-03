// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

type Locals interface{}

type Params map[string]string

func (p Params) Get(key string) string { return p[key] }
