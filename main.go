package main

import (
	"log"
	"net/http"
	"time"
)

func main() {

	log.Println("Creating ServeMux to route requests...")
	ServeMux := http.NewServeMux()

	log.Println("Creating server struct...")
	server := &http.Server{
		Addr:           ":8080",
		Handler:        ServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fileHandler := http.FileServer(http.Dir("."))
	ServeMux.Handle("/", fileHandler)

	log.Println("Starting server...")
	log.Fatal(server.ListenAndServe())

}
