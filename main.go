package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /healthz", readiness)
	mux.HandleFunc("GET /metrics", apiCfg.hitsCount)
	mux.HandleFunc("POST /reset", apiCfg.hitsReset)

	//start server
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Issue starting server: %v\n", err)
	}

}

func readiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		statusCode := 405
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	statusCode := 200
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte("OK"))
}

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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) hitsReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		statusCode := 405
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	statusCode := 200
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Hits reset."))
}
