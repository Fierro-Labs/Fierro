package main

import (
	"fmt"
	"net/http"
	"testing"
)

// func TestGetRecord(t *testing.T) {
// 	resp := make(map[string]string)
// 	cid := "/ipfs/QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH"

// 	// use shell to generate key in ipfs node
// 	key, err := genKey("temp")
// 	check(err)
// 	fmt.Println("key name: ", key.Name)

// 	// use shell to publish to ipfs node
// 	pubResp, err := publishToIPNS(cid, key.Name)
// 	check(err)
// 	ipnskey := pubResp.Name

// 	// send request to API
// 	req, err := http.NewRequest("GET", "/getRecord?ipnskey="+ipnskey, nil)
// 	check(err)
// 	rr := executeRequest(GetRecord, req)

// 	// check if OK
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("Error: %v\n Status Code: %v",
// 			rr.Body, status)
// 	}

// 	err = json.Unmarshal(rr.Body.Bytes(), &resp)
// 	fmt.Println(resp["message"], resp["value"])
// 	// use shell to delete key from ipfs node
// 	deleteKey("temp")
// }

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
