package main

import (
	"encoding/json"
	"log"
	"main/internal/auth"
	"main/internal/database"
	"net/http"

	"github.com/google/uuid"
)

type polkaRequest struct {
	Event string           `json:"event"`
	Data  polkaRequestData `json:"data"`
}
type polkaRequestData struct {
	User_id uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) polkaWebhook(w http.ResponseWriter, r *http.Request) {
	APICheck, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Couldn't retrieve API: %s", err)
		respondWithError(w, "Unauthorized", 401)
		return
	}
	if APICheck != cfg.polkaKey {
		respondWithError(w, "Unauthorized", 401)
		return
	}

	decoder := json.NewDecoder(r.Body)
	req := polkaRequest{}
	err = decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, "Invalid request", 500)
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	_, err = cfg.db.UpdateChirpyRed(r.Context(), database.UpdateChirpyRedParams{IsChirpyRed: true, ID: req.Data.User_id})
	if err != nil {
		respondWithError(w, "User not Found", 404)
		return
	}

	w.WriteHeader(204)
}
