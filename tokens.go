package main

import (
	"net/http"
	"time"

	"github.com/mgmaster24/chirpy/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get token", err)
		return
	}

	rt, err := cfg.db.GetToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get token from db", err)
		return
	}

	if time.Now().After(rt.ExpiresAt.Time) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is expired", nil)
		return
	}

	if rt.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token has been revoked", nil)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), rt.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user from token", err)
		return
	}

	token, err = auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create token", err)
		return
	}

	type ret struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, ret{
		Token: token,
	})
}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get token", err)
		return
	}

	rt, err := cfg.db.GetToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get token from db", err)
		return
	}

	if time.Now().After(rt.ExpiresAt.Time) {
		respondWithError(w, http.StatusUnauthorized, "Token is expired", nil)
		return
	}

	if rt.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token has been revoked", nil)
		return
	}

	_, err = cfg.db.RevokeToken(r.Context(), rt.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token was NOT revoked", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
