package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Healthz(writer http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) AdminMetrics(writer http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "tex/html; charset=utf-8")
	writer.WriteHeader(200)
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
	writer.WriteHeader(200)
	cfg.fileServerHits.Store(0)
	writer.Write([]byte("Reset Successful"))
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		body string
	}

	type response struct {
		error string
		valid bool
	}

	decoder := json.NewDecoder(r.Body)
	chrp := chirp{}
	resp := response{}

	w.Header().Add("Content-Type", "application/json")
	err := decoder.Decode(&chrp)
	if err != nil {
		log.Printf("Error while decoding request. Error: e%", err)
		w.WriteHeader(500)
		resp.error = "Something went wrong"
		b, e := json.Marshal(resp)
		if e != nil {
			panic(e)
		}

		w.Write(b)
		return
	}

	if len(chrp.body) <= 0 || len(chrp.body) > 140 {
		w.WriteHeader(500)
		resp.error = "Chirp is too long"
		b, e := json.Marshal(resp)
		if e != nil {
			panic(e)
		}

		w.Write(b)
	}

	w.WriteHeader(200)
	resp.valid = true

	b, e := json.Marshal(resp)
	if e != nil {
		panic(e)
	}

	w.Write(b)
}
