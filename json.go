package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	type errStruct struct {
		Error string `json:"error"`
	}

	if err != nil {
		fmt.Println(err)
	}

	if code > 499 {
		fmt.Printf("Responding with 5XX error: %s", msg)
	}

	respondWithJSON(w, code, errStruct{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
