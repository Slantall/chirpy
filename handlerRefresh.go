package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't retrieve token: %s", err)
		respondWithError(w, "Unauthorized", 401)
		return
	}
	usertkn, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("Couldn't retrieve refresh token: %s", err)
		respondWithError(w, "Unauthorized", 401)
		return
	}
	if usertkn.ExpiresAt.Before(time.Now()) || usertkn.RevokedAt.Valid == true {
		log.Printf("Token expired: %s", err)
		respondWithError(w, "Token expired", 401)
		return
	}
	newAccessTkn, err := auth.MakeJWT(usertkn.UserID, cfg.jwtS, time.Duration(3600)*time.Second)
	if err != nil {
		log.Printf("Error creating token: %s", err)
		respondWithError(w, "Error creating token", 401)
		return
	}

	respBody := struct {
		Token string `json:"token"`
	}{Token: newAccessTkn}
	fmt.Println(respBody)
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
