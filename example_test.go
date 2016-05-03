// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv_test

import (
	"fmt"
	"github.com/gotschmarcel/goserv"
	"log"
	"net/http"
)

func ExampleServer_simple() {
	// A simple server example.
	//
	// First an access logging function is registered which gets invoked
	// before the request is forwarded to the home handler. After that
	// the home handler is registered which is the final handler writing
	// a simple message to the response body.
	//
	// As a last step server.Listen is called to start listening for incoming
	// requests.
	server := goserv.NewServer()

	server.Use(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Access %s %s", r.Method, r.URL.String())
	}).Get("/", func(w http.ResponseWriter, r *http.Request) {
		goserv.WriteString(w, "Welcome Home")
	})

	log.Fatalln(server.Listen(":12345"))
}

func ExampleServer_static() {
	// Example file server:
	server := goserv.NewServer()

	server.UseHandler(http.FileServer(http.Dir("/files")))

	log.Fatalln(server.Listen(":12345"))
}

func ExampleServer_context() {
	// Share data between handlers:
	//
	// The middleware stores a shared value in the RequestContext under the name "shared".
	// The GET handler is the next handler in line and retrieves the value from the
	// context. Since a RequestContext can store arbitrary types a type assertion
	// is necessary to get the value in it's real type.
	server := goserv.NewServer()

	server.Use(func(w http.ResponseWriter, r *http.Request) {
		goserv.Context(r).Set("shared", "something to share")
	})

	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		shared := goserv.Context(r).Get("shared").(string)
		goserv.WriteString(w, shared)
	})

	log.Fatalln(server.Listen(":12345"))
}

func ExampleServer_json() {
	// Example server showing how to read and write JSON body:
	//
	// Since WriteJSON and ReadJSONBody are based on the encoding/json
	// package of the standard library the usage is very similar.
	// One thing to notice is that occuring errors are passed to
	// the RequestContext which stops further processing and passes
	// the error to the server's error handler.
	server := goserv.NewServer()

	// Send a simple JSON response.
	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// JSON data to send.
		data := &struct{ Title string }{"My First Todo"}

		// Try to write the data.
		// In case of an error pass it to the RequestContext
		// so it gets forwarded to the next error handler.
		if err := goserv.WriteJSON(w, data); err != nil {
			goserv.Context(r).Error(err, http.StatusInternalServerError)
			return
		}
	})

	// Handle send JSON data.
	server.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var data struct{ Title string }

		// Read and decode the request's body.
		// In case of an error pass it to the RequestContext
		// so it gets forwarded to the next error handler.
		if err := goserv.ReadJSONBody(r, &data); err != nil {
			goserv.Context(r).Error(err, http.StatusBadRequest)
			return
		}

		log.Println(data)
	})

	log.Fatalln(server.Listen(":12345"))
}

func ExampleServer_parameters() {
	// Use URL parameters:
	//
	// URL parameters can be specified by prefixing the name with a ":" in the handler path.
	// The captured value can be retrieved from the RequestContext using the .Param method and
	// the parameter's name.
	//
	// Servers and Routers both support parameter handlers which can be added using the
	// .Param method, i.e. server.Param(...). The first argument is the name of the parameter
	// (without the leading ":"). The parameter handlers are always invoked once before
	// the request handlers get invoked.
	//
	server := goserv.NewServer()

	// This route captures a single URL parameter named "resource_id".
	server.Get("/resource/:resource_id", func(w http.ResponseWriter, r *http.Request) {
		id := goserv.Context(r).Param("resource_id")
		goserv.WriteStringf(w, "Requested resource: %s", id)
	})

	// Registers a parameter handler for the "resource_id" parameter.
	server.Param("resource_id", func(w http.ResponseWriter, r *http.Request, id string) {
		// Some sort of validation.
		if len(id) < 12 {
			goserv.Context(r).Error(fmt.Errorf("Invalid id"), http.StatusBadRequest)
			return
		}

		log.Printf("Requesting resource: %s", id)
	})

	log.Fatalln(server.Listen(":12345"))
}

func ExampleServer_errorHandling() {
	// Custom error handling:
	//
	// Every Router can have its own error handler. In this example
	// a custom error handler is set on the API sub router to handler
	// all errors occured on the /api route.
	//
	// Note: it is also possible to overwrite the default error handler of
	// the server.
	server := goserv.NewServer()

	server.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Errorf("a server error")
		goserv.Context(r).Error(err, http.StatusInternalServerError)
	})

	api := server.SubRouter("/api")
	api.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Errorf("a API error")
		goserv.Context(r).Error(err, http.StatusInternalServerError)
	})

	// Handle errors occured on the API router.
	api.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err *goserv.ContextError) {
		log.Printf("API Error: %s", err)

		w.WriteHeader(err.Code)
		goserv.WriteString(w, err.String())
	}

	log.Fatalln(server.Listen(":12345"))
}
