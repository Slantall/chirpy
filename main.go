package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	//start server
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Issue starting server: %v\n", err)
	}

}
