package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/JohnstoneDev/chirpy/internal/helpers"
	"github.com/go-chi/chi/v5"
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
	adminFs := http.FileServer(http.Dir("/admin"))

	// serve assets
	assets := http.FileServer(http.Dir("/assets"))
	fsHandler := http.StripPrefix("/app/", fs)
	appHandler := http.StripPrefix("/app", fs)

	adminRootHandler := http.StripPrefix("/admin/", adminFs)
	adminHandler := http.StripPrefix("/admin", adminFs)

	adminRouter.Handle("/metrics", apiCfg.ReportMetrics(adminHandler))
	adminRouter.Handle("/admin", apiCfg.ReportMetrics(adminRootHandler))

	apiRouter.Get("/healthz", handlerReadiness) 	// Restricted to GET only
	apiRouter.Post("/validate_chirp", validateChirpHandler)
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

// validate chirp POST handler
func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type Received struct {
		Body string `json:"body"`
	}

	type ErrorResponse struct {
		Error string
	}

	type validResp struct {
		Valid bool
	}

	resp := ErrorResponse{
		Error: "Something went wrong",
	}

	decoder := json.NewDecoder(r.Body)
	parameters := Received{
		Body: "",
	}

 	err := decoder.Decode(&parameters)
	if err != nil {
		log.Println("Error decoding parameters")

		w.WriteHeader(http.StatusBadRequest)
		data, _  := json.Marshal(resp)
		w.Write(data)

		return
	}

	// parameters is populated successfully
	if len(parameters.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(ErrorResponse{
			Error: "Chirp is too long",
		})

		w.Write(data)
		return

	}

	data, _ := json.Marshal(validResp{
		Valid : true,
	})

	// the request body to return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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