package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) metrics(writer http.ResponseWriter, req *http.Request) {
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

func (cfg *apiConfig) reset(writer http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	req.Header.Add("Cache-Control", "no-cache")
	if cfg.platform != "dev" {
		writer.WriteHeader(http.StatusForbidden)
		return
	}
	writer.WriteHeader(http.StatusOK)
	cfg.fileServerHits.Store(0)
	cfg.db.ResetUsers(req.Context())
	writer.Write([]byte("Reset Successful"))
}
