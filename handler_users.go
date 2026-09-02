package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/i-am-ok-15/chirpy/internal/auth"
	"github.com/i-am-ok-15/chirpy/internal/database"
)

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type newUser struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userDetails := newUser{}
	err := decoder.Decode(&userDetails)
	if err != nil {
		log.Printf("Error decoding user email: %s", err)
		respondWithError(w, http.StatusBadRequest, "Error decoding user email")
		return
	}

	hashedPassword, err := auth.HashPassword(userDetails.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		return
	}

	createdUser, err := a.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          userDetails.Email,
		HashedPassword: hashedPassword,
	})
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
