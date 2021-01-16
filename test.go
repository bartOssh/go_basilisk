package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	db "github.com/bartOssh/go_basilisk/services"
)

func TestHealthCheckHandler(t *testing.T) {
	dbURI := os.Getenv("MONGODB_ADDON_URI")
	dbName := os.Getenv("MONGODB_ADDON_DB")
	dbClient, _ = db.NewMongoClient(dbURI, dbName)
	appToken, _ = dbClient.GetToken()

	uri := fmt.Sprintf("/url?token=%s", appToken)

	req, err := http.NewRequest("GET", uri, bytes.NewBuffer([]byte(`{"url" : "http://google.com"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(scanPng)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	fmt.Printf("status: %v\n", status)
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
