package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error parsing Chirp UUID: %s", err)
		respondWithError(w, "Failed to retrieve chirp", 404)
	}

	chirp, err := cfg.db.GetChirp(r.Context(), uuid.UUID(chirpID))
	if err != nil {
		log.Printf("Error retrieving chrip: %s", err)
		respondWithError(w, "Failed to retrieve chirp", 404)
	}

	chirpJson := Chirps{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_id:   chirp.UserID,
	}

	dat, err := json.Marshal(chirpJson)
	if err != nil {
		log.Printf("Error marshalling JSON for pulledChirps: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
