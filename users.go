package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/mgmaster24/chirpy/internal/auth"
	"github.com/mgmaster24/chirpy/internal/database"
)

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type userInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (config *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	in := userInput{}
	w.Header().Add("Content-Type", "application/json")
	err := decoder.Decode(&in)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPW, err := auth.HsshPassword(in.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash the password", err)
		return
	}

	user, err := config.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          in.Email,
		HashedPassword: hashedPW,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, UserResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})
}

func (config *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	in := userInput{}
	w.Header().Add("Content-Type", "application/json")
	err := decoder.Decode(&in)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := config.db.GetUserByEmail(r.Context(), in.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to get user by email", err)
		return
	}

	err = auth.CheckPasswordHash(in.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Inccorrect Password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, config.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create token", err)
		return
	}

	rt, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token", err)
		return
	}

	_, err = config.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  rt,
		UserID: user.ID,
		ExpiresAt: sql.NullTime{
			Time:  time.Now().AddDate(0, 2, 0),
			Valid: true,
		},
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to create refresh token in db",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		IsChirpyRed:  user.IsChirpyRed.Bool,
		Token:        token,
		RefreshToken: rt,
	})
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.Authenticate(cfg.tokenSecret, r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could authenticate user", err)
		return
	}

	in := userInput{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&in)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPW, err := auth.HsshPassword(in.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash the password", err)
		return
	}

	user, err := cfg.db.UpdateUserEmailPass(r.Context(), database.UpdateUserEmailPassParams{
		ID:             userId,
		Email:          in.Email,
		HashedPassword: hashedPW,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
	}

	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})
}
