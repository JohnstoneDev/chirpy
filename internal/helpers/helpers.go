package helpers

import (
	"fmt"
	"net/http"
	"strings"
	_ "text/template"
)

type ApiConfig struct {
	FileServerHits int
}

// middleware to report metrics
func (config *ApiConfig) ReportMetrics (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// html template
		htmlContent := `
			<html>
				<body>
					<h1>Welcome, Chirpy Admin</h1>
					<p>Chirpy has been visited %d times!</p>
				</body>
			</html>
		`

		htmlContent = fmt.Sprintf(htmlContent, config.FileServerHits)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlContent))

		next.ServeHTTP(w, r)
	})
}

// middleware to increment fileserver hits
func (config *ApiConfig) MiddlewareMetricsInc (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.FileServerHits++
		hits := config.FileServerHits
		text := fmt.Sprintf("Hits: %v \n", hits)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
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


func ReplaceProfanity (input string) string {
	profane := []string{"kerfuffle", "sharbert", "fornax"}
	var returnString string

	words := strings.Split(input, " ")

	for _, word := range words {
		for _, prof := range profane {
			if strings.EqualFold(prof, word) {
				word = "****"
			}
		}
		returnString += word + " "
	}

	return strings.TrimSpace(returnString)
}