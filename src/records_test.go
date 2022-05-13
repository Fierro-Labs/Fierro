package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestPostRecord(t *testing.T) {
	resp := make(map[string]string)

	// Create and add HelloWorld file to IPFS
	path := "/tmp/dat"
	addFileResponse := addTmpFileToIPFS(path)
	err := json.Unmarshal(addFileResponse.Body.Bytes(), &resp)
	cid := resp["value"][6:] // remove leading `/ipfs/` characters
	fmt.Println(cid)

	// request key from getKey endpoint
	reqkeyResponse := requestKey()

	disp := reqkeyResponse.Header().Get("Content-Disposition")
	line := strings.Split(disp, "=")
	filename := line[1]
	fmt.Println("filename: ", filename)

	//Store bytes in a file
	keyPath := "/tmp/temp.key"
	file := createTmpFile(reqkeyResponse, keyPath)

	// Create writer
	w, body := createWriter(file)

	// Call postRecord endpoint
	req, err := http.NewRequest("POST", "/PostRecord?CID="+cid, body)
	check(err)
	req.Header.Add("Content-Type", w.FormDataContentType())

	// execute request
	rr := executeRequest(PostRecord, req)

	// check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Result().Body, status)
	}

	// parse key name from response
	rb, err := ioutil.ReadAll(rr.Body)
	check(err)
	err = json.Unmarshal(rb, &resp)
	fmt.Println(resp)

	// delete tmp file
	err = os.Remove(path)
	check(err)
	// delete key file
	err = os.Remove(keyPath)
	check(err)
	// use shell to delete key from ipfs node
	deleteKey(resp["keyname"])
}

func TestGetRecord(t *testing.T) {
	resp := make(map[string]string)
	cid := "/ipfs/QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH"

	// use shell to generate key in ipfs node
	key, err := genKey("temp")
	check(err)
	fmt.Println("key name: ", key.Name)

	// use shell to publish to ipfs node
	pubResp, err := publishToIPNS(cid, key.Name)
	check(err)
	ipnskey := pubResp.Name

	// send request to API
	req, err := http.NewRequest("GET", "/getRecord?ipnskey="+ipnskey, nil)
	check(err)
	rr := executeRequest(GetRecord, req)

	// check if OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, status)
	}

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	fmt.Println(resp["message"], resp["value"])
	// use shell to delete key from ipfs node
	deleteKey("temp")
}

func TestStartFollowing(t *testing.T) {
	ipnskey := "k51qzi5uqu5diir8lcwn6n9o4k2dhohp0e16ur82vw55abv5xfi91mp00ie0ml"

	// send request to API
	req, err := http.NewRequest("GET", "/startFollowing?ipnskey="+ipnskey, nil)
	check(err)
	rr := executeRequest(StartFollowing, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, status)
	}
	fmt.Println(rr.Body.String())
}

func TestStopFollowing(t *testing.T) {
	ipnskey := "k51qzi5uqu5diir8lcwn6n9o4k2dhohp0e16ur82vw55abv5xfi91mp00ie0ml"

	// First we must add a key or else we get "Queue is empty" error
	req, err := http.NewRequest("GET", "/startFollowing?ipnskey="+ipnskey, nil)
	check(err)
	executeRequest(StartFollowing, req)

	// send request to API
	req, err = http.NewRequest("GET", "/stopFollowing?ipnskey="+ipnskey, nil)
	check(err)
	rr := executeRequest(StopFollowing, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, status)
	}
	fmt.Println(rr.Body.String())
}
