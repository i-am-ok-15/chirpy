package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
