package main

import (
	"encoding/json"
	"log"
	"main/internal/auth"
	"main/internal/database"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type chrip struct {
		Body    string    `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	chrp := chrip{}
	err := decoder.Decode(&chrp)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, "Invalid request", 500)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't retrieve token: %s", err)
		respondWithError(w, "Unauthorized", 401)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtS)
	if err != nil {
		log.Printf("incorrect token: %s", err)
		respondWithError(w, "Unauthorized", 401)
		return
	}

	if len(chrp.Body) >= 140 {
		respondWithError(w, "Chirp is too long", 400)
		return
	}

	type returnValid struct {
		Cleaned_body string `json:"cleaned_body"`
	}
	filtered := filterSwears(chrp.Body)
	/*
		respBody := returnValid{
			Cleaned_body: filtered,
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
	*/
	chirpParams := database.CreateChirpParams{
		Body:   filtered,
		UserID: userID,
	}
	createChirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating chrip: %s", err)
		respondWithError(w, "Failed to create chirp in database", 500)
	}

	chirpJSON := Chirps{
		ID:        createChirp.ID,
		CreatedAt: createChirp.CreatedAt,
		UpdatedAt: createChirp.UpdatedAt,
		Body:      createChirp.Body,
		User_id:   createChirp.UserID,
	}

	dat, err := json.Marshal(chirpJSON)
	if err != nil {
		log.Printf("Error marshalling JSON for createChirp: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)

}
