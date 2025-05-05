package main

import (
	"encoding/json"
	"log"
	"main/internal/auth"
	"main/internal/database"
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
	//expires := time.Duration(account.Expires_in_seconds) * time.Second  //used if `expires_in_seconds` will be used
	//if account.Expires_in_seconds <= 0 || account.Expires_in_seconds > 3600 {
	expires := time.Duration(3600) * time.Second
	//}

	token, err := auth.MakeJWT(user.ID, cfg.jwtS, expires)
	if err != nil {
		log.Printf("Error creating token: %s", err)
		respondWithError(w, "Error creating token", 401)
		return
	}
	refshtkn, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		respondWithError(w, "Error creating refresh token", 401)
		return
	}
	err = cfg.db.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{Token: refshtkn, UserID: user.ID})
	if err != nil {
		log.Printf("Error storing refresh token: %s", err)
		respondWithError(w, "Error creating refresh token", 401)
		return
	}

	respBody := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refshtkn,
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
