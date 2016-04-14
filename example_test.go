// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv_test

import (
	"github.com/gotschmarcel/goserv"
	"io"
	"log"
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

	server.UseFunc(func(res goserv.ResponseWriter, req *goserv.Request) {
		log.Printf("Access %s %s", req.Method, req.URL.String())
	}).GetFunc("/", func(res goserv.ResponseWriter, req *goserv.Request) {
		io.WriteString(res, "Home")
	})

	log.Fatalln(server.Listen(":12345"))
}

func ExampleServer_subrouter() {
	// Example server with API sub router:
	server := goserv.NewServer()

	apiRouter := server.SubRouter("/api")

	apiRouter.GetFunc("/users", func(res goserv.ResponseWriter, req *goserv.Request) {
		// ...
	})

	apiRouter.GetFunc("/users/:user_id", func(res goserv.ResponseWriter, req *goserv.Request) {
		// ...
	})

	apiRouter.ParamFunc("user_id", func(res goserv.ResponseWriter, req *goserv.Request, val string) {
		// ...
	})

	log.Fatalln(server.Listen(":8080"))
}

func ExampleServer_static() {
	// Example file server:
	server := goserv.NewServer()

	server.Static("/", "/usr/share/doc")
	log.Fatalln(server.Listen(":12345"))
}

func ExampleServer_templates() {
	// Example server rendering template files with Go's html/template package.
	//
	// The template files are placed inside the views folder in the current working
	// directory. Inside of the views folder is a file called home.tpl with the following
	// content:
	//
	//	<html>
	//		<head><title>{{.Title}}</title></head>
	//		<body>Welcome Home</body>
	//	</html>
	server := goserv.NewServer()

	server.TemplateEngine = goserv.NewStdTemplateEngine(".tpl", true /* enable caches */)
	server.TemplateRoot = "views" // Relative folder

	server.GetFunc("/", func(res goserv.ResponseWriter, req *goserv.Request) {
		res.Render("home", &struct{ Title string }{Title: "Home"})
	})

	log.Fatalln(server.Listen(":12345"))
}
