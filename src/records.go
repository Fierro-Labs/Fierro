package main

import (
	"fmt"
	"net/http"
	// "errors"
)

// This function will accept a ipns key and resolve it
// returns the ipfs path it resolved.
func GetRecord(w http.ResponseWriter, r *http.Request) {
	ipnsKey, ok := HasParam(r, "ipnskey") // grab ipnskey from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting ipnskey", nil)
		return
	}

	ipfsPath, err := resolve(ipnsKey) // download content and return ipfs path
	if err != nil {
		writeJSONError(w, "Error in resolve", err)
		return
	}

	writeJSONSuccess(w, "Success - GetRecord", ipfsPath)
}

// This function will take a ipnskey and add it to the queue
// Returns 200 response
func StartFollowing(w http.ResponseWriter, r *http.Request) {
	key, ok := HasParam(r, "ipnskey") // grab key from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting key: "+key, nil)
		return
	}
	fmt.Print("working?")
	q.PushBack(key)

	writeJSONSuccess(w, "Success - Started following", key)
}

func StopFollowing(w http.ResponseWriter, r *http.Request) {
	key, ok := HasParam(r, "ipnskey") // grab key from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting key: "+key, nil)
		return
	}

	err := stopFollow(key)
	if err != nil {
		writeJSONError(w, "Error in StopFollow", err)
		return
	}
	writeJSONSuccess(w, "Success - Stopped Following", key)
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
