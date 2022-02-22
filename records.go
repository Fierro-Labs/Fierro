package main

import (
	"fmt"
	"net/http"
	"strings"
)
// This function will accept a ipns key and resolve it
// returns the ipfs path it resolved.
func GetRecord(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting IPNS Key...")
	ipnsKey, ok := GetParam(r, "ipnskey") // grab ipnskey from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting ipnskey", nil)
		return
	}

	fmt.Println("Resolving IPNS Record...")
	ipfsPath, err := resolve(ipnsKey) // download content and return ipfs path
	if err != nil {
		writeJSONError(w, "Error in resolve", err)
		return
	}

	q.PushBack(ipnsKey)

	writeJSONSuccess(w, "Success - GetRecord", ipfsPath)
}

// This function takes a CID and file.key and publishes brand new IPNS records to IPFS
// IPFS Node handles republishing automatically in the background as long as it is up and running
// Returns ACK & IPNS path
func PostRecord(w http.ResponseWriter, r *http.Request) {
	const dir = "KeyStore"
	
	fmt.Println("Getting CID...")
	CID, ok := GetParam(r, "CID") // grab CID from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting CID: "+CID, nil)
		return
	}

	FileName, err := saveFile(r, dir, 32 << 10) // grab uploaded .key file
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}
	name := strings.Split(FileName, ".")[0]

	fmt.Println("Importing Key...")
	err = importKey(name, dir+"/"+FileName) //import key to local node keystore
	if err != nil {
		writeJSONError(w, "Error in importKey", err)
		return
	}

	fmt.Println("Publishing to IPNS...")
	pubResp, err := publishToIPNS(ipfsURI + CID, name) //publish IPNS record to IPFS
	if err != nil {
		writeJSONError(w, "Error in publishToIPNS", err)
		return
	}

	fmt.Println("Deleting exported key...")
	err = diskDelete(dir+"/"+FileName) // delete key from disk
	if err != nil {
		writeJSONError(w, "Error in diskDelete", err)
		return
	}

	fmt.Printf("\nresponse Name: %s\nresponse Value: %s\n", pubResp.Name, pubResp.Value)
	writeJSONSuccess(w, "Success - PostKey", ipnsURI+pubResp.Name)
}

// This function takes an IPNS Key and file.key and resolves IPNS record
// IPFS Node handles republishing automatically in the background as long as it is up and running
// Returns ACK & resolved content
func PutRecord(w http.ResponseWriter, r *http.Request) {
	const dir = "KeyStore"
	
	fmt.Println("Getting IPNS Key...")
	key, ok := GetParam(r, "ipnskey") // grab key from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting key: "+key, nil)
		return
	}

	FileName, err := saveFile(r, dir, 32 << 10) // grab uploaded .key file
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}
	name := strings.Split(FileName, ".")[0]

	fmt.Println("Importing Key...")
	err = importKey(name, dir+"/"+FileName) //import key to local node keystore
	if err != nil {
		writeJSONError(w, "Error in importKey", err)
		return
	}

	fmt.Println("Resolving IPNS Record...")
	ipfsPath, err := resolve(key) // download content and return ipfs path
	if err != nil {
		writeJSONError(w, "Error in resolve", err)
		return
	}

	fmt.Println("Deleting saved key from disk...")
	err = diskDelete(dir+"/"+FileName) // delete key from disk
	if err != nil {
		writeJSONError(w, "Error in deleteKey", err)
		return
	}
	writeJSONSuccess(w, "Success - PutRecord", ipfsPath)
}

// This function will return the first element of the queue and add it to the back
func FollowRecord(w http.ResponseWriter, r *http.Request) {
	ipnsKey := follow()
	if ipnsKey == "" {
		writeJSONError(w, "Queue is empty", nil)
		return
	}
	writeJSONSuccess(w, "Success - PopFront", ipnsKey)
}

// Helper function so that we can call internally and externally
// returns the key at front of queue
func follow() string {
	defer func() { // <- executes only if there is a panic when using the queue
		if (recover() != nil){
			fmt.Sprintf("Error with queue operations %s")
			return
		}
	}()

	ipnsKeyInt := q.PopFront()
	ipnsKey := fmt.Sprintf("%s", ipnsKeyInt) // Used to convert interface to string
	// fmt.Printf("Success - top key is: %s\n", ipnsKey)
	q.PushBack(ipnsKey)
	ipfsPath, err := resolve(ipnsKey)
	if err != nil {
		fmt.Printf("Error in resolve %s", err)
		return ""
	}
	return ipfsPath
}
