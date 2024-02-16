Test Signer Service
Overview

The Test Signer Service is a RESTful API designed to sign user-submitted answers to a set of questions with a unique signature and verify these signatures at a later time. This service ensures that users' test submissions are authenticated and can be reliably verified.
Getting Started

To run the Test Signer Service on your local machine, follow these steps:
Prerequisites

    Go installed on your machine (version 1.15 or later recommended).

Installation

Clone this repository to your local machine:

bash

git clone <repository-url>

Navigate to the project directory:

bash

cd test-signer

Run the service:

bash

go run cmd/test-signer/main.go

The service will start running on http://localhost:8080.
API Endpoints
Sign Answers

    URL: /sign
    Method: POST
    Description: Accepts a JSON payload containing a user JWT, questions, and their answers. Returns a unique signature for the submission.
    Request Body:

    json

{
  "jwt": "UserJWT",
  "questions": [
    {
      "id": "q1",
      "question": "What is 2+2?",
      "answer": "4"
    },
    {
      "id": "q2",
      "question": "What is the capital of France?",
      "answer": "Paris"
    }
  ]
}

Response:

    Success:

    json

{
  "testSignature": "GeneratedSignature"
}

Error:

json

        {
          "error": "ErrorDescription"
        }

Verify Signature

    URL: /verify
    Method: POST
    Description: Accepts a JSON payload with a user identifier and a signature to verify the authenticity of a previously signed test submission.
    Request Body:

    json

{
  "user": "UsernameOrUserID",
  "testSignature": "SignatureToVerify"
}

Response:

    Success:

    json

{
  "status": "OK",
  "timestamp": "SignatureTimestamp",
  "answers": [
    // The answers associated with the signature
  ]
}

Error:

json

{
  "error": "ErrorDescription"
}



Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.