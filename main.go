package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	fmt.Println("Creating ServeMux to route requests...")
	ServeMux := http.NewServeMux()

	fmt.Println("Creating server struct...")
	server := &http.Server{
		Addr:           ":8080",
		Handler:        ServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Starting server...")
	server.ListenAndServe()

}
