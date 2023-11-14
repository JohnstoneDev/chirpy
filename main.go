package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/JohnstoneDev/chirpy/internal/helpers"
)


func main(){
	const filePathRoot = "."

	// Create config for caching
	apiCfg := &helpers.ApiConfig{
		FileServerHits: 0,
	}

	// chi router
	router := chi.NewRouter()

	// API router for non-website requests
	apiRouter := chi.NewRouter()

	// admin router
	adminRouter := chi.NewRouter()

	// create a file server
	fs := http.FileServer(http.Dir(filePathRoot))

	// serve assets
	assets := http.FileServer(http.Dir("/assets"))
	fsHandler := http.StripPrefix("/app/", fs)
	appHandler := http.StripPrefix("/app", fs)
	adminHandler := http.StripPrefix("/admin", fs)


	adminRouter.Handle("/metrics", apiCfg.ReportMetrics(adminHandler))
	apiRouter.Get("/healthz", handlerReadiness) 	// Restricted to GET only
	apiRouter.Handle("/reset", apiCfg.MiddlewareResetInfo(fs))

	// Mount the api router to the /api path
	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)

	router.Handle("/app", apiCfg.MiddlewareMetricsInc(appHandler))
	router.Handle("/app/*", apiCfg.MiddlewareMetricsInc(fsHandler))
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