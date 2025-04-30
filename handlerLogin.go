package main

import (
	"encoding/json"
	"log"
	"main/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	account := LoginAccount{}
	err := decoder.Decode(&account)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, "Invalid request", 500)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), account.Email)
	if err != nil {
		respondWithError(w, "Incorrect email or password", 401)
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, account.Password)
	if err != nil {
		respondWithError(w, "Incorrect email or password", 401)
		return
	}
	expires := time.Duration(account.Expires_in_seconds) * time.Second
	if account.Expires_in_seconds <= 0 || account.Expires_in_seconds > 3600 {
		expires = time.Duration(3600) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtS, expires)
	if err != nil {
		log.Printf("Error creating token: %s", err)
		respondWithError(w, "Error creating token", 401)
		return
	}

	respBody := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
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
