package main

import (
    "github.com/dgrijalva/jwt-go"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// Question is a struct to hold a question and its answer
type Question struct {
    ID       string `json:"id"`
    Question string `json:"question"`
    Answer   string `json:"answer"`
}

// SignRequest is the expected format of the /sign request body
type SignRequest struct {
    JWT       string      `json:"jwt"`
    Questions []Question  `json:"questions"`
}

// SignResponse is the format of the /sign response body
type SignResponse struct {
    TestSignature string `json:"testSignature,omitempty"`
    Error         string `json:"error,omitempty"`
}

// VerifyRequest is the expected format of the /verify request body
type VerifyRequest struct {
    User          string `json:"user"`
    TestSignature string `json:"testSignature"`
}

// VerifyResponse is the format of the /verify response body
type VerifyResponse struct {
    Status    string     `json:"status,omitempty"`
    Timestamp time.Time  `json:"timestamp,omitempty"`
    Answers   []Question `json:"answers,omitempty"`
    Error     string     `json:"error,omitempty"`
}

// signAnswers handles the /sign endpoint
func signAnswers(w http.ResponseWriter, r *http.Request) {
    var req SignRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Here you would implement JWT verification and answer signing logic

    // For example purposes, we are generating a mock signature
    mockSignature := fmt.Sprintf("signature-%d", time.Now().Unix())

    resp := SignResponse{
        TestSignature: mockSignature,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// verifySignature handles the /verify endpoint
func verifySignature(w http.ResponseWriter, r *http.Request) {
    var req VerifyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Here you would implement signature verification logic

    // For example purposes, we are assuming the signature is valid
    if req.TestSignature == "known-valid-signature" {
        resp := VerifyResponse{
            Status:    "OK",
            Timestamp: time.Now(), // The timestamp would be retrieved from the stored signature data
            Answers:   []Question{}, // The answers would be retrieved from the stored signature data
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    } else {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
    }
}

func main() {
    http.HandleFunc("/sign", signAnswers)
    http.HandleFunc("/verify", verifySignature)

    fmt.Println("Server is running on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}
