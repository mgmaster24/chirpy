package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"

	"github.com/mgmaster24/chirpy/internal/auth"
	"github.com/mgmaster24/chirpy/internal/database"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.Authenticate(cfg.tokenSecret, r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not authenticate user", err)
		return
	}

	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chrp := chirp{}
	w.Header().Add("Content-Type", "application/json")
	err = decoder.Decode(&chrp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(chrp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanedBody := cleanBody(chrp.Body, profanities)
	c, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save chirp to DB", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error
	if len(authorId) == 0 {
		chirps, err = cfg.db.GetGhirps(r.Context())
	} else {
		userId, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
			return
		}
		chirps, err = cfg.db.GetChirpsForUserById(r.Context(), userId)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
		return
	}

	sortOrder := r.URL.Query().Get("sort")
	if len(sortOrder) > 0 && sortOrder == "desc" {
		slices.Reverse(chirps)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func cleanBody(body string, profanities map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowered := strings.ToLower(word)
		if _, ok := profanities[lowered]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	uuid, err := getUUIDFromPath(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get uuid from path", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to retrieve chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.Authenticate(cfg.tokenSecret, r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not authenticate user", err)
		return
	}

	uuid, err := getUUIDFromPath(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get uuid from path", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find chirp", err)
		return
	}

	if chirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "User can't delete this chirp", err)
		return
	}

	_, err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		UserID: userId,
		ID:     uuid,
	})
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Unable to delete chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func getUUIDFromPath(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(r.PathValue("chirp_id"))
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
