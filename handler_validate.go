package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respBody := returnVals{
		CleanedBody: cleanProfanity(params.Body),
	}

	respondWithJSON(w, http.StatusOK, respBody)
}

func cleanProfanity(body string) string {
	profanity := map[string] bool{
		"kerfuffle": true,
		"sharbert": true,
		"fornax": true,
	}

	words := strings.Split(body, " ")
	
	for i, word := range words {
		if profanity[strings.ToLower(word)] {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
