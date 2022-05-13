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

func TestForceErrors(t *testing.T) {
	tests := []struct {
		method string
		fx     func() error
	}{
		{
			method: "genKey",
			fx:     func() error { _, err := genKey(""); return err },
		},
		{
			method: "deleteKey",
			fx:     func() error { return deleteKey("") },
		},
		{
			method: "diskDelete",
			fx:     func() error { return diskDelete("/tmp/dne8943phqbtnu4ijher") },
		},
		{
			method: "exportKey",
			fx:     func() error { return exportKey("") },
		},
		{
			method: "importKey",
			fx:     func() error { return importKey("", "") },
		},
	}

	for _, tt := range tests {
		err := tt.fx()
		if err == nil {
			t.Errorf("function %s did not error out.", tt.method)
		}
		fmt.Printf("Failed as Expected: %s\n", err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
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
