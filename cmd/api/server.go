package main

import (
	"fmt"
	"net/http"
)

func (app *application) server() {
	server := &http.Server{
		Addr:    fmt.Sprint(":8080"),
		Handler: app.routes(),
	}

	fmt.Print("Starting server...")
	server.ListenAndServe()
}
