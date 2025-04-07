package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cfg.fileserverHits.Add(1)
		fmt.Printf("Incrementing, hits: %d\n", cfg.fileserverHits.Load())

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitsCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		statusCode := 405
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	statusCode := 200
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())))
}
