package main

import (
	"sync/atomic"
	"net/http"
	"github.com/wNRG/chirpy/internal/database"
)

type apiConfig struct {
	dbQueries *database.Queries
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}
