package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	log.Println("Go Server Listening on Port 8081")
	server.ListenAndServe()
}
