package main

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/mgmaster24/chirpy/internal/database"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	config := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
		tokenSecret:    secret,
		polkaKey:       polkaKey,
	}

	httpServerMux := http.NewServeMux()
	httpServerMux.Handle(
		"/app/",
		config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))),
	)
	httpServerMux.HandleFunc("GET /api/healthz", healthz)
	httpServerMux.HandleFunc("POST /api/polka/webhooks", config.polka_ep)
	httpServerMux.HandleFunc("POST /api/chirps", config.createChirp)
	httpServerMux.HandleFunc("GET /api/chirps", config.getChirps)
	httpServerMux.HandleFunc("GET /api/chirps/{chirp_id}", config.getChirpById)
	httpServerMux.HandleFunc("DELETE /api/chirps/{chirp_id}", config.deleteChirp)
	httpServerMux.HandleFunc("POST /api/users", config.createUser)
	httpServerMux.HandleFunc("PUT /api/users", config.updateUser)
	httpServerMux.HandleFunc("POST /api/login", config.login)
	httpServerMux.HandleFunc("POST /api/refresh", config.refresh)
	httpServerMux.HandleFunc("POST /api/revoke", config.revoke)
	httpServerMux.HandleFunc("POST /admin/reset", config.reset)
	httpServerMux.HandleFunc("GET /admin/metrics/", config.metrics)

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: httpServerMux,
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
