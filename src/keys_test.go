package main

import (
	"net/http"
	"net/http/httptest"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func executeRequest(handle func(w http.ResponseWriter, r *http.Request), req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle)
	handler.ServeHTTP(rr, req)
	return rr
}
