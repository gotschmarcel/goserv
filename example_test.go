// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv_test

import (
	"fmt"
	"github.com/gotschmarcel/goserv"
	"log"
	"net/http"
	"os"
)

func accessLogger(res goserv.ResponseWriter, req *goserv.Request) {
	log.Printf("Access %s %s", req.Method, req.URL.String())
}

func ExampleSimpleServer() {
	server := goserv.NewRouter()

	server.UseFunc(accessLogger)
	server.GetFunc("/", func(res goserv.ResponseWriter, req *goserv.Request) {
		fmt.Fprint(res, "Home")
	})

	// Everything else is a 404 error.

	log.Fatalln(http.ListenAndServe(":8080", server))
}

type MyController struct {
	AppName string
	Logger  *log.Logger
}

func (m *MyController) logName(res goserv.ResponseWriter, req *goserv.Request) {
	m.Logger.Printf("%s %s", req.Method, req.URL.String())
}

func (m *MyController) getUsers(res goserv.ResponseWriter, req *goserv.Request) {
	m.Logger.Println("Requesting all users")
	fmt.Fprint(res, "Alex, Peter, Marc")
}

func (m *MyController) getUser(res goserv.ResponseWriter, req *goserv.Request) {
	m.Logger.Printf("Requesting user: %s", req.Params.Get("user_id"))
	fmt.Fprint(res, req.Context.Get("user").(string))
}

func (m *MyController) paramUserID(res goserv.ResponseWriter, req *goserv.Request, id string) {
	m.Logger.Printf("Handling user ID: %s", id)
	req.Context.Set("user", fmt.Sprintf("User (id: %s)", id))
}

func ExampleAPISubrouter() {
	controller := &MyController{"MyApp", log.New(os.Stderr, "[main] ", log.LstdFlags)}
	server := goserv.NewRouter()

	server.UseFunc(controller.logName)
	apiRouter := server.Router("/api")

	apiRouter.GetFunc("/users", controller.getUsers)
	apiRouter.GetFunc("/users/:user_id", controller.getUser)
	apiRouter.ParamFunc("user_id", controller.paramUserID)

	log.Fatalln(http.ListenAndServe(":8080", server))
}
