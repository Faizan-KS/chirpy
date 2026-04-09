package main

import (
	"encoding/json"
	"http_adv/internal/auth"
	"http_adv/internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type postBody struct {
	Body string `json:"body"`
}
type errorResponse struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	// 1. Read Bearer token
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "Unauthorized"})
		return
	}

	// 2. Validate JWT and get user ID
	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "Unauthorized"})
		return
	}

	// 3. Decode request body
	var pBody postBody
	if err := json.NewDecoder(r.Body).Decode(&pBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Invalid JSON"})
		return
	}

	// 4. Censor words
	words := strings.Fields(pBody.Body)
	for i, word := range words {
		if strings.EqualFold(word, "kerfuffle") ||
			strings.EqualFold(word, "sharbert") ||
			strings.EqualFold(word, "fornax") {
			words[i] = "****"
		}
	}

	// 5. Store chirp with authenticated user ID
	chirp, err := cfg.database.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      strings.Join(words, " "),
		UserID:    userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "Unauthorized"})
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "Unauthorized"})
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirpIDStr := parts[3]
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_,err = cfg.database.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: chirpID,
		UserID: userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerFollowChirp(w http.ResponseWriter, r *http.Request){
	
}
