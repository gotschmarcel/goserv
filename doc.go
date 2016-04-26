// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package goserv provides a fast, easy and minimalistic framework for
// web applications in Go.
//
//      goserv requires at least Go v1.6
//
// Getting Started
//
// The first thing to do is to import the goserv package:
//
//      import "github.com/gotschmarcel/goserv"
//
// After that we need a goserv.Server as our entry point for incoming requests.
//
//      server := goserv.NewServer()
//
// To start handling things we must register handlers to paths using one of the Server's
// embedded Router functions, like Get.
//
//      server.Get("/", func (res goserv.ResponseWriter, req *goserv.Request) {
//              res.WriteString("Welcome Home")
//      })
//
// The first argument in the Get() call is the path for which the
// handler gets registered and the second argument is the handler function itself.
// To learn more about the path syntax take a look at the "Path Syntax" section.
//
// As the name of the function suggests the requests are only getting dispatched to the handler
// if the request method is "GET". There are a lot more methods like this one available, just take
// a look at the documentation of the Router or Route.
//
// At this point it is worth noting that goserv provides it's own Handler and HandlerFunc types using the
// goserv.ResponseWriter and goserv.Request. This means that all http.Handlers must be wrapped before
// they are passed to a goserv.Router. To make this more convenient two functions called WrapHTTPHandler
// and WrapHTTPHandlerFunc exist.
//
//      server.Get("/", WrapHTTPHandler(httpHandler))
//
// Sometimes it is useful to have handlers that are invoked all the time, also known as middleware.
// To register middleware use the Use() function.
//
//      server.Use(func (res goserv.ResponseWriter, req *goserv.Request) {
//              log.Printf("Access %s %s", req.Method, req.URL.String())
//      })
//
// After we've registered all handlers there is only one thing left to do, which is to start
// listening for incoming requests. The easiest way to do that, is to use the Listen method
// of goserv.Server which is a convenience wrapper around http.ListenAndServe.
//
//      err := server.Listen(":12345")
//
// Now we have a running HTTP server which automatically dispatches incoming requests to the
// right handler functions. This was of course just a simple example of what can be achieved
// with goserv. To get deeper into it take a look at the examples or read the reference documentation
// below.
//
// Path Syntax
//
// All Routes and Routers are registered under a path. This path is matched against the path
// of incoming requests and decides wether or not a handler will be processed.
// The following examples show the features of the path syntax supported by goserv.
//
//      NOTE: Paths must start with a "/". Also query strings are not part of a path.
//
// This simple route will match request to "/mypath":
//
//      server.Get("/mypath", handler)
//
// To match everything starting with "/mypath" use an asterisk as wildcard:
//
//      server.Get("/mypath*", handler)
//
// The wildcard can be positioned anywhere. Multiple wildcards are also possible. The
// following route matches request to "/mypath" or anything starting with "/my" and ending
// with "path", e.g. "/myfunnypath":
//
//      server.Get("/my*path", handler)
//
// The next route matches requests to "/abc" or "/ac" by using the "?" expression:
//
//      server.Get("/ab?c", handler)
//
// To make multiple characters optional wrap them into parentheses:
//
//      server.Get("/a(bc)?d", handler)
//
// Sometimes it is necessary to capture values from parts of the request path, so called parameters.
// To include parameters in a Route the path must contain a named parameter:
//
//	server.Get("/users/:user_id", handler)
//
// Parameters always start with a ":". The name (without the leading ":") can contain
// alphanumeric symbols as well as "_" and "-". By default a parameter captures everything
// until the next slash. This behavior can be changed by providing a custom matching pattern:
//
//      server.Get("/users/:user_id(\\d+)", handler)
//
// Remember to escape the backslash when using custom patterns.
//
// Strict vs non-strict Slash
//
// A Route can have either strict slash or non-strict slash behavior. In non-strict mode paths with or
// without a trailing slash are considered to be the same, i.e. a Route registered with "/mypath" in
// non-strict mode matches both "/mypath" and "/mypath/". In strict mode both
// paths are considered to be different.
// The behavior can be modified by changing a Router's .StrictSlash property. Sub routers automatically
// inherit the strict slash behavior from their parent.
//
// Order matters
//
// The order in which handlers are registered does matter, since incoming requests go through the
// exact same order. After each handler the Router checks wether an error was set on the
// ResponseWriter or if a response was written and ends the processing if necessary. In case of an error
// the Router forwards the request along with the error to its ErrorHandler, but only if one is available.
// All sub Routers have no ErrorHandler by default, so all errors are handled by the top level Server. It
// is possible though to handle errors in a sub Router by setting a custom ErrorHandler.
//
package goserv
