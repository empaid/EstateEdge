package main

import (
	"log"
	"net/http"
)

type application struct {
	config cfg
}

type cfg struct {
	addr string
}

func HandleRequest(w http.ResponseWriter, r *http.ResponseWriter) {

}

func (app *application) Run() error {

	server := http.Server{
		Addr: app.config.addr,
	}
	log.Printf("HTTP Server Started on: %s", app.config.addr)
	return server.ListenAndServe()

}
