package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	dat := []byte("OK")
	w.Write(dat)
}
