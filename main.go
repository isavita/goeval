package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
)

type SyntaxCheckRequest struct {
	Code string `json:"code"`
}

type SyntaxCheckResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

func checkGoSyntaxHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-Api-Key")
	expectedApiKey := os.Getenv("API_KEY")
	if apiKey != expectedApiKey || apiKey == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SyntaxCheckRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, "", req.Code, parser.AllErrors)

	response := SyntaxCheckResponse{
		Valid: err == nil,
	}
	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func privacyHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent := `<!DOCTYPE html>
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

	fmt.Fprintf(w, htmlContent)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("default to port %s", port)
	}

	http.HandleFunc("/privacy", privacyHandler)
	http.HandleFunc("/check/gosyntax", checkGoSyntaxHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
