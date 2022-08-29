package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// This function will set the appropiate headers for when there is an error.
func writeJSONError(w http.ResponseWriter, msg string, err error) {
	resp := make(map[string]interface{})
	resp["message"] = msg
	resp["error"] = err
	jsonResp, err := json.Marshal(resp)
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
func writeJSONSuccess(w http.ResponseWriter, msg string, val string) {
	resp := make(map[string]string)
	resp["message"] = msg
	resp["value"] = val
	jsonResp, err := json.Marshal(resp)
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
func HasParam(r *http.Request, parameter string) (string, bool) {
	params, ok := r.URL.Query()[parameter]
	// Query()[parameter] will return an array of items,
	// we only want the single item.
	if !ok || len(params[0]) < 1 {
		return "Missing " + parameter + " parameter", ok
	}
	return params[0], ok
}
