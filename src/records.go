package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Fierro-Labs/Fierro/src/models"
	"github.com/Fierro-Labs/Fierro/src/restutils"
	"github.com/google/uuid"
)

// This function will accept a object ID & search DB for it
// returns the PinStatus object found.
func GetRecord(w http.ResponseWriter, r *http.Request) {
	reqToken := restutils.GetTokenFromContext(r.Context())

	status := getPins(reqToken)
	if len(status) == 0 {
		fmt.Println("User has no pins.")
		restutils.WriteJSONError(w)
		return
	}

	restutils.WriteJSONSuccessResults(w, status)
}

// This function will take a ipns-name and add it to the queue
// Returns 202 response
// Might need to think about adding pinstatus objects into queue instead of just names
func StartFollowing(w http.ResponseWriter, r *http.Request) {
	reqToken := restutils.GetTokenFromContext(r.Context())
	var status models.Status = models.QUEUED

	// create Pin object from request
	pin, ok := getPin(r)
	if ok != nil {
		fmt.Println("problem with request")
		restutils.WriteJSONError(w)
		return
	}
	// Create uuid
	requestid := uuid.NewString()

	// Create pinstatus object to store & return in response
	pinstatus := models.PinStatus{Requestid: requestid, Pin: &pin, Created: time.Now(), Status: &status, Delegates: &MLTRADRS}

	// Update users' saved pins
	pinRes := users[reqToken]
	pinRes.Count++
	pinRes.Results = append(pinRes.Results, pinstatus)
	users[reqToken] = pinRes

	// Add name to queue
	q.PushBack(pin.Cid)
	// Return result
	restutils.WriteJSONSuccessStatus(w, pinstatus)
}

// This function will take an ipns name and delete from queue
// Continue here
func StopFollowing(w http.ResponseWriter, r *http.Request) {
	reqToken := restutils.GetTokenFromContext(r.Context())

	// grab key from query parameter
	rid, ok := hasParam(r, "requestid")
	if !ok {
		restutils.WriteJSONError(w)
		return
	}

	// Find user by their authToken and check if they have a pin with that requestID
	idx := searchUser(reqToken, rid)
	if idx < 0 {
		fmt.Println("request id not found")
		restutils.WriteJSONError(w)
	}

	// Pass string requestID to remove from queue
	err := removeFromQueue(users[reqToken].Results[idx].Pin.Cid)
	if err != nil {
		restutils.WriteJSONError(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
}

// This function will delete a key from the queue
// return success or error
func removeFromQueue(ipnsKey string) error {
	if q.Len() < 1 {
		fmt.Printf("Queue is empty")
		return fmt.Errorf("queue is empty")
	}

	var ipnsKeyInt interface{} = ipnsKey
	index := q.Index(func(item interface{}) bool {
		return item == ipnsKeyInt
	})

	// if q.Index() returns -1 aka value not found
	if index < 0 {
		fmt.Printf("Key %s not in queue", ipnsKey)
		return fmt.Errorf("key %s not in queue", ipnsKey)
	}
	q.Remove(index)
	return nil
}

// This function resolves and rotates keys in queue
// returns the key at front of queue
func follow() string {
	// if q is empty
	if q.Len() < 1 {
		fmt.Printf("Queue is empty")
		return ""
	}

	ipnsKey := fmt.Sprintf("%s", q.Front()) // Used to convert interface to string

	q.Rotate(1) //moves front elem to the back

	ipfsPath, err := resolve(ipnsKey)
	if err != nil {
		fmt.Printf("Error in resolve %s", err)
		return ""
	}
	return ipfsPath
}
