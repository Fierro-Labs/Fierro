package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
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
