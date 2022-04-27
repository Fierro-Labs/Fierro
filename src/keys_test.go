package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/getKey", nil)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetKey)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, status)
	}
}

func TestGetRecord(t *testing.T) {
	resp := make(map[string]string)
	req, err := http.NewRequest("GET", "/getRecord?ipnskey=k51qzi5uqu5diuu5c8uiuzhwk05gycgg5hakjle664382artgxwkw93na4lgmt", nil)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetRecord)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, status)
	}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["value"] == "/ipfs/QmUXTtySmd7LD4p6RG6rZW6RuUuPZXTtNMmRQ6DSQo3aMw" {
		fmt.Println(resp["message"], resp["value"])
	}
}
func TestGet(t *testing.T) {
	// resp := make(map[string]string)
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
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
		rr := httptest.NewRecorder()
		handler := tt.fx
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.statusCode {
			t.Errorf("Error: %v\n Status Code: %v",
				rr.Body, rr.Code)
		}
		// err = json.Unmarshal(rr.Body.Bytes(), &resp)
	}
}
