package main

import (
	"database/sql"
	"fmt"
	"main/internal/database"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	plat           string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Issue opening server: %v\n", err)
	}
	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		db:   dbQueries,
		plat: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	//start server
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Issue starting server: %v\n", err)
	}

}
