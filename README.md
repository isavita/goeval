# Go Syntax Checking API

This API allows for checking the syntax correctness of Go code snippets. It's designed for simplicity and ease of use with basic authentication for secure access.

## Usage

### Checking Syntax Validity

To check if a Go program has valid syntax, send a POST request with the Go code in JSON format. The response will indicate whether the syntax is valid and, if not, provide the syntax error details.

**Endpoint**: `/check/gosyntax`

**Method**: POST

**Content-Type**: `application/json`

**Request Body**:
- `code` (string): The Go code to be checked for syntax correctness.

**Request Format**:
```json
{
    "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}"
}
```
**Successful Response:**
```json
{
    "valid": true
}
```
**Response When Syntax Error Exists:**
```json
{
    "valid": false,
    "error": "6:29: missing ',' before newline in argument list (and 6 more errors)"
}
```

**Example Request:**
```bash
curl -X POST https://goeval-production.up.railway.app/check/gosyntax \
-H 'X-Api-Key: XXXXXX' \
-H "Content-Type: application/json" \
-d '{"code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}"}'
```

## Authentication
The API uses an API key for authentication. When making requests to the /check/gosyntax endpoint, you need to include an X-Api-Key header with your API key. The header should be in the following format:
```text
X-Api-Key: your_api_key_here
```
