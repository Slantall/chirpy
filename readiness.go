package main

import "net/http"

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
