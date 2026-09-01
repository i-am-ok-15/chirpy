package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		Valid bool `json:"valid"`
	}

	type cleanedResponse struct {
		CleanedString string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder((r.Body))
	message := chirp{}
	err := decoder.Decode(&message)
	if err != nil {
		log.Printf("Error decoding message: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding message")
		return
	}

	if lengthChecker(message.Body) == false {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	CleanedBody := profanityChecker(message.Body)

	respondWithJSON(w, http.StatusOK, cleanedResponse{
		CleanedString: CleanedBody,
	})
}
