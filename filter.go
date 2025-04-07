package main

import "strings"

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
