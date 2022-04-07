package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	// "log"
	// "strconv"
	// "encoding/json"
)

func TestGenKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/getKey", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetKey)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: want %v got %v",
			http.StatusOK, status)
	}
	expected := `[{}]`
}
