package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/i-am-ok-15/chirpy/internal/database"
)

func profanityChecker(messageBody string) string {
	profanity := [3]string{"kerfuffle", "sharbert", "fornax"}
	lowersWords := strings.Split(messageBody, " ")
	cleanWords := []string{}

	for _, word := range lowersWords {
		lowerWord := strings.ToLower(word)
		newWord := word
		for _, profanityWord := range profanity {
			if lowerWord == profanityWord {
				newWord = "****"
			}
		}
		cleanWords = append(cleanWords, newWord)
	}

	return strings.Join(cleanWords, " ")
}

func lengthChecker(messageBody string) bool {
	if len(messageBody) > 140 {
		log.Printf("Chirp is too long.")
		return false
	} else {
		return true
	}
}

func (a *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding message: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding message")
		return
	}

	if lengthChecker(params.Body) == false {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	params.Body = profanityChecker(params.Body)

	createdChirp, err := a.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusBadRequest, "Error creating user")
		return
	}

	chirp := Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (a *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpArray, err := a.dbQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error fetching chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error fetching chirps")
		return
	}

	APIChirps := []Chirp{}

	for _, chirp := range chirpArray {
		chirpStruct := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		APIChirps = append(APIChirps, chirpStruct)
	}

	respondWithJSON(w, http.StatusOK, APIChirps)
}

func (a *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	requestedID := r.PathValue("chirpID")
	parsedUUID, err := uuid.Parse(requestedID)
	if err != nil {
		log.Printf("Error parsing id to uuid: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error parsing id to uuid")
		return
	}

	singleChirp, err := a.dbQueries.GetChirp(r.Context(), parsedUUID)
	if err != nil {
		log.Printf("Error fetching chirp: %s", err)
		respondWithError(w, http.StatusNotFound, "Error fetching chirp")
		return
	}

	APIChirp := Chirp{
		ID:        singleChirp.ID,
		CreatedAt: singleChirp.CreatedAt,
		UpdatedAt: singleChirp.UpdatedAt,
		Body:      singleChirp.Body,
		UserID:    singleChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, APIChirp)
}
