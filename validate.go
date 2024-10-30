package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type validStruct struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Error deconding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	maxLength := 140
	if len(params.Body) > maxLength {
		fmt.Println("Error chirp too long")
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	words := strings.Split(params.Body, " ")
	profane_words := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range words {
		if slices.Contains(profane_words, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	cleaned_body := strings.Join(words, " ")
	respondWithJSON(w, 200, validStruct{
		CleanedBody: cleaned_body,
	})
}
