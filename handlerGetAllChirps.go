package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error retrieving chrips: %s", err)
		respondWithError(w, "Failed to retrieve chirps", 500)
	}
	pulledChirps := []Chirps{}
	for _, chirp := range chirps {
		c := Chirps{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			User_id:   chirp.UserID,
		}
		pulledChirps = append(pulledChirps, c)
	}

	dat, err := json.Marshal(pulledChirps)
	if err != nil {
		log.Printf("Error marshalling JSON for pulledChirps: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
