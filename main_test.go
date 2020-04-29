package main

import (
	"testing"
	"net/http"
)

func TestMain(t *testing.T) {
	go main()
	resp, err := http.Get("http://localhost:8080/") 
	if err != nil { 
        t.Fatalf("Request failed") 
    }
    if resp.StatusCode != 200 {
    	t.Fatalf("Bad response code: %d", resp.StatusCode) 
    } 
}