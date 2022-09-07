package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// This function will set the appropiate headers for when there is an error.
func writeJSONError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
}

// This function will set the appropiate headers for when there is an error.
func writeJSONSuccessStatus(w http.ResponseWriter, Pin PinStatus) {

	jsonResp, err := json.Marshal(Pin)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func writeJSONSuccessResults(w http.ResponseWriter, Pins []PinStatus) {

	jsonResp, err := json.Marshal(Pins)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// This function will set the appropiate headers for when there is an error.
func writeJSONSuccess(w http.ResponseWriter, Pin PinResults) {

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

// This function will create a pin object from incoming request
// will return Pin to use in DB
func getPin(r *http.Request) (Pin, error) {
	var pin Pin
	bodyBuffer, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBuffer, &pin); err != nil {
		return pin, err
	}
	return pin, nil
}

// This function grabs the specified parameter value out the URL
// returns true or false if `parameter` exists
func hasParam(r *http.Request, parameter string) (string, bool) {
	category := mux.Vars(r)[parameter]
	if parameter == "requestcid" {
		if len(category) != 62 {
			return "not a base36 cid", false
		}

	}
	if category == "" {
		return "Missing " + parameter + " parameter", false
	}
	return category, true
}

func getAuthToken(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	return reqToken
}

func checkToken(reqToken string) bool {
	if _, ok := users[reqToken]; !ok {
		return false
	}
	return true
}

func searchUser(reqToken string, rid string) int {
	for i := range users[reqToken].Results {
		if users[reqToken].Results[i].Requestid == rid {
			return i
		}
	}
	return -1
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func getPins(reqToken string) []PinStatus {
	return users[reqToken].Results
}
