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
	fmt.Fprintf(w, "<html><body><p>We do not store any personal data or information from our users.</p></body></html>")
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
