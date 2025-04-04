package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/database"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.hitsReset)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	//start server
	err = server.ListenAndServe()
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
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())))
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

func respondWithError(w http.ResponseWriter, ErrString string, statusCode int) {
	type returnError struct {
		ErrString string `json:"error"`
	}

	respBody := returnError{
		ErrString: ErrString,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}

func filterSwears(body string) string {
	forbidden := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	lowerWords := strings.Split(strings.ToLower(body), " ")
	for i := 0; i < len(words); i++ {
		for _, forbid := range forbidden {
			if lowerWords[i] == forbid {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}
