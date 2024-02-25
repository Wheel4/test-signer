package main

import (
    "github.com/dgrijalva/jwt-go"
    "encoding/json"
    "os"
    "fmt"
    "net/http"
    "time"
    "strings"
)
const jwtSecretKey = "test-signer"
const dataFile = "submissions.json"

var submissions []Submission

// Structs for handling request and response data
type Question struct {
    ID       string `json:"id"`
    Question string `json:"question"`
    Answer   string `json:"answer"`
}

type SignRequest struct {
    JWT       string     `json:"jwt"`
    Questions []Question `json:"questions"`
}

type SignResponse struct {
    TestSignature string `json:"testSignature,omitempty"`
    Error         string `json:"error,omitempty"`
}

type VerifyRequest struct {
    User          string `json:"user"`
    TestSignature string `json:"testSignature"`
}

type VerifyResponse struct {
    Status    string     `json:"status,omitempty"`
    Timestamp time.Time  `json:"timestamp,omitempty"`
    Answers   []Question `json:"answers,omitempty"`
    Error     string     `json:"error,omitempty"`
}
type Submission struct {
    User      string     `json:"user"`
    Signature string     `json:"signature"`
    Timestamp time.Time  `json:"timestamp"`
    Questions []Question `json:"questions"`
}


// Function to verify JWT from the Authorization header
func verifyJWT(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(jwtSecretKey), nil
    })

    if err != nil {
        return nil, err
    }
    return token, nil
}

func signAnswers(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Authorization header is required", http.StatusUnauthorized)
        return
    }
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")

    _, err := verifyJWT(tokenString)
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    var req SignRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    mockSignature := fmt.Sprintf("signature-%d", time.Now().Unix())

    submissions = append(submissions, Submission{
        User:      "exampleUser", // Extract from JWT in real scenarios
        Signature: mockSignature,
        Timestamp: time.Now(),
        Questions: req.Questions,
    })

    err = saveData(dataFile)
    if err != nil {
        fmt.Println("Error saving data:", err)
        http.Error(w, "Failed to save submission", http.StatusInternalServerError)
        return
    }

    resp := SignResponse{TestSignature: mockSignature}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func verifySignature(w http.ResponseWriter, r *http.Request) {
    var req VerifyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    for _, submission := range submissions {
        if submission.User == req.User && submission.Signature == req.TestSignature {
            resp := VerifyResponse{
                Status:    "OK",
                Timestamp: submission.Timestamp,
                Answers:   submission.Questions,
            }
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(resp)
            return
        }
    }

    http.Error(w, "Invalid signature", http.StatusUnauthorized)
}

func loadData(filename string) error {
    file, err := os.ReadFile(filename)
    if err != nil {
        if os.IsNotExist(err) {
            submissions = []Submission{}
            return nil
        }
        return err
    }

    err = json.Unmarshal(file, &submissions)
    if err != nil {
        return err
    }

    return nil
}

func saveData(filename string) error {
    data, err := json.MarshalIndent(submissions, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(filename, data, 0644)
}

func main() {
    err := loadData(dataFile)
    if err != nil {
        fmt.Println("Error loading data:", err)
        return
    }

    http.HandleFunc("/sign", signAnswers)
    http.HandleFunc("/verify", verifySignature)

    fmt.Println("Server is running on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}
