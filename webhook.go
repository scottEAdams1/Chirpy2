package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/scottEAdams1/Chirpy2/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error getting apiKey", err)
		return
	}

	if apiKey != cfg.polka_key {
		respondWithError(w, 401, "Error wrong apiKey", err)
		return
	}

	type userid struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  userid `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, struct{}{})
		return
	}

	_, err = cfg.db.UpdateUserToRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "Error upgrading user", err)
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
