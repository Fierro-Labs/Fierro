package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
