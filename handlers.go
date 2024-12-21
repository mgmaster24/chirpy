package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func Healthz(writer http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) AdminMetrics(writer http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "tex/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	hits := cfg.fileServerHits.Load()
	html := `
		<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
		</html>`

	html = fmt.Sprintf(html, hits)
	writer.Write([]byte(html))
}

func (cfg *apiConfig) Reset(writer http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	req.Header.Add("Cache-Control", "no-cache")
	writer.WriteHeader(http.StatusOK)
	cfg.fileServerHits.Store(0)
	writer.Write([]byte("Reset Successful"))
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type cleaned struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chrp := chirp{}
	w.Header().Add("Content-Type", "application/json")
	err := decoder.Decode(&chrp)
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
	respondWithJSON(w, http.StatusOK, cleaned{
		CleanedBody: cleanedBody,
	})
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
