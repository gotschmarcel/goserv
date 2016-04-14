// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package goserv provides a fast, easy and minimalistic framework for
// web applications in Go.
//
// Getting Started
//
// Every goserv application starts with a single instance of goserv.Server. This
// instance is the entry point for all incoming requests.
//
//      server := goserv.NewServer()
//
// To start handling paths you must register handlers to paths using one of the Server's
// embedded Router functions, like Get.
//
//      server.Get("/", homeHandler)
//
// The first argument in the Get() call is the path for which the
// handler gets registered. To learn more about paths take a look at the Path section.
// As the name of the function suggests the requests are only getting dispatched to the handler
// if the request method is "GET". There are a lot more methods like this one available, just take
// a look at the documentation of the Router.
//
// At this point it is worth noting that goserv provides it's own Handler type using the
// goserv.ResponseWriter and goserv.Request. This means that all http.Handlers cannot
// be used with goserv without wrapping them in a goserv.Handler. To make this more convenient
// a function called WrapHTTPHandler exists.
//
//      server.Get("/", WrapHTTPHandler(httpHomeHandler))
//
// Sometimes it is useful to have handlers that are invoked all the time, also known as middleware.
// To register middleware use the Use() function.
//
//      server.Use(accessLogger)
//
// Path Syntax
//
// All Routes and Routers are registered under a path. This path is matched against the path
// of incoming requests.
//
// A Route created from "/mypath" will only match the request path if it is exactly
// "/mypath".
//
// In case that the Route should match everything starting with "/mypath" a wildcard
// can be appended to the path, i.e.
//	/mypath*
//
// The wildcard can be at any position in the path. It also possible to just use the wildcard
// as path.
//
// Sometimes it is necessary to capture values from parts of the request path, so called parameters.
// To include parameters in a Route the path must contain a named parameter, e.g.
//	/users/:user_id
//
// Parameters always start with a ":" after a "/". The name (without the leading ":") can contain
// alphanumeric symbols as well as "_" and "-". The number of parameters in a Route is not limited, but
// they must be separated by at least a single "/".
//	/:value1/:value2
//
// Routers allow you to register handlers for each parameter in a Route. Each handler gets invoked
// once per parameter and request. That means even though a Router may invoked multiple Routes the
// parameter handlers are invoked only once.
//
//      router.Get("/users/:user_id", handler)
//      router.Param("user_id", userIDHandler)
//
// When a Route processes a Request it automatically extracts the captured parameter values from the path
// and stores the values under their name in the .Param field of the Request.
//
// Strict vs non-strict Slash
//
// A Route can have either strict slash or non-strict slash behavior. In non-strict mode paths with or
// without a trailing slash are considered to be the same, i.e. a Route registered with "/mypath" in
// non-strict mode matches both "/mypath" and "/mypath/". In strict mode both
// paths are considered to be different.
//
// Order matters
//
// The order in which Handlers are registered does matter. An incoming request goes through in the
// exact same order. After each handler the Router checks wether an error was set on the
// ResponseWriter or if a response was written and ends the processing if necessary. In case of an error
// the Router forwards the request along with the error to its ErrorHandler, but only if one is available.
// All sub Routers have no ErrorHandler by default, so all errors are handled by the top level Server. It
// is possible though to handle errors in a sub Router by setting a custom ErrorHandler.
//
// A small example shows how requests are processed in a Router:
//
//      server.Use(A)
//      server.Get("/", B)
//      server.Use(C)
//      server.SubRouter("/sub", D)
//
// All incoming request will follow the order "ABCD". The
package goserv
