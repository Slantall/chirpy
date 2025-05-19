package main

import (
	"log"
	"main/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	//Parse chirpID
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error parsing Chirp UUID: %s", err)
		respondWithError(w, "Failed to retrieve chirp", 404)
	}
	//try to find chirp
	chirp, err := cfg.db.GetChirp(r.Context(), uuid.UUID(chirpID))
	if err != nil {
		log.Printf("Error retrieving chrip: %s", err)
		respondWithError(w, "Failed to retrieve chirp", 404)
	}
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
	if chirp.UserID != userID {
		respondWithError(w, "You are not the author of this chrip.", 403)
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, "Couldn't delete chirp", 401)
		return
	}
	w.WriteHeader(204)
}
