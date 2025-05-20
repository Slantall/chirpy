package main

import (
	"encoding/json"
	"log"
	"main/internal/database"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")
	var err error
	var chirps []database.Chirp
	if author != "" {
		authorID, err := uuid.Parse(author)
		if err != nil {
			log.Printf("Error Parse author: %s", err)
			respondWithError(w, "Failed to retrieve chirps", 500)
		}
		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorID)
	} else {
		chirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			log.Printf("Error retrieving chrips: %s", err)
			respondWithError(w, "Failed to retrieve chirps", 500)
		}
	}
	srt := r.URL.Query().Get("sort")
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
	if srt == "desc" {
		sort.Slice(pulledChirps, func(i, j int) bool { return pulledChirps[i].CreatedAt.After(pulledChirps[j].CreatedAt) })
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
