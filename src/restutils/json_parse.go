package restutils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Fierro-Labs/Fierro/src/models"
)

const contentType = "application/json"

// WriteJSONError will set the appropiate headers for when there is an error.
func WriteJSONError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", contentType)
}

// WriteJSONSuccessStatus will set the appropiate headers for when there is an error.
func WriteJSONSuccessStatus(w http.ResponseWriter, Pin models.PinStatus) {
	jsonResp, err := json.Marshal(Pin)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", contentType)
	w.Write(jsonResp)
}

func WriteJSONSuccessResults(w http.ResponseWriter, Pins []models.PinStatus) {
	jsonResp, err := json.Marshal(Pins)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", contentType)
	w.Write(jsonResp)
}

// WriteJSONSuccess will set the appropiate headers for when there is an error.
func WriteJSONSuccess(w http.ResponseWriter, Pin models.PinResults) {
	jsonResp, err := json.Marshal(Pin)
	if err != nil {
		log.Fatalf("error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", contentType)
	w.Write(jsonResp)
}
