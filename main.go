package main

import (
	"log"
	"net/http"
)

// middleware function that adds CORS headers to the response
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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


func main(){

	// create a file server
	fs := http.FileServer(http.Dir("."))
	mux := http.NewServeMux()

	// serve assets
	assets := http.FileServer(http.Dir("/assets"))

	mux.Handle("/", fs)
	mux.Handle("/assets", assets)

	indexServer := http.Server {
		Addr: ":8080",
		Handler: middlewareCors(mux) ,
	}

	log.Println("Listening on port 8080")
	if err := indexServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}