package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main(){
	const filePathRoot = "."

	// Create config for caching
	apiCfg := &apiConfig{
		fileserverHits: 0,
	}

	// chi router
	router := chi.NewRouter()

	// create a file server
	fs := http.FileServer(http.Dir(filePathRoot))

	// serve assets
	assets := http.FileServer(http.Dir("/assets"))
	fsHandler := http.StripPrefix("/app/", fs)
	appHandler := http.StripPrefix("/app", fs)

	router.Get("/healthz", handlerReadiness) 	// Restricted to GET only
	router.Handle("/metrics", apiCfg.reportMetrics(fsHandler))
	router.Handle("/reset", apiCfg.middlewareResetInfo(fs))

	router.Handle("/app", apiCfg.middlewareMetricsInc(appHandler))
	router.Handle("/app/*", apiCfg.middlewareMetricsInc(fsHandler))
	router.Handle("/assets/", http.StripPrefix("/assets/", assets))

	indexServer := &http.Server {
		Addr: ":8080",
		Handler: middlewareCors(router) ,
	}

	log.Println("Listening on port 8080")
	if err := indexServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// readiness Handler : checks if the server is ready to receive requests
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits int
}

// middleware to report metrics
func (config *apiConfig) reportMetrics (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics := config.fileserverHits

		value := fmt.Sprintf("Hits: %v \n", metrics)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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

// middleware function that adds CORS headers to the response
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}