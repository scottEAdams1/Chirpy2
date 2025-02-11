package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/scottEAdams1/Chirpy2/internal/auth"
	"github.com/scottEAdams1/Chirpy2/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token string", err)
		return
	}

	id, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Invalid token", err)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	cleaned_body := validate_chirp(w, params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: id,
	})
	if err != nil {
		fmt.Println("Error creating chirp")
		respondWithError(w, 400, "Error creating chirp", err)
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("author_id")
	sort1 := r.URL.Query().Get("sort")
	var chirps []database.Chirp
	var err error
	if id == "" {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, 400, "Error retrieving chirps", err)
			return
		}
	} else {
		uuid, err := uuid.Parse(id)
		if err != nil {
			respondWithError(w, 400, "Unable to parse", err)
			return
		}
		chirps, err = cfg.db.GetChirpByAuthorID(r.Context(), uuid)
		if err != nil {
			fmt.Println("Error retrieving chirps")
			respondWithError(w, 404, "Error retrieving chirps", err)
			return
		}
	}
	var structChirps []Chirp

	for _, chirp := range chirps {
		structChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		structChirps = append(structChirps, structChirp)
	}
	if sort1 == "" {
		respondWithJSON(w, 200, structChirps)
	} else if sort1 == "asc" {
		sort.Slice(structChirps, func(i, j int) bool {
			return structChirps[i].CreatedAt.Before(structChirps[j].CreatedAt)
		})
		respondWithJSON(w, 200, structChirps)
	} else if sort1 == "desc" {
		sort.Slice(structChirps, func(i, j int) bool {
			return structChirps[i].CreatedAt.After(structChirps[j].CreatedAt)
		})
		respondWithJSON(w, 200, structChirps)
	}
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpID")
	if path == "" {
		respondWithError(w, 400, "No ID given", errors.New("no id given"))
		return
	}

	id, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, 400, "Unable to parse", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		fmt.Println("Error creating chirp")
		respondWithError(w, 404, "Error creating chirp", err)
		return
	}

	respondWithJSON(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token string", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Invalid token", err)
		return
	}

	path := r.PathValue("chirpID")
	if path == "" {
		respondWithError(w, 400, "No ID given", errors.New("no id given"))
		return
	}

	chirpID, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, 400, "Unable to parse", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Error getting chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, 403, "Error chirp doesn't belong to user", err)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Error deleting chirp", err)
		return
	}

	respondWithJSON(w, 204, struct{}{})
}

func validate_chirp(w http.ResponseWriter, body string) string {
	maxLength := 140
	if len(body) > maxLength {
		fmt.Println("Error chirp too long")
		respondWithError(w, 400, "Chirp is too long", nil)
		return ""
	}

	words := strings.Split(body, " ")
	profane_words := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range words {
		if slices.Contains(profane_words, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	cleaned_body := strings.Join(words, " ")
	return cleaned_body
}
