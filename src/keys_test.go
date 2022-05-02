package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// AddFile to ipfs through API /addFile endpoint
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
	fmt.Println(resp["key"])
}

func TestPostKey(t *testing.T) {
	resp := make(map[string]string)

	// request key from getKey endpoint
	rr := requestKey()

	// store key at path and execute request to post key
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
	deleteKey("temp")
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

	req, err := http.NewRequest("GET", "/getRecord?ipnskey="+ipnskey, nil)
	check(err)

	// send request to API
	rr := executeRequest(GetRecord, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, status)
	}

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	fmt.Println(resp["message"], resp["value"])
}

func TestRecords(t *testing.T) {
	resp := make(map[string]string)
	tests := []struct {
		method     string
		request    string
		fx         http.HandlerFunc
		statusCode int
		payload    []byte
		expected   []byte
	}{
		{
			method:     "POST",
			request:    "/addFile",
			fx:         http.HandlerFunc(AddFile),
			statusCode: 200,
			payload:    []byte("Hello World"),
		},
	}

	for i, tt := range tests {
		// create file to send to server
		pl := tt.payload
		file, err := os.Create(fmt.Sprintf("%s%d", "/tmp/dat", i))
		check(err)
		defer file.Close()
		_, err = file.Write(pl)

		// create writer
		body := &bytes.Buffer{}
		w := multipart.NewWriter(body)
		fw, err := w.CreateFormFile("file", filepath.Base(file.Name()))
		check(err)
		io.Copy(fw, file)
		w.Close()

		// Create request
		req, err := http.NewRequest("POST", tt.request, body)
		check(err)
		req.Header.Add("Content-Type", w.FormDataContentType())

		// execute request
		rr := httptest.NewRecorder()
		handler := tt.fx
		handler.ServeHTTP(rr, req)

		// check response
		if rr.Code != tt.statusCode {
			t.Errorf("Error: %v\n Status Code: %v",
				rr.Body, rr.Code)
		}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		fmt.Println(resp)
	}
}

func TestKeyOperations(t *testing.T) {
	resp := make(map[string]string)
	tests := []struct {
		method     string
		request    string
		fx         http.HandlerFunc
		statusCode int
		expected   []byte
	}{
		{
			method:     "GET",
			fx:         http.HandlerFunc(GetKey),
			request:    "/getKey",
			statusCode: 200,
		},
		{
			method:     "GET",
			fx:         http.HandlerFunc(GetRecord),
			request:    "/getRecord?ipnskey=k51qzi5uqu5diuu5c8uiuzhwk05gycgg5hakjle664382artgxwkw93na4lgmt",
			expected:   []byte("/ipfs/QmUXTtySmd7LD4p6RG6rZW6RuUuPZXTtNMmRQ6DSQo3aMw"),
			statusCode: 200,
		},
		{
			method:     "POST",
			fx:         http.HandlerFunc(StartFollowing),
			request:    "/startFollowing?ipnskey=k51qzi5uqu5diuu5c8uiuzhwk05gycgg5hakjle664382artgxwkw93na4lgmt",
			statusCode: 200,
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.request, nil)
		check(err)
		rr := httptest.NewRecorder()
		handler := tt.fx
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.statusCode {
			t.Errorf("Error: %v\n Status Code: %v",
				rr.Body, rr.Code)
		}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
	}
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
