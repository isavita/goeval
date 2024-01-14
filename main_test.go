package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	os.Setenv("API_KEY", "test_api_key")
}

func createJSONRequest(t *testing.T, code string) *http.Request {
	jsonData, err := json.Marshal(SyntaxCheckRequest{Code: code})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/check/gosyntax", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "test_api_key")
	return req
}

func TestValidGoProgramWithApiKey(t *testing.T) {
	validProgram := `package main
                     func main() {
                         println("Hello, world!")
                     }`
	req := createJSONRequest(t, validProgram)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(checkGoSyntaxHandler)
	handler.ServeHTTP(rr, req)

	var resp SyntaxCheckResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if !resp.Valid {
		t.Errorf("Expected valid program, got invalid")
	}
}

func TestInvalidApiKey(t *testing.T) {
	invalidProgram := `package main
                       func main() {
                           println("Hello, world!")
                       }`
	req, err := http.NewRequest("POST", "/check/gosyntax", bytes.NewBufferString(invalidProgram))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "invalid_api_key")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(checkGoSyntaxHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestInvalidGoProgram(t *testing.T) {
	invalidProgram := `package main
                       func main() {
                           println("Hello, world!"
                       }`
	req := createJSONRequest(t, invalidProgram)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(checkGoSyntaxHandler)
	handler.ServeHTTP(rr, req)

	var resp SyntaxCheckResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp.Valid || resp.Error == "" {
		t.Errorf("Expected invalid program, got valid")
	}
}
