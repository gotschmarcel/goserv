# goserv

[![GoDoc](https://godoc.org/github.com/gotschmarcel/goserv?status.svg)](https://godoc.org/github.com/gotschmarcel/goserv)
[![Build Status](https://travis-ci.org/gotschmarcel/goserv.svg?branch=dev)](https://travis-ci.org/gotschmarcel/goserv)

A fast, easy and minimalistic framework for
web applications in Go.

> goserv requires at least Go v1.6.0

## Features

- Robust and fast routing
- Middleware handlers
- Nested routers
- Static file serving
- Template rendering
- Request context to share data between handlers
- URL parameters
- Improved response building
- Support for http.Handler

## Getting Started

Here's a small example showing some of the features supported by **goserv**:

```go
import (
	"io"
	"log"
	"net/http"
	"github.com/gotschmarcel/goserv"
)

func accessLogger(res goserv.ResponseWriter, req *goserv.Request) {
	log.Printf("Access: %s", req.URL.String())
}

func homeHandler(res goserv.ResponseWriter, req *goserv.Request) {
	io.WriteString(res, "Welcome Home")
}

func secureMiddleware(res goserv.ResponseWriter, req *goserv.Request) {
	// Authenticate ...
}

func userHandler(res goserv.ResponseWriter, req *goserv.Request) {
	id := req.Param.Get("user_id")
	// ...
}

func main() {
	server := goserv.NewServer()

	server.UseFunc(accessLogger)
	server.GetFunc("/", homeHandler)

	api := server.SubRouter("/api")
	api.UseFunc(secureMiddleware)
	api.GetFunc("/users/:user_id", userHandler)

	log.Fatalln(server.Listen(":8080"))
}

```

## License

MIT licensed. See the LICENSE for more information.
