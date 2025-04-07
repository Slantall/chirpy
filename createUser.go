package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type Email struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	email := Email{}
	err := decoder.Decode(&email)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, "Invalid request", 500)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), email.Email)
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
