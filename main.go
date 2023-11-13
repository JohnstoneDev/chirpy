package main

import (
	_ "fmt"
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
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	// server that uses the corsMux as the handler
	server :=  http.Server{
		Addr: ":3000",
		Handler: corsMux,
	}

	// listen for request on the server
	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
		return
	}
}