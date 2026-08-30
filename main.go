package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	log.Println("Creating ServeMux to route requests...")
	mux := http.NewServeMux()

	apiCfg := apiConfig{}

	log.Println("Creating server struct...")
	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHitCount)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileHandler))

	log.Println("Starting server...")
	log.Fatal(server.ListenAndServe())
}
