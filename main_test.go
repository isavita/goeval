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

func TestPrivacyEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/privacy", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(privacyHandler)
	handler.ServeHTTP(rr, req)

	expectedContentType := "text/html; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	expectedBody := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Privacy Policy</title>
	<style>
		body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; }
		p { font-size: 1.2em; }
	</style>
</head>
<body>
	<p>Your privacy is respected here. No personal data is collected or stored.</p>
</body>
</html>`

	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}
