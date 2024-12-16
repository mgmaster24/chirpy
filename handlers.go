package main

import (
	"fmt"
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
