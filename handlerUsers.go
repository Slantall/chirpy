package main

import (
	"encoding/json"
	"log"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	//Get token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't retrieve token: %s", err)
		respondWithError(w, "Unauthorized", 401)
		return
	}
	//Get user from token
	userID, err := auth.ValidateJWT(token, cfg.jwtS)
	if err != nil {
		log.Printf("Couldn't validate token: %s", err)
		respondWithError(w, "Unauthorized", 401)
		return
	}
	//get the request

	decoder := json.NewDecoder(r.Body)
	account := LoginAccount{}
	err = decoder.Decode(&account)
	if err != nil {
		log.Printf("Error decoding account parameters: %s", err)
		respondWithError(w, "Invalid request", 500)
		return
	}
	//hash the password they requested
	hashedPass, err := auth.HashPassword(account.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, "Failed to create account", 500)
		return
	}

	updatedUser, err := cfg.db.UpdateEmailAndPass(r.Context(), database.UpdateEmailAndPassParams{Email: account.Email, HashedPassword: hashedPass, ID: userID})
	if err != nil {
		log.Printf("Error updating information: %s", err)
		respondWithError(w, "Failed to update account", 500)
		return
	}

	respBody := User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
