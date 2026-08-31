package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type newUser struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userDetails := newUser{}
	err := decoder.Decode(&userDetails)
	if err != nil {
		log.Printf("Error decoding user email: %s", err)
		respondWithError(w, http.StatusBadRequest, "Error decoding user email")
		return
	}

	createdUser, err := a.dbQueries.CreateUser(r.Context(), userDetails.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusBadRequest, "Error creating user")
		return
	}

	user := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)
}
