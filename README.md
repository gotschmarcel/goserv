# goserv

![GoServ](logo/Goserv_Logo_300.png)

A fast, easy and minimalistic framework for
web applications in Go.

> goserv requires at least Go v1.6.0

[![GoDoc](https://godoc.org/github.com/gotschmarcel/goserv?status.svg)](https://godoc.org/github.com/gotschmarcel/goserv)
[![Build Status](https://travis-ci.org/gotschmarcel/goserv.svg?branch=dev)](https://travis-ci.org/gotschmarcel/goserv)

**Read all about it at [goserv.it](http://goserv.it)**

```go
package main

import (
	"github.com/gotschmarcel/goserv"
	"net/http"
	"log"
)

func main() {
	server := goserv.NewServer()
	server.Get("/", func (w http.ResponseWriter, r *http.Request) {
		goserv.WriteString(w, "Welcome Home")
	}
	log.Fatalln(server.Listen(":12345"))
}
```

## Installation

```go
$ go get github.com/gotschmarcel/goserv
```

## Features

- Fully compatible with net/http
- Robust and fast routing
- Middleware handlers
- Nested routers
- Improved file server
- Request context
- URL parameters
- Response and request helpers
- Centralized error handling


## Examples

Examples can be found in `example_test.go`

## License

MIT licensed. See the LICENSE file for more information.
