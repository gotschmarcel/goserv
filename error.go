// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	errNotFound = errors.New(http.StatusText(http.StatusNotFound))
)

func defaultErrorHandler(res ResponseWriter, req *Request, err error) {
	status := http.StatusInternalServerError

	if err == errNotFound {
		status = http.StatusNotFound
	}

	res.WriteHeader(status)
	fmt.Fprintf(res, err.Error())
}
