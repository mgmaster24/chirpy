package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/mgmaster24/chirpy/internal/auth"
	"github.com/mgmaster24/chirpy/internal/database"
)

func healthz(writer http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) polka_ep(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetAPIKeys(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "key not found", err)
		return
	}

	if key != cfg.polkaKey {
		fmt.Println("Key", key)
		fmt.Println("Server Key", cfg.polkaKey)
		respondWithError(w, http.StatusUnauthorized, "Incorrect key provided", err)
		return
	}

	in := input{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&in)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if in.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	_, err = cfg.db.MakeChirpyRed(r.Context(), database.MakeChirpyRedParams{
		ID: in.Data.UserId,
		IsChirpyRed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
