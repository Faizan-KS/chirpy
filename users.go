package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Faizan-KS/chirpy/internal/auth"
	"github.com/Faizan-KS/chirpy/internal/database"

	"github.com/google/uuid"
)

type chirpyRedUser struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	defer r.Body.Close()
	postIntake := request{}
	if err := json.NewDecoder(r.Body).Decode(&postIntake); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	hashedPwd, err := auth.HashPassword(postIntake.Password)
	if err != nil {
		return
	}
	user, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          postIntake.Email,
		HashedPassword: hashedPwd,
	})
	if err != nil {
		http.Error(w, "Couldn't create user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	type getUserEmail struct {
		Email string `json:"email"`
	}
	getUsers, err := cfg.database.GetAllUsers(r.Context())
	if err != nil {
		return
	}
	for _, emails := range getUsers {
		w.Write([]byte(emails))
	}
}

func (cfg *apiConfig) upgradeUserHook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var redUser chirpyRedUser
	if err := json.NewDecoder(r.Body).Decode(&redUser); err != nil {
		return
	}
	userID, err := uuid.Parse(redUser.Data.UserID)
	if err != nil {
		return
	}
	if redUser.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		cfg.database.UpgradeToRedByID(r.Context(), userID)
		w.WriteHeader(http.StatusNoContent)
		//json.NewEncoder(w).Encode("Payment Success. User upgraded")
	}
}
