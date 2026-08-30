package main

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hitCount := a.fileserverHits.Load()

	hitCountStr := fmt.Sprintf(
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`, hitCount)

	dat := []byte(hitCountStr)
	w.Write(dat)
}

func (a *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	a.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
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
	dat := []byte("OK")
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	respBody := errorResponse{
		Error: msg,
	}

	respondWithJSON(w, code, respBody)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder((r.Body))
	message := chirp{}
	err := decoder.Decode(&message)
	if err != nil {
		log.Printf("Error decoding message: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding message")
		return
	}

	if len(message.Body) > 140 {
		log.Printf("Chirp is too long.")
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, validResponse{
		Valid: true,
	})
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
