package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, errorMsg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Println("Responding with 5xx error: :s", errorMsg)
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResp{
		Error: errorMsg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}
