package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		statusCode := 405
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	if cfg.plat != "dev" {
		statusCode := 403
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode)
		w.Write([]byte("Forbidden"))
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Failed to delete all users: %v", err) // Logs the detailed error for server-side debuggings
		statusCode := 500
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode)
		w.Write([]byte("Failed to delete users"))
		return
	}
	statusCode := 200
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Hits and Users reset."))
}
