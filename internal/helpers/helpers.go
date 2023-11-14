package helpers

import (
	"fmt"
	"net/http"
)

type ApiConfig struct {
	FileServerHits int
}

// middleware to report metrics
func (config *ApiConfig) ReportMetrics (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics := config.FileServerHits

		value := fmt.Sprintf("Hits: %v \n", metrics)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(value))

		next.ServeHTTP(w, r)
	})
}

// middleware to increment fileserver hits
func (config *ApiConfig) MiddlewareMetricsInc (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.FileServerHits++
		hits := config.FileServerHits
		text := fmt.Sprintf("Hits: %v \n", hits)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(text))

		next.ServeHTTP(w, r)
	})
}

// middleware tpo reset fileserver metrics
func (config *ApiConfig) MiddlewareResetInfo (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.FileServerHits = 0

		w.WriteHeader(http.StatusOK)
		next.ServeHTTP(w, r)
	})
}
