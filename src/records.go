package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
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

// This function takes a CID and file.key and publishes brand new IPNS records to IPFS
// IPFS Node handles republishing automatically in the background as long as it is up and running
// Returns ACK & IPNS path
func PostRecord(w http.ResponseWriter, r *http.Request) {
	var dir = abs + "/KeyStore"

	CID, ok := HasParam(r, "CID") // grab CID from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting CID: "+CID, nil)
		return
	}

	FilePath, err := saveFile(r, dir, 32<<10) // grab uploaded .key file
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}

	name := removeExtenstion(path.Base(FilePath))

	err = importKey(name, FilePath) //import key to local node keystore
	if err != nil {
		writeJSONError(w, "Error in importKey", err)
		return
	}

	pubResp, err := publishToIPNS(ipfsURI+CID, name) //publish IPNS record to IPFS
	if err != nil {
		writeJSONError(w, "Error in publishToIPNS", err)
		return
	}

	err = diskDelete(FilePath) // delete key from disk
	if err != nil {
		writeJSONError(w, "Error in diskDelete", err)
		return
	}

	resp := make(map[string]interface{})
	resp["message"] = "Success - PostRecord"
	resp["value"] = ipnsURI + pubResp.Name
	resp["keyname"] = name
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
	// writeJSONSuccess(w, "Success - PostKey", ipnsURI+pubResp.Name)
}

// This function takes an IPNS Key and file.key and resolves IPNS record
// IPFS Node handles republishing automatically in the background as long as it is up and running
// Returns ACK & resolved content
func PutRecord(w http.ResponseWriter, r *http.Request) {
	var dir = abs + "/KeyStore"

	key, ok := HasParam(r, "ipnskey") // grab key from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting key: "+key, nil)
		return
	}

	FilePath, err := saveFile(r, dir, 32<<10) // grab uploaded .key file
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}
	name := removeExtenstion(path.Base(FilePath))

	err = importKey(name, FilePath) //import key to local node keystore
	if err != nil {
		writeJSONError(w, "Error in importKey", err)
		return
	}

	ipfsPath, err := resolve(key) // download content and return ipfs path
	if err != nil {
		writeJSONError(w, "Error in resolve", err)
		return
	}

	err = diskDelete(FilePath) // delete key from disk
	if err != nil {
		writeJSONError(w, "Error in deleteKey", err)
		return
	}

	// add custom return function to include generated keyname
	writeJSONSuccess(w, "Success - PutRecord", ipfsPath)
}

// This function will take a ipnskey and add it to the queue
// Returns 200 response
func StartFollowing(w http.ResponseWriter, r *http.Request) {
	key, ok := HasParam(r, "ipnskey") // grab key from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting key: "+key, nil)
		return
	}

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
