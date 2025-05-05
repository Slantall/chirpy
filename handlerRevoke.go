package main

import (
	"log"
	"main/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't retrieve token: %s", err)
		respondWithError(w, "Unauthorized", 401)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("Couldn't revoke token: %s", err)
		respondWithError(w, "Unauthorized", 401)
		return
	}
	w.WriteHeader(204)
}
