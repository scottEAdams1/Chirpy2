package main

import (
	"net/http"

	"github.com/scottEAdams1/Chirpy2/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Error with refresh token", err)
		return
	}

	databaseToken, err := cfg.db.GetTokenByTokenString(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Token doesn't exist", err)
		return
	}

	accessToken, err := auth.MakeJWT(databaseToken.UserID, cfg.secret)
	if err != nil {
		respondWithError(w, 400, "Error making access token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "Error with refresh token", err)
		return
	}

	databaseToken, err := cfg.db.GetTokenByTokenString(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Token doesn't exist", err)
		return
	}

	err = cfg.db.UpdateRevokeField(r.Context(), databaseToken.Token)
	if err != nil {
		respondWithError(w, 500, "Error revoking token", err)
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
