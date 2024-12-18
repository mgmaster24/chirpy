package main

import (
	"net/http"
	"sync/atomic"
)

func main() {
	config := apiConfig{
		fileServerHits: atomic.Int32{},
	}

	httpServerMux := http.NewServeMux()
	httpServerMux.Handle(
		"/app/",
		config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))),
	)
	httpServerMux.HandleFunc(
		"/admin/metrics/", config.AdminMetrics)
	httpServerMux.HandleFunc("GET /api/healthz", Healthz)
	httpServerMux.HandleFunc("POST /admin/reset", config.Reset)
	httpServerMux.HandleFunc("POST /api/validate_chirp", ValidateChirp)
	httpServer := http.Server{
		Addr:    ":8080",
		Handler: httpServerMux,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
