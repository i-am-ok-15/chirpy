package main

import (
	"log"
	"net/http"
	"time"
)

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	log.Println("Creating ServeMux to route requests...")
	mux := http.NewServeMux()

	log.Println("Creating server struct...")
	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.HandleFunc("/healthz", handlerReadiness)
	mux.Handle("/app/", fileHandler)

	log.Println("Starting server...")
	log.Fatal(server.ListenAndServe())
}
