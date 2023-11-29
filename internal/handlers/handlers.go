package handlers

import (
	"net/http"
	"log"
	"encoding/json"

	"github.com/JohnstoneDev/chirpy/internal/helpers"
)

// readiness Handler : checks if the server is ready to receive requests
func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// validate chirp POST handler
func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type Received struct {
		Body string `json:"body"`
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	type validResp struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	resp := ErrorResponse{
		Error: "Something went wrong",
	}

	decoder := json.NewDecoder(r.Body)
	parameters := Received{}

 	err := decoder.Decode(&parameters)
	if err != nil {
		log.Println("Error decoding parameters:", err)
		w.WriteHeader(http.StatusBadRequest)

		if data, err  := json.Marshal(resp); err == nil {
			w.Write(data)
		} else {
			log.Println("Error marshaling response:", err)
		}

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

	cleanResponse := helpers.ReplaceProfanity(parameters.Body)

	data, _ := json.Marshal(validResp{
		Cleaned_body: cleanResponse,
	})

	// the request body to return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}