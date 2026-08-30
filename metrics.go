package main

import (
	"fmt"
	"net/http"
)

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
