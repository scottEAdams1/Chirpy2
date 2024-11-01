package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/scottEAdams1/Chirpy2/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Error deconding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if params.ExpiresInSeconds > 3600 || params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = 3600
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(time.Duration(params.ExpiresInSeconds)*time.Second))
	if err != nil {
		respondWithError(w, 500, "Error making token", err)
		return
	}

	respondWithJSON(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
