package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddFile(t *testing.T) {
	resp := make(map[string]string)

	// create file
	path := "/tmp/dat"
	rr := addTmpFileToIPFS(path)

	// check response
	if rr.Code != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, rr.Code)
	}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	fmt.Println("response value:", resp["value"])

	// delete tmp file
	err = os.Remove(path)
	check(err)
}

func TestGenKey(t *testing.T) {
	resp := make(map[string][]byte)

	// create and execute request
	rr := requestKey()

	// check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, status)
	}

	resp["key"] = rr.Body.Bytes()
	fmt.Println("Success - GetKey")
}

func TestPostKey(t *testing.T) {
	port := "8082"
	ln, err := net.Listen("tcp", ":"+port)
	if err == nil {
		ln.Close()
		t.Skip("skipping test. Service needs to be running on localhost.")
	}

	resp := make(map[string]string)

	// request key from getKey endpoint
	rr := requestKey()

	// store key at path and execute request to post key
	// TODO: fix fileName to be what I get back from requestKey()
	fileName := "temp.key"
	response := submitKey(rr, fileName)

	rb, err := ioutil.ReadAll(response.Body)
	check(err)

	// check response
	if status := response.StatusCode; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			string(rb), status)
	}
	err = json.Unmarshal(rb, &resp)
	fmt.Println("Key Name stored in remote node:", resp["value"])
	response.Body.Close()

	// use shell to delete key from ipfs node
	deleteKey(resp["value"])
}

func TestDeleteKey(t *testing.T) {
	// use shell to generate key in ipfs node
	key, err := genKey("temp")

	// create new request
	req, err := http.NewRequest("DELETE", "/deleteKey?keyName="+key.Name, nil)
	check(err)

	// execute request
	rr := executeRequest(DeleteKey, req)

	// check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, status)
	}
	fmt.Println(rr.Body)
}

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func addTmpFileToIPFS(path string) *httptest.ResponseRecorder {
	pl := []byte("Hello World")
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	_, err = file.Write(pl)
	file.Seek(0, 0) //reset pointer to start of file

	// create writer to send to API
	w, body := createWriter(file)

	// Create request
	req, err := http.NewRequest("POST", "/addFile", body)
	check(err)
	req.Header.Add("Content-Type", w.FormDataContentType())

	// execute request
	return executeRequest(AddFile, req)
}

func requestKey() *httptest.ResponseRecorder {
	// create request
	req, err := http.NewRequest("GET", "/getKey", nil)
	check(err)

	// execute request
	return executeRequest(GetKey, req)
}

func submitKey(rr *httptest.ResponseRecorder, fileName string) *http.Response {
	// create file
	path := "/tmp/" + fileName
	file := createTmpFile(rr, path)

	// create writer
	w, body := createWriter(file)

	// Create request
	req, err := http.NewRequest("POST", "http://localhost:8082/postKey", bytes.NewReader(body.Bytes()))
	check(err)
	req.Header.Add("Content-Type", w.FormDataContentType())
	client := &http.Client{}

	// execute request
	response, err := client.Do(req)
	check(err)

	// delete tmp file
	err = os.Remove(path)
	check(err)

	return response
}

func createWriter(file *os.File) (*multipart.Writer, *bytes.Buffer) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	fw, err := w.CreateFormFile("file", filepath.Base(file.Name()))
	check(err)
	_, err = io.Copy(fw, file)
	check(err)
	w.Close()
	return w, body
}

func createTmpFile(rr *httptest.ResponseRecorder, path string) *os.File {
	file, err := os.Create(path)
	check(err)
	_, err = io.Copy(file, rr.Body)
	check(err)
	file.Seek(0, 0) //reset pointer to start of file

	return file
}

func executeRequest(handle func(w http.ResponseWriter, r *http.Request), req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle)
	handler.ServeHTTP(rr, req)
	return rr
}
