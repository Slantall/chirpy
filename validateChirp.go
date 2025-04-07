package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chrip struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chrp := chrip{}
	err := decoder.Decode(&chrp)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, "Invalid request", 500)
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

}
