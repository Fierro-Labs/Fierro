package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// This function will set the appropiate headers for when there is an error.
func writeJSONError(w http.ResponseWriter, Pin PinStatus) {
	jsonResp, err := json.Marshal(Pin)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// This function will set the appropiate headers for when there is an error.
func writeJSONSuccess(w http.ResponseWriter, Pin PinStatus) {

	jsonResp, err := json.Marshal(Pin)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// This function grabs the specified parameter value out the URL
// returns true or false if `parameter` exists
func hasParam(r *http.Request, parameter string) (string, bool) {
	vars := mux.Vars(r)
	category := vars[parameter]
	if category == "" {
		return "Missing " + parameter + " parameter", false
	}
	return category, true
}
