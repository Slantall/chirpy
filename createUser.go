package main

import (
	"encoding/json"
	"log"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	account := LoginAccount{}
	err := decoder.Decode(&account)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, "Invalid request", 500)
		return
	}

	hashedPass, err := auth.HashPassword(account.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, "Failed to create account", 500)
		return
	}

	userParams := database.CreateUserParams{
		Email:          account.Email,
		HashedPassword: hashedPass,
	}

	user, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		log.Printf("Error Creating User: %s", err)
		respondWithError(w, "Invalid request", 500)
		return
	}

	respBody := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)

}
