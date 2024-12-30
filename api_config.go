package main

import (
	"sync/atomic"

	"github.com/mgmaster24/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	tokenSecret    string
	polkaKey       string
}
