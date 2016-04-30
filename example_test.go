// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv_test

import (
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

func ExampleServer_subrouter() {
	// Example server with API sub router:
	server := goserv.NewServer()

	apiRouter := server.SubRouter("/api")

	apiRouter.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		// ...
	})

	apiRouter.Get("/users/:user_id", func(w http.ResponseWriter, r *http.Request) {
		// ...
	})

	apiRouter.Param("user_id", func(w http.ResponseWriter, r *http.Request, val string) {
		// ...
	})

	log.Fatalln(server.Listen(":8080"))
}

func ExampleServer_static() {
	// Example file server:
	server := goserv.NewServer()

	server.SubRouter("/").UseHandler(goserv.FileServer("/home/myfile", "/", "index.html"))
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
