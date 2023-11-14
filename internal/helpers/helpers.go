package helpers

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

// middleware to report metrics
func (config *apiConfig) reportMetrics (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics := config.fileserverHits

		value := fmt.Sprintf("Hits: %v \n", metrics)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(value))

		next.ServeHTTP(w, r)
	})
}

// middleware to increment fileserver hits
func (config *apiConfig) middlewareMetricsInc (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.fileserverHits++
		hits := config.fileserverHits
		text := fmt.Sprintf("Hits: %v \n", hits)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(text))

		next.ServeHTTP(w, r)
	})
}

// middleware tpo reset fileserver metrics
func (config *apiConfig) middlewareResetInfo (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.fileserverHits = 0

		w.WriteHeader(http.StatusOK)
		next.ServeHTTP(w, r)
	})
}
