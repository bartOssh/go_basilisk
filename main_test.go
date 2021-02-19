package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	db "github.com/bartOssh/go_basilisk/services"
	"go.uber.org/goleak"
	// "go.uber.org/goleak"
)

// TestHealthCheckHandler for screenshot/jpeg route
func TestHealthCheckHandler(t *testing.T) {
	dbURI := os.Getenv("MONGODB_ADDON_URI")
	dbName := os.Getenv("MONGODB_ADDON_DB")
	dbClient, _ = db.NewMongoClient(dbURI, dbName)
	appToken, _ = dbClient.GetToken()

	uri := fmt.Sprintf("/screenshot/jpeg?token=%s", appToken)

	req, err := http.NewRequest("GET", uri, bytes.NewBuffer([]byte(`{"url" : "http://google.com"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(screenshotJpeg)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	fmt.Printf("status: %v\n", status)
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Test1"},
		{name: "Test2"},
		{name: "Test3"},
		{name: "Test4"},
		{name: "Test5"},
	}
	for _, tt := range tests {
		fmt.Printf(" test: %s \n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
