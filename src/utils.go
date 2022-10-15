package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Fierro-Labs/Fierro/src/models"
	"github.com/gorilla/mux"
)

// This function will create a pin object from incoming request
// will return Pin to use in DB
func getPin(r *http.Request) (models.Pin, error) {
	var pin models.Pin
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

func getPins(reqToken string) []models.PinStatus {
	return users[reqToken].Results
}
