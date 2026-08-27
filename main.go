package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (a *apiConfig) handlerHitCount(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hitCount := a.fileserverHits.Load()
	hitCountStr := fmt.Sprintf("Hits: %d", hitCount)
	w.Write([]byte(hitCountStr))
}

func (a *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	a.fileserverHits.Store(0)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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

	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCfg.handlerHitCount)
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileHandler))

	log.Println("Starting server...")
	log.Fatal(server.ListenAndServe())
}
