package main

import (
	"fmt"
	"net/http"
	"time"

	guuid "github.com/google/uuid"
)

// This function will accept a ipns key and resolve it
// returns the ipfs path it resolved.
// func GetRecord(w http.ResponseWriter, r *http.Request) {
// 	ipnsKey, ok := hasParam(r, "requestcid") // grab ipnskey from query parameter
// 	if ok != true {
// 		writeJSONError(w, PinStatus{})
// 		return
// 	}
//
// ipfsPath, err := resolve(ipnsKey) // download content and return ipfs path
// if err != nil {
// 	writeJSONError(w, PinStatus{})
// 	return
// }
//
// 	writeJSONSuccess(w, PinStatus{})
// }

// This function will take a ipnskey and add it to the queue
// Returns 200 response
func StartFollowing(w http.ResponseWriter, r *http.Request) {
	var pin *Pin
	var status Status

	key, ok := hasParam(r, "requestcid") // grab key from query parameter
	if ok != true {
		writeJSONError(w, PinStatus{})
		return
	}

	requestid := guuid.NewString()

	pin = new(Pin)
	pin.Cid = key

	status = QUEUED

	q.PushBack(key)

	writeJSONSuccess(w, PinStatus{Requestid: requestid, Pin: pin, Created: time.Now(), Status: &status, Delegates: &MLTRADRS})
}

func StopFollowing(w http.ResponseWriter, r *http.Request) {
	key, ok := hasParam(r, "requestcid") // grab key from query parameter
	if ok != true {
		writeJSONError(w, PinStatus{})
		return
	}

	err := stopFollow(key)
	if err != nil {
		writeJSONError(w, PinStatus{})
		return
	}
	writeJSONSuccess(w, PinStatus{})
}

// This function will delete a key from the queue
// return success or error
func stopFollow(ipnsKey string) error {
	if q.Len() < 1 {
		fmt.Printf("Queue is empty")
		return fmt.Errorf("Queue is empty")
	}

	var ipnsKeyInt interface{} = ipnsKey
	index := q.Index(func(item interface{}) bool {
		return item == ipnsKeyInt
	})

	// if q.Index() returns -1 aka value not found
	if index < 0 {
		fmt.Printf("Key %s not in queue", ipnsKey)
		return fmt.Errorf("Key %s not in queue", ipnsKey)
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
